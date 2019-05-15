package digitalocean

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDigitalOceanRecord_importBasic(t *testing.T) {
	resourceName := "digitalocean_record.foobar"
	domainName := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanRecordConfig_basic, domainName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Requires passing both the ID and domain
				ImportStateIdPrefix: fmt.Sprintf("%s,", domainName),
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateVerify:   false,
				ImportStateIdPrefix: fmt.Sprintf("%s,", "nonexistent.com"),
				ExpectError:         regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}
