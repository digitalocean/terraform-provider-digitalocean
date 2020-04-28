package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceDigitalOceanSpacesBucketObjects_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigResources(rInt), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectsDataSourceExists("data.digitalocean_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.0", "arch/navajo/north_window"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.1", "arch/navajo/sand_dune"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObjects_all(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigResources(rInt), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigAll(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectsDataSourceExists("data.digitalocean_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.#", "7"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.0", "arch/courthouse_towers/landscape"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.1", "arch/navajo/north_window"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.2", "arch/navajo/sand_dune"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.3", "arch/partition/park_avenue"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.4", "arch/rubicon"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.5", "arch/three_gossips/broken"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.6", "arch/three_gossips/turret"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObjects_prefixes(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigResources(rInt), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigPrefixes(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectsDataSourceExists("data.digitalocean_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.0", "arch/rubicon"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "common_prefixes.#", "4"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "common_prefixes.0", "arch/courthouse_towers/"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "common_prefixes.1", "arch/navajo/"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "common_prefixes.2", "arch/partition/"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "common_prefixes.3", "arch/three_gossips/"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObjects_encoded(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigExtraResource(rInt), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigEncoded(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectsDataSourceExists("data.digitalocean_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.0", "arch%2Fru%20b%20ic%20on"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.1", "arch%2Frubicon"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObjects_maxKeys(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigResources(rInt), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigMaxKeys(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectsDataSourceExists("data.digitalocean_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.0", "arch/courthouse_towers/landscape"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.1", "arch/navajo/north_window"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObjects_startAfter(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigResources(rInt), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigStartAfter(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectsDataSourceExists("data.digitalocean_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.0", "arch/three_gossips/turret"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObjects_fetchOwner(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigResources(rInt), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigOwners(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectsDataSourceExists("data.digitalocean_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "owners.#", "2"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanSpacesObjectsDataSourceExists(addr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("Can't find S3 objects data source: %s", addr)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("S3 objects data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigResources(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "objects_bucket" {
  name          = "tf-acc-objects-test-bucket-%d"
  region        = "nyc3"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object1" {
  bucket  = digitalocean_spaces_bucket.objects_bucket.name
  region  = digitalocean_spaces_bucket.objects_bucket.region
  key     = "arch/three_gossips/turret"
  content = "Delicate"
}

resource "digitalocean_spaces_bucket_object" "object2" {
  bucket  = digitalocean_spaces_bucket.objects_bucket.name
  region  = digitalocean_spaces_bucket.objects_bucket.region
  key     = "arch/three_gossips/broken"
  content = "Dark Angel"
}

resource "digitalocean_spaces_bucket_object" "object3" {
  bucket  = digitalocean_spaces_bucket.objects_bucket.name
  region  = digitalocean_spaces_bucket.objects_bucket.region
  key     = "arch/navajo/north_window"
  content = "Balanced Rock"
}

resource "digitalocean_spaces_bucket_object" "object4" {
  bucket  = digitalocean_spaces_bucket.objects_bucket.name
  region  = digitalocean_spaces_bucket.objects_bucket.region
  key     = "arch/navajo/sand_dune"
  content = "Queen Victoria Rock"
}

resource "digitalocean_spaces_bucket_object" "object5" {
  bucket  = digitalocean_spaces_bucket.objects_bucket.name
  region  = digitalocean_spaces_bucket.objects_bucket.region
  key     = "arch/partition/park_avenue"
  content = "Double-O"
}

resource "digitalocean_spaces_bucket_object" "object6" {
  bucket  = digitalocean_spaces_bucket.objects_bucket.name
  region  = digitalocean_spaces_bucket.objects_bucket.region
  key     = "arch/courthouse_towers/landscape"
  content = "Fiery Furnace"
}

resource "digitalocean_spaces_bucket_object" "object7" {
  bucket  = digitalocean_spaces_bucket.objects_bucket.name
  region  = digitalocean_spaces_bucket.objects_bucket.region
  key     = "arch/rubicon"
  content = "Devils Garden"
}
`, randInt)
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigBasic(randInt int) string {
	return fmt.Sprintf(`
%s

data "digitalocean_spaces_bucket_objects" "yesh" {
  bucket    = digitalocean_spaces_bucket.objects_bucket.name
  region    = digitalocean_spaces_bucket.objects_bucket.region
  prefix    = "arch/navajo/"
  delimiter = "/"
}
`, testAccDataSourceDigitalOceanSpacesObjectsConfigResources(randInt))
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigAll(randInt int) string {
	return fmt.Sprintf(`
%s

data "digitalocean_spaces_bucket_objects" "yesh" {
  bucket = digitalocean_spaces_bucket.objects_bucket.name
  region = digitalocean_spaces_bucket.objects_bucket.region
}
`, testAccDataSourceDigitalOceanSpacesObjectsConfigResources(randInt))
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigPrefixes(randInt int) string {
	return fmt.Sprintf(`
%s

data "digitalocean_spaces_bucket_objects" "yesh" {
  bucket    = digitalocean_spaces_bucket.objects_bucket.name
  region    = digitalocean_spaces_bucket.objects_bucket.region
  prefix    = "arch/"
  delimiter = "/"
}
`, testAccDataSourceDigitalOceanSpacesObjectsConfigResources(randInt))
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigExtraResource(randInt int) string {
	return fmt.Sprintf(`
%s

resource "digitalocean_spaces_bucket_object" "object8" {
  bucket  = digitalocean_spaces_bucket.objects_bucket.name
  region  = digitalocean_spaces_bucket.objects_bucket.region
  key     = "arch/ru b ic on"
  content = "Goose Island"
}
`, testAccDataSourceDigitalOceanSpacesObjectsConfigResources(randInt))
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigEncoded(randInt int) string {
	return fmt.Sprintf(`
%s

data "digitalocean_spaces_bucket_objects" "yesh" {
  bucket        = digitalocean_spaces_bucket.objects_bucket.name
  region        = digitalocean_spaces_bucket.objects_bucket.region
  encoding_type = "url"
  prefix        = "arch/ru"
}
`, testAccDataSourceDigitalOceanSpacesObjectsConfigExtraResource(randInt))
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigMaxKeys(randInt int) string {
	return fmt.Sprintf(`
%s

data "digitalocean_spaces_bucket_objects" "yesh" {
  bucket   = digitalocean_spaces_bucket.objects_bucket.name
  region   = digitalocean_spaces_bucket.objects_bucket.region
  max_keys = 2
}
`, testAccDataSourceDigitalOceanSpacesObjectsConfigResources(randInt))
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigStartAfter(randInt int) string {
	return fmt.Sprintf(`
%s

data "digitalocean_spaces_bucket_objects" "yesh" {
  bucket      = digitalocean_spaces_bucket.objects_bucket.name
  region      = digitalocean_spaces_bucket.objects_bucket.region
  start_after = "arch/three_gossips/broken"
}
`, testAccDataSourceDigitalOceanSpacesObjectsConfigResources(randInt))
}

func testAccDataSourceDigitalOceanSpacesObjectsConfigOwners(randInt int) string {
	return fmt.Sprintf(`
%s

data "digitalocean_spaces_bucket_objects" "yesh" {
  bucket      = digitalocean_spaces_bucket.objects_bucket.name
  region      = digitalocean_spaces_bucket.objects_bucket.region
  prefix      = "arch/three_gossips/"
  fetch_owner = true
}
`, testAccDataSourceDigitalOceanSpacesObjectsConfigResources(randInt))
}
