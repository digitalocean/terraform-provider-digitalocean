package loadbalancer

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
	resource.AddTestSweepers("digitalocean_loadbalancer", &resource.Sweeper{
		Name: "digitalocean_loadbalancer",
		F:    sweepLoadbalancer,
	})

}

func sweepLoadbalancer(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	lbs, _, err := client.LoadBalancers.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, l := range lbs {
		if strings.HasPrefix(l.Name, "loadbalancer-") {
			log.Printf("Destroying loadbalancer %s", l.Name)

			if _, err := client.LoadBalancers.Delete(context.Background(), l.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
