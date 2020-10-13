package digitalocean

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseDB_importBasic(t *testing.T) {
	resourceName := "digitalocean_database_db.foobar_db"
	databaseClusterName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))
	databaseDBName := fmt.Sprintf("foobar-test-db-terraform-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseDBConfigBasic, databaseClusterName, databaseDBName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Requires passing both the cluster ID and DB name
				ImportStateIdFunc: testAccDatabaseDBImportID(resourceName),
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     fmt.Sprintf("%s,%s", "this-cluster-id-does-not-exist", databaseDBName),
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}

func testAccDatabaseDBImportID(n string) resource.ImportStateIdFunc {
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
