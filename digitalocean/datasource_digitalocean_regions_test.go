package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccCheckDataSourceDigitalOceanRegionsConfig_basic = `
data "digitalocean_regions" "foobar" {
}
`

func TestAccDigitalOceanRegions_importBasic(t *testing.T) {
	resourceName := "digitalocean_regions.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanRegionsConfig_basic),
			}, {
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
