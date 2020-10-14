package digitalocean

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanImage_Basic(t *testing.T) {
	var droplet godo.Droplet
	var snapshotsId []int
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					takeSnapshotsOfDroplet(rInt, &droplet, &snapshotsId),
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
					deleteDropletSnapshots(&snapshotsId),
				),
			},
		},
	})
}

func TestAccDigitalOceanImage_PublicSlug(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
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

func takeSnapshotsOfDroplet(rInt int, droplet *godo.Droplet, snapshotsId *[]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()
		for i := 0; i < 3; i++ {
			err := takeSnapshotOfDroplet(rInt, i%2, droplet)
			if err != nil {
				return err
			}
		}
		retrieveDroplet, _, err := client.Droplets.Get(context.Background(), (*droplet).ID)
		if err != nil {
			return err
		}
		*snapshotsId = retrieveDroplet.SnapshotIDs
		return nil
	}
}

func takeSnapshotOfDroplet(rInt, sInt int, droplet *godo.Droplet) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()
	action, _, err := client.DropletActions.Snapshot(context.Background(), (*droplet).ID, fmt.Sprintf("snap-%d-%d", rInt, sInt))
	if err != nil {
		return err
	}
	waitForAction(client, action)
	return nil
}

func deleteDropletSnapshots(snapshotsId *[]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("Deleting Droplet snapshots")

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		snapshots := *snapshotsId
		for _, value := range snapshots {
			log.Printf("Deleting %d", value)
			_, err := client.Images.Delete(context.Background(), value)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckDigitalOceanImageConfig_basic(rInt, sInt int) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  name               = "snap-%d-%d"
}
`, rInt, sInt)
}

func testAccCheckDigitalOceanImageConfig_nonexisting(rInt int) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  name               = "snap-%d-nonexisting"
}
`, rInt)
}

func testAccCheckDigitalOceanImageConfig_slug(slug string) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  slug               = "%s"
}
`, slug)
}
