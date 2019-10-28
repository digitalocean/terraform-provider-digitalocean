package digitalocean

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanDatabase_Basic(t *testing.T) {
	var database godo.DatabaseDB
	databaseClusterName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))
	databaseName := fmt.Sprintf("foobar-test-db-terraform-%s", acctest.RandString(10))
	databaseNameUpdated := databaseName + "-up"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseConfigBasic, databaseClusterName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseExists("digitalocean_database.foobar_db", &database),
					testAccCheckDigitalOceanDatabaseAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database.foobar_db", "name", databaseName),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseConfigBasic, databaseClusterName, databaseNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseExists("digitalocean_database.foobar_db", &database),
					testAccCheckDigitalOceanDatabaseNotExists("digitalocean_database.foobar_db", databaseName),
					testAccCheckDigitalOceanDatabaseAttributes(&database, databaseNameUpdated),
					resource.TestCheckResourceAttr(
						"digitalocean_database.foobar_db", "name", databaseNameUpdated),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database" {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		// Try to find the database
		_, _, err := client.Databases.GetReplica(context.Background(), clusterID, name)

		if err == nil {
			return fmt.Errorf("Database still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanDatabaseExists(n string, database *godo.DatabaseDB) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		foundDatabase, _, err := client.Databases.GetDB(context.Background(), clusterID, name)

		if err != nil {
			return err
		}

		if foundDatabase.Name != name {
			return fmt.Errorf("Database not found")
		}

		*database = *foundDatabase

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseNotExists(n string, databaseName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()
		clusterID := rs.Primary.Attributes["cluster_id"]

		_, resp, err := client.Databases.GetDB(context.Background(), clusterID, databaseName)

		if err != nil && resp.StatusCode != http.StatusNotFound {
			return err
		}

		if err == nil {
			return fmt.Errorf("Database %s still exists", databaseName)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseAttributes(database *godo.DatabaseDB, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if database.Name != name {
			return fmt.Errorf("Bad name: %s", database.Name)
		}

		return nil
	}
}

const testAccCheckDigitalOceanDatabaseConfigBasic = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
	node_count = 1

	maintenance_window {
        day  = "friday"
        hour = "13:00:00"
	}
}

resource "digitalocean_database" "foobar_db" {
  cluster_id = "${digitalocean_database_cluster.foobar.id}"
  name       = "%s"
}`
