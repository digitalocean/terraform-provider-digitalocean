package digitalocean

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_database_cluster", &resource.Sweeper{
		Name: "digitalocean_database_cluster",
		F:    testSweepDatabaseCluster,
	})

}

func testSweepDatabaseCluster(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*CombinedConfig).godoClient()

	opt := &godo.ListOptions{PerPage: 200}
	databases, _, err := client.Databases.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, db := range databases {
		if strings.HasPrefix(db.Name, testNamePrefix) {
			log.Printf("Destroying database cluster %s", db.Name)

			if _, err := client.Databases.Delete(context.Background(), db.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanDatabaseCluster_Basic(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
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
						"digitalocean_database_cluster.foobar", "uri"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "private_uri"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "urn"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "tags.#", "1"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "private_network_uuid"),
					testAccCheckDigitalOceanDatabaseClusterURIPassword(
						"digitalocean_database_cluster.foobar", "uri"),
					testAccCheckDigitalOceanDatabaseClusterURIPassword(
						"digitalocean_database_cluster.foobar", "private_uri"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_WithUpdate(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
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
					testAccCheckDigitalOceanDatabaseClusterURIPassword(
						"digitalocean_database_cluster.foobar", "uri"),
					testAccCheckDigitalOceanDatabaseClusterURIPassword(
						"digitalocean_database_cluster.foobar", "private_uri"),
				),
				ExpectError: regexp.MustCompile(`The argument "version" is required, but no definition was found.`),
			},
		},
	})
}

// For backwards compatibility the API allows for POST requests that specify "5"
// for the version, but a Redis 6 cluster is actually created. The response body
// specifies "6" for the version. This should be handled without Terraform
// attempting to recreate the cluster.
func TestAccDigitalOceanDatabaseCluster_oldRedisVersion(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterRedis, databaseName, "5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "engine", "redis"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "version", "6"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_RedisWithEvictionPolicy(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
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
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterRedis, databaseName, "6"),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicyError, databaseName),
				ExpectError: regexp.MustCompile(`eviction_policy is only supported for Redis`),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_TagUpdate(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "tags.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigTagUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_WithVPC(t *testing.T) {
	var database godo.Database
	vpcName := randomTestName()
	databaseName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithVPC, vpcName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "private_network_uuid"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_MongoDBPassword(t *testing.T) {
	var database godo.Database
	databaseName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigMongoDB, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists(
						"digitalocean_database_cluster.foobar", &database),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "password"),
					testAccCheckDigitalOceanDatabaseClusterURIPassword(
						"digitalocean_database_cluster.foobar", "uri"),
					testAccCheckDigitalOceanDatabaseClusterURIPassword(
						"digitalocean_database_cluster.foobar", "private_uri"),
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

// testAccCheckDigitalOceanDatabaseClusterURIPassword checks that the password in
// a database cluster's URI or private URI matches the password value stored in
// its password attribute.
func testAccCheckDigitalOceanDatabaseClusterURIPassword(name string, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		uri, ok := rs.Primary.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("%s not set", attributeName)
		}

		u, err := url.Parse(uri)
		if err != nil {
			return err
		}

		password, ok := u.User.Password()
		if !ok || password == "" {
			return fmt.Errorf("password not set in %s: %s", attributeName, uri)
		}

		return resource.TestCheckResourceAttr(name, "password", password)(s)
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
	version    = "8"
	size       = "db-s-1vcpu-1gb"
	region     = "lon1"
    node_count = 1
    sql_mode   = "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ZERO_DATE,NO_ZERO_IN_DATE"
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithSQLModeUpdate = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "mysql"
	version    = "8"
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

const testAccCheckDigitalOceanDatabaseClusterRedis = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "redis"
	version    = "%s"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
    node_count = 1
	tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithEvictionPolicy = `
resource "digitalocean_database_cluster" "foobar" {
	name            = "%s"
	engine          = "redis"
	version         = "5"
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
	version         = "5"
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
	engine          = "pg"
	version         = "11"
	size            = "db-s-1vcpu-1gb"
	region          = "nyc1"
    node_count      = 1
	eviction_policy = "allkeys_lru"
}
`

const testAccCheckDigitalOceanDatabaseClusterConfigTagUpdate = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
    node_count = 1
	tags       = ["production", "foo"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithVPC = `
resource "digitalocean_vpc" "foobar" {
  name        = "%s"
  region      = "nyc1"
}

resource "digitalocean_database_cluster" "foobar" {
	name                 = "%s"
	engine               = "pg"
	version              = "11"
	size                 = "db-s-1vcpu-1gb"
	region               = "nyc1"
	node_count           = 1
	tags                 = ["production"]
	private_network_uuid = digitalocean_vpc.foobar.id
}`

const testAccCheckDigitalOceanDatabaseClusterConfigMongoDB = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "mongodb"
	version    = "4"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
    node_count = 1
}`
