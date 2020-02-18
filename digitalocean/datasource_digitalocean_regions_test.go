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
	filter {
        key = "available"
        values = ["true"]
    }
    sort {
		key = "slug"
    }
}
`

	configFeaturesFilter := `
data "digitalocean_regions" "filtered" {
	filter {
        key = "features"
        values = ["private_networking", "backups"]
    }
    sort {
		key = "available"
		direction = "desc"
    }
}
`

	configAllFilters := `
data "digitalocean_regions" "filtered" {
	filter {
        key = "available"
        values = ["true"]
    }
	filter {
        key = "features"
        values = ["private_networking", "backups"]
    }
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configNoFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.all", "regions.#"),
				),
			},
			{
				Config: configAvailableFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.filtered", "regions.#"),
				),
			},
			{
				Config: configFeaturesFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.filtered", "regions.#"),
				),
			},
			{
				Config: configAllFilters,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.filtered", "regions.#"),
				),
			},
		},
	})
}
