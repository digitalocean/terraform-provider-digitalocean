package account_test

import (
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanAccount_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanAccountConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_account.foobar", "uuid"),
				),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanAccountConfig_basic = `
data "digitalocean_account" "foobar" {
}`
