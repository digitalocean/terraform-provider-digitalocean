package digitalocean

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDomain_importBasic(t *testing.T) {
	resourceName := "digitalocean_domain.foobar"
	domainName := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDomainConfig_basic, domainName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ip_address"}, //we ignore the IP Address as we do not set to state
			},
		},
	})
}
