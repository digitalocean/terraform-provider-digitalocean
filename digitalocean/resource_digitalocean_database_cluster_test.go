package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanDatabaseCluster_Basic(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

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
						"digitalocean_database_cluster.foobar", "private_host"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "port"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "user"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "password"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "urn"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "tags.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_WithUpdate(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

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
	databaseName := randomTestName()

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
	databaseName := randomTestName()

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

func TestAccDigitalOceanDatabaseCluster_WithSQLMode(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithSQLMode, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "sql_mode",
						"ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ZERO_DATE,NO_ZERO_IN_DATE"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithSQLModeUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "sql_mode",
						"ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_CheckSQLModeSupport(t *testing.T) {
	databaseName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithRedisSQLModeError, databaseName),
				ExpectError: regexp.MustCompile(`sql_mode is only supported for MySQL`),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_RedisNoVersion(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterRedisNoVersion, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "engine", "redis"),
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
			// Add eviction policy when not initially set
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicyUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "eviction_policy", "allkeys_lru"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_RedisWithEvictionPolicy(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			// Create with an eviction policy
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicy, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "eviction_policy", "volatile_random"),
				),
			},
			// Update eviction policy
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicyUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "eviction_policy", "allkeys_lru"),
				),
			},
			// Remove eviction policy
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterRedisNoVersion, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_CheckEvictionPolicySupport(t *testing.T) {
	databaseName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicyError, databaseName),
				ExpectError: regexp.MustCompile(`eviction_policy is only supported for Redis`),
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
	tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithUpdate = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-2gb"
	region     = "nyc1"
    node_count = 1
	tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithMigration = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "lon1"
    node_count = 1
	tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithMaintWindow = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
	node_count = 1
	tags       = ["production"]

	maintenance_window {
        day  = "friday"
        hour = "13:00:00"
	}
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithSQLMode = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "mysql"
	size       = "db-s-1vcpu-1gb"
	region     = "lon1"
    node_count = 1
    sql_mode   = "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ZERO_DATE,NO_ZERO_IN_DATE"
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithSQLModeUpdate = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "mysql"
	size       = "db-s-1vcpu-1gb"
	region     = "lon1"
    node_count = 1
    sql_mode   = "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE"
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithRedisSQLModeError = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "redis"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
    node_count = 1
    sql_mode   = "ANSI"
}`

const testAccCheckDigitalOceanDatabaseClusterRedisNoVersion = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "redis"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
    node_count = 1
	tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicy = `
resource "digitalocean_database_cluster" "foobar" {
	name            = "%s"
	engine          = "redis"
	size            = "db-s-1vcpu-1gb"
	region          = "nyc1"
    node_count      = 1
	tags            = ["production"]
	eviction_policy = "volatile_random"
}
`

const testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicyUpdate = `
resource "digitalocean_database_cluster" "foobar" {
	name            = "%s"
	engine          = "redis"
	size            = "db-s-1vcpu-1gb"
	region          = "nyc1"
    node_count      = 1
	tags            = ["production"]
	eviction_policy = "allkeys_lru"
}
`

const testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicyError = `
resource "digitalocean_database_cluster" "foobar" {
	name            = "%s"
	engine          = "psql"
	size            = "db-s-1vcpu-1gb"
	region          = "nyc1"
    node_count      = 1
	eviction_policy = "allkeys_lru"
}
`
