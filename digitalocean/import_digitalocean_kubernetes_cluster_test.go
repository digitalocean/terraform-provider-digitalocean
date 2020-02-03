package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
