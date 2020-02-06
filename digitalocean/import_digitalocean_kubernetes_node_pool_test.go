package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanKubernetesNodePool_Import(t *testing.T) {
	testName1 := randomTestName()
	testName2 := randomTestName()

	config := fmt.Sprintf(`
resource "digitalocean_kubernetes_cluster" "foobar" {
  name = "%s"
  region = "lon1"
  version = "%s"

  node_pool {
    name = "default"
	size = "s-1vcpu-2gb"
	node_count = 1
  }
}

resource digitalocean_kubernetes_node_pool "barfoo" {
  cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"
  name = "%s"
  size = "s-1vcpu-2gb"
  node_count = 1
}
`, testName1, testClusterVersion16, testName2)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
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
					clusterState, nodePoolState := s[0], s[1]

					if clusterState.Attributes["name"] != testName1 {
						return fmt.Errorf("expected name attribute for cluster to match: expected=%s, actual=%s",
							clusterState.Attributes["name"], testName1)
					}

					if nodePoolState.Attributes["name"] != testName2 {
						return fmt.Errorf("expected name attribute for node pool to match: expected=%s, actual=%s",
							nodePoolState.Attributes["name"], testName2)
					}

					return nil
				},
			},
		},
	})
}
