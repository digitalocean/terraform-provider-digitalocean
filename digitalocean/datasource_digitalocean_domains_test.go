package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDigitalOceanDomains_Basic(t *testing.T) {
	name1 := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))
	name2 := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resourcesConfig := fmt.Sprintf(`
resource "digitalocean_domain" "foo" {
  name     = "%s"
}

resource "digitalocean_domain" "bar" {
  name     = "%s"
}
`, name1, name2)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_domains" "result" {
  filter {
    key = "name"
    values = ["%s"]
  }
}
`, name1)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
					// skip checking `ttl` because digitalocean_domain does not expose the default TTL yet
					//resource.TestCheckResourceAttrPair("data.digitalocean_domains.result", "domains.0.ttl", "digitalocean_domain.foo", "ttl"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
