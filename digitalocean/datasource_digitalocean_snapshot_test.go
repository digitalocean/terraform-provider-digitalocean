package digitalocean

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDigitalOceanSnapshotDataSource_droplet(t *testing.T) {
	var (
		droplet             godo.Droplet
		dropletSnapshotsIds []int
	)

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					takeSnapshotsOfDroplet(rInt, &droplet, true, &dropletSnapshotsIds),
				),
			},
			{
				Config: testAccCheckDigitalOceanDropletSnapshotDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_snapshot.droplet_test_snap", "name", fmt.Sprintf("snap-%d-2", rInt)),
					resource.TestCheckResourceAttr("data.digitalocean_snapshot.droplet_test_snap", "min_disk_size", "20"),
					resource.TestCheckResourceAttr("data.digitalocean_snapshot.droplet_test_snap", "regions.0", "nyc3"),
				),
			},
			{
				Config:      testAccCheckDigitalOceanDropletSnapshotDataSourceConfig_fail,
				ExpectError: regexp.MustCompile(`.*Your query returned more than one result.*`),
			},
			{
				Config: " ",
				Check: resource.ComposeTestCheckFunc(
					deleteDropletSnapshots(&dropletSnapshotsIds),
				),
			},
		},
	})
}

func TestAccDigitalOceanSnapshotDataSource_volume(t *testing.T) {
	var volumeSnapshotsIds []string

	rInt := acctest.RandInt()
	name := fmt.Sprintf("volume-%v", rInt)
	volume := godo.Volume{
		Name: name,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeConfig_basic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					takeSnapshotsOfVolume(rInt, &volume, &volumeSnapshotsIds),
				),
			},
			{
				Config: testAccCheckDigitalOceanVolumeSnapshotDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_snapshot.volume_test_snap", "name", fmt.Sprintf("vol-snap-%d-2", rInt)),
					resource.TestCheckResourceAttr("data.digitalocean_snapshot.volume_test_snap", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.digitalocean_snapshot.volume_test_snap", "regions.0", "nyc1"),
				),
			},
			{
				Config:      testAccCheckDigitalOceanVolumeSnapshotDataSourceConfig_fail,
				ExpectError: regexp.MustCompile(`.*Your query returned more than one result.*`),
			},
			{
				Config: " ",
				Check: resource.ComposeTestCheckFunc(
					deleteVolumeSnapshots(&volumeSnapshotsIds),
				),
			},
		},
	})
}

func takeSnapshotsOfVolume(rInt int, volume *godo.Volume, snapshotsId *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*godo.Client)
		for i := 0; i < 3; i++ {
			createRequest := &godo.SnapshotCreateRequest{
				VolumeID: (*volume).ID,
				Name:     fmt.Sprintf("vol-snap-%d-%d", rInt, i),
			}
			volume, _, err := client.Storage.CreateSnapshot(context.Background(), createRequest)
			if err != nil {
				return err
			}
			*snapshotsId = append(*snapshotsId, volume.ID)
		}
		return nil
	}
}

func deleteVolumeSnapshots(snapshotsId *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("XXX Deleting volume snapshots")
		client := testAccProvider.Meta().(*godo.Client)
		snapshots := *snapshotsId
		for _, value := range snapshots {
			log.Printf("XXX Deleting %v", value)
			_, err := client.Snapshots.Delete(context.Background(), value)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

const testAccCheckDigitalOceanDropletSnapshotDataSourceConfig = `
data "digitalocean_snapshot" "droplet_test_snap" {
  most_recent = true
  resource_type = "droplet"
  name_regex = "^snap"
}
`

const testAccCheckDigitalOceanDropletSnapshotDataSourceConfig_fail = `
data "digitalocean_snapshot" "droplet_test_snap" {
  most_recent = false
  resource_type = "droplet"
  name_regex = "^snap"
}
`

const testAccCheckDigitalOceanVolumeSnapshotDataSourceConfig = `
data "digitalocean_snapshot" "volume_test_snap" {
    most_recent = true
    resource_type = "volume"
    name_regex = "^vol-snap"
}
`
const testAccCheckDigitalOceanVolumeSnapshotDataSourceConfig_fail = `
data "digitalocean_snapshot" "volume_test_snap" {
    most_recent = false
    resource_type = "volume"
    name_regex = "^vol-snap"
}
`
