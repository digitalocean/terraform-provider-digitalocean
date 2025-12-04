package nfs

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
	resource.AddTestSweepers("digitalocean_nfs", &resource.Sweeper{
		Name: "digitalocean_nfs",
		F:    sweepNfs,
	})

}

func sweepNfs(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}

	shares, _, err := client.Nfs.List(context.Background(), opt, "atl1")
	if err != nil {
		return err
	}

	for _, s := range shares {
		if strings.HasPrefix(s.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying share %s", s.Name)

			if _, err := client.Nfs.Delete(context.Background(), s.ID, "atl1"); err != nil {
				return err
			}
		}
	}

	return nil
}
