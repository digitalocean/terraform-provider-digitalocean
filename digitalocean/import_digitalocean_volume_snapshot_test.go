package digitalocean

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDigitalOceanVolumeSnapshot_importBasic(t *testing.T) {
	resourceName := "digitalocean_volume_snapshot.foobar"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeSnapshotConfig_basic, rInt, rInt),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
