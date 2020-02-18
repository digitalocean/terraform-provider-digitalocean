package digitalocean

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDigitalOceanRegion_Basic(t *testing.T) {
	config := `
data "digitalocean_region" "lon1" {
	slug = "lon1"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_region.lon1", "slug", "lon1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_region.lon1", "name"),
					resource.TestCheckResourceAttrSet("data.digitalocean_region.lon1", "available"),
					resource.TestCheckResourceAttrSet("data.digitalocean_region.lon1", "sizes.#"),
					resource.TestCheckResourceAttrSet("data.digitalocean_region.lon1", "features.#"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRegion_MissingSlug(t *testing.T) {
	config := `
data "digitalocean_region" "xyz5" {
	slug = "xyz5"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Region does not exist: xyz5"),
			},
		},
	})
}
