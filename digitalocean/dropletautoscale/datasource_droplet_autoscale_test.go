package dropletautoscale_test

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDropletAutoscale_Static(t *testing.T) {
	var autoscalePool godo.DropletAutoscalePool
	name := acceptance.RandomTestName()

	createConfig := testAccCheckDigitalOceanDropletAutoscaleConfig_static(name, 1)
	dataSourceIDConfig := `
data "digitalocean_droplet_autoscale" "foo" {
  id = digitalocean_droplet_autoscale.foobar.id
}`
	dataSourceNameConfig := `
data "digitalocean_droplet_autoscale" "foo" {
  name = digitalocean_droplet_autoscale.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletAutoscaleDestroy,
		Steps: []resource.TestStep{
			{
				// Test create
				Config: createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("digitalocean_droplet_autoscale.foobar", &autoscalePool),
				),
			},
			{
				// Import by id
				Config: createConfig + dataSourceIDConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("data.digitalocean_droplet_autoscale.foo", &autoscalePool),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_autoscale.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_number_instances", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.image"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "updated_at"),
				),
			},
			{
				// Import by name
				Config: createConfig + dataSourceNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("data.digitalocean_droplet_autoscale.foo", &autoscalePool),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_autoscale.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.min_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.max_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_cpu_utilization", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_memory_utilization", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.cooldown_minutes", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_number_instances", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.image"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "updated_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDropletAutoscale_Dynamic(t *testing.T) {
	var autoscalePool godo.DropletAutoscalePool
	name := acceptance.RandomTestName()

	createConfig := testAccCheckDigitalOceanDropletAutoscaleConfig_dynamic(name, 1)
	dataSourceIDConfig := `
data "digitalocean_droplet_autoscale" "foo" {
  id = digitalocean_droplet_autoscale.foobar.id
}`
	dataSourceNameConfig := `
data "digitalocean_droplet_autoscale" "foo" {
  name = digitalocean_droplet_autoscale.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletAutoscaleDestroy,
		Steps: []resource.TestStep{
			{
				// Test create
				Config: createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("digitalocean_droplet_autoscale.foobar", &autoscalePool),
				),
			},
			{
				// Import by id
				Config: createConfig + dataSourceIDConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("data.digitalocean_droplet_autoscale.foo", &autoscalePool),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_autoscale.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.min_instances", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.max_instances", "3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_cpu_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_memory_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.cooldown_minutes", "5"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_number_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.image"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "updated_at"),
				),
			},
			{
				// Import by name
				Config: createConfig + dataSourceNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanDropletAutoscaleExists("data.digitalocean_droplet_autoscale.foo", &autoscalePool),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_autoscale.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.min_instances", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.max_instances", "3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_cpu_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_memory_utilization", "0.5"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.cooldown_minutes", "5"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "config.0.target_number_instances", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.size", "c-2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.region", "nyc3"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.image"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.with_droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.user_data", "\n#cloud-config\nruncmd:\n- apt-get update\n- apt-get install -y stress-ng\n"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "droplet_template.0.ssh_keys.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet_autoscale.foo", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_droplet_autoscale.foo", "updated_at"),
				),
			},
		},
	})
}
