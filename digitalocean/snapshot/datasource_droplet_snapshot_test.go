package snapshot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanDropletSnapshot_basic(t *testing.T) {
	var snapshot godo.Snapshot
	testName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanDropletSnapshot_basic, testName, testName)
	dataSourceConfig := `
data "digitalocean_droplet_snapshot" "foobar" {
  most_recent = true
  name = digitalocean_droplet_snapshot.foo.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletSnapshotExists("data.digitalocean_droplet_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "name", fmt.Sprintf("%s-snapshot", testName)),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "min_disk_size", "25"),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_snapshot.foobar", "droplet_id"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDropletSnapshot_regex(t *testing.T) {
	var snapshot godo.Snapshot
	testName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanDropletSnapshot_basic, testName, testName)
	dataSourceConfig := fmt.Sprintf(`
data "digitalocean_droplet_snapshot" "foobar" {
  most_recent = true
  name_regex = "^%s"
}`, testName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletSnapshotExists("data.digitalocean_droplet_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "name", fmt.Sprintf("%s-snapshot", testName)),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "min_disk_size", "25"),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_snapshot.foobar", "droplet_id"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDropletSnapshot_region(t *testing.T) {
	var snapshot godo.Snapshot
	testName := acceptance.RandomTestName()
	nycResourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanDropletSnapshot_basic, testName, testName)
	lonResourceConfig := fmt.Sprintf(`
resource "digitalocean_droplet" "bar" {
  region = "lon1"
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
}

resource "digitalocean_droplet_snapshot" "bar" {
  name = "%s-snapshot"
  droplet_id = "${digitalocean_droplet.bar.id}"
}`, testName, testName)
	dataSourceConfig := `
data "digitalocean_droplet_snapshot" "foobar" {
  name = digitalocean_droplet_snapshot.bar.name
  region = "lon1"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: nycResourceConfig + lonResourceConfig,
			},
			{
				Config: nycResourceConfig + lonResourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletSnapshotExists("data.digitalocean_droplet_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "name", fmt.Sprintf("%s-snapshot", testName)),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "min_disk_size", "25"),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_snapshot.foobar", "droplet_id"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanDropletSnapshotExists(n string, snapshot *godo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No snapshot ID is set")
		}

		foundSnapshot, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundSnapshot.ID != rs.Primary.ID {
			return fmt.Errorf("Snapshot not found")
		}

		*snapshot = *foundSnapshot

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanDropletSnapshot_basic = `
resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_droplet_snapshot" "foo" {
  name = "%s-snapshot"
  droplet_id = digitalocean_droplet.foo.id
}
`
