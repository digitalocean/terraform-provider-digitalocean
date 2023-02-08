package domain_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDomains_Basic(t *testing.T) {
	name1 := acceptance.RandomTestName() + ".com"
	name2 := acceptance.RandomTestName() + ".com"

	resourcesConfig := fmt.Sprintf(`
resource "digitalocean_domain" "foo" {
  name = "%s"
}

resource "digitalocean_domain" "bar" {
  name = "%s"
}
`, name1, name2)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_domains" "result" {
  filter {
    key    = "name"
    values = ["%s"]
  }
}
`, name1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_domains.result", "domains.#", "1"),
					resource.TestCheckResourceAttrPair("data.digitalocean_domains.result", "domains.0.name", "digitalocean_domain.foo", "name"),
					resource.TestCheckResourceAttrPair("data.digitalocean_domains.result", "domains.0.urn", "digitalocean_domain.foo", "urn"),
					resource.TestCheckResourceAttrPair("data.digitalocean_domains.result", "domains.0.ttl", "digitalocean_domain.foo", "ttl"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
