package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanTag_importBasic(t *testing.T) {
	resourceName := "digitalocean_tag.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanTagConfig_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
