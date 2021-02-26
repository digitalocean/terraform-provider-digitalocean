package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanFirewalls_Basic(t *testing.T) {

	config := `
resource "digitalocean_firewall" "test" {
	name = "tf-test"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		Steps:           []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_firewalls.tf-test", "firewalls.#"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.test", "name", "tf-test"),
				),
			},
		},
	})
}