package vpcpeering_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanVPCPeering_ByName(t *testing.T) {
	vpcPeeringName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVPCPeeringConfig_Basic, vpcPeeringName)
	dataSourceConfig := `
data "digitalocean_vpcpeering" "foobar" {
  name = digitalocean_vpcpeering.foobar.name
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
					testAccCheckDigitalOceanVPCPeeringExists("data.digitalocean_vpcpeering.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "name", vpcPeeringName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_vpcpeering.foobar", "vpc_ids.0", "digitalocean_vpc.tf-vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_vpcpeering.foobar", "vpc_ids.1", "digitalocean_vpc.tf-vpc2", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpcpeering.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpcpeering.foobar", "status"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVPCPeering_ExpectErrors(t *testing.T) {
	vpcPeeringName := acceptance.RandomTestName()
	vpcPeeringNotExist := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVPCPeeringConfig_DoesNotExist, vpcPeeringName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      vpcPeeringNotExist,
				ExpectError: regexp.MustCompile(`no VPC Peerings found with name`),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanVPCPeeringConfig_Basic = `
resource "digitalocean_vpc" "tf-vpc1" {
  name   = "tf-vpc1"
  region = "s2r1"
}

resource "digitalocean_vpc" "tf-vpc2" {
  name   = "tf-vpc2"
  region = "s2r1"
}

resource "digitalocean_vpcpeering" "foobar" {
  name = "%s"
  vpc_ids = [
    digitalocean_vpc.tf-vpc1.id,
    digitalocean_vpc.tf-vpc2.id
  ]
  depends_on = [
    digitalocean_vpc.tf-vpc1,
    digitalocean_vpc.tf-vpc2
  ]
}
`

const testAccCheckDataSourceDigitalOceanVPCPeeringConfig_DoesNotExist = `
data "digitalocean_vpcpeering" "foobar" {
  name = "%s"
}
`
