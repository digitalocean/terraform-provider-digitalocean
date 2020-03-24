package digitalocean

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDigitalOceanVPC_importBasic(t *testing.T) {
	resourceName := "digitalocean_vpc.foobar"
	vpcName := randomTestName()
	vpcCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanVPCConfig_Basic, vpcName, "A description for the VPC")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: vpcCreateConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: fix when godo is updated
				ImportStateVerifyIgnore: []string{"description"},
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
