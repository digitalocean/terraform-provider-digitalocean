package image_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanImage_Basic(t *testing.T) {
	var droplet godo.Droplet
	var snapshotsId []int
	snapName := acceptance.RandomTestName()
	dropletName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			// Creates a Droplet and takes multiple snapshots of it.
			// One will have the suffix -1 and two will have -0
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(dropletName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					acceptance.TakeSnapshotsOfDroplet(snapName, &droplet, &snapshotsId),
				),
			},
			// Find snapshot with suffix -1
			{
				Config: testAccCheckDigitalOceanImageConfig_basic(snapName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "name", snapName+"-1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "min_disk_size", "25"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "private", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "type", "snapshot"),
				),
			},
			// Expected error with  suffix -0 as multiple exist
			{
				Config:      testAccCheckDigitalOceanImageConfig_basic(snapName, 0),
				ExpectError: regexp.MustCompile(`.*too many images found with name tf-acc-test-.*\ .found 2, expected 1.`),
			},
			{
				Config:      testAccCheckDigitalOceanImageConfig_nonexisting(snapName),
				Destroy:     false,
				ExpectError: regexp.MustCompile(`.*no image found with name tf-acc-test-.*-nonexisting`),
			},
			{
				Config: " ",
				Check: resource.ComposeTestCheckFunc(
					acceptance.DeleteDropletSnapshots(&snapshotsId),
				),
			},
		},
	})
}

func TestAccDigitalOceanImage_PublicSlug(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanImageConfig_slug("ubuntu-22-04-x64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "slug", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "min_disk_size", "7"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "private", "false"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "type", "snapshot"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "distribution", "Ubuntu"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanImageConfig_basic(name string, sInt int) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  name = "%s-%d"
}
`, name, sInt)
}

func testAccCheckDigitalOceanImageConfig_nonexisting(name string) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  name = "%s-nonexisting"
}
`, name)
}

func testAccCheckDigitalOceanImageConfig_slug(slug string) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  slug = "%s"
}
`, slug)
}
