package reservedipv6_test

import (
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanReservedIPV6_importBasicRegion(t *testing.T) {
	resourceName := "digitalocean_reserved_ipv6.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPV6Config_regionSlug,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDigitalOceanReservedIPV6_importBasicDroplet(t *testing.T) {
	resourceName := "digitalocean_reserved_ipv6.foobar"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPV6Config_droplet(name),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
