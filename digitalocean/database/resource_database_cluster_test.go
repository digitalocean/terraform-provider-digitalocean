package database_test

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseCluster_Basic(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "project_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "storage_size_mib"),
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
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "size", "db-s-1vcpu-2gb"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName),
				Check: resource.TestCheckFunc(
					func(s *terraform.State) error {
						time.Sleep(30 * time.Second)
						return nil
					},
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithUpdate, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "size", "db-s-2vcpu-4gb"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_WithAdditionalStorage(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "storage_size_mib", "30720"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithAdditionalStorage, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "storage_size_mib", "61440"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_WithMigration(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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

// DigitalOcean only supports one version of Redis. For backwards compatibility
// the API allows for POST requests that specifies a previous version, but new
// clusters are created with the latest/only supported version, regardless of
// the version specified in the config.
// The provider suppresses diffs when the config version is <= to the latest
// version. New clusters is always created with the latest version .
func TestAccDigitalOceanDatabaseCluster_oldRedisVersion(t *testing.T) {
	var (
		database godo.Database
	)

	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_cluster.foobar", "version"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_RedisWithEvictionPolicy(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	vpcName := acceptance.RandomTestName()
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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

func TestAccDigitalOceanDatabaseCluster_WithBackupRestore(t *testing.T) {
	var originalDatabase godo.Database
	var backupDatabase godo.Database

	originalDatabaseName := acceptance.RandomTestName()
	backupDatabasename := acceptance.RandomTestName()

	originalDatabaseConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, originalDatabaseName)
	backUpRestoreConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithBackupRestore, backupDatabasename, originalDatabaseName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: originalDatabaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &originalDatabase),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&originalDatabase, originalDatabaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "region", "nyc1"),
					func(s *terraform.State) error {
						err := waitForDatabaseBackups(originalDatabaseName)
						return err
					},
				),
			},
			{
				Config: originalDatabaseConfig + backUpRestoreConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar_backup", &backupDatabase),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&backupDatabase, backupDatabasename),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar_backup", "region", "nyc1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_MongoDBPassword(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
			// Pause before running CheckDestroy
			{
				Config: " ",
				Check: resource.TestCheckFunc(
					func(s *terraform.State) error {
						time.Sleep(30 * time.Second)
						return nil
					},
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_Upgrade(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()
	previousPGVersion := "14"
	latestPGVersion := "15"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				// TODO: Hardcoding the versions here is not ideal.
				// We will need to determine a better way to fetch the last and latest versions dynamically.
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigCustomVersion, databaseName, "pg", previousPGVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists(
						"digitalocean_database_cluster.foobar", &database),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "engine", "pg"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "version", previousPGVersion),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigCustomVersion, databaseName, "pg", latestPGVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "version", latestPGVersion),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseCluster_nonDefaultProject(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()
	projectName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigNonDefaultProject, projectName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttrPair(
						"digitalocean_project.foobar", "id", "digitalocean_database_cluster.foobar", "project_id"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseClusterDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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

func waitForDatabaseBackups(originalDatabaseName string) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	var (
		tickerInterval = 10 * time.Second
		timeoutSeconds = 300.0
		timeout        = int(timeoutSeconds / tickerInterval.Seconds())
		n              = 0
		ticker         = time.NewTicker(tickerInterval)
	)

	databases, _, err := client.Databases.List(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error retrieving backups from original cluster")
	}

	// gets original database's ID
	var originalDatabaseID string
	for _, db := range databases {
		if db.Name == originalDatabaseName {
			originalDatabaseID = db.ID
		}
	}

	if originalDatabaseID == "" {
		return fmt.Errorf("Error retrieving backups from cluster")
	}

	for range ticker.C {
		backups, resp, err := client.Databases.ListBackups(context.Background(), originalDatabaseID, nil)
		if resp.StatusCode == 412 {
			continue
		}

		if err != nil {
			ticker.Stop()
			return fmt.Errorf("Error retrieving backups from cluster")
		}

		if len(backups) >= 1 {
			ticker.Stop()
			return nil
		}

		if n > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return fmt.Errorf("Timeout waiting for database cluster to have a backup to be restored from")
}

const testAccCheckDigitalOceanDatabaseClusterConfigBasic = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithBackupRestore = `
resource "digitalocean_database_cluster" "foobar_backup" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]

  backup_restore {
    database_name = "%s"
  }
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithUpdate = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-2vcpu-4gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithAdditionalStorage = `
resource "digitalocean_database_cluster" "foobar" {
  name             = "%s"
  engine           = "pg"
  version          = "15"
  size             = "db-s-1vcpu-2gb"
  region           = "nyc1"
  node_count       = 1
  tags             = ["production"]
  storage_size_mib = 61440
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithMigration = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "lon1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithMaintWindow = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
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

const testAccCheckDigitalOceanDatabaseClusterKafka = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "kafka"
  version    = "%s"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 3
  tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseClusterMySQL = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
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
  version         = "15"
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
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production", "foo"]
}`

const testAccCheckDigitalOceanDatabaseClusterConfigWithVPC = `
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "nyc1"
}

resource "digitalocean_database_cluster" "foobar" {
  name                 = "%s"
  engine               = "pg"
  version              = "15"
  size                 = "db-s-1vcpu-2gb"
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
  region     = "nyc3"
  node_count = 1
}`

const testAccCheckDigitalOceanDatabaseClusterConfigCustomVersion = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "%s"
  version    = "%s"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc3"
  node_count = 1
}`

const testAccCheckDigitalOceanDatabaseClusterConfigNonDefaultProject = `
resource "digitalocean_project" "foobar" {
  name = "%s"
}

resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  project_id = digitalocean_project.foobar.id
}`
