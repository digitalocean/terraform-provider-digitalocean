package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanKubernetesNodePool_Basic(t *testing.T) {
	t.Parallel()
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
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
	t.Parallel()
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasicWithNodePool(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "tags.#", "2"),
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
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesNodePool_CreateWithAutoScale(t *testing.T) {
	t.Parallel()
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create without auto-scaling.
			{
				Config: fmt.Sprintf(`
					resource "digitalocean_kubernetes_cluster" "foobar" {
						name = "%s"
						region = "lon1"
						version = "%s"

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
				`, rName, testClusterVersion, rName),
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
				Config: fmt.Sprintf(`
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
						auto_scale = true
						min_nodes = 1
						max_nodes = 3
					}
				`, rName, testClusterVersion, rName),
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
				Config: fmt.Sprintf(`
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
						node_count = 2
						auto_scale = true
						min_nodes = 1
						max_nodes = 3
					}
				`, rName, testClusterVersion, rName),
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
				Config: fmt.Sprintf(`
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
						node_count = 2
					}
				`, rName, testClusterVersion, rName),
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
	t.Parallel()
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create without auto-scaling.
			{
				Config: fmt.Sprintf(`
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
				`, rName, testClusterVersion, rName),
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
				Config: fmt.Sprintf(`
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
						auto_scale = true
						min_nodes = 1
						max_nodes = 3
					}
				`, rName, testClusterVersion, rName),
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
				Config: fmt.Sprintf(`
					resource "digitalocean_kubernetes_cluster" "foobar" {
						name = "%s"
						region = "lon1"
						version = "%s"

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
				`, rName, testClusterVersion, rName),
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

func TestAccDigitalOceanKubernetesNodePool_WithEmptyNodePool(t *testing.T) {
	t.Parallel()
	rName := randomTestName()
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigWithEmptyNodePool(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "name", fmt.Sprintf("%s-pool", rName)),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "actual_node_count", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "nodes.#", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "min_nodes", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_node_pool.barfoo", "max_nodes", "3"),
				),
			},
		},
	})
}

func testAccDigitalOceanKubernetesConfigBasicWithNodePool(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = "%s"
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
`, rName, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigBasicWithNodePool2(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = "%s"
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
}
`, rName, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigWithEmptyNodePool(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = "%s"

	node_pool {
		name       = "default"
		size       = "s-1vcpu-2gb"
		node_count = 1
	}
}

resource digitalocean_kubernetes_node_pool "barfoo" {
  cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"

	name       = "%s-pool"
	size       = "s-1vcpu-2gb"
	auto_scale = true
	min_nodes  = 0
	max_nodes  = 3
}
`, rName, testClusterVersion, rName)
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
