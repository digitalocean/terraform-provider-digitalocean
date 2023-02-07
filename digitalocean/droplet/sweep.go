package droplet

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
	resource.AddTestSweepers("digitalocean_droplet", &resource.Sweeper{
		Name: "digitalocean_droplet",
		F:    sweepDroplets,
	})
}

func sweepDroplets(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	droplets, _, err := client.Droplets.List(context.Background(), opt)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Found %d droplets to sweep", len(droplets))

	for _, d := range droplets {
		if strings.HasPrefix(d.Name, "foo-") || strings.HasPrefix(d.Name, "bar-") || strings.HasPrefix(d.Name, "baz-") || strings.HasPrefix(d.Name, "tf-acc-test-") || strings.HasPrefix(d.Name, "foobar-") {
			log.Printf("Destroying Droplet %s", d.Name)

			if _, err := client.Droplets.Delete(context.Background(), d.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
