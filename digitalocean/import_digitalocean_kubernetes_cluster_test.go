package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanKubernetesCluster_ImportBasic(t *testing.T) {
	clusterName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(clusterName, testClusterVersion16),
			},
			{
				// Remove the default node pool tag so that the import code which infers
				// the need to add the tag gets triggered.
				Check:       testAccDigitalOceanKubernetesRemoveDefaultNodePoolTag(clusterName),
				ExpectError: regexp.MustCompile("No default node pool was found"),
			},
			{
				ResourceName:      "digitalocean_kubernetes_cluster.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"kube_config",            // because kube_config was completely different for imported state
					"node_pool.0.node_count", // because import test failed before DO had started the node in pool
				},
			},
		},
	})
}

func testAccDigitalOceanKubernetesRemoveDefaultNodePoolTag(clusterName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		clusters, resp, err := client.Kubernetes.List(context.Background(), &godo.ListOptions{})
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return fmt.Errorf("No clusters found")
			}

			return fmt.Errorf("Error listing Kubernetes clusters: %s", err)
		}

		var cluster *godo.KubernetesCluster
		for _, c := range clusters {
			if c.Name == clusterName {
				cluster = c
				break
			}
		}
		if cluster == nil {
			return fmt.Errorf("Unable to find Kubernetes cluster with name: %s", clusterName)
		}

		for _, nodePool := range cluster.NodePools {
			tags := make([]string, 0)
			for _, tag := range nodePool.Tags {
				if tag != digitaloceanKubernetesDefaultNodePoolTag {
					tags = append(tags, tag)
				}
			}

			if len(tags) != len(nodePool.Tags) {
				nodePoolUpdateRequest := &godo.KubernetesNodePoolUpdateRequest{
					Tags: tags,
				}

				_, _, err := client.Kubernetes.UpdateNodePool(context.Background(), cluster.ID, nodePool.ID, nodePoolUpdateRequest)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}
