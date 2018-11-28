package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceDigitalOceanDropletSnapshot_basic(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanDropletSnapshot_basic, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletSnapshotExists("data.digitalocean_droplet_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "min_disk_size", "20"),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_snapshot.foobar", "droplet_id"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDropletSnapshot_regex(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanDropletSnapshot_regex, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletSnapshotExists("data.digitalocean_droplet_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "min_disk_size", "20"),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_snapshot.foobar", "droplet_id"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDropletSnapshot_region(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanDropletSnapshot_region, rInt, rInt, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletSnapshotExists("data.digitalocean_droplet_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "min_disk_size", "20"),
					resource.TestCheckResourceAttr("data.digitalocean_droplet_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet_snapshot.foobar", "droplet_id"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanDropletSnapshotExists(n string, snapshot *godo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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
  name   = "%d"
  size   = "512mb"
  image  = "centos-7-x64"
  region = "nyc3"
  ipv6   = true
}

resource "digitalocean_droplet_snapshot" "foo" {
  name = "snapshot-%d"
  droplet_id = "${digitalocean_droplet.foo.id}"
}

data "digitalocean_droplet_snapshot" "foobar" {
  most_recent = true
  name = "${digitalocean_droplet_snapshot.foo.name}"
}`

const testAccCheckDataSourceDigitalOceanDropletSnapshot_regex = `
resource "digitalocean_droplet" "foo" {
  region      = "nyc1"
  name        = "foo-%d"
  size        = "512mb"
  image  			= "centos-7-x64"
  ipv6   			= true
}

resource "digitalocean_droplet_snapshot" "foo" {
  name = "snapshot-%d"
  droplet_id = "${digitalocean_droplet.foo.id}"
}

data "digitalocean_droplet_snapshot" "foobar" {
  most_recent = true
  name_regex = "^${digitalocean_droplet_snapshot.foo.name}"
}`

const testAccCheckDataSourceDigitalOceanDropletSnapshot_region = `
resource "digitalocean_droplet" "foo" {
  region      = "nyc1"
  name        = "foo-nyc-%d"
  size        = "512mb"
  image  			= "centos-7-x64"
  ipv6   			= true
}

resource "digitalocean_droplet" "bar" {
  region      = "lon1"
  name        = "bar-lon-%d"
  size        = "512mb"
  image  			= "centos-7-x64"
  ipv6   			= true
}

resource "digitalocean_droplet_snapshot" "foo" {
  name = "snapshot-%d"
  droplet_id = "${digitalocean_droplet.foo.id}"
}

resource "digitalocean_droplet_snapshot" "bar" {
  name = "snapshot-%d"
  droplet_id = "${digitalocean_droplet.bar.id}"
}

data "digitalocean_droplet_snapshot" "foobar" {
  name = "${digitalocean_droplet_snapshot.bar.name}"
  region = "lon1"
}`
