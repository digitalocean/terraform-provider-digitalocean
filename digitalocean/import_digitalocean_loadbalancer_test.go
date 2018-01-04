package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDigitalOceanLoadBalancer_importBasic(t *testing.T) {
	resourceName := "digitalocean_loadbalancer.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
