package app

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
	resource.AddTestSweepers("digitalocean_app", &resource.Sweeper{
		Name: "digitalocean_app",
		F:    sweepApp,
	})

}

func sweepApp(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	apps, _, err := client.Apps.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, app := range apps {
		if strings.HasPrefix(app.Spec.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying app %s", app.Spec.Name)

			if _, err := client.Apps.Delete(context.Background(), app.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
