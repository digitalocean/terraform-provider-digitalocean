package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseOpensearchConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterOpensearch, name, "2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseOpensearchConfigConfigBasic, dbConfig, true, 10, "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "enable_security_audit", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "ism_enabled", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "ism_history_enabled", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "ism_history_max_age_hours", "10"),
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "ism_history_max_docs", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseOpensearchConfigConfigBasic, dbConfig, false, 1, "1000000000000000000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "enable_security_audit", "false"),
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "ism_enabled", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "ism_history_enabled", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "ism_history_max_age_hours", "1"),
					resource.TestCheckResourceAttr("digitalocean_database_opensearch_config.foobar", "ism_history_max_docs", "1000000000000000000"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseOpensearchConfigConfigBasic = `
%s

resource "digitalocean_database_opensearch_config" "foobar" {
  cluster_id                = digitalocean_database_cluster.foobar.id
  enable_security_audit     = %t
  ism_enabled               = true
  ism_history_enabled       = true
  ism_history_max_age_hours = %d
  ism_history_max_docs      = %s
}`
