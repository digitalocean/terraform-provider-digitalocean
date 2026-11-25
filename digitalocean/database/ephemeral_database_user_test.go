package database_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseUser_Ephemeral(t *testing.T) {
	databaseClusterName := acceptance.RandomTestName()
	databaseUserName := acceptance.RandomTestName()
	var databaseUserPassword string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acceptance.TestAccProtoV5ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"postgresql": {
				Source:            "cyrilgdn/postgresql",
				VersionConstraint: "1.22.0",
			},
		},
		CheckDestroy: testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDigitalOceanDatabaseUserEphemeralClusterConfig, databaseClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCreateDatabaseUser("digitalocean_database_cluster.foobar", databaseUserName, &databaseUserPassword),
				),
			},
			{
				Config: fmt.Sprintf(testAccDigitalOceanDatabaseUserEphemeralConfig, databaseClusterName, databaseUserName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("data.postgresql_schemas.app", "schemas.*", "public"),
					testAccCheckStateDoesNotContain(&databaseUserPassword),
				),
			},
		},
	})
}

func testAccCreateDatabaseUser(clusterResourceName, userName string, password *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[clusterResourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", clusterResourceName)
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		user, _, err := client.Databases.CreateUser(context.Background(), rs.Primary.ID, &godo.DatabaseCreateUserRequest{Name: userName})
		if err != nil {
			return fmt.Errorf("creating database user %q: %w", userName, err)
		}
		if user.Password == "" {
			return fmt.Errorf("API returned an empty password for database user %q", userName)
		}

		*password = user.Password
		return nil
	}
}

func testAccCheckStateDoesNotContain(secret *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *secret == "" {
			return fmt.Errorf("sensitive value was not captured before the state check")
		}

		for resourceName, rs := range s.RootModule().Resources {
			for attribute, value := range rs.Primary.Attributes {
				if strings.Contains(value, *secret) {
					return fmt.Errorf("sensitive value found in state at %s.%s", resourceName, attribute)
				}
			}
		}

		return nil
	}
}

const testAccDigitalOceanDatabaseUserEphemeralClusterConfig = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
`

const testAccDigitalOceanDatabaseUserEphemeralConfig = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

ephemeral "digitalocean_database_user" "app" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}

provider "postgresql" {
  alias    = "app"
  host     = digitalocean_database_cluster.foobar.host
  port     = digitalocean_database_cluster.foobar.port
  database = digitalocean_database_cluster.foobar.database
  username = ephemeral.digitalocean_database_user.app.name
  password = ephemeral.digitalocean_database_user.app.password
  sslmode  = "require"
}

data "postgresql_schemas" "app" {
  provider = postgresql.app
  database = digitalocean_database_cluster.foobar.database
}
`
