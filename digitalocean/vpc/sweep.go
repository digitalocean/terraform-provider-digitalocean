package vpc

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
	resource.AddTestSweepers("digitalocean_vpc", &resource.Sweeper{
		Name: "digitalocean_vpc",
		F:    sweepVPC,
		Dependencies: []string{
			"digitalocean_droplet",
		},
	})
}

func sweepVPC(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	vpcs, _, err := client.VPCs.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, v := range vpcs {
		if strings.HasPrefix(v.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying VPC %s", v.Name)
			resp, err := client.VPCs.Delete(context.Background(), v.ID)
			if err != nil {
				if resp.StatusCode == http.StatusForbidden {
					log.Printf("[DEBUG] Skipping VPC %s; still contains resources", v.Name)
				} else {
					return err
				}
			}
		}
	}

	return nil
}
