package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDigitalOceanDatabaseCluster_Basic(t *testing.T) {
	var database godo.Database
	databaseName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "engine", "pg"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "host"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "port"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "user"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "password"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "urn"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_WithUpdate(t *testing.T) {
	var database godo.Database
	databaseName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "size", "db-s-1vcpu-1gb"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "size", "db-s-1vcpu-2gb"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_WithMigration(t *testing.T) {
	var database godo.Database
	databaseName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "region", "nyc1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithMigration, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "region", "lon1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_WithMaintWindow(t *testing.T) {
	var database godo.Database
	databaseName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithMaintWindow, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "maintenance_window.0.day"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "maintenance_window.0.hour"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_cluster" {
			continue
		}

		// Try to find the database
		_, _, err := client.Databases.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("DatabaseCluster still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanDatabaseClusterAttributes(database *godo.Database, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if database.Name != name {
			return fmt.Errorf("Bad name: %s", database.Name)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseClusterExists(n string, database *godo.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DatabaseCluster ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundDatabaseCluster, _, err := client.Databases.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDatabaseCluster.ID != rs.Primary.ID {
			return fmt.Errorf("DatabaseCluster not found")
		}

		*database = *foundDatabaseCluster

		return nil
	}
}

const testAccCheckDigitalOceanDatabaseClusterConfigBasic = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
    node_count = 1
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithUpdate = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-2gb"
	region     = "nyc1"
    node_count = 1
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithMigration = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "lon1"
    node_count = 1
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithMaintWindow = `
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
}`
