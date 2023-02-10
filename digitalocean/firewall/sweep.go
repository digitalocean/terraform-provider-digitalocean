package firewall

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
	resource.AddTestSweepers("digitalocean_firewall", &resource.Sweeper{
		Name: "digitalocean_firewall",
		F:    sweepFirewall,
	})

}

func sweepFirewall(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	fws, _, err := client.Firewalls.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, f := range fws {
		if strings.HasPrefix(f.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying firewall %s", f.Name)

			if _, err := client.Firewalls.Delete(context.Background(), f.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
