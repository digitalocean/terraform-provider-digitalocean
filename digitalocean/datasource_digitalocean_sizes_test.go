package digitalocean

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
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

func TestFilterDigitalOceanSizes(t *testing.T) {
	testCases := []struct {
		name         string
		filter       commonFilter
		expectations []string // Expectations are filled with the expected size slugs in order
	}{
		{"BySlug", commonFilter{"slug", []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByMemory", commonFilter{"memory", []string{"1024", "8192"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByCPU", commonFilter{"vcpus", []string{"1", "4"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByDisk", commonFilter{"disk", []string{"25", "160"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByTransfer", commonFilter{"transfer", []string{"1.0", "5.0"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByPriceMonthly", commonFilter{"price_monthly", []string{"5.0", "40.0"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByPriceHourly", commonFilter{"price_hourly", []string{"0.00744", "0.05952"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByRegions", commonFilter{"regions", []string{"sgp1", "ams2"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByAvailable", commonFilter{"available", []string{"true"}}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sizes := filterDigitalOceanSizes(sizesTestData(), []commonFilter{testCase.filter})
			if len(sizes) != len(testCase.expectations) {
				t.Fatalf("Expecting %d size results, found %d size results instead", len(testCase.expectations), len(sizes))
			}
			for i, expectedSlug := range testCase.expectations {
				if sizes[i].Slug != expectedSlug {
					t.Fatalf("Expecting size index %d to be %s, found %s instead", i, expectedSlug, sizes[i].Slug)
				}
			}
		})
	}
}

func TestSortDigitalOceanSizes(t *testing.T) {
	testCases := []struct {
		name        string
		key         string
		expectedAsc []string // Expected sizes if sorted ascendingly
	}{
		{"BySlug", "slug", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByMemory", "memory", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByCPU", "vcpus", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByDisk", "disk", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByTransfer", "transfer", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByPriceMonthly", "price_monthly", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByPriceHourly", "price_hourly", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Test ascending order
			sizes := sortDigitalOceanSizes(sizesTestData(), []commonSort{{testCase.key, "asc"}})
			if len(sizes) != len(testCase.expectedAsc) {
				t.Fatalf("Expecting %d size results, found %d size results instead", len(testCase.expectedAsc), len(sizes))
			}
			for i, expectedSlug := range testCase.expectedAsc {
				if sizes[i].Slug != expectedSlug {
					t.Fatalf("Expecting size index %d to be %s, found %s instead", i, expectedSlug, sizes[i].Slug)
				}
			}

			// Test descending order
			sizes = sortDigitalOceanSizes(sizesTestData(), []commonSort{{testCase.key, "desc"}})
			if len(sizes) != len(testCase.expectedAsc) {
				t.Fatalf("Expecting %d size results, found %d size results instead", len(testCase.expectedAsc), len(sizes))
			}
			for i, expectedSlug := range testCase.expectedAsc {
				if sizes[len(sizes)-i-1].Slug != expectedSlug {
					t.Fatalf("Expecting size index %d to be %s, found %s instead", i, expectedSlug, sizes[i].Slug)
				}
			}
		})
	}
}

func TestSortMultipleDigitalOceanSizes(t *testing.T) {
	// Test ascending order
	sizes := sortDigitalOceanSizes(
		sizesTestDataForTestMultipleSort(),
		[]commonSort{
			{"memory", "desc"}, // Sort by memory descendingly first
			{"disk", "asc"},    // Then for sizes with same memory, sort by disk ascendingly
		},
	)

	if len(sizes) != 3 {
		t.Fatalf("Expecting 3 size results, found %d size results instead", len(sizes))
	}

	// s-2vcpu-2gb 	(Memory = 2048)
	// s-1vcpu-1gb 	(Memory = 1024, Disk = 25)
	// 1gb			(Memory = 1024, Disk = 30)
	if sizes[0].Slug != "s-2vcpu-2gb" ||
		sizes[1].Slug != "s-1vcpu-1gb" ||
		sizes[2].Slug != "1gb" {
		t.Fatalf("Expecting sizes to be sorted by memory in descending order, then by disk in ascending order")
	}
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

func sizesTestData() []godo.Size {
	return []godo.Size{
		godo.Size{
			Slug:         "s-1vcpu-1gb",
			Memory:       1024,
			Vcpus:        1,
			Disk:         25,
			Transfer:     1.0,
			PriceMonthly: 5.0,
			PriceHourly:  0.007439999841153622,
			Regions:      []string{"sgp1", "sgp2"},
			Available:    true,
		},
		godo.Size{
			Slug:         "s-2vcpu-2gb",
			Memory:       2048,
			Vcpus:        2,
			Disk:         60,
			Transfer:     3.0,
			PriceMonthly: 15.0,
			PriceHourly:  0.02232000045478344,
			Regions:      []string{"nyc1", "nyc2"},
			Available:    false,
		},
		godo.Size{
			Slug:         "s-4vcpu-8gb",
			Memory:       8192,
			Vcpus:        4,
			Disk:         160,
			Transfer:     5.0,
			PriceMonthly: 40.0,
			PriceHourly:  0.05951999872922897,
			Regions:      []string{"ams1", "ams2"},
			Available:    true,
		},
	}
}

func sizesTestDataForTestMultipleSort() []godo.Size {
	return []godo.Size{
		godo.Size{
			Slug:         "s-1vcpu-1gb",
			Memory:       1024,
			Vcpus:        1,
			Disk:         25,
			Transfer:     1.0,
			PriceMonthly: 5.0,
			PriceHourly:  0.007439999841153622,
			Regions:      []string{"sgp1", "sgp2"},
			Available:    true,
		},
		godo.Size{
			Slug:         "1gb",
			Memory:       1024,
			Vcpus:        1,
			Disk:         30,
			Transfer:     2.0,
			PriceMonthly: 10.0,
			PriceHourly:  0.01487999968230724,
			Regions:      []string{"sgp1", "sgp2"},
			Available:    true,
		},
		godo.Size{
			Slug:         "s-2vcpu-2gb",
			Memory:       2048,
			Vcpus:        2,
			Disk:         60,
			Transfer:     3.0,
			PriceMonthly: 15.0,
			PriceHourly:  0.02232000045478344,
			Regions:      []string{"nyc1", "nyc2"},
			Available:    false,
		},
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
