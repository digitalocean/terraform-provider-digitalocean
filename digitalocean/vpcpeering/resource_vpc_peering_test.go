package vpcpeering_test

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

func TestAccDigitalOceanVPCPeering_Basic(t *testing.T) {
	var vpcPeering godo.VPCPeering
	vpcPeeringName := acceptance.RandomTestName()
	vpcPeeringCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCPeeringConfig_Basic, vpcPeeringName)

	updateVPCPeeringName := acceptance.RandomTestName()
	vpcPeeringUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCPeeringConfig_Basic, updateVPCPeeringName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVPCPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcPeeringCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCPeeringExists("digitalocean_vpc_peering.foobar", &vpcPeering),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_peering.foobar", "name", vpcPeeringName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_peering.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_vpc_peering.foobar", "vpc_ids.0", "digitalocean_vpc.vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_vpc_peering.foobar", "vpc_ids.1", "digitalocean_vpc.vpc2", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_peering.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_peering.foobar", "status"),
				),
			},
			{
				Config: vpcPeeringUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCPeeringExists("digitalocean_vpc_peering.foobar", &vpcPeering),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_peering.foobar", "name", updateVPCPeeringName),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVPCPeeringDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "digitalocean_vpc_peering" {
			continue
		}

		_, _, err := client.VPCs.GetVPCPeering(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC Peering resource still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanVPCPeeringExists(resource string, vpcPeering *godo.VPCPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundVPCPeering, _, err := client.VPCs.GetVPCPeering(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundVPCPeering.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		*vpcPeering = *foundVPCPeering

		return nil
	}
}

const testAccCheckDigitalOceanVPCPeeringConfig_Basic = `
resource "digitalocean_vpc" "vpc1" {
  name   = "vpc1"
  region = "nyc3"
}

resource "digitalocean_vpc" "vpc2" {
  name   = "vpc2"
  region = "nyc3"
}

resource "digitalocean_vpc_peering" "foobar" {
  name = "%s"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc2.id
  ]
  depends_on = [
    digitalocean_vpc.vpc1,
    digitalocean_vpc.vpc2
  ]
}
`
