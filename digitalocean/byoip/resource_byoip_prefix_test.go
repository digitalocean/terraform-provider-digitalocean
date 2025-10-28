package byoip_test

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

func TestAccDigitalOceanBYOIPPrefix_Basic(t *testing.T) {
	var byoipPrefix godo.BYOIPPrefix

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBYOIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanBYOIPPrefixConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBYOIPPrefixExists("digitalocean_byoip_prefix.test", &byoipPrefix),
					resource.TestCheckResourceAttr("digitalocean_byoip_prefix.test", "prefix", "203.0.113.0/24"),
					resource.TestCheckResourceAttr("digitalocean_byoip_prefix.test", "region", "nyc3"),
				),
			},
		},
	})
}

func TestAccDigitalOceanBYOIPPrefix_CreateReadUpdate(t *testing.T) {
	var byoipPrefix godo.BYOIPPrefix

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBYOIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanBYOIPPrefixConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBYOIPPrefixExists("digitalocean_byoip_prefix.test", &byoipPrefix),
					resource.TestCheckResourceAttr("digitalocean_byoip_prefix.test", "prefix", "203.0.113.0/24"),
					resource.TestCheckResourceAttr("digitalocean_byoip_prefix.test", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanBYOIPPrefixConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBYOIPPrefixExists("digitalocean_byoip_prefix.test", &byoipPrefix),
					resource.TestCheckResourceAttr("digitalocean_byoip_prefix.test", "advertised", "true"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanBYOIPPrefixDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_byoip_prefix" {
			continue
		}

		_, _, err := client.BYOIPPrefixes.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("BYOIP Prefix still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanBYOIPPrefixExists(n string, byoipPrefix *godo.BYOIPPrefix) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		found, _, err := client.BYOIPPrefixes.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		*byoipPrefix = *found
		return nil
	}
}

var testAccCheckDigitalOceanBYOIPPrefixConfig_basic = `
resource "digitalocean_byoip_prefix" "test" {
  prefix    = "203.0.113.0/24"
  region    = "nyc3"
  signature = "your-signature-here"
}
`

var testAccCheckDigitalOceanBYOIPPrefixConfig_update = `
resource "digitalocean_byoip_prefix" "test" {
  prefix     = "203.0.113.0/24"
  region     = "nyc3"
  signature  = "your-signature-here"
  advertised = true
}
`
