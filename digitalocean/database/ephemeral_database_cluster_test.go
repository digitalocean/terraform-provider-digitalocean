package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseCluster_Ephemeral(t *testing.T) {
	databaseClusterName := acceptance.RandomTestName()

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
				Config: fmt.Sprintf(testAccDigitalOceanDatabaseClusterEphemeralConfig, databaseClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("postgresql_schema.app", "name", "app_schema"),
				),
			},
		},
	})
}

const testAccDigitalOceanDatabaseClusterEphemeralConfig = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

ephemeral "digitalocean_database_cluster" "foobar" {
  id = digitalocean_database_cluster.foobar.id
}

provider "postgresql" {
  host     = ephemeral.digitalocean_database_cluster.foobar.host
  port     = ephemeral.digitalocean_database_cluster.foobar.port
  database = ephemeral.digitalocean_database_cluster.foobar.database
  username = ephemeral.digitalocean_database_cluster.foobar.user
  password = ephemeral.digitalocean_database_cluster.foobar.password
  sslmode  = "require"
}

resource "postgresql_schema" "app" {
  database = digitalocean_database_cluster.foobar.database
  name     = "app_schema"
}
`
