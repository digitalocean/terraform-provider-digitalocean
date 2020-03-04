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
					testAccCheckDataSourceDigitalOceanSizesExist("data.digitalocean_sizes.foobar"),
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
					testAccCheckDataSourceDigitalOceanSizesExist("data.digitalocean_sizes.foobar"),
					testAccCheckDataSourceDigitalOceanSizesFilteredAndSorted("data.digitalocean_sizes.foobar"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanSizesExist(n string) resource.TestCheckFunc {
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

func testAccCheckDataSourceDigitalOceanSizesFilteredAndSorted(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		rawTotal := rs.Primary.Attributes["sizes.#"]
		total, err := strconv.Atoi(rawTotal)
		if err != nil {
			return err
		}

		stringInSlice := func(value string, slice []string) bool {
			for _, item := range slice {
				if item == value {
					return true
				}
			}
			return false
		}

		var prevSlug string
		var prevPriceMonthly float64
		for i := 0; i < total; i++ {
			slug := rs.Primary.Attributes[fmt.Sprintf("sizes.%d.slug", i)]
			if !stringInSlice(slug, []string{"s-1vcpu-1gb", "s-1vcpu-2gb", "s-2vcpu-2gb", "s-3vcpu-1gb"}) {
				return fmt.Errorf("Slug is not in expected test filter values")
			}
			if prevSlug != "" && prevSlug < slug {
				return fmt.Errorf("Sizes is not sorted by slug in descending order")
			}
			prevSlug = slug

			vcpus := rs.Primary.Attributes[fmt.Sprintf("sizes.%d.vcpus", i)]
			if !stringInSlice(vcpus, []string{"1", "2"}) {
				return fmt.Errorf("Virtual CPU is not in expected test filter values")
			}

			priceMonthly, _ := strconv.ParseFloat(rs.Primary.Attributes[fmt.Sprintf("sizes.%d.price_monthly", i)], 64)
			if prevPriceMonthly > 0 && prevPriceMonthly < priceMonthly {
				return fmt.Errorf("Sizes is not sorted by price monthly in descending order")
			}
			prevPriceMonthly = priceMonthly
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
