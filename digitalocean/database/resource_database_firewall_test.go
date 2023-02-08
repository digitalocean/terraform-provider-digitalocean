package database_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseFirewall_Basic(t *testing.T) {
	databaseClusterName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	dbName := acceptance.RandomTestName()
	dropletName := acceptance.RandomTestName()
	tagName := acceptance.RandomTestName()
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseFirewallConfigMultipleResourceTypes,
					dbName, dropletName, tagName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_firewall.example", "rule.#", "4"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseFirewallDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_tag" "foobar" {
  name = "%s"
}

resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
    region = "nyc"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "basic-xxs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }
    }
  }
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

  rule {
    type  = "app"
    value = digitalocean_app.foobar.id
  }
}
`
