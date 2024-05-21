package vpcpeering

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
	resource.AddTestSweepers("digitalocean_vpcpeering", &resource.Sweeper{
		Name: "digitalocean_vpcpeering",
		F:    sweepVPCPeering,
	})
}

func sweepVPCPeering(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()
	ctx := context.Background()

	opt := &godo.ListOptions{PerPage: 200}
	vpcPeerings, _, err := client.VPCs.ListVPCPeerings(ctx, opt)
	if err != nil {
		return err
	}

	for _, v := range vpcPeerings {
		if strings.HasPrefix(v.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying VPC Peering %s", v.Name)
			resp, err := client.VPCs.DeleteVPCPeering(ctx, v.ID)
			if err != nil {
				if resp.StatusCode == http.StatusForbidden {
					log.Printf("[DEBUG] Skipping VPC Peering %s; still contains resources", v.Name)
				} else {
					return err
				}
			}
			log.Printf("[DEBUG] Waiting for VPC Peering (%s) to be deleted", v.Name)
			stateConf := &retry.StateChangeConf{
				Pending:    []string{"DELETING"},
				Target:     []string{http.StatusText(http.StatusNotFound)},
				Refresh:    vpcPeeringStateRefreshFunc(client, v.ID),
				Timeout:    10 * time.Minute,
				MinTimeout: 2 * time.Second,
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("error waiting for VPC Peering (%s) to be deleted: %s", v.Name, err)
			}
		}
	}

	return nil
}
