package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseLogsink_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterPostgreSQL, name, "15")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkBasic, dbConfig, "lname", "opensearch", "https://user:passwd@192.168.0.1:25060", "logs", 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_logsink.logsink", "opensearch_config.0.url", "https://user:passwd@192.168.0.1:25060"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink.logsink", "opensearch_config.0.index_prefix", "logs"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink.logsink", "opensearch_config.0.index_days_max", "5"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink.logsink", "sink_type", "opensearch"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink.logsink", "sink_name", "lname"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkBasic, dbConfig, "new-lname", "opensearch", "https://user:passwd@192.168.0.1:25060", "logs", 4),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_logsink.logsink", "sink_name", "new-lname"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink.logsink", "opensearch_config.0.index_days_max", "4"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseLogsinkBasic = `
%s

resource "digitalocean_database_logsink" "logsink" {
  cluster_id = digitalocean_database_cluster.foobar.id
  sink_name  = "%s"
  sink_type  = "%s"

  opensearch_config {
    url            = "%s"
    index_prefix   = "%s"
    index_days_max = %d
    timeout        = 10
  }
}`
