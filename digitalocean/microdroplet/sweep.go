package microdroplet

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
	resource.AddTestSweepers("digitalocean_microdroplet", &resource.Sweeper{
		Name: "digitalocean_microdroplet",
		F:    sweepMicroDroplets,
	})
	resource.AddTestSweepers("digitalocean_microdroplet_image", &resource.Sweeper{
		Name: "digitalocean_microdroplet_image",
		// Images may be referenced by MicroDroplets, so sweep them second.
		Dependencies: []string{"digitalocean_microdroplet"},
		F:            sweepMicroDropletImages,
	})
}

func sweepMicroDroplets(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	for {
		mds, resp, err := client.MicroDroplets.List(context.Background(), opt)
		if err != nil {
			return err
		}
		for _, m := range mds {
			if !strings.HasPrefix(m.Name, sweep.TestNamePrefix) {
				continue
			}
			log.Printf("[DEBUG] Destroying MicroDroplet %s (%s)", m.Name, m.ID)
			if _, err := client.MicroDroplets.Delete(context.Background(), m.ID); err != nil {
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

func sweepMicroDropletImages(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}
	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	for {
		images, resp, err := client.MicroDropletImages.List(context.Background(), opt)
		if err != nil {
			return err
		}
		for _, img := range images {
			if !strings.HasPrefix(img.Name, sweep.TestNamePrefix) {
				continue
			}
			log.Printf("[DEBUG] Destroying MicroDroplet image %s (%s)", img.Name, img.ID)
			if _, err := client.MicroDropletImages.Delete(context.Background(), img.ID); err != nil {
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
