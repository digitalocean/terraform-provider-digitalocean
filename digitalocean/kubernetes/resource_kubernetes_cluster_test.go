package kubernetes_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/kubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testClusterVersionPrevious = `data "digitalocean_kubernetes_versions" "latest" {
}

locals {
  previous_version = format("%s.",
    join(".", [
      split(".", data.digitalocean_kubernetes_versions.latest.latest_version)[0],
      tostring(parseint(split(".", data.digitalocean_kubernetes_versions.latest.latest_version)[1], 10) - 1)
    ])
  )
}

data "digitalocean_kubernetes_versions" "test" {
  version_prefix = local.previous_version
}`

	testClusterVersionLatest = `data "digitalocean_kubernetes_versions" "test" {
}`
)

func TestAccDigitalOceanKubernetesCluster_Basic(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster
	expectedURNRegEx, _ := regexp.Compile(`do:kubernetes:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "surge_upgrade", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "ha", "false"),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "cluster_subnet"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "service_subnet"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "endpoint"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "tags.#", "3"),
					resource.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "tags.*", "foo"),
					resource.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "tags.*", "foo"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "status"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "updated_at"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.actual_node_count", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.tags.*", "one"),
					resource.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.tags.*", "two"),
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
					resource.TestMatchResourceAttr("digitalocean_kubernetes_cluster.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.day"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.start_time"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "registry_integration", "false"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "destroy_all_associated_resources", "false"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_utilization_threshold"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_unneeded_time"),
				),
			},
			// Update: remove default node_pool taints
			{
				Config: fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  version       = data.digitalocean_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar", "one"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
  }
}`, testClusterVersionLatest, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.tags.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.nodes.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.taint.#", "0"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_CreateWithHAControlPlane(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "nyc1"
  ha      = true
  version = data.digitalocean_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
				`, testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "ha", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "ipv4_address", ""),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "status"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "endpoint"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_CreateWithRegistry(t *testing.T) {
	var (
		rName          = acceptance.RandomTestName()
		k8s            godo.KubernetesCluster
		registryConfig = fmt.Sprintf(`
resource "digitalocean_container_registry" "foobar" {
  name                   = "%s"
  region                 = "nyc3"
  subscription_tier_slug = "starter"
}`, rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create container registry
			{
				Config: registryConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_container_registry.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_container_registry.foobar", "endpoint", "registry.digitalocean.com/"+rName),
					resource.TestCheckResourceAttr("digitalocean_container_registry.foobar", "server_url", "registry.digitalocean.com"),
					resource.TestCheckResourceAttr("digitalocean_container_registry.foobar", "subscription_tier_slug", "starter"),
					resource.TestCheckResourceAttr("digitalocean_container_registry.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttrSet("digitalocean_container_registry.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_container_registry.foobar", "storage_usage_bytes"),
				),
			},
			// Create cluster with registry integration enabled
			{
				Config: fmt.Sprintf(`%s

%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name                 = "%s"
  region               = "nyc3"
  registry_integration = true
  version              = data.digitalocean_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
				`, testClusterVersionLatest, registryConfig, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "registry_integration", "true"),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "status"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "endpoint"),
				),
			},
			// Disable registry integration
			{
				Config: fmt.Sprintf(`%s

%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "nyc3"
  version = data.digitalocean_kubernetes_versions.test.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
				`, testClusterVersionLatest, registryConfig, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "registry_integration", "false"),
				),
			},
			// Re-enable registry integration
			{
				Config: fmt.Sprintf(`%s

%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name                 = "%s"
  region               = "nyc3"
  version              = data.digitalocean_kubernetes_versions.test.latest_version
  registry_integration = true

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
				`, testClusterVersionLatest, registryConfig, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "registry_integration", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_UpdateCluster(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "ha", "false"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigBasic4(testClusterVersionLatest, rName+"-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "tags.*", "one"),
					resource.TestCheckTypeSetElemAttr("digitalocean_kubernetes_cluster.foobar", "tags.*", "two"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.%", "0"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "surge_upgrade", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "ha", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_utilization_threshold", "0.8"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_unneeded_time", "2m"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_MaintenancePolicy(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	policy := `
	maintenance_policy {
		day = "monday"
		start_time = "00:00"
	}
`

	updatedPolicy := `
	maintenance_policy {
		day = "any"
		start_time = "04:00"
	}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigMaintenancePolicy(testClusterVersionLatest, rName, policy),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.day", "monday"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.start_time", "00:00"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.duration"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigMaintenancePolicy(testClusterVersionLatest, rName, updatedPolicy),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.day", "any"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.start_time", "04:00"),
					resource.TestCheckResourceAttrSet("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.duration"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_ControlPlaneFirewall(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	firewall := `
	control_plane_firewall {
		enabled = true
		allowed_addresses = ["1.2.3.4/16"]
	}
`

	firewallUpdate := `
	control_plane_firewall {
		enabled = false
		allowed_addresses = ["1.2.3.4/16", "5.6.7.8/16"]
	}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigControlPlaneFirewall(testClusterVersionLatest, rName, firewall),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "control_plane_firewall.0.enabled", "true"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "control_plane_firewall.0.allowed_addresses.0", "1.2.3.4/16"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigControlPlaneFirewall(testClusterVersionLatest, rName, firewallUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "control_plane_firewall.0.enabled", "false"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "control_plane_firewall.0.allowed_addresses.0", "1.2.3.4/16"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "control_plane_firewall.0.allowed_addresses.1", "5.6.7.8/16"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_ClusterAutoscalerConfiguration(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	clusterAutoscalerConfiguration := `
	cluster_autoscaler_configuration {
		scale_down_utilization_threshold = 0.5
		scale_down_unneeded_time = "1m30s"
	}
`

	updatedClusterAutoscalerConfiguration := `
	cluster_autoscaler_configuration {
		scale_down_utilization_threshold = 0.8
		scale_down_unneeded_time = "2m"
	}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigClusterAutoscalerConfiguration(testClusterVersionLatest, rName, clusterAutoscalerConfiguration),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_utilization_threshold", "0.5"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_unneeded_time", "1m30s"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigClusterAutoscalerConfiguration(testClusterVersionLatest, rName, updatedClusterAutoscalerConfiguration),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.day", "any"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_utilization_threshold", "0.8"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_unneeded_time", "2m"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_UpdatePoolDetails(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersionLatest, rName),
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
				Config: testAccDigitalOceanKubernetesConfigBasic2(testClusterVersionLatest, rName),
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
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersionLatest, rName),
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
				Config: testAccDigitalOceanKubernetesConfigBasic3(testClusterVersionLatest, rName),
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
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			// Create with auto-scaling and explicit node_count.
			{
				Config: fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name         = "%s"
  region       = "lon1"
  version      = data.digitalocean_kubernetes_versions.test.latest_version
  auto_upgrade = true

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
  maintenance_policy {
    start_time = "05:00"
    day        = "sunday"
  }
}
				`, testClusterVersionLatest, rName),
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
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.day", "sunday"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.start_time", "05:00"),
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
    name       = "default"
    size       = "s-1vcpu-2gb"
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
}
				`, testClusterVersionLatest, rName),
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
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 2
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
}
				`, testClusterVersionLatest, rName),
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
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 2
  }
}
				`, testClusterVersionLatest, rName),
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
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
			`, testClusterVersionLatest, rName),
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
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
}
				`, testClusterVersionLatest, rName),
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
    name       = "default"
    size       = "s-1vcpu-2gb"
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 3
  }
}
				`, testClusterVersionLatest, rName),
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
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"kubernetes": {
				Source:            "hashicorp/kubernetes",
				VersionConstraint: "2.0.1",
			},
		},
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfig_KubernetesProviderInteroperability(testClusterVersionLatest, rName),
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
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersionPrevious, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
				),
			},
			{
				Config: testAccDigitalOceanKubernetesConfigBasic(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPtr("digitalocean_kubernetes_cluster.foobar", "id", &k8s.ID),
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_DestroyAssociated(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigDestroyAssociated(testClusterVersionPrevious, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttrPair("digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "destroy_all_associated_resources", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKubernetesCluster_VPCNative(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigVPCNative(testClusterVersionLatest, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "cluster_subnet", "192.168.0.0/20"),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "service_subnet", "192.168.16.0/22"),
				),
			},
		},
	})
}

func testAccDigitalOceanKubernetesConfigBasic(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "nyc1"
  version       = data.digitalocean_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar", "one"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
    taint {
      key    = "key1"
      value  = "val1"
      effect = "PreferNoSchedule"
    }
  }
	
	cluster_autoscaler_configuration {
		scale_down_utilization_threshold = 0.5
		scale_down_unneeded_time = "1m30s"
	}
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigClusterAutoscalerConfiguration(testClusterVersion string, rName string, clusterAutoscalerConfiguration string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  version       = data.digitalocean_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar", "one"]

%s

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
    taint {
      key    = "key1"
      value  = "val1"
      effect = "PreferNoSchedule"
    }
  }
}
`, testClusterVersion, rName, clusterAutoscalerConfiguration)
}

