package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDigitalOceanSizes_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanSizesConfigBasic),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSizes_WithFilterAndSort(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanSizesConfigWithFilterAndSort),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanSizesConfigBasic = `
data "digitalocean_sizes" "foobar" {
}`

const testAccCheckDataSourceDigitalOceanSizesConfigWithFilterAndSort = `
data "digitalocean_sizes" "foobar" {
	filter {
		key 	= "slug"
		values 	= ["s-1vcpu-1gb", "s-1vcpu-2gb", "s-2vcpu-2gb", "s-3vcpu-1gb"]
	}

	filter {
		key 	= "vcpus"
		values 	= ["1", "2"]
	}

	sort {
		key 		= "price_monthly"
		direction 	= "desc"
	}

	sort {
		key 		= "slug"
		direction 	= "desc"
	}
}`
