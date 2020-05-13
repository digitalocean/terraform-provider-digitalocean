package digitalocean

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDigitalOceanKubernetesNodePool_Import(t *testing.T) {
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
`, testClusterVersion16, testName1, testName2)
	resourceName := "digitalocean_kubernetes_node_pool.barfoo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "this-is-not-a-valid-ID",
				ExpectError:       regexp.MustCompile("Did not find the cluster owning the node pool"),
			},
		},
	})
}
