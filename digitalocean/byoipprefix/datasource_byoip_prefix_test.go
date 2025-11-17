package byoipprefix_test

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

func TestAccDataSourceDigitalOceanBYOIPPrefix_Basic(t *testing.T) {
	var prefix godo.BYOIPPrefix
	prefixCIDR := "192.0.2.0/24"
	region := "nyc3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanBYOIPPrefixConfig_basic(prefixCIDR, region, false),
			},
			{
				Config: testAccCheckDataSourceDigitalOceanBYOIPPrefixConfig_basic(prefixCIDR, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanBYOIPPrefixExists("data.digitalocean_byoip_prefix.foobar", &prefix),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_byoip_prefix.foobar", "uuid",
						"digitalocean_byoip_prefix.foo", "uuid"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_byoip_prefix.foobar", "prefix", prefixCIDR),
					resource.TestCheckResourceAttr(
						"data.digitalocean_byoip_prefix.foobar", "region", region),
					resource.TestCheckResourceAttr(
						"data.digitalocean_byoip_prefix.foobar", "advertised", "false"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_byoip_prefix.foobar", "status"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanBYOIPPrefixExists(n string, prefix *godo.BYOIPPrefix) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BYOIP prefix ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundPrefix, _, err := client.BYOIPPrefixes.Get(context.Background(), rs.Primary.Attributes["uuid"])
		if err != nil {
			return err
		}

		*prefix = *foundPrefix

		return nil
	}
}

func testAccCheckDataSourceDigitalOceanBYOIPPrefixConfig_basic(
	prefixCIDR, region string,
	includeDataSource bool,
) string {
	config := fmt.Sprintf(`
resource "digitalocean_byoip_prefix" "foo" {
  prefix     = "%s"
  signature  = "test-signature-data"
  region     = "%s"
  advertised = false
}
`, prefixCIDR, region)

	if includeDataSource {
		config += `
data "digitalocean_byoip_prefix" "foobar" {
  uuid = digitalocean_byoip_prefix.foo.uuid
}
`
	}

	return config
}
