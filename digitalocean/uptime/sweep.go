package uptime

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
	resource.AddTestSweepers("digitalocean_uptime_check", &resource.Sweeper{
		Name: "digitalocean_uptime_check",
		F:    sweepUptimeCheck,
	})

	// Note: Deleting the check will delete associated alerts. So no sweeper is
	// needed for digitalocean_uptime_alert
}

func sweepUptimeCheck(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	checks, _, err := client.UptimeChecks.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, c := range checks {
		if strings.HasPrefix(c.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Deleting uptime check %s", c.Name)

			if _, err := client.UptimeChecks.Delete(context.Background(), c.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
