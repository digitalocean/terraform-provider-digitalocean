package vpcpeering

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("digitalocean_vpcpeering", &resource.Sweeper{
		Name: "digitalocean_vpcpeering",
		F:    sweepVPCPeering,
		Dependencies: []string{
			"digitalocean_droplet",
			"digitalocean_vpc",
		},
	})
}

func sweepVPCPeering(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	vpcPeerings, _, err := client.VPCs.ListVPCPeerings(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, v := range vpcPeerings {
		if strings.HasPrefix(v.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying VPC Peering %s", v.Name)
			resp, err := client.VPCs.DeleteVPCPeering(context.Background(), v.ID)
			if err != nil {
				if resp.StatusCode == http.StatusForbidden {
					log.Printf("[DEBUG] Skipping VPC Peering %s; still contains resources", v.Name)
				} else {
					return err
				}
			}
		}
	}

	return nil
}
