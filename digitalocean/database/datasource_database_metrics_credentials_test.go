package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDatabaseMetricsCredentials(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()
	databaseConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
				),
			},
			{
				Config: databaseConfig + testAccCheckDigitalOceanDatasourceMetricsCredentialsConfigNew,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_database_metrics_credentials.creds", "username"),
					resource.TestCheckResourceAttrSet("data.digitalocean_database_metrics_credentials.creds", "password"),
				),
			},
		},
	})
}

const (
	testAccCheckDigitalOceanDatasourceMetricsCredentialsConfigNew = `
data "digitalocean_database_metrics_credentials" "creds" {}`
)
