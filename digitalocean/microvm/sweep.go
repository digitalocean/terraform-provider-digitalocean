package microvm

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
	resource.AddTestSweepers("digitalocean_microvm", &resource.Sweeper{
		Name: "digitalocean_microvm",
		F:    sweepMicroVMs,
	})
}

func sweepMicroVMs(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	for {
		mds, resp, err := client.MicroVMs.List(context.Background(), opt)
		if err != nil {
			return err
		}
		for _, m := range mds {
			if !strings.HasPrefix(m.Name, sweep.TestNamePrefix) {
				continue
			}
			log.Printf("[DEBUG] Destroying MicroVM %s (%s)", m.Name, m.ID)
			if _, err := client.MicroVMs.Delete(context.Background(), m.ID); err != nil {
				return err
			}
		}
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}
		opt.Page = page + 1
	}
	return nil
}
