package digitalocean

import (
	"fmt"
	"reflect"
	"sort"
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
