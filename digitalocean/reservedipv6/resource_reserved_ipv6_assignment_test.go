package reservedipv6_test

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

func TestAccDigitalOceanReservedIPV6Assignment(t *testing.T) {
	var reservedIPv6 godo.ReservedIPV6

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPV6AssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6AttachmentExists("digitalocean_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},

			{
				Config: testAccCheckDigitalOceanReservedIPV6AssignmentDeleteAssignment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6Exists("digitalocean_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6.foobar", "ip", regexp.MustCompile(ipv6Regex)),
				),
			},
		},
	})
}

func TestAccDigitalOceanReservedIPV6Assignment_createBeforeDestroy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPV6AssignmentConfig_createBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6AttachmentExists("digitalocean_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPV6AssignmentConfig_createBeforeDestroyReassign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6AttachmentExists("digitalocean_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
		},
	})
}

func TestAccDigitalOceanReservedIPV6Assignment_unassignAndAssignToNewDroplet(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPV6AssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6AttachmentExists("digitalocean_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPV6AssignmentUnAssign,
			},
			{
				Config: testAccCheckDigitalOceanReservedIPV6ReAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6AttachmentExists("digitalocean_reserved_ipv6_assignment.foobar"),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "ip", regexp.MustCompile(ipv6Regex)),
					resource.TestMatchResourceAttr(
						"digitalocean_reserved_ipv6_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanReservedIPV6AttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["ip"] == "" {
			return fmt.Errorf("No reserved IPv6 is set")
		}
		fipID := rs.Primary.Attributes["ip"]
		dropletID, err := strconv.Atoi(rs.Primary.Attributes["droplet_id"])
		if err != nil {
			return err
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		// Try to find the ReservedIPv6
		foundReservedIP, _, err := client.ReservedIPV6s.Get(context.Background(), fipID)
		if err != nil {
			return err
		}

		if foundReservedIP.IP != fipID || foundReservedIP.Droplet.ID != dropletID {
			return fmt.Errorf("wrong floating IP attachment found")
		}

		return nil
	}
}

var testAccCheckDigitalOceanReservedIPV6AssignmentConfig = `
resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count  = 1
  name   = "tf-acc-test-assign"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}

resource "digitalocean_reserved_ipv6_assignment" "foobar" {
  ip         = digitalocean_reserved_ipv6.foobar.ip
  droplet_id = digitalocean_droplet.foobar.0.id
}
`

var testAccCheckDigitalOceanReservedIPV6AssignmentUnAssign = `
resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count  = 1
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test-assign"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ipv6   = true
}

`

var testAccCheckDigitalOceanReservedIPV6ReAssignmentConfig = `
resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "digitalocean_droplet" "foobar1" {
  count  = 1
  name   = "tf-acc-test-reassign"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}
resource "digitalocean_droplet" "foobar" {
  count  = 1
  name   = "tf-acc-test-assign"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}

resource "digitalocean_reserved_ipv6_assignment" "foobar" {
  ip         = digitalocean_reserved_ipv6.foobar.ip
  droplet_id = digitalocean_droplet.foobar1.0.id
}
`

var testAccCheckDigitalOceanReservedIPV6AssignmentDeleteAssignment = `
resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  count              = 1
  name               = "tf-acc-test-${count.index}"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
}
`

var testAccCheckDigitalOceanReservedIPV6AssignmentConfig_createBeforeDestroy = `
resource "digitalocean_droplet" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ipv6   = true

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "digitalocean_reserved_ipv6_assignment" "foobar" {
  ip         = digitalocean_reserved_ipv6.foobar.ip
  droplet_id = digitalocean_droplet.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`
var testAccCheckDigitalOceanReservedIPV6AssignmentConfig_createBeforeDestroyReassign = `
resource "digitalocean_droplet" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test-01"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ipv6   = true

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "digitalocean_reserved_ipv6_assignment" "foobar" {
  ip         = digitalocean_reserved_ipv6.foobar.ip
  droplet_id = digitalocean_droplet.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
`
