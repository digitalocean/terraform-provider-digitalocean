package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseFirewall_Basic(t *testing.T) {
	databaseClusterName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseFirewallConfigBasic, databaseClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_firewall.example", "rule.#", "1"),
				),
			},
			// Add a new rule
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseFirewallConfigAddRule, databaseClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_firewall.example", "rule.#", "2"),
				),
			},
			// Remove an existing rule
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseFirewallConfigBasic, databaseClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_firewall.example", "rule.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseFirewall_MultipleResourceTypes(t *testing.T) {
	dbName := randomTestName()
	dropletName := randomTestName()
	tagName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseFirewallConfigMultipleResourceTypes,
					dbName, dropletName, tagName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_firewall.example", "rule.#", "3"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseFirewallDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_firewall" {
			continue
		}

		clusterId := rs.Primary.Attributes["cluster_id"]

		_, _, err := client.Databases.GetFirewallRules(context.Background(), clusterId)
		if err == nil {
			return fmt.Errorf("DatabaseFirewall still exists")
		}
	}

	return nil
}

const testAccCheckDigitalOceanDatabaseFirewallConfigBasic = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
	node_count = 1
}

resource "digitalocean_database_firewall" "example" {
	cluster_id = digitalocean_database_cluster.foobar.id

	rule {
		type  = "ip_addr"
		value = "192.168.1.1"
	}
}
`

const testAccCheckDigitalOceanDatabaseFirewallConfigAddRule = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
	node_count = 1
}

resource "digitalocean_database_firewall" "example" {
	cluster_id = digitalocean_database_cluster.foobar.id

	rule {
		type  = "ip_addr"
		value = "192.168.1.1"
	}

	rule {
		type  = "ip_addr"
		value = "192.0.2.0"
	}
}
`

const testAccCheckDigitalOceanDatabaseFirewallConfigMultipleResourceTypes = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
	node_count = 1
}

resource "digitalocean_droplet" "foobar" {
	name      = "%s"
	size      = "s-1vcpu-1gb"
	image     = "centos-7-x64"
	region    = "nyc3"
}

resource "digitalocean_tag" "foobar" {
	name = "%s"
}

resource "digitalocean_database_firewall" "example" {
	cluster_id = digitalocean_database_cluster.foobar.id

	rule {
		type  = "ip_addr"
		value = "192.168.1.1"
	}

	rule {
		type  = "droplet"
		value = digitalocean_droplet.foobar.id
	}

	rule {
		type  = "tag"
		value = digitalocean_tag.foobar.name
	}
}
`
