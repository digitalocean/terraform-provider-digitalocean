package digitalocean

import (
	"testing"

	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseFirewall_importBasic(t *testing.T) {
	resourceName := "digitalocean_database_firewall.example"
	databaseClusterName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseFirewallConfigBasic, databaseClusterName),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// Needs the cluster's ID, not the firewall's
				ImportStateIdFunc: testAccDatabaseFirewallImportID(resourceName),
				// We do not have a way to match IDs as we the cluster's ID
				// with resource.PrefixedUniqueId() for the DatabaseFirewall resource.
				// ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return fmt.Errorf("expected 1 state: %#v", s)
					}
					rs := s[0]

					clusterId := rs.Attributes["cluster_id"]
					if !strings.HasPrefix(rs.ID, clusterId) {
						return fmt.Errorf("expected ID to be set and begin with %s-, received: %s\n %#v",
							clusterId, rs.ID, rs.Attributes)
					}

					if rs.Attributes["rule.#"] != "1" {
						return fmt.Errorf("expected 1 rule, received: %s\n %#v",
							rs.Attributes["rule.#"], rs.Attributes)
					}

					return nil
				},
			},
		},
	})
}

func testAccDatabaseFirewallImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		return rs.Primary.Attributes["cluster_id"], nil
	}
}
