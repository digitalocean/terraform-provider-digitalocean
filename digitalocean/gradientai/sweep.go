package gradientai

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
	resource.AddTestSweepers("digitalocean_gradientai_custom_model", &resource.Sweeper{
		Name: "digitalocean_gradientai_custom_model",
		F:    sweepCustomModel,
	})
}

func sweepCustomModel(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.CustomModelListOptions{
		ListOptions: godo.ListOptions{Page: 1, PerPage: 200},
	}

	var all []*godo.CustomModel
	for {
		listResp, resp, err := client.GradientAI.ListCustomModels(context.Background(), opts)
		if err != nil {
			return err
		}
		if listResp != nil {
			all = append(all, listResp.Models...)
		}
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}
		opts.Page = page + 1
	}

	log.Printf("[DEBUG] Found %d custom models to consider sweeping", len(all))

	var swept int
	for _, m := range all {
		if !strings.HasPrefix(m.Name, sweep.TestNamePrefix) {
			continue
		}
		// STATUS_FAILED models cannot currently be deleted via the API
		// (DELETE returns 404); skip them rather than aborting the sweep.
		if m.Status == godo.CustomModelStatusFailed {
			log.Printf("Skipping custom model %s (%s) in %s", m.Name, m.Uuid, m.Status)
			continue
		}
		log.Printf("Destroying custom model %s (%s)", m.Name, m.Uuid)
		if _, _, err := client.GradientAI.DeleteCustomModel(context.Background(), m.Uuid); err != nil {
			log.Printf("Error destroying custom model %s (%s): %s", m.Name, m.Uuid, err)
			continue
		}
		swept++
	}
	log.Printf("[DEBUG] Deleted %d of %d custom models", swept, len(all))

	return nil
}
