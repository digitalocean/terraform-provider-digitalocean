package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanKubernetesNodePool_Basic(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasicWithNodePool(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesNodePool_Update(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasicWithNodePool(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "tags.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "labels.%", "0"),
					resource.TestCheckNoResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "labels.priority"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "1"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigBasicWithNodePool2(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "tags.#", "3"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "labels.%", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "labels.priority", "high"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesNodePool_CreateWithAutoScale(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create without auto-scaling.
			{
				Config: fmt.Sprintf(`%s

					resource "digitalocean_kubernetes_cluster" "foobar" {
						name = "%s"
						region = "lon1"
						version = data.digitalocean_kubernetes_versions.test.latest_version

						node_pool {
							name = "default"
							size  = "s-1vcpu-2gb"
							node_count = 1
						}
					}

					resource digitalocean_kubernetes_node_pool "barfoo" {
						cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"
						name = "%s"
						size = "s-1vcpu-2gb"
						node_count = 1
						auto_scale = true
						min_nodes = 1
						max_nodes = 5
					}
				`, testClusterVersion19, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "max_nodes", "5"),
				),
			},
			// Remove node count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

					resource digitalocean_kubernetes_node_pool "barfoo" {
						cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"
						name = "%s"
						size = "s-1vcpu-2gb"
						auto_scale = true
						min_nodes = 1
						max_nodes = 3
					}
				`, testClusterVersion19, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
			// Update node count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

					resource digitalocean_kubernetes_node_pool "barfoo" {
						cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"
						name = "%s"
						size = "s-1vcpu-2gb"
						node_count = 2
						auto_scale = true
						min_nodes = 1
						max_nodes = 3
					}
				`, testClusterVersion19, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
			// Disable auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

					resource digitalocean_kubernetes_node_pool "barfoo" {
						cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"
						name = "%s"
						size = "s-1vcpu-2gb"
						node_count = 2
					}
				`, testClusterVersion19, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "auto_scale", "false"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "min_nodes", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "max_nodes", "0"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesNodePool_UpdateWithAutoScale(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create without auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

					resource digitalocean_kubernetes_node_pool "barfoo" {
						cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"
						name = "%s"
						size = "s-1vcpu-2gb"
						node_count = 1
					}
				`, testClusterVersion19, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "auto_scale", "false"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "min_nodes", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "max_nodes", "0"),
				),
			},
			// Update to enable auto-scaling.
			{
				Config: fmt.Sprintf(`%s

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

					resource digitalocean_kubernetes_node_pool "barfoo" {
						cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"
						name = "%s"
						size = "s-1vcpu-2gb"
						node_count = 1
						auto_scale = true
						min_nodes = 1
						max_nodes = 3
					}
				`, testClusterVersion19, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
			// Remove node count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

					resource "digitalocean_kubernetes_cluster" "foobar" {
						name = "%s"
						region = "lon1"
						version = data.digitalocean_kubernetes_versions.test.latest_version

						node_pool {
							name = "default"
							size  = "s-1vcpu-2gb"
							node_count = 1
						}
					}

					resource digitalocean_kubernetes_node_pool "barfoo" {
						cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"
						name = "%s"
						size = "s-1vcpu-2gb"
						auto_scale = true
						min_nodes = 1
						max_nodes = 3
					}
				`, testClusterVersion19, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
		},
	})
}

func testAccDigitalOceanKubernetesConfigBasicWithNodePool(rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = data.digitalocean_kubernetes_versions.test.latest_version
	tags    = ["foo","bar"]

	node_pool {
		name = "default"
		size  = "s-1vcpu-2gb"
		node_count = 1
		tags  = ["one","two"]
	}
}

resource digitalocean_kubernetes_node_pool "barfoo" {
  cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"

	name    = "%s"
	size  = "s-1vcpu-2gb"
	node_count = 1
	tags  = ["three","four"]
}
`, testClusterVersion19, rName, rName)
}

func testAccDigitalOceanKubernetesConfigBasicWithNodePool2(rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = data.digitalocean_kubernetes_versions.test.latest_version
	tags    = ["foo","bar"]

	node_pool {
		name = "default"
		size  = "s-1vcpu-2gb"
		node_count = 1
		tags  = ["one","two"]
	}
}

resource digitalocean_kubernetes_node_pool "barfoo" {
	cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"

	name    = "%s-updated"
	size  = "s-1vcpu-2gb"
	node_count = 2
	tags  = ["one","two", "three"]
	labels = {
      priority = "high"
	}
}
`, testClusterVersion19, rName, rName)
}

func testAccCheckDigitalOceanKubernetesNodePoolExists(n string, cluster *godo.KubernetesCluster, pool *godo.KubernetesNodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundPool, _, err := client.Kubernetes.GetNodePool(context.Background(), cluster.ID, rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundPool.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*pool = *foundPool

		return nil
	}
}
