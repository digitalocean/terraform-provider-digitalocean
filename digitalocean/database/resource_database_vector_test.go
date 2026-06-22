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

func TestAccDigitalOceanDatabaseVector_Basic(t *testing.T) {
	var vectorDB godo.VectorDB
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseVectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseVectorConfigBasic, vectorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseVectorExists("digitalocean_database_vector.foobar", &vectorDB),
					testAccCheckDigitalOceanDatabaseVectorAttributes(&vectorDB, vectorName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_vector.foobar", "name", vectorName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_vector.foobar", "region", "tor1"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_vector.foobar", "size", "db-s-1vcpu-1gb"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_vector.foobar", "status"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_vector.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_vector.foobar", "tags.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseVector_Update(t *testing.T) {
	var vectorDB godo.VectorDB
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseVectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseVectorConfigBasic, vectorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseVectorExists("digitalocean_database_vector.foobar", &vectorDB),
					resource.TestCheckResourceAttr(
						"digitalocean_database_vector.foobar", "size", "db-s-1vcpu-1gb"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseVectorConfigResize, vectorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseVectorExists("digitalocean_database_vector.foobar", &vectorDB),
					resource.TestCheckResourceAttr(
						"digitalocean_database_vector.foobar", "size", "db-s-2vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_vector.foobar", "config.0.enable_auto_schema", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_vector.foobar", "config.0.default_quantization", "none"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseVector_Import(t *testing.T) {
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseVectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseVectorConfigBasic, vectorName),
			},
			{
				ResourceName:      "digitalocean_database_vector.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseVectorDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_vector" {
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

func testAccCheckDigitalOceanDatabaseVectorExists(n string, vectorDB *godo.VectorDB) resource.TestCheckFunc {
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

func testAccCheckDigitalOceanDatabaseVectorAttributes(vectorDB *godo.VectorDB, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vectorDB.Name != name {
			return fmt.Errorf("Bad name: %s", vectorDB.Name)
		}

		return nil
	}
}

const testAccCheckDigitalOceanDatabaseVectorConfigBasic = `
resource "digitalocean_database_vector" "foobar" {
  name   = "%s"
  region = "tor1"
  size   = "db-s-1vcpu-1gb"
  tags   = ["production"]
}`

const testAccCheckDigitalOceanDatabaseVectorConfigResize = `
resource "digitalocean_database_vector" "foobar" {
  name   = "%s"
  region = "tor1"
  size   = "db-s-2vcpu-2gb"
  tags   = ["production"]

  config {
    enable_auto_schema   = true
    default_quantization = "none"
  }
}`
