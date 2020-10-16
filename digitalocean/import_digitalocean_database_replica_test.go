package digitalocean

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseReplica_importBasic(t *testing.T) {
	var database godo.Database
	resourceName := "digitalocean_database_replica.read-01"
	databaseName := randomTestName()
	databaseReplicaName := randomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseReplicaConfigBasic, databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
				),
			},
			{
				Config: databaseConfig + replicaConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Requires passing both the cluster ID and replica name
				ImportStateIdFunc: testAccDatabaseReplicaImportID(resourceName),
				// The DO API does not return the size on read, but it is required on create
				ImportStateVerifyIgnore: []string{"size"},
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     fmt.Sprintf("%s,%s", "this-cluster-id-does-not-exist", databaseReplicaName),
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}

func testAccDatabaseReplicaImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		return fmt.Sprintf("%s,%s", clusterId, name), nil
	}
}
