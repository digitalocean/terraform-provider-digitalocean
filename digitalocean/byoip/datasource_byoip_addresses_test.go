package byoip_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanBYOIPAddresses_Basic(t *testing.T) {
	var prefix godo.BYOIPPrefix
	prefixCIDR := "192.0.2.0/24"
	region := "nyc3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanBYOIPAddressesConfig_basic(prefixCIDR, region, false),
			},
			{
				Config: testAccCheckDataSourceDigitalOceanBYOIPAddressesConfig_basic(prefixCIDR, region, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanBYOIPPrefixExists("data.digitalocean_byoip_prefix.foobar", &prefix),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_byoip_addresses.foobar", "addresses.#"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanBYOIPAddressesConfig_basic(
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

data "digitalocean_byoip_prefix" "foobar" {
  uuid = digitalocean_byoip_prefix.foo.uuid
}
`, prefixCIDR, region)

	if includeDataSource {
		config += `
data "digitalocean_byoip_addresses" "foobar" {
  byoip_prefix_uuid = data.digitalocean_byoip_prefix.foobar.uuid
}
`
	}

	return config
}
