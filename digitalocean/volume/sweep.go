package volume

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("digitalocean_volume", &resource.Sweeper{
		Name:         "digitalocean_volume",
		F:            testSweepVolumes,
		Dependencies: []string{"digitalocean_droplet"},
	})
}

func testSweepVolumes(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListVolumeParams{
		ListOptions: &godo.ListOptions{PerPage: 200},
	}
	volumes, _, err := client.Storage.ListVolumes(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, v := range volumes {
		if strings.HasPrefix(v.Name, "volume-") || strings.HasPrefix(v.Name, "tf-acc-test-") {

			if len(v.DropletIDs) > 0 {
				log.Printf("Detaching volume %v from Droplet %v", v.ID, v.DropletIDs[0])

				action, _, err := client.StorageActions.DetachByDropletID(context.Background(), v.ID, v.DropletIDs[0])
				if err != nil {
					return fmt.Errorf("Error resizing volume (%s): %s", v.ID, err)
				}

				if err = util.WaitForAction(client, action); err != nil {
					return fmt.Errorf(
						"Error waiting for volume (%s): %s", v.ID, err)
				}
			}

			log.Printf("Destroying Volume %s", v.Name)

			if _, err := client.Storage.DeleteVolume(context.Background(), v.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
