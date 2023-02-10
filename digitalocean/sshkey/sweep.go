package sshkey

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
	resource.AddTestSweepers("digitalocean_ssh_key", &resource.Sweeper{
		Name: "digitalocean_ssh_key",
		F:    sweepSSHKey,
	})

}

func sweepSSHKey(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	keys, _, err := client.Keys.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, k := range keys {
		if strings.HasPrefix(k.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying SSH key %s", k.Name)

			if _, err := client.Keys.DeleteByID(context.Background(), k.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
