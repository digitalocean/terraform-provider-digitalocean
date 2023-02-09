package acceptance

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCheckDigitalOceanDropletDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_droplet" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Droplet
		_, _, err = client.Droplets.Get(context.Background(), id)

		// Wait

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for droplet (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func TestAccCheckDigitalOceanDropletExists(n string, droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Droplet ID is set")
		}

		client := TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Droplet
		retrieveDroplet, _, err := client.Droplets.Get(context.Background(), id)

		if err != nil {
			return err
		}

		if strconv.Itoa(retrieveDroplet.ID) != rs.Primary.ID {
			return fmt.Errorf("Droplet not found")
		}

		*droplet = *retrieveDroplet

		return nil
	}
}

func TestAccCheckDigitalOceanDropletConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-22-04-x64"
  region    = "nyc3"
  user_data = "foobar"
}`, name)
}

// TakeSnapshotsOfDroplet takes three snapshots of the given Droplet. One will have the suffix -1 and two will have -0.
func TakeSnapshotsOfDroplet(snapName string, droplet *godo.Droplet, snapshotsIDs *[]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		for i := 0; i < 3; i++ {
			err := takeSnapshotOfDroplet(snapName, i%2, droplet)
			if err != nil {
				return err
			}
		}
		retrieveDroplet, _, err := client.Droplets.Get(context.Background(), (*droplet).ID)
		if err != nil {
			return err
		}
		*snapshotsIDs = retrieveDroplet.SnapshotIDs
		return nil
	}
}

func takeSnapshotOfDroplet(snapName string, intSuffix int, droplet *godo.Droplet) error {
	client := TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
	action, _, err := client.DropletActions.Snapshot(context.Background(), (*droplet).ID, fmt.Sprintf("%s-%d", snapName, intSuffix))
	if err != nil {
		return err
	}
	util.WaitForAction(client, action)
	return nil
}

func DeleteDropletSnapshots(snapshotsId *[]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("Deleting Droplet snapshots")

		client := TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		snapshots := *snapshotsId
		for _, value := range snapshots {
			log.Printf("Deleting %d", value)
			_, err := client.Images.Delete(context.Background(), value)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
