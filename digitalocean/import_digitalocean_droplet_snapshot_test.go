package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDigitalOceanDropletSnapshot_importBasic(t *testing.T) {
	resourceName := "digitalocean_droplet_snapshot.foobar"
	rInt1 := acctest.RandInt()
	rInt2 := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDropletSnapshotConfig_basic, rInt1, rInt2),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// image_id oposite than image isnt store if no slug
				ImportStateVerifyIgnore: []string{
					"image_id"},
			},
		},
	})
}
