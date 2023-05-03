package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDatabaseUser_Basic(t *testing.T) {
	var user godo.DatabaseUser

	databaseName := acceptance.RandomTestName()
	userName := acceptance.RandomTestName()

	resourceConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseUserConfigBasic, databaseName, userName)
	datasourceConfig := fmt.Sprintf(testAccCheckDigitalOceanDatasourceDatabaseUserConfigBasic, userName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseUserExists("digitalocean_database_user.foobar_user", &user),
					testAccCheckDigitalOceanDatabaseUserAttributes(&user, userName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_user.foobar_user", "name", userName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_user.foobar_user", "password"),
				),
			},
			{
				Config: resourceConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("digitalocean_database_user.foobar_user", "name",
						"data.digitalocean_database_user.foobar_user", "name"),
					resource.TestCheckResourceAttrPair("digitalocean_database_user.foobar_user", "role",
						"data.digitalocean_database_user.foobar_user", "role"),
					resource.TestCheckResourceAttrPair("digitalocean_database_user.foobar_user", "password",
						"data.digitalocean_database_user.foobar_user", "password"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatasourceDatabaseUserConfigBasic = `
data "digitalocean_database_user" "foobar_user" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}`
