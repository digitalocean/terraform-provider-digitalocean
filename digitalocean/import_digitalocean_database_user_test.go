package digitalocean

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseUser_importBasic(t *testing.T) {
	resourceName := "digitalocean_database_user.foobar_user"
	databaseClusterName := randomTestName()
	databaseUserName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigBasic, databaseClusterName, databaseUserName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Requires passing both the cluster ID and user name
				ImportStateIdFunc: testAccDatabaseUserImportID(resourceName),
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     fmt.Sprintf("%s,%s", "this-cluster-id-does-not-exist", databaseUserName),
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}

func testAccDatabaseUserImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		return fmt.Sprintf("%s,%s", clusterID, name), nil
	}
}
