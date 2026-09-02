package microdroplet_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanMicroDroplet_importBasic(t *testing.T) {
	resourceName := "digitalocean_microdroplet.foobar"
	name := acceptance.RandomTestName()
	config := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletImage)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			{Config: config},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// `state` is settable-with-default; the API does not persist
				// user intent so we can't verify it round-trips exactly.
				ImportStateVerifyIgnore: []string{"state"},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "not-a-real-uuid",
				ExpectError:       regexp.MustCompile(`(not found|Cannot import non-existent remote object)`),
			},
		},
	})
}

func TestAccDigitalOceanMicroDropletImage_importBasic(t *testing.T) {
	resourceName := "digitalocean_microdroplet_image.foobar"
	name := acceptance.RandomTestName()
	config := fmt.Sprintf(testAccMicroDropletImageConfig_Basic, name, testMicroDropletImage)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletImageDestroy,
		Steps: []resource.TestStep{
			{Config: config},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
