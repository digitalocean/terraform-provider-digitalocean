package digitalocean

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceDigitalOceanSizes_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanSizesConfigBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanSizesHasAtLeastOne("data.digitalocean_sizes.foobar"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanSizesHasAtLeastOne("data.digitalocean_sizes.foobar"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanSizesHasAtLeastOne(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		rawTotal := rs.Primary.Attributes["sizes.#"]
		total, err := strconv.Atoi(rawTotal)
		if err != nil {
			return err
		}

		if total < 1 {
			return fmt.Errorf("No digital ocean sizes retrieved")
		}

		return nil
	}
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
