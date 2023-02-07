package kubernetes

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
	resource.AddTestSweepers("digitalocean_kubernetes_cluster", &resource.Sweeper{
		Name: "digitalocean_kubernetes_cluster",
		F:    sweepKubernetesClusters,
	})

}

func sweepKubernetesClusters(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	clusters, _, err := client.Kubernetes.List(context.Background(), opt)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Found %d Kubernetes clusters to sweep", len(clusters))

	for _, c := range clusters {
		if strings.HasPrefix(c.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying Kubernetes cluster %s", c.Name)
			if _, err := client.Kubernetes.Delete(context.Background(), c.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
