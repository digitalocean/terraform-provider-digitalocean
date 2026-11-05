package dedicatedinference

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
	resource.AddTestSweepers("digitalocean_dedicated_inference", &resource.Sweeper{
		Name: "digitalocean_dedicated_inference",
		F:    sweepDedicatedInference,
	})
}

func sweepDedicatedInference(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.DedicatedInferenceListOptions{
		ListOptions: godo.ListOptions{Page: 1, PerPage: 200},
	}

	var allItems []godo.DedicatedInferenceListItem
	for {
		items, resp, err := client.DedicatedInference.List(context.Background(), opts)
		if err != nil {
			return err
		}
		allItems = append(allItems, items...)
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}
		opts.Page = page + 1
	}

	log.Printf("[DEBUG] Found %d dedicated inference endpoints to sweep", len(allItems))

	var swept int
	for _, di := range allItems {
		if strings.HasPrefix(di.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying Dedicated Inference endpoint %s (%s)", di.Name, di.ID)

			if _, err := client.DedicatedInference.Delete(context.Background(), di.ID); err != nil {
				return err
			}
			swept++
		}
	}
	log.Printf("[DEBUG] Deleted %d of %d dedicated inference endpoints", swept, len(allItems))

	return nil
}
