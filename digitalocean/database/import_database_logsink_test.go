package database_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseLogsinkRsyslog_ImportBasic(t *testing.T) {
	resourceName := "digitalocean_database_logsink_rsyslog.test"
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigBasic, clusterName, logsinkName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Import requires cluster_id,logsink_id format
				ImportStateIdFunc: testAccDatabaseLogsinkImportID(resourceName),
			},
			// Test importing non-existent resource provides expected error
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     fmt.Sprintf("%s,%s", "this-cluster-id-does-not-exist", "non-existent-logsink-id"),
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
			// Test invalid ID format
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "invalid-id-format",
				ExpectError:       regexp.MustCompile("must use the format 'cluster_id,logsink_id' for import"),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseLogsinkOpensearch_ImportBasic(t *testing.T) {
	resourceName := "digitalocean_database_logsink_opensearch.test"
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigBasic, clusterName, logsinkName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Import requires cluster_id,logsink_id format
				ImportStateIdFunc: testAccDatabaseLogsinkImportID(resourceName),
			},
			// Test importing non-existent resource provides expected error
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     fmt.Sprintf("%s,%s", "this-cluster-id-does-not-exist", "non-existent-logsink-id"),
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
			// Test invalid ID format
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "invalid-id-format",
				ExpectError:       regexp.MustCompile("must use the format 'cluster_id,logsink_id' for import"),
			},
		},
	})
}

func testAccDatabaseLogsinkImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		logsinkID := rs.Primary.Attributes["logsink_id"]

		if clusterID == "" {
			return "", fmt.Errorf("cluster_id not found in resource state")
		}

		if logsinkID == "" {
			return "", fmt.Errorf("logsink_id not found in resource state")
		}

		return fmt.Sprintf("%s,%s", clusterID, logsinkID), nil
	}
}
