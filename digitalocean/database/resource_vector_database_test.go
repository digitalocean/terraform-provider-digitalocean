package database_test

import (
	"context"
	"fmt"
	"testing"

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
				Config: fmt.Sprintf(testAccCheckDigitalOceanVectorDatabaseConfigBasic, vectorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVectorDatabaseExists("digitalocean_vector_database.foobar", &vectorDB),
					testAccCheckDigitalOceanVectorDatabaseAttributes(&vectorDB, vectorName),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "name", vectorName),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "region", "tor1"),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "size", "db-s-1vcpu-1gb"),
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
				Config: fmt.Sprintf(testAccCheckDigitalOceanVectorDatabaseConfigBasic, vectorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVectorDatabaseExists("digitalocean_vector_database.foobar", &vectorDB),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "size", "db-s-1vcpu-1gb"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVectorDatabaseConfigResize, vectorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVectorDatabaseExists("digitalocean_vector_database.foobar", &vectorDB),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "size", "db-s-2vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "config.0.enable_auto_schema", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_vector_database.foobar", "config.0.default_quantization", "none"),
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
				Config: fmt.Sprintf(testAccCheckDigitalOceanVectorDatabaseConfigBasic, vectorName),
			},
			{
				ResourceName:      "digitalocean_vector_database.foobar",
				ImportState:       true,
				ImportStateVerify: true,
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

		// Try to find the vector database
		_, _, err := client.VectorDBs.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Vector database still exists")
		}
	}

	return nil
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
  name   = "%s"
  region = "tor1"
  size   = "db-s-1vcpu-1gb"
  tags   = ["production"]
}`

const testAccCheckDigitalOceanVectorDatabaseConfigResize = `
resource "digitalocean_vector_database" "foobar" {
  name   = "%s"
  region = "tor1"
  size   = "db-s-2vcpu-2gb"
  tags   = ["production"]

  config {
    enable_auto_schema   = true
    default_quantization = "none"
  }
}`
