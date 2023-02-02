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

func TestResourceInstanceState(name string, check func(*terraform.InstanceState) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := s.RootModule()
		if rs, ok := m.Resources[name]; ok {
			is := rs.Primary
			if is == nil {
				return fmt.Errorf("No primary instance: %s", name)
			}

			return check(is)
		} else {
			return fmt.Errorf("Not found: %s", name)
		}

	}
}

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

func TestAccCheckDigitalOceanDropletConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-22-04-x64"
  region    = "nyc3"
  user_data = "foobar"
}`, rInt)
}

func TakeSnapshotsOfDroplet(rInt int, droplet *godo.Droplet, snapshotsId *[]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		for i := 0; i < 3; i++ {
			err := takeSnapshotOfDroplet(rInt, i%2, droplet)
			if err != nil {
				return err
			}
		}
		retrieveDroplet, _, err := client.Droplets.Get(context.Background(), (*droplet).ID)
		if err != nil {
			return err
		}
		*snapshotsId = retrieveDroplet.SnapshotIDs
		return nil
	}
}

func takeSnapshotOfDroplet(rInt, sInt int, droplet *godo.Droplet) error {
	client := TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
	action, _, err := client.DropletActions.Snapshot(context.Background(), (*droplet).ID, fmt.Sprintf("snap-%d-%d", rInt, sInt))
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
