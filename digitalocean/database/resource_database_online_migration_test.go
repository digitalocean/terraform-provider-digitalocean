package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseOnlineMigration_Basic(t *testing.T) {
	source := "source-" + acceptance.RandomTestName()
	destination := "destination-" + acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		//CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseOnlineMigrationBasic, source, "8", destination, "8"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("digitalocean_database_online_migration.foobar", "id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_online_migration.foobar", "status"),
					resource.TestCheckResourceAttrSet("digitalocean_database_online_migration.foobar", "created_at"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseOnlineMigrationBasic = `
resource "digitalocean_database_cluster" "source" {
	name       = "%s"
	engine     = "mysql"
	version    = "%s"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
	node_count = 1
	tags       = ["production"]
}

resource "digitalocean_database_cluster" "destination" {
	name       = "%s"
	engine     = "mysql"
	version    = "%s"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
	node_count = 1
	tags       = ["production"]
}

resource "digitalocean_database_db" "source_db" {
	cluster_id = digitalocean_database_cluster.source.id
	name       = "terraform-db-om-source"
}

resource "digitalocean_database_online_migration" "foobar" {
	cluster_id = digitalocean_database_cluster.destination.id
	source {
		host = digitalocean_database_cluster.source.host
		db_name = digitalocean_database_db.source_db.name
		port = digitalocean_database_cluster.source.port
		username = digitalocean_database_cluster.source.user
		password = digitalocean_database_cluster.source.password
  	}
  depends_on = [digitalocean_database_cluster.destination, digitalocean_database_cluster.source, digitalocean_database_db.source_db]
}`
