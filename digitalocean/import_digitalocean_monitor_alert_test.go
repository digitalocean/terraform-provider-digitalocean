package digitalocean

// import (
// 	"fmt"
// 	"regexp"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// )

// func TestAccDigitalOceanDatabaseMonitorAlert_importBasic(t *testing.T) {
// 	resourceName := "digitalocean_monitor_alert.cpu_alert"

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheck(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckDigitalOceanDatabaseDBDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseDBConfigBasic, databaseClusterName, databaseDBName),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 				// Requires passing both the cluster ID and DB name
// 				ImportStateIdFunc: testAccDatabaseDBImportID(resourceName),
// 			},
// 			// Test importing non-existent resource provides expected error.
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: false,
// 				ImportStateId:     fmt.Sprintf("%s,%s", "this-cluster-id-does-not-exist", databaseDBName),
// 				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
// 			},
// 		},
// 	})
// }
