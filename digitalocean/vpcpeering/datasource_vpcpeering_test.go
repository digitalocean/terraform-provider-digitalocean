package vpcpeering_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanVPCPeering_ByID(t *testing.T) {
	var vpcPeering godo.VPCPeering
	vpcPeeringName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVPCPeeringConfig_Basic, vpcName1, vpcName2, vpcPeeringName)
	dataSourceConfig := `
data "digitalocean_vpcpeering" "foobar" {
  id = digitalocean_vpcpeering.foobar.id
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
					testAccCheckDigitalOceanVPCPeeringExists("data.digitalocean_vpcpeering.foobar", &vpcPeering),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpcpeering.foobar", "name", vpcPeeringName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpcpeering.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_vpcpeering.foobar", "vpc_ids.0", "digitalocean_vpc.vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_vpcpeering.foobar", "vpc_ids.1", "digitalocean_vpc.vpc2", "id"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpcpeering.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpcpeering.foobar", "status"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVPCPeering_ByName(t *testing.T) {
	var vpcPeering godo.VPCPeering
	vpcPeeringName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVPCPeeringConfig_Basic, vpcName1, vpcName2, vpcPeeringName)
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
					testAccCheckDigitalOceanVPCPeeringExists("data.digitalocean_vpcpeering.foobar", &vpcPeering),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "name", vpcPeeringName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpcpeering.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_vpcpeering.foobar", "vpc_ids.0", "digitalocean_vpc.vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_vpcpeering.foobar", "vpc_ids.1", "digitalocean_vpc.vpc2", "id"),
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
				ExpectError: regexp.MustCompile(`Error retrieving VPC Peering`),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanVPCPeeringConfig_Basic = `
resource "digitalocean_vpc" "vpc1" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_vpc" "vpc2" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_vpcpeering" "foobar" {
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

const testAccCheckDataSourceDigitalOceanVPCPeeringConfig_DoesNotExist = `
data "digitalocean_vpcpeering" "foobar" {
  id = "%s"
}
`
