package database_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanVectorDatabase_Basic(t *testing.T) {
	var vectorDB godo.VectorDB
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVectorDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVectorDatabaseConfigBasic, vectorName, os.Getenv("DIGITALOCEAN_PROJECT_ID")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVectorDatabaseExists("digitalocean_vector_database.foobar", &vectorDB),
					testAccCheckDigitalOceanVectorDatabaseAttributes(&vectorDB, vectorName),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "name", vectorName),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "region", "tor1"),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "size", "small"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vector_database.foobar", "status"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vector_database.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "tags.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanVectorDatabase_Update(t *testing.T) {
	var vectorDB godo.VectorDB
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVectorDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVectorDatabaseConfigBasic, vectorName, os.Getenv("DIGITALOCEAN_PROJECT_ID")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVectorDatabaseExists("digitalocean_vector_database.foobar", &vectorDB),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "size", "small"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVectorDatabaseConfigResize, vectorName, os.Getenv("DIGITALOCEAN_PROJECT_ID")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVectorDatabaseExists("digitalocean_vector_database.foobar", &vectorDB),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "size", "medium"),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "config.0.enable_auto_schema", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "config.0.default_quantization", "pq"),
				),
			},
		},
	})
}

func TestAccDigitalOceanVectorDatabase_Import(t *testing.T) {
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVectorDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVectorDatabaseConfigBasic, vectorName, os.Getenv("DIGITALOCEAN_PROJECT_ID")),
			},
			{
				ResourceName:      "digitalocean_vector_database.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				// project_id is write-only on the API (not returned on read),
				// so it cannot be verified against imported state.
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func testAccCheckDigitalOceanVectorDatabaseDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_vector_database" {
			continue
		}

		// A vector database that is still in the "deleting" state continues to
		// return a 200 from Get, so poll until the API reports it gone (404)
		// rather than treating any successful response as a failure.
		if err := waitForVectorDatabaseDestroyed(client, rs.Primary.ID); err != nil {
			return err
		}
	}

	return nil
}

func waitForVectorDatabaseDestroyed(client *godo.Client, id string) error {
	const (
		tickerInterval = 10 * time.Second
		maxAttempts    = 60
	)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		_, resp, err := client.VectorDBs.Get(context.Background(), id)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("error checking vector database destroy status: %s", err)
		}

		time.Sleep(tickerInterval)
	}

	return fmt.Errorf("Vector database still exists")
}

func testAccCheckDigitalOceanVectorDatabaseExists(n string, vectorDB *godo.VectorDB) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No vector database ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundVectorDB, _, err := client.VectorDBs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundVectorDB.ID != rs.Primary.ID {
			return fmt.Errorf("Vector database not found")
		}

		*vectorDB = *foundVectorDB

		return nil
	}
}

func testAccCheckDigitalOceanVectorDatabaseAttributes(vectorDB *godo.VectorDB, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vectorDB.Name != name {
			return fmt.Errorf("Bad name: %s", vectorDB.Name)
		}

		return nil
	}
}

const testAccCheckDigitalOceanVectorDatabaseConfigBasic = `
resource "digitalocean_vector_database" "foobar" {
  name       = "%[1]s"
  region     = "tor1"
  size       = "small"
  project_id = "%[2]s"
  tags       = ["production"]
}`

const testAccCheckDigitalOceanVectorDatabaseConfigResize = `
resource "digitalocean_vector_database" "foobar" {
  name       = "%[1]s"
  region     = "tor1"
  size       = "medium"
  project_id = "%[2]s"
  tags       = ["production"]

  config {
    enable_auto_schema   = true
    default_quantization = "pq"
  }
}`