func testAccDigitalOceanKubernetesConfigMaintenancePolicy(testClusterVersion string, rName string, policy string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  version       = data.digitalocean_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar", "one"]

%s

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
    taint {
      key    = "key1"
      value  = "val1"
      effect = "PreferNoSchedule"
    }
  }
}
`, testClusterVersion, rName, policy)
}

func testAccDigitalOceanKubernetesConfigControlPlaneFirewall(testClusterVersion string, rName string, controlPlaneFirewall string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  version       = data.digitalocean_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar", "one"]

%s

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
    taint {
      key    = "key1"
      value  = "val1"
      effect = "PreferNoSchedule"
    }
  }
}
`, testClusterVersion, rName, controlPlaneFirewall)
}

func testAccDigitalOceanKubernetesConfigBasic2(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  version       = data.digitalocean_kubernetes_versions.test.latest_version
  surge_upgrade = true
  tags          = ["foo", "bar"]

  node_pool {
    name       = "default-rename"
    size       = "s-1vcpu-2gb"
    node_count = 2
    tags       = ["one", "two", "three"]
    labels = {
      priority = "high"
      purpose  = "awesome"
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
  tags    = ["foo", "bar"]

  node_pool {
    name       = "default"
    size       = "s-2vcpu-4gb"
    node_count = 1
    tags       = ["one", "two"]
  }
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigBasic4(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name          = "%s"
  region        = "lon1"
  surge_upgrade = true
  ha            = true
  version       = data.digitalocean_kubernetes_versions.test.latest_version
  tags          = ["one", "two"]

  node_pool {
    name       = "default"
    size       = "s-2vcpu-4gb"
    node_count = 1
    tags       = ["foo", "bar"]
  }

	cluster_autoscaler_configuration {
		scale_down_utilization_threshold = 0.8
		scale_down_unneeded_time = "2m"
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
    name       = "default"
    size       = "s-2vcpu-4gb"
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

resource "kubernetes_namespace" "example" {
  metadata {
    name = "example-namespace"
  }
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigDestroyAssociated(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name                             = "%s"
  region                           = "nyc1"
  version                          = data.digitalocean_kubernetes_versions.test.latest_version
  destroy_all_associated_resources = true

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
`, testClusterVersion, rName)
}

func testAccDigitalOceanKubernetesConfigVPCNative(testClusterVersion string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foobar" {
  name           = "%s"
  region         = "nyc1"
  version        = data.digitalocean_kubernetes_versions.test.latest_version
  cluster_subnet = "192.168.0.0/20"
  service_subnet = "192.168.16.0/22"
  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}
`, testClusterVersion, rName)
}

func testAccCheckDigitalOceanKubernetesClusterDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
		filteredTags := kubernetes.FilterTags(tt.have)
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
	kubeConfigRendered, err := kubernetes.RenderKubeconfig("test-cluster", "lon1", &creds)
	if err != nil {
		t.Errorf("error calling renderKubeconfig: %s", err)

	}
	got := string(kubeConfigRendered)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("renderKubeconfig returned %+v\n, expected %+v\n", got, expected)
	}
}
