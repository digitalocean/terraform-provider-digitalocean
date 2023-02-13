package monitoring

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
	resource.AddTestSweepers("digitalocean_monitor_alert", &resource.Sweeper{
		Name: "digitalocean_monitor_alert",
		F:    sweepMonitoringAlerts,
	})

}

func sweepMonitoringAlerts(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	alerts, _, err := client.Monitoring.ListAlertPolicies(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, a := range alerts {
		if strings.HasPrefix(a.Description, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying alert %s", a.Description)

			if _, err := client.Monitoring.DeleteAlertPolicy(context.Background(), a.UUID); err != nil {
				return err
			}
		}
	}

	return nil
}
