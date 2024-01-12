package database_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseUser_Basic(t *testing.T) {
	var databaseUser godo.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()
	databaseUserNameUpdated := databaseUserName + "-up"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigBasic, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &databaseUser),
					testAccCheckDigitalOceanDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "password"),
					resource.TestCheckNoResourceAttr(
						"digitalocean_database_user.foobar_user", "mysql_auth_plugin"),
					resource.TestCheckNoResourceAttr(
						"digitalocean_database_user.foobar_user", "access_cert"),
					resource.TestCheckNoResourceAttr(
						"digitalocean_database_user.foobar_user", "access_key"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigBasic, databaseClusterName, databaseUserNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &databaseUser),
					testAccCheckDigitalOceanDatabaseUserNotExists("digitalocean_database_user.foobar_user", databaseUserName),
					testAccCheckDigitalOceanDatabaseUserAttributes(&databaseUser, databaseUserNameUpdated),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", databaseUserNameUpdated),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseUser_MongoDB(t *testing.T) {
	var databaseUser godo.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigMongo, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &databaseUser),
					testAccCheckDigitalOceanDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "password"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseUser_MongoDBMultiUser(t *testing.T) {
	databaseClusterName := acceptance.RandomTestName()
	users := []string{"foo", "bar", "baz", "one", "two"}
	config := fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigMongoMultiUser,
		databaseClusterName,
		users[0], users[0],
		users[1], users[1],
		users[2], users[2],
		users[3], users[3],
		users[4], users[4],
	)
	userResourceNames := make(map[string]string, len(users))
	for _, u := range users {
		userResourceNames[u] = fmt.Sprintf("digitalocean_database_user.%s", u)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						userResourceNames[users[0]], "name", users[0]),
					resource.TestCheckResourceAttr(
						userResourceNames[users[1]], "name", users[1]),
					resource.TestCheckResourceAttr(
						userResourceNames[users[2]], "name", users[2]),
					resource.TestCheckResourceAttr(
						userResourceNames[users[3]], "name", users[3]),
					resource.TestCheckResourceAttr(
						userResourceNames[users[4]], "name", users[4]),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseUser_MySQLAuth(t *testing.T) {
	var databaseUser godo.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigMySQLAuth, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &databaseUser),
					testAccCheckDigitalOceanDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "mysql_auth_plugin", "mysql_native_password"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigMySQLAuthUpdate, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &databaseUser),
					testAccCheckDigitalOceanDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "mysql_auth_plugin", "caching_sha2_password"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigMySQLAuthRemoved, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &databaseUser),
					testAccCheckDigitalOceanDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "mysql_auth_plugin", "caching_sha2_password"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseUser_KafkaACLs(t *testing.T) {
	var databaseUser godo.DatabaseUser
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigKafkaACL, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &databaseUser),
					testAccCheckDigitalOceanDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "access_cert"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "access_key"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "settings.0.acl.0.id"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.0.topic", "topic-1"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.0.permission", "admin"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "settings.0.acl.1.id"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.1.topic", "topic-2"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.1.permission", "produceconsume"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "settings.0.acl.2.id"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.2.topic", "topic-*"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.2.permission", "produce"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "settings.0.acl.3.id"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.3.topic", "topic-*"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.3.permission", "consume"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigKafkaACLUpdate, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &databaseUser),
					testAccCheckDigitalOceanDatabaseUserAttributes(&databaseUser, databaseUserName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", databaseUserName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "password"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "access_cert"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "access_key"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "settings.0.acl.0.id"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.0.topic", "topic-1"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "settings.0.acl.0.permission", "produceconsume"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseUserDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_user" {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		// Try to find the database
		_, _, err := client.Databases.GetUser(context.Background(), clusterID, name)

		if err == nil {
			return fmt.Errorf("Database User still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanDatabaseUserExists(n string, databaseUser *godo.DatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database User ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		foundDatabaseUser, _, err := client.Databases.GetUser(context.Background(), clusterID, name)

		if err != nil {
			return err
		}

		if foundDatabaseUser.Name != name {
			return fmt.Errorf("Database user not found")
		}

		*databaseUser = *foundDatabaseUser

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseUserNotExists(n string, databaseUserName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Database User ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		clusterID := rs.Primary.Attributes["cluster_id"]

		_, resp, err := client.Databases.GetDB(context.Background(), clusterID, databaseUserName)

		if err != nil && resp.StatusCode != http.StatusNotFound {
			return err
		}

		if err == nil {
			return fmt.Errorf("Database User %s still exists", databaseUserName)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseUserAttributes(databaseUser *godo.DatabaseUser, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if databaseUser.Name != name {
			return fmt.Errorf("Bad name: %s", databaseUser.Name)
		}

		return nil
	}
}

const testAccCheckDigitalOceanDatabaseUserConfigBasic = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1

  maintenance_window {
    day  = "friday"
    hour = "13:00:00"
  }
}

resource "digitalocean_database_user" "foobar_user" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}`

const testAccCheckDigitalOceanDatabaseUserConfigMongo = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mongodb"
  version    = "4"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1

  maintenance_window {
    day  = "friday"
    hour = "13:00:00"
  }
}

resource "digitalocean_database_user" "foobar_user" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}`

const testAccCheckDigitalOceanDatabaseUserConfigMongoMultiUser = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mongodb"
  version    = "4"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_user" "%s" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}

resource "digitalocean_database_user" "%s" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}

resource "digitalocean_database_user" "%s" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}

resource "digitalocean_database_user" "%s" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}

resource "digitalocean_database_user" "%s" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}`

const testAccCheckDigitalOceanDatabaseUserConfigMySQLAuth = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_user" "foobar_user" {
  cluster_id        = digitalocean_database_cluster.foobar.id
  name              = "%s"
  mysql_auth_plugin = "mysql_native_password"
}`

const testAccCheckDigitalOceanDatabaseUserConfigKafkaACL = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "kafka"
  version    = "3.5"
  size       = "db-s-2vcpu-2gb"
  region     = "nyc1"
  node_count = 3
}

resource "digitalocean_database_user" "foobar_user" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
  settings {
    acl {
      topic      = "topic-1"
      permission = "admin"
    }
    acl {
      topic      = "topic-2"
      permission = "produceconsume"
    }
    acl {
      topic      = "topic-*"
      permission = "produce"
    }
    acl {
      topic      = "topic-*"
      permission = "consume"
    }
  }
}`

const testAccCheckDigitalOceanDatabaseUserConfigKafkaACLUpdate = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "kafka"
  version    = "3.5"
  size       = "db-s-2vcpu-2gb"
  region     = "nyc1"
  node_count = 3
}

resource "digitalocean_database_user" "foobar_user" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
  settings {
    acl {
      topic      = "topic-1"
      permission = "produceconsume"
    }
  }
}`

const testAccCheckDigitalOceanDatabaseUserConfigMySQLAuthUpdate = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_user" "foobar_user" {
  cluster_id        = digitalocean_database_cluster.foobar.id
  name              = "%s"
  mysql_auth_plugin = "caching_sha2_password"
}`

const testAccCheckDigitalOceanDatabaseUserConfigMySQLAuthRemoved = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_user" "foobar_user" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}`
