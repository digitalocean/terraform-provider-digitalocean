package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanFloatingIP_importBasicRegion(t *testing.T) {
	resourceName := "digitalocean_floating_ip.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_region,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDigitalOceanFloatingIP_importBasicDroplet(t *testing.T) {
	resourceName := "digitalocean_floating_ip.foobar"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_droplet(rInt),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
