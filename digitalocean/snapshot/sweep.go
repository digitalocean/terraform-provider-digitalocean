package snapshot

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
	resource.AddTestSweepers("digitalocean_droplet_snapshot", &resource.Sweeper{
		Name:         "digitalocean_droplet_snapshot",
		F:            testSweepDropletSnapshots,
		Dependencies: []string{"digitalocean_droplet"},
	})

	resource.AddTestSweepers("digitalocean_volume_snapshot", &resource.Sweeper{
		Name:         "digitalocean_volume_snapshot",
		F:            testSweepVolumeSnapshots,
		Dependencies: []string{"digitalocean_volume"},
	})
}

func testSweepDropletSnapshots(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	snapshots, _, err := client.Snapshots.ListDroplet(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, s := range snapshots {
		if strings.HasPrefix(s.Name, "snapshot-") || strings.HasPrefix(s.Name, "snap-") {
			log.Printf("Destroying Droplet Snapshot %s", s.Name)

			if _, err := client.Snapshots.Delete(context.Background(), s.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func testSweepVolumeSnapshots(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	snapshots, _, err := client.Snapshots.ListVolume(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, s := range snapshots {
		if strings.HasPrefix(s.Name, "snapshot-") {
			log.Printf("Destroying Volume Snapshot %s", s.Name)

			if _, err := client.Snapshots.Delete(context.Background(), s.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
