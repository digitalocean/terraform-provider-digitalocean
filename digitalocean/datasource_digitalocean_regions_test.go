package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDigitalOceanRegions_Basic(t *testing.T) {
	configNoFilter := `
data "digitalocean_regions" "all" {
}
`
	configAvailableFilter := `
data "digitalocean_regions" "filtered" {
	available = true
}
`

	configFeaturesFilter := `
data "digitalocean_regions" "filtered" {
	features = ["private_networking", "backups"]
}
`

	configAllFilters := `
data "digitalocean_regions" "filtered" {
	available = true
	features = ["private_networking", "backups"]
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configNoFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.all", "slugs.#"),
				),
			},
			{
				Config: configAvailableFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.filtered", "slugs.#"),
				),
			},
			{
				Config: configFeaturesFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.filtered", "slugs.#"),
				),
			},
			{
				Config: configAllFilters,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.filtered", "slugs.#"),
				),
			},
		},
	})
}
