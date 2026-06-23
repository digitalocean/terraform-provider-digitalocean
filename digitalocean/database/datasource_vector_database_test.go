package database_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanVectorDatabase_ByName(t *testing.T) {
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVectorDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanVectorDatabaseConfigByName, vectorName, os.Getenv("DIGITALOCEAN_PROJECT_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_vector_database.foobar", "name", vectorName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vector_database.foobar", "region", "tor1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vector_database.foobar", "size", "small"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vector_database.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vector_database.foobar", "status"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVectorDatabase_ByID(t *testing.T) {
	vectorName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVectorDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanVectorDatabaseConfigByID, vectorName, os.Getenv("DIGITALOCEAN_PROJECT_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_vector_database.foobar", "name", vectorName),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vector_database.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vector_database.foobar", "status"),
				),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanVectorDatabaseConfigByName = `
resource "digitalocean_vector_database" "foobar" {
  name       = "%[1]s"
  region     = "tor1"
  size       = "small"
  project_id = "%[2]s"
}

data "digitalocean_vector_database" "foobar" {
  name = digitalocean_vector_database.foobar.name
}`

const testAccCheckDataSourceDigitalOceanVectorDatabaseConfigByID = `
resource "digitalocean_vector_database" "foobar" {
  name       = "%[1]s"
  region     = "tor1"
  size       = "small"
  project_id = "%[2]s"
}

data "digitalocean_vector_database" "foobar" {
  id = digitalocean_vector_database.foobar.id
}`
