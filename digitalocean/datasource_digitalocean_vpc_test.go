package digitalocean

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanVPC_ByName(t *testing.T) {
	vpcName := randomTestName()
	vpcDesc := "A description for the VPC"
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVPCConfig_Basic, vpcName, vpcDesc)
	dataSourceConfig := `
data "digitalocean_vpc" "foobar" {
  name = digitalocean_vpc.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCExists("data.digitalocean_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc.foobar", "name", vpcName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc.foobar", "description", vpcDesc),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc.foobar", "default"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc.foobar", "ip_range"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc.foobar", "urn"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVPC_RegionDefault(t *testing.T) {
	vpcDropletName := randomTestName()
	vpcConfigRegionDefault := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVPCConfig_RegionDefault, vpcDropletName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: vpcConfigRegionDefault,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCExists("data.digitalocean_vpc.foobar"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc.foobar", "name"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc.foobar", "default", "true"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVPC_ExpectErrors(t *testing.T) {
	vpcName := randomTestName()
	vpcNotExist := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVPCConfig_DoesNotExist, vpcName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataSourceDigitalOceanVPCConfig_MissingRegionDefault,
				ExpectError: regexp.MustCompile(`unable to find default VPC in foo region`),
			},
			{
				Config:      vpcNotExist,
				ExpectError: regexp.MustCompile(`no VPCs found with name`),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanVPCConfig_Basic = `
resource "digitalocean_vpc" "foobar" {
	name        = "%s"
	description = "%s"
	region      = "nyc3"
}`

const testAccCheckDataSourceDigitalOceanVPCConfig_RegionDefault = `
// Create Droplet to ensure default VPC exists
resource "digitalocean_droplet" "foo" {
	image  = "ubuntu-18-04-x64"
	name   = "%s"
	region = "nyc3"
	size   = "s-1vcpu-1gb"
	private_networking = "true"
}

data "digitalocean_vpc" "foobar" {
	region = "nyc3"
}
`

const testAccCheckDataSourceDigitalOceanVPCConfig_MissingRegionDefault = `
data "digitalocean_vpc" "foobar" {
	region = "foo"
}
`

const testAccCheckDataSourceDigitalOceanVPCConfig_DoesNotExist = `
data "digitalocean_vpc" "foobar" {
	name = "%s"
}
`
