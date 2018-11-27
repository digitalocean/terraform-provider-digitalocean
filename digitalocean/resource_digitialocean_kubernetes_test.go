package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDigitalOceanKubernetes_Basic(t *testing.T) {
	rName := acctest.RandString(10)
	var k8s godo.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "region", "lon1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "version", "1.12.1-do.2"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "ipv4_address"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "cluster_subnet"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "service_subnet"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "endpoint"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "tags.2356372769", "foo"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "tags.1996459178", "bar"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "status"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "updated_at"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.name", "pool1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.count", "3"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.tags.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.tags.2053932785", "one"), // Currently tags are being copied from parent this will fail
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.tags.298486374", "two"),  // requires API update
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.nodes.#", "3"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.nodes.0.name"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.nodes.0.status"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.nodes.0.created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.nodes.0.updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.nodes.1.name"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.nodes.2.name"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.raw_config"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.host"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.client_certificate"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetes_MultiplePools(t *testing.T) {
	rName := acctest.RandString(10)
	var k8s godo.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigMultipleNodePools(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.nodes.#", "3"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.513870257.nodes.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.name", "pool1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.513870257.name", "pool2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetes_AddPool(t *testing.T) {
	rName := acctest.RandString(10)
	var k8s godo.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.name", "pool1"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigMultipleNodePools(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.name", "pool1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.513870257.name", "pool2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetes_RemovePool(t *testing.T) {
	rName := acctest.RandString(10)
	var k8s godo.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigMultipleNodePools(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.name", "pool1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.513870257.name", "pool2"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.2275956747.name", "pool1"),
				),
			},
		},
	})
}

func testAccDigitalOceanKubernetesConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = "1.12.1-do.2"
	tags    = ["foo","bar"]

	node_pool {
		name  = "pool1"
		size  = "s-1vcpu-2gb"
		count = 3
		tags  = ["one","two"]
	}
}
`, rName)
}

func testAccDigitalOceanKubernetesConfigMultipleNodePools(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = "1.12.1-do.2"
	tags    = ["foo","bar"]

	node_pool {
		name  = "pool1"
		size  = "s-1vcpu-2gb"
		count = 3
		tags  = ["foo","bar"]
	}
	
	node_pool {
		name  = "pool2"
		size  = "s-1vcpu-2gb"
		count = 1
		tags  = ["foo","bar"]
	}
}
`, rName)
}

func testAccCheckDigitalOceanKubernetesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*godo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_kubernetes_cluster" {
			continue
		}

		// Try to find the cluster
		_, _, err := client.Kubernetes.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("K8s Cluster still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanKubernetesExists(n string, cluster *godo.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*godo.Client)

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
