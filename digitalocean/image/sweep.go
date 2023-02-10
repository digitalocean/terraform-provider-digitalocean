package image

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
	resource.AddTestSweepers("digitalocean_custom_image", &resource.Sweeper{
		Name: "digitalocean_custom_image",
		F:    sweepCustomImage,
	})

}

func sweepCustomImage(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	images, _, err := client.Images.ListUser(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, i := range images {
		if strings.HasPrefix(i.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying image %s", i.Name)

			if _, err := client.Images.Delete(context.Background(), i.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
