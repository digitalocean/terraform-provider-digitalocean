package reservedip_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanReservedIPAssignment(t *testing.T) {
	var reservedIP godo.ReservedIP

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPAttachmentExists("digitalocean_reserved_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPAssignmentReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPAttachmentExists("digitalocean_reserved_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPAssignmentDeleteAssignment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPExists("digitalocean_reserved_ip.foobar", &reservedIP),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip.foobar", "ip_address", regexp.MustCompile("[0-9.]+")),
				),
			},
		},
	})
}

func TestAccDigitalOceanReservedIPAssignment_createBeforeDestroy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPAssignmentConfig_createBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPAttachmentExists("digitalocean_reserved_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPAssignmentConfig_createBeforeDestroyReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPAttachmentExists("digitalocean_reserved_ip_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanReservedIPAttachmentExists(n string) resource.TestCheckFunc {
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

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		// Try to find the ReservedIP
		foundReservedIP, _, err := client.ReservedIPs.Get(context.Background(), fipID)
		if err != nil {
			return err
		}

		if foundReservedIP.IP != fipID || foundReservedIP.Droplet.ID != dropletID {
			return fmt.Errorf("wrong floating IP attachment found")
		}

		return nil
	}
}

var testAccCheckDigitalOceanReservedIPAssignmentConfig = `
resource "digitalocean_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_reserved_ip_assignment" "foobar" {
  ip_address = digitalocean_reserved_ip.foobar.ip_address
  droplet_id = digitalocean_droplet.foobar.0.id
}
`

var testAccCheckDigitalOceanReservedIPAssignmentReassign = `
resource "digitalocean_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_reserved_ip_assignment" "foobar" {
  ip_address = digitalocean_reserved_ip.foobar.ip_address
  droplet_id = digitalocean_droplet.foobar.1.id
}
`

var testAccCheckDigitalOceanReservedIPAssignmentDeleteAssignment = `
resource "digitalocean_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count              = 2
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}
`

var testAccCheckDigitalOceanReservedIPAssignmentConfig_createBeforeDestroy = `
resource "digitalocean_droplet" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test"
  region = "nyc3"
  size   = "s-1vcpu-1gb"

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "digitalocean_reserved_ip_assignment" "foobar" {
  ip_address = digitalocean_reserved_ip.foobar.id
  droplet_id = digitalocean_droplet.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`

var testAccCheckDigitalOceanReservedIPAssignmentConfig_createBeforeDestroyReassign = `
resource "digitalocean_droplet" "foobar" {
  image  = "ubuntu-18-04-x64"
  name   = "tf-acc-test-01"
  region = "nyc3"
  size   = "s-1vcpu-1gb"

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_reserved_ip" "foobar" {
  region = "nyc3"
}

resource "digitalocean_reserved_ip_assignment" "foobar" {
  ip_address = digitalocean_reserved_ip.foobar.id
  droplet_id = digitalocean_droplet.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`
