package digitalocean

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanKubernetesCluster_ImportBasic(t *testing.T) {
	clusterName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersion19, clusterName),
				// Remove the default node pool tag so that the import code which infers
				// the need to add the tag gets triggered.
				Check: testAccDigitalOceanKubernetesRemoveDefaultNodePoolTag(clusterName),
			},
			{
				ResourceName:      "digitalocean_kubernetes_cluster.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"kube_config",            // because kube_config was completely different for imported state
					"node_pool.0.node_count", // because import test failed before DO had started the node in pool
					"updated_at",             // because removing default tag updates the resource outside of Terraform
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

func TestAccDigitalOceanKubernetesCluster_ImportNonDefaultNodePool(t *testing.T) {
	testName1 := randomTestName()
	testName2 := randomTestName()

	config := fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name = "%s"
  region = "lon1"
  version = data.digitalocean_kubernetes_versions.test.latest_version

  node_pool {
    name = "default"
	size = "s-1vcpu-2gb"
	node_count = 1
  }
}

resource "digitalocean_kubernetes_node_pool" "barfoo" {
  cluster_id = digitalocean_kubernetes_cluster.foobar.id
  name = "%s"
  size = "s-1vcpu-2gb"
  node_count = 1
}
`, testClusterVersion19, testName1, testName2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName: "digitalocean_kubernetes_cluster.foobar",
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 2 {
						return fmt.Errorf("expected 2 states: %#v", s)
					}

					actualNames := []string{s[0].Attributes["name"], s[1].Attributes["name"]}
					expectedNames := []string{testName1, testName2}
					sort.Strings(actualNames)
					sort.Strings(expectedNames)

					if !reflect.DeepEqual(actualNames, expectedNames) {
						return fmt.Errorf("expected name attributes for cluster and node pools to match: expected=%#v, actual=%#v, s=%#v",
							expectedNames, actualNames, s)
					}

					return nil
				},
			},
		},
	})
}
