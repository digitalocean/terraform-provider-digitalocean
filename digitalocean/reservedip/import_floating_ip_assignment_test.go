package reservedip_test

import (
	"context"
	"strconv"
	"testing"

	"fmt"
	"regexp"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanFloatingIPAssignment_importBasic(t *testing.T) {
	resourceName := "digitalocean_floating_ip_assignment.foobar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPAttachmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccFLIPAssignmentImportID(resourceName),
				// floating_ip_assignments are "virtual" resources that have unique, timestamped IDs.
				// As the imported one will have a different ID that the initial one, the states will not match.
				// Verify the attachment is correct using an ImportStateCheck function instead.
				ImportStateVerify: false,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return fmt.Errorf("expected 1 state: %+v", s)
					}

					rs := s[0]
					flipID := rs.Attributes["ip_address"]
					dropletID, err := strconv.Atoi(rs.Attributes["droplet_id"])
					if err != nil {
						return err
					}

					client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
					foundFloatingIP, _, err := client.FloatingIPs.Get(context.Background(), flipID)
					if err != nil {
						return err
					}

					if foundFloatingIP.IP != flipID || foundFloatingIP.Droplet.ID != dropletID {
						return fmt.Errorf("wrong floating IP attachment found")
					}

					return nil
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "192.0.2.1",
				ExpectError:       regexp.MustCompile("joined with a comma"),
			},
		},
	})
}

func testAccFLIPAssignmentImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		flip := rs.Primary.Attributes["ip_address"]
		droplet := rs.Primary.Attributes["droplet_id"]

		return fmt.Sprintf("%s,%s", flip, droplet), nil
	}
}
