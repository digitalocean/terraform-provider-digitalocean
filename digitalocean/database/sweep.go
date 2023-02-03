package database

import (
	"context"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("digitalocean_database_cluster", &resource.Sweeper{
		Name: "digitalocean_database_cluster",
		F:    testSweepDatabaseCluster,
	})

}

func testSweepDatabaseCluster(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	databases, _, err := client.Databases.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, db := range databases {
		if strings.HasPrefix(db.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying database cluster %s", db.Name)

			if _, err := client.Databases.Delete(context.Background(), db.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
