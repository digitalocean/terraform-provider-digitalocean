package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanFloatingIPAssignment(t *testing.T) {
	var floatingIP godo.FloatingIP

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPAttachmentExists("digitalocean_floating_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanFloatingIPAssignmentReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPAttachmentExists("digitalocean_floating_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanFloatingIPAssignmentDeleteAssignment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip.foobar", "ip_address", regexp.MustCompile("[0-9.]+")),
				),
			},
		},
	})
}

func TestAccDigitalOceanFloatingIPAssignment_createBeforeDestroy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPAssignmentConfig_createBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPAttachmentExists("digitalocean_floating_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanFloatingIPAssignmentConfig_createBeforeDestroyReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPAttachmentExists("digitalocean_floating_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanFloatingIPAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["ip_address"] == "" {
			return fmt.Errorf("No floating IP is set")
		}
		fipID := rs.Primary.Attributes["ip_address"]
		dropletID, err := strconv.Atoi(rs.Primary.Attributes["droplet_id"])
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		// Try to find the FloatingIP
		foundFloatingIP, _, err := client.FloatingIPs.Get(context.Background(), fipID)
		if err != nil {
			return err
		}

		if foundFloatingIP.IP != fipID || foundFloatingIP.Droplet.ID != dropletID {
			return fmt.Errorf("wrong floating IP attachment found")
		}

		return nil
	}
}

var testAccCheckDigitalOceanFloatingIPAssignmentConfig = `
resource "digitalocean_floating_ip" "foobar" {
  region = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "centos-7-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip_assignment" "foobar" {
  ip_address = "${digitalocean_floating_ip.foobar.ip_address}"
  droplet_id = "${digitalocean_droplet.foobar.0.id}"
}
`

var testAccCheckDigitalOceanFloatingIPAssignmentReassign = `
resource "digitalocean_floating_ip" "foobar" {
  region = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "centos-7-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip_assignment" "foobar" {
  ip_address = "${digitalocean_floating_ip.foobar.ip_address}"
  droplet_id = "${digitalocean_droplet.foobar.1.id}"
}
`

var testAccCheckDigitalOceanFloatingIPAssignmentDeleteAssignment = `
resource "digitalocean_floating_ip" "foobar" {
  region = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "centos-7-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}
`

var testAccCheckDigitalOceanFloatingIPAssignmentConfig_createBeforeDestroy = `
resource "digitalocean_droplet" "foobar" {
  image = "centos-7-x64"
  name = "tf-acc-test"
  region = "nyc3"
  size = "s-1vcpu-1gb"

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_floating_ip" "foobar" {
  region     = "nyc3"
}

resource "digitalocean_floating_ip_assignment" "foobar" {
  ip_address = "${digitalocean_floating_ip.foobar.id}"
  droplet_id = "${digitalocean_droplet.foobar.id}"

  lifecycle {
    create_before_destroy = true
  }
}
`

var testAccCheckDigitalOceanFloatingIPAssignmentConfig_createBeforeDestroyReassign = `
resource "digitalocean_droplet" "foobar" {
  image = "ubuntu-18-04-x64"
  name = "tf-acc-test"
  region = "nyc3"
  size = "s-1vcpu-1gb"

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_floating_ip" "foobar" {
  region     = "nyc3"
}

resource "digitalocean_floating_ip_assignment" "foobar" {
  ip_address = "${digitalocean_floating_ip.foobar.id}"
  droplet_id = "${digitalocean_droplet.foobar.id}"

  lifecycle {
    create_before_destroy = true
  }
}
`
