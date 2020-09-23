package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDroplet_importBasic(t *testing.T) {
	resourceName := "digitalocean_droplet.foobar"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ssh_keys", "user_data", "resize_disk"}, //we ignore the ssh_keys, resize_disk and user_data as we do not set to state
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "123",
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_ImportWithNoImageSlug(t *testing.T) {
	rInt := acctest.RandInt()
	var droplet godo.Droplet
	var snapshotId []int

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					takeDropletSnapshot(rInt, &droplet, &snapshotId),
				),
			},
			{
				Config: testAccCheckDigitalOceanDropletConfig_fromSnapshot(rInt),
			},
			{
				ResourceName:      "digitalocean_droplet.from-snapshot",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ssh_keys", "user_data", "resize_disk"}, //we ignore the ssh_keys, resize_disk and user_data as we do not set to state
			},
			{
				Config: " ",
				Check: resource.ComposeTestCheckFunc(
					deleteDropletSnapshots(&snapshotId),
				),
			},
		},
	})
}

func takeDropletSnapshot(rInt int, droplet *godo.Droplet, snapshotId *[]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		action, _, err := client.DropletActions.Snapshot(context.Background(), (*droplet).ID, fmt.Sprintf("snap-%d", rInt))
		if err != nil {
			return err
		}
		waitForAction(client, action)

		retrieveDroplet, _, err := client.Droplets.Get(context.Background(), (*droplet).ID)
		if err != nil {
			return err
		}
		*snapshotId = retrieveDroplet.SnapshotIDs
		return nil
	}
}

func testAccCheckDigitalOceanDropletConfig_fromSnapshot(rInt int) string {
	return fmt.Sprintf(`
data "digitalocean_image" "snapshot" {
  name = "snap-%d"
}

resource "digitalocean_droplet" "from-snapshot" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "${data.digitalocean_image.snapshot.id}"
  region    = "nyc3"
  user_data = "foobar"
}`, rInt, rInt)
}
