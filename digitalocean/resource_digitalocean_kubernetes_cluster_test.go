package digitalocean

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/setutil"

	"context"
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testClusterVersion18 = `data "digitalocean_kubernetes_versions" "test" {
  version_prefix = "1.18."
}`
	testClusterVersion19 = `data "digitalocean_kubernetes_versions" "test" {
  version_prefix = "1.19."
}`
)

func init() {
	resource.AddTestSweepers("digitalocean_kubernetes_cluster", &resource.Sweeper{
		Name: "digitalocean_kubernetes_cluster",
		F:    testSweepKubernetesClusters,
	})

}

func testSweepKubernetesClusters(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*CombinedConfig).godoClient()

	opt := &godo.ListOptions{PerPage: 200}
	clusters, _, err := client.Kubernetes.List(context.Background(), opt)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Found %d Kubernetes clusters to sweep", len(clusters))

	for _, c := range clusters {
		if strings.HasPrefix(c.Name, testNamePrefix) {
			log.Printf("Destroying Kubernetes cluster %s", c.Name)
			if _, err := client.Kubernetes.Delete(context.Background(), c.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanKubernetesCluster_Basic(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "region", "lon1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "surge_upgrade", "true"),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "ipv4_address"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "cluster_subnet"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "service_subnet"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "endpoint"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "tags.#", "3"),
					setutil.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "tags.*", "foo"),
					setutil.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "tags.*", "foo"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "status"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "updated_at"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.tags.#", "2"),
					setutil.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.tags.*", "one"),
					setutil.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.tags.*", "two"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.%", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.priority", "high"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.nodes.#", "1"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.0.nodes.0.created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "node_pool.0.nodes.0.updated_at"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.taint.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.taint.0.effect", "PreferNoSchedule"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.raw_config"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.host"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.token"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.expires_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "auto_upgrade"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_UpdateCluster(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigBasic4(testClusterVersion19, rName+"-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "tags.#", "2"),
					setutil.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "tags.*", "one"),
					setutil.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "tags.*", "two"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.%", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "surge_upgrade", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_UpdatePoolDetails(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.name", "default"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigBasic2(testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.name", "default-rename"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.tags.#", "3"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.%", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.priority", "high"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.purpose", "awesome"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_UpdatePoolSize(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigBasic3(testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-2vcpu-4gb"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_CreatePoolWithAutoScale(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create with auto-scaling and explicit node_count.
			{
				Config: fmt.Sprintf(`%s

					resource "digitalocean_kubernetes_cluster" "foobar" {
						name    	 = "%s"
						region  	 = "lon1"
						version 	 = data.digitalocean_kubernetes_versions.test.latest_version
						auto_upgrade = true

						node_pool {
							name = "default"
							size  = "s-1vcpu-2gb"
							node_count = 1
							auto_scale = true
							min_nodes = 1
							max_nodes = 3
						}
					}
				`, testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "auto_upgrade", "true"),
				),
			},
			// Remove node_count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

					resource "digitalocean_kubernetes_cluster" "foobar" {
						name    = "%s"
						region  = "lon1"
						version = data.digitalocean_kubernetes_versions.test.latest_version

						node_pool {
							name = "default"
							size  = "s-1vcpu-2gb"
							auto_scale = true
							min_nodes = 1
							max_nodes = 3
						}
					}
				`, testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
				),
			},
			// Update node_count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

					resource "digitalocean_kubernetes_cluster" "foobar" {
						name    = "%s"
						region  = "lon1"
						version = data.digitalocean_kubernetes_versions.test.latest_version

						node_pool {
							name = "default"
							size  = "s-1vcpu-2gb"
							node_count = 2
							auto_scale = true
							min_nodes = 1
							max_nodes = 3
						}
					}
				`, testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
				),
			},
			// Disable auto-scaling.
			{
				Config: fmt.Sprintf(`%s

					resource "digitalocean_kubernetes_cluster" "foobar" {
						name    = "%s"
						region  = "lon1"
						version = data.digitalocean_kubernetes_versions.test.latest_version

						node_pool {
							name = "default"
							size  = "s-1vcpu-2gb"
							node_count = 2
						}
					}
				`, testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "false"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "0"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_UpdatePoolWithAutoScale(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create with auto-scaling disabled.
			{
				Config: fmt.Sprintf(`%s

				resource "digitalocean_kubernetes_cluster" "foobar" {
					name    = "%s"
					region  = "lon1"
					version = data.digitalocean_kubernetes_versions.test.latest_version

					node_pool {
						name = "default"
						size  = "s-1vcpu-2gb"
						node_count = 1
					}
				}
			`, testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "false"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "0"),
				),
			},
			// Enable auto-scaling with explicit node_count.
			{
				Config: fmt.Sprintf(`%s

					resource "digitalocean_kubernetes_cluster" "foobar" {
						name    = "%s"
						region  = "lon1"
						version = data.digitalocean_kubernetes_versions.test.latest_version

						node_pool {
							name = "default"
							size  = "s-1vcpu-2gb"
							node_count = 1
							auto_scale = true
							min_nodes = 1
							max_nodes = 3
						}
					}
				`, testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
				),
			},
			// Remove node_count, keep auto-scaling.
			{
				Config: fmt.Sprintf(`%s

					resource "digitalocean_kubernetes_cluster" "foobar" {
						name    = "%s"
						region  = "lon1"
						version = data.digitalocean_kubernetes_versions.test.latest_version

						node_pool {
							name = "default"
							size  = "s-1vcpu-2gb"
							auto_scale = true
							min_nodes = 1
							max_nodes = 3
						}
					}
				`, testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.auto_scale", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.min_nodes", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.max_nodes", "3"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_KubernetesProviderInteroperability(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"kubernetes": {
				Source:            "hashicorp/kubernetes",
				VersionConstraint: "2.0.1",
			},
		},
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfig_KubernetesProviderInteroperability(testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s), resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.raw_config"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.host"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "kube_config.0.token"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_UpgradeVersion(t *testing.T) {
	rName := randomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersion18, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersion19, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPtr("digitalocean_kubernetes_cluster.foobar", "id", &k8s.ID),
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
				),
			},
		},
	})
}

