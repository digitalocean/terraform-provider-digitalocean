package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_floating_ip", &resource.Sweeper{
		Name: "digitalocean_floating_ip",
		F:    testSweepFloatingIps,
	})

}

func testSweepFloatingIps(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*CombinedConfig).godoClient()

	ips, _, err := client.FloatingIPs.List(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		if _, err := client.FloatingIPs.Delete(context.Background(), ip.IP); err != nil {
			return err
		}
	}

	return nil
}

func TestAccDigitalOceanFloatingIP_Region(t *testing.T) {
	var floatingIP godo.FloatingIP

	expectedURNRegEx, _ := regexp.Compile(`do:floatingip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_region,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("digitalocean_floating_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccDigitalOceanFloatingIP_Droplet(t *testing.T) {
	var floatingIP godo.FloatingIP
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_droplet(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_Reassign(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_Unassign(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip.foobar", "region", "nyc3"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanFloatingIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_floating_ip" {
			continue
		}

		// Try to find the key
		_, _, err := client.FloatingIPs.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Floating IP still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanFloatingIPExists(n string, floatingIP *godo.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		// Try to find the FloatingIP
		foundFloatingIP, _, err := client.FloatingIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundFloatingIP.IP != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*floatingIP = *foundFloatingIP

		return nil
	}
}

var testAccCheckDigitalOceanFloatingIPConfig_region = `
resource "digitalocean_floating_ip" "foobar" {
  region = "nyc3"
}`

func testAccCheckDigitalOceanFloatingIPConfig_droplet(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name               = "tf-acc-test-%d"
  size               = "s-1vcpu-1gb"
  image              = "centos-7-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip" "foobar" {
  droplet_id = "${digitalocean_droplet.foobar.id}"
  region     = "${digitalocean_droplet.foobar.region}"
}`, rInt)
}

func testAccCheckDigitalOceanFloatingIPConfig_Reassign(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "baz" {
  name               = "tf-acc-test-%d"
  size               = "s-1vcpu-1gb"
  image              = "centos-7-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip" "foobar" {
  droplet_id = "${digitalocean_droplet.baz.id}"
  region     = "${digitalocean_droplet.baz.region}"
}`, rInt)
}

func testAccCheckDigitalOceanFloatingIPConfig_Unassign(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "baz" {
  name               = "tf-acc-test-%d"
  size               = "s-1vcpu-1gb"
  image              = "centos-7-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip" "foobar" {
  region     = "nyc3"
}`, rInt)
}
