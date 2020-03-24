package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceDigitalOceanKubernetesCluster_Basic(t *testing.T) {
	t.Parallel()
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigWithDataSource(rName, testClusterVersion16),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanKubernetesClusterExists("data.digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "region", "lon1"),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "version", testClusterVersion16),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.priority", "high"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "vpc_uuid"),
				),
			},
		},
	})
}

func testAccDigitalOceanKubernetesConfigWithDataSource(rName string, version string) string {
	return fmt.Sprintf(`
resource "digitalocean_kubernetes_cluster" "foo" {
	name    = "%s"
	region  = "lon1"
	version = "%s"
	tags    = ["foo","bar"]

	node_pool {
	    name = "default"
		size  = "s-1vcpu-2gb"
		node_count = 1
		tags  = ["one","two"]
        labels = {
          priority = "high"
        }
	}
}

data "digitalocean_kubernetes_cluster" "foobar" {
	name = "${digitalocean_kubernetes_cluster.foo.name}"
}
`, rName, version)
}

func testAccCheckDataSourceDigitalOceanKubernetesClusterExists(n string, cluster *godo.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundCluster, _, err := client.Kubernetes.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundCluster.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*cluster = *foundCluster

		return nil
	}
}