func testAccDigitalOceanKubernetesConfigBasic(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = data.digitalocean_kubernetes_versions.test.latest_version
	surge_upgrade = true
	tags    = ["foo","bar", "one"]

	node_pool {
	  name = "default"
      size  = "s-1vcpu-2gb"
      node_count = 1
      tags  = ["one","two"]
      labels = {
        priority = "high"
      }
      taint {
        key = "key1"
        value = "val1"
        effect = "PreferNoSchedule"
      }
	}
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigBasic2(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = data.digitalocean_kubernetes_versions.test.latest_version
	surge_upgrade = true
	tags    = ["foo","bar"]

	node_pool {
	  name = "default-rename"
		size  = "s-1vcpu-2gb"
		node_count = 2
		tags  = ["one","two","three"]
        labels = {
            priority = "high"
            purpose = "awesome"
        }
	}
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigBasic3(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = data.digitalocean_kubernetes_versions.test.latest_version
	tags    = ["foo","bar"]

	node_pool {
	  name = "default"
		size  = "s-2vcpu-4gb"
		node_count = 1
		tags  = ["one","two"]
	}
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigBasic4(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
    surge_upgrade = true
	version = data.digitalocean_kubernetes_versions.test.latest_version
	tags    = ["one","two"]

	node_pool {
	  name = "default"
		size  = "s-2vcpu-4gb"
		node_count = 1
		tags  = ["foo","bar"]
	}
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigBasic5(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = data.digitalocean_kubernetes_versions.test.latest_version
	tags    = ["one","two"]
	version = true

	node_pool {
	  name = "default"
		size  = "s-2vcpu-4gb"
		node_count = 1
		tags  = ["foo","bar"]
	}
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfig_KubernetesProviderInteroperability(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = data.digitalocean_kubernetes_versions.test.latest_version

	node_pool {
	  name = "default"
		size  = "s-2vcpu-4gb"
		node_count = 1
	}
}

provider "kubernetes" {
  host = digitalocean_kubernetes_cluster.foobar.endpoint
  cluster_ca_certificate = base64decode(
    digitalocean_kubernetes_cluster.foobar.kube_config[0].cluster_ca_certificate
  )
  token = digitalocean_kubernetes_cluster.foobar.kube_config[0].token
}

resource "kubernetes_service_account" "tiller" {
  metadata {
    name      = "tiller"
    namespace = "kube-system"
  }

  automount_service_account_token = true
}
`, testClusterVersion, rName)
}

func testAccCheckDigitalOceanKubernetesClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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

func testAccCheckDigitalOceanKubernetesClusterExists(n string, cluster *godo.KubernetesCluster) resource.TestCheckFunc {
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

func Test_filterTags(t *testing.T) {
	tests := []struct {
		have []string
		want []string
	}{
		{
			have: []string{"k8s", "foo"},
			want: []string{"foo"},
		},
		{
			have: []string{"k8s", "k8s:looks-like-a-uuid", "bar"},
			want: []string{"bar"},
		},
		{
			have: []string{"k8s", "k8s:looks-like-a-uuid", "bar", "k8s-this-is-ok"},
			want: []string{"bar", "k8s-this-is-ok"},
		},
		{
			have: []string{"k8s", "k8s:looks-like-a-uuid", "terraform:default-node-pool", "baz"},
			want: []string{"baz"},
		},
	}

	for _, tt := range tests {
		filteredTags := filterTags(tt.have)
		if !reflect.DeepEqual(filteredTags, tt.want) {
			t.Errorf("filterTags returned %+v, expected %+v", filteredTags, tt.want)
		}
	}
}

func Test_renderKubeconfig(t *testing.T) {
	certAuth := []byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURKekNDQWWlOQT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K")
	expected := fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: %v
    server: https://6a37a0f6-c355-4527-b54d-521beffd9817.k8s.ondigitalocean.com
  name: do-lon1-test-cluster
contexts:
- context:
    cluster: do-lon1-test-cluster
    user: do-lon1-test-cluster-admin
  name: do-lon1-test-cluster
current-context: do-lon1-test-cluster
users:
- name: do-lon1-test-cluster-admin
  user:
    token: 97ae2bbcfd85c34155a56b822ffa73909d6770b28eb7e5dfa78fa83e02ffc60f
`, base64.StdEncoding.EncodeToString(certAuth))

	creds := godo.KubernetesClusterCredentials{
		Server:                   "https://6a37a0f6-c355-4527-b54d-521beffd9817.k8s.ondigitalocean.com",
		CertificateAuthorityData: certAuth,
		Token:                    "97ae2bbcfd85c34155a56b822ffa73909d6770b28eb7e5dfa78fa83e02ffc60f",
		ExpiresAt:                time.Now(),
	}
	kubeConfigRenderd, err := renderKubeconfig("test-cluster", "lon1", &creds)
	if err != nil {
		t.Errorf("error calling renderKubeconfig: %s", err)

	}
	got := string(kubeConfigRenderd)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("renderKubeconfig returned %+v\n, expected %+v\n", got, expected)
	}
}
