package image_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanImage_Basic(t *testing.T) {
	var droplet godo.Droplet
	var snapshotsId []int
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					acceptance.TakeSnapshotsOfDroplet(rInt, &droplet, &snapshotsId),
				),
			},
			{
				Config: testAccCheckDigitalOceanImageConfig_basic(rInt, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "name", fmt.Sprintf("snap-%d-1", rInt)),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "min_disk_size", "20"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "private", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "type", "snapshot"),
				),
			},
			{
				Config:      testAccCheckDigitalOceanImageConfig_basic(rInt, 0),
				ExpectError: regexp.MustCompile(`.*too many images found with name snap-.*\ .found 2, expected 1.`),
			},
			{
				Config:      testAccCheckDigitalOceanImageConfig_nonexisting(rInt),
				Destroy:     false,
				ExpectError: regexp.MustCompile(`.*no image found with name snap-.*-nonexisting`),
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
				Config: testAccCheckDigitalOceanImageConfig_slug("ubuntu-18-04-x64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "slug", "ubuntu-18-04-x64"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_image.foobar", "min_disk_size", "15"),
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

func testAccCheckDigitalOceanImageConfig_basic(rInt, sInt int) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  name = "snap-%d-%d"
}
`, rInt, sInt)
}

func testAccCheckDigitalOceanImageConfig_nonexisting(rInt int) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  name = "snap-%d-nonexisting"
}
`, rInt)
}

func testAccCheckDigitalOceanImageConfig_slug(slug string) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  slug = "%s"
}
`, slug)
}
