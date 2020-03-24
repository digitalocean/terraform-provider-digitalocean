package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanVPC_Basic(t *testing.T) {
	vpcName := randomTestName()
	vpcDesc := "A description for the VPC"
	vpcCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCConfig_Basic, vpcName, vpcDesc)

	updatedVPCName := randomTestName()
	// updatedVPVDesc := "A brand new updated description for the VPC"
	vpcUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCConfig_Basic, updatedVPCName, vpcDesc) //updatedVPVDesc)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCExists("digitalocean_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "name", vpcName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "default", "false"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "description", vpcDesc),
				),
			},
			{
				Config: vpcUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCExists("digitalocean_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "name", updatedVPCName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "description", vpcDesc), //updatedVPVDesc),
				),
			},
		},
	})
}

func TestAccDigitalOceanVPC_IPRange(t *testing.T) {
	vpcName := randomTestName()
	vpcCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCConfig_IPRange, vpcName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCExists("digitalocean_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "name", vpcName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "ip_range", "10.00.10.0/24"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVPCDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "digitalocean_vpc" {
			continue
		}

		_, _, err := client.VPCs.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC resource still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanVPCExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundVPC, _, err := client.VPCs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundVPC.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		return nil
	}
}

const testAccCheckDigitalOceanVPCConfig_Basic = `
resource "digitalocean_vpc" "foobar" {
	name        = "%s"
	description = "%s"
	region      = "s2r1" # "nyc3"
}
`
const testAccCheckDigitalOceanVPCConfig_IPRange = `
resource "digitalocean_vpc" "foobar" {
	name     = "%s"
	region   = "s2r1" # "nyc3"
	ip_range = "10.00.10.0/24"
}
`
