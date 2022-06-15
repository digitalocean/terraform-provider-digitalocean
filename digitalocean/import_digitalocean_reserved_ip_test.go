package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanReservedIP_importBasicRegion(t *testing.T) {
	resourceName := "digitalocean_reserved_ip.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPConfig_region,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDigitalOceanReservedIP_importBasicDroplet(t *testing.T) {
	resourceName := "digitalocean_reserved_ip.foobar"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPConfig_droplet(rInt),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
