package dropletautoscale

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
	resource.AddTestSweepers("digitalocean_droplet_autoscale", &resource.Sweeper{
		Name: "digitalocean_droplet_autoscale",
		F:    sweepDropletAutoscale,
	})
}

func sweepDropletAutoscale(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GodoClient()
	pools, _, err := client.DropletAutoscale.List(context.Background(), &godo.ListOptions{PerPage: 200})
	if err != nil {
		return err
	}
	for _, pool := range pools {
		if strings.HasPrefix(pool.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying droplet autoscale pool %s", pool.Name)
			if _, err = client.DropletAutoscale.DeleteDangerous(context.Background(), pool.ID); err != nil {
				return err
			}
		}
	}
	return nil
}
