package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDigitalOceanKubernetesClusterNodePool_Import(t *testing.T) {
	rName := randomTestName()
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
`, rName, testClusterVersion16, rName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "digitalocean_kubernetes_cluster_node_pool.barfoo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"node_count", // because import test failed before DO had started the node in pool
				},
			},
		},
	})
}
