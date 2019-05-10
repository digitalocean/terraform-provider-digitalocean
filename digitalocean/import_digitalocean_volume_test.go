package digitalocean

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDigitalOceanVolume_importBasic(t *testing.T) {
	resourceName := "digitalocean_volume.foobar"
	volumeName := fmt.Sprintf("volume-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeConfig_basic, volumeName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "123abc",
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}
