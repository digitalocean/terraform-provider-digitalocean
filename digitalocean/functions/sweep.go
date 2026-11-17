package functions

import (
	"context"
	"log"
	"strings"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("digitalocean_functions_namespace", &resource.Sweeper{
		Name: "digitalocean_functions_namespace",
		F:    sweepFunctionsNamespaces,
	})
}

// Deleting a namespace also deletes its triggers, so no separate trigger
// sweeper is needed.
func sweepFunctionsNamespaces(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	namespaces, _, err := client.Functions.ListNamespaces(context.Background())
	if err != nil {
		return err
	}

	for _, ns := range namespaces {
		if strings.HasPrefix(ns.Label, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying Functions namespace %s (%s)", ns.Label, ns.Namespace)
			if _, err := client.Functions.DeleteNamespace(context.Background(), ns.Namespace); err != nil {
				return err
			}
		}
	}

	return nil
}
