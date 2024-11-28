package dropletautoscale_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDropletAutoscale_Static(t *testing.T) {
	var autoscalePool godo.DropletAutoscalePool
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletAutoscaleDestroy,
		Steps: []resource.TestStep{
			{
				// Test create
				Config: testAccCheckDigitalOceanDropletAutoscaleConfig_static(name, 1, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("digitalocean_droplet_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("digitalocean_droplet_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_number_instances", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.region", "s2r1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.image", "547864"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "updated_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "members.#", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "history_events.#", "0"),
				),
			},
			{
				// Test update (static scale up)
				Config: testAccCheckDigitalOceanDropletAutoscaleConfig_static(name, 2, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("digitalocean_droplet_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("digitalocean_droplet_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_number_instances", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.region", "s2r1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.image", "547864"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "updated_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "members.#", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "history_events.#", "0"),
				),
			},
			{
				// Test listing members
				Config: testAccCheckDigitalOceanDropletAutoscaleConfig_static(name, 2, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("digitalocean_droplet_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("digitalocean_droplet_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_number_instances", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.region", "s2r1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.image", "547864"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "members.#"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "members.0.droplet_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "members.0.created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "members.0.updated_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "members.0.health_status"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "members.0.status"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "history_events.#", "0"),
				),
			},
			{
				// Test listing history events
				Config: testAccCheckDigitalOceanDropletAutoscaleConfig_static(name, 2, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("digitalocean_droplet_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("digitalocean_droplet_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_number_instances", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.region", "s2r1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.image", "547864"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "updated_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "members.#", "0"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "history_events.#"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "history_events.0.history_event_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "history_events.0.current_instance_count"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "history_events.0.desired_instance_count"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "history_events.0.reason"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "history_events.0.status"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "history_events.0.created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "history_events.0.updated_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDropletAutoscale_Dynamic(t *testing.T) {
	var autoscalePool godo.DropletAutoscalePool
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletAutoscaleDestroy,
		Steps: []resource.TestStep{
			{
				// Test create
				Config: testAccCheckDigitalOceanDropletAutoscaleConfig_dynamic(name, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("digitalocean_droplet_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("digitalocean_droplet_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.min_instances", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.max_instances", "3"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_cpu_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_memory_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.cooldown_minutes", "5"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_number_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.region", "s2r1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.image", "547864"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "updated_at"),
				),
			},
			{
				// Test update (dynamic scale up)
				Config: testAccCheckDigitalOceanDropletAutoscaleConfig_dynamic(name, 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("digitalocean_droplet_autoscale.foobar", &autoscalePool),
					resource.TestCheckResourceAttrSet("digitalocean_droplet_autoscale.foobar", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.min_instances", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.max_instances", "3"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_cpu_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_memory_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.cooldown_minutes", "5"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "config.0.target_number_instances", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.region", "s2r1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.image", "547864"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_autoscale.foobar", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_autoscale.foobar", "updated_at"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDropletAutoscaleExists(n string, autoscalePool *godo.DropletAutoscalePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %v", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set")
		}
		// Check for valid ID response to validate that the resource has been created
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		pool, _, err := client.DropletAutoscale.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if pool.ID != rs.Primary.ID {
			return fmt.Errorf("Droplet autoscale pool not found")
		}
		*autoscalePool = *pool
		return nil
	}
}

func testAccCheckDigitalOceanDropletAutoscaleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_droplet_autoscale" {
			continue
		}
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		_, _, err := client.DropletAutoscale.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("autoscale group with id %s not found", rs.Primary.ID)) {
				return nil
			}
			return fmt.Errorf("Droplet autoscale pool still exists")
		}
	}
	return nil
}

func testAccCheckDigitalOceanDropletAutoscaleConfig_static(name string, size int, getMembers, getEvents bool) string {
	pubKey1, _, err := acctest.RandSSHKeyPair("digitalocean@acceptance-test")
	if err != nil {
		fmt.Println("Unable to generate public key", err)
		return ""
	}

	pubKey2, _, err := acctest.RandSSHKeyPair("digitalocean@acceptance-test")
	if err != nil {
		fmt.Println("Unable to generate public key", err)
		return ""
	}

	var memberOpts, eventOpts string
	if getMembers {
		memberOpts = `
  list_member_opts {
    page     = 1
    per_page = 10
  }
`
	}
	if getEvents {
		memberOpts = `
  list_history_opts {
    page     = 1
    per_page = 10
  }
`
	}

	return fmt.Sprintf(`
resource "digitalocean_ssh_key" "foo" {
  name       = "%s"
  public_key = "%s"
}

resource "digitalocean_ssh_key" "bar" {
  name       = "%s"
  public_key = "%s"
}

resource "digitalocean_tag" "foo" {
  name = "%s"
}

resource "digitalocean_tag" "bar" {
  name = "%s"
}

resource "digitalocean_droplet_autoscale" "foobar" {
  name = "%s"

  config {
    target_number_instances = %d
  }

  droplet_template {
    size               = "c-2"
    region             = "s2r1"
    image              = "547864"
    tags               = [digitalocean_tag.foo.id, digitalocean_tag.bar.id]
    ssh_keys           = [digitalocean_ssh_key.foo.id, digitalocean_ssh_key.bar.id]
    with_droplet_agent = true
    ipv6               = true
    user_data          = "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"
  }

  # embed pagination parameters for listing members
  %s

  # embed pagination parameters for listing history events
  %s
}`,
		acceptance.RandomTestName("sshKey1"), pubKey1,
		acceptance.RandomTestName("sshKey2"), pubKey2,
		acceptance.RandomTestName("tag1"),
		acceptance.RandomTestName("tag2"),
		name, size, memberOpts, eventOpts)
}

func testAccCheckDigitalOceanDropletAutoscaleConfig_dynamic(name string, size int) string {
	pubKey1, _, err := acctest.RandSSHKeyPair("digitalocean@acceptance-test")
	if err != nil {
		fmt.Println("Unable to generate public key", err)
		return ""
	}

	pubKey2, _, err := acctest.RandSSHKeyPair("digitalocean@acceptance-test")
	if err != nil {
		fmt.Println("Unable to generate public key", err)
		return ""
	}

	return fmt.Sprintf(`
resource "digitalocean_ssh_key" "foo" {
  name       = "%s"
  public_key = "%s"
}

resource "digitalocean_ssh_key" "bar" {
  name       = "%s"
  public_key = "%s"
}

resource "digitalocean_tag" "foo" {
  name = "%s"
}

resource "digitalocean_tag" "bar" {
  name = "%s"
}

resource "digitalocean_droplet_autoscale" "foobar" {
  name = "%s"

  config {
    min_instances             = %d
    max_instances             = 3
    target_cpu_utilization    = 0.5
    target_memory_utilization = 0.5
    cooldown_minutes          = 5
  }

  droplet_template {
    size               = "c-2"
    region             = "s2r1"
    image              = "547864"
    tags               = [digitalocean_tag.foo.id, digitalocean_tag.bar.id]
    ssh_keys           = [digitalocean_ssh_key.foo.id, digitalocean_ssh_key.bar.id]
    with_droplet_agent = true
    ipv6               = true
    user_data          = "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"
  }
}`,
		acceptance.RandomTestName("sshKey1"), pubKey1,
		acceptance.RandomTestName("sshKey2"), pubKey2,
		acceptance.RandomTestName("tag1"),
		acceptance.RandomTestName("tag2"),
		name, size)
}
