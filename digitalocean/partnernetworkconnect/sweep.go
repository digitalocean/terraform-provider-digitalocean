package partnernetworkconnect

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("digitalocean_partner_network_connect", &resource.Sweeper{
		Name: "digitalocean_partner_network_connect",
		F:    sweepPartnerNetworkConnect,
	})
}

func sweepPartnerNetworkConnect(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()
	ctx := context.Background()

	opt := &godo.ListOptions{PerPage: 200}
	partnerNetworkConnects, _, err := client.PartnerNetworkConnect.List(ctx, opt)
	if err != nil {
		return err
	}

	for _, p := range partnerNetworkConnects {
		if strings.HasPrefix(p.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying Partner Network Connect %s", p.Name)
			resp, err := client.PartnerNetworkConnect.Delete(ctx, p.ID)
			if err != nil {
				if resp.StatusCode == http.StatusForbidden {
					log.Printf("[DEBUG] Skipping Partner Network Connect %s; still contains resources", p.Name)
				} else {
					return err
				}
			}
			log.Printf("[DEBUG] Waiting for Partner Network Connect (%s) to be deleted", p.Name)
			stateConf := &retry.StateChangeConf{
				Pending:    []string{"DELETING"},
				Target:     []string{http.StatusText(http.StatusNotFound)},
				Refresh:    partnerNetworkConnectStateRefreshFunc(client, p.ID),
				Timeout:    10 * time.Minute,
				MinTimeout: 2 * time.Second,
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("error waiting for Partner Network Connect (%s) to be deleted: %s", p.Name, err)
			}
		}
	}

	return nil
}
