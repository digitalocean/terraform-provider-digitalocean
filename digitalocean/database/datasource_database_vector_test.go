package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDatabaseVector_ByName(t *testing.T) {
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseVectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanDatabaseVectorConfigByName, vectorName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_vector.foobar", "name", vectorName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_vector.foobar", "region", "tor1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_vector.foobar", "size", "db-s-1vcpu-1gb"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_vector.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_vector.foobar", "status"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDatabaseVector_ByID(t *testing.T) {
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseVectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanDatabaseVectorConfigByID, vectorName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_vector.foobar", "name", vectorName),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_vector.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_vector.foobar", "status"),
				),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanDatabaseVectorConfigByName = `
resource "digitalocean_database_vector" "foobar" {
  name   = "%s"
  region = "tor1"
  size   = "db-s-1vcpu-1gb"
}

data "digitalocean_database_vector" "foobar" {
  name = digitalocean_database_vector.foobar.name
}`

const testAccCheckDataSourceDigitalOceanDatabaseVectorConfigByID = `
resource "digitalocean_database_vector" "foobar" {
  name   = "%s"
  region = "tor1"
  size   = "db-s-1vcpu-1gb"
}

data "digitalocean_database_vector" "foobar" {
  id = digitalocean_database_vector.foobar.id
}`
