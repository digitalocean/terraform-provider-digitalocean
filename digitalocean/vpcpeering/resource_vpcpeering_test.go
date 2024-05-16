package vpcpeering_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanVPCPeering_Basic(t *testing.T) {
	vpcPeeringName := acceptance.RandomTestName()
	vpcIDs := `["foo", "bar"]`
	vpcPeeringCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCPeeringConfig_Basic, vpcPeeringName, vpcIDs)

	updateVPCPeeringName := acceptance.RandomTestName()
	vpcPeeringUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCPeeringConfig_Basic, updateVPCPeeringName, vpcIDs)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVPCPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcPeeringCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCPeeringExists("digitalocean_vpcpeering.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "name", vpcPeeringName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "vpc_ids.0", "foo"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "vpc_ids.1", "bar"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpcpeering.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpcpeering.foobar", "status"),
				),
			},
			{
				Config: vpcPeeringUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCPeeringExists("digitalocean_vpcpeering.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "name", updateVPCPeeringName),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVPCPeeringDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "digitalocean_vpcpeering" {
			continue
		}

		_, _, err := client.VPCs.GetVPCPeering(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC Peering resource still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanVPCPeeringExists(resource string) resource.TestCheckFunc {
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

		return nil
	}
}

const testAccCheckDigitalOceanVPCPeeringConfig_Basic = `
resource "digitalocean_vpcpeering" "foobar" {
  name        = "%s"
  vpc_ids     = %s
}
`
