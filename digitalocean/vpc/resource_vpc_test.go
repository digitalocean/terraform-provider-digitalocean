package vpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanVPC_Basic(t *testing.T) {
	vpcName := acceptance.RandomTestName()
	vpcDesc := "A description for the VPC"
	vpcCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCConfig_Basic, vpcName, vpcDesc)

	updatedVPCName := acceptance.RandomTestName()
	updatedVPVDesc := "A brand new updated description for the VPC"
	vpcUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCConfig_Basic, updatedVPCName, updatedVPVDesc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVPCDestroy,
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
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc.foobar", "ip_range"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc.foobar", "urn"),
				),
			},
			{
				Config: vpcUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCExists("digitalocean_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "name", updatedVPCName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "description", updatedVPVDesc),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "default", "false"),
				),
			},
		},
	})
}

func TestAccDigitalOceanVPC_IPRange(t *testing.T) {
	vpcName := acceptance.RandomTestName()
	vpcCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCConfig_IPRange, vpcName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCExists("digitalocean_vpc.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "name", vpcName),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "ip_range", "10.10.10.0/24"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foobar", "default", "false"),
				),
			},
		},
	})
}

// https://github.com/digitalocean/terraform-provider-digitalocean/issues/551
func TestAccDigitalOceanVPC_IPRangeRace(t *testing.T) {
	vpcNameOne := acceptance.RandomTestName()
	vpcNameTwo := acceptance.RandomTestName()
	vpcCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCConfig_IPRangeRace, vpcNameOne, vpcNameTwo)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVPCExists("digitalocean_vpc.foo"),
					testAccCheckDigitalOceanVPCExists("digitalocean_vpc.bar"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.foo", "name", vpcNameOne),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc.foo", "ip_range"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc.bar", "name", vpcNameTwo),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc.bar", "ip_range"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVPCDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
  region      = "nyc3"
}
`
const testAccCheckDigitalOceanVPCConfig_IPRange = `
resource "digitalocean_vpc" "foobar" {
  name     = "%s"
  region   = "nyc3"
  ip_range = "10.10.10.0/24"
}
`

const testAccCheckDigitalOceanVPCConfig_IPRangeRace = `
resource "digitalocean_vpc" "foo" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_vpc" "bar" {
  name   = "%s"
  region = "nyc3"
}
`
