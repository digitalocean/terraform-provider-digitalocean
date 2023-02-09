package region_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanRegions_Basic(t *testing.T) {
	configNoFilter := `
data "digitalocean_regions" "all" {
}
`
	configAvailableFilter := `
data "digitalocean_regions" "filtered" {
  filter {
    key    = "available"
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
    key    = "features"
    values = ["private_networking", "backups"]
  }
  sort {
    key       = "available"
    direction = "desc"
  }
}
`

	configAllFilters := `
data "digitalocean_regions" "filtered" {
  filter {
    key    = "available"
    values = ["true"]
  }
  filter {
    key    = "features"
    values = ["private_networking", "backups"]
  }
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configNoFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.all", "regions.#"),
					resource.TestCheckResourceAttrSet("data.digitalocean_regions.all", "regions.#"),
					acceptance.TestResourceInstanceState("data.digitalocean_regions.all", func(is *terraform.InstanceState) error {
						n, err := strconv.Atoi(is.Attributes["regions.#"])
						if err != nil {
							return err
						}

						for i := 0; i < n; i++ {
							key := fmt.Sprintf("regions.%d.slug", i)
							if _, ok := is.Attributes[key]; !ok {
								return fmt.Errorf("missing key in instance state for %s in %s", key, "data.digitalocean_regions.all")
							}
						}

						return nil
					}),
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
