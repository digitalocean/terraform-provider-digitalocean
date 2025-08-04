package kubernetes_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanKubernetesCluster_Basic(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster
	expectedURNRegEx, _ := regexp.Compile(`do:kubernetes:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)
	resourceConfig := testAccDigitalOceanKubernetesConfigForDataSource(testClusterVersionLatest, rName)
	dataSourceConfig := `
data "digitalocean_kubernetes_cluster" "foobar" {
  name = digitalocean_kubernetes_cluster.foo.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"kubernetes": {
				Source:            "hashicorp/kubernetes",
				VersionConstraint: "1.13.2",
			},
		},
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanKubernetesClusterExists("data.digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "name", rName),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "region", "lon1"),
					resource.TestCheckResourceAttrPair("data.digitalocean_kubernetes_cluster.foobar", "version", "data.digitalocean_kubernetes_versions.test", "latest_version"),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "node_pool.0.labels.priority", "high"),
					resource.TestCheckResourceAttrSet("data.digitalocean_kubernetes_cluster.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttrSet("data.digitalocean_kubernetes_cluster.foobar", "auto_upgrade"),
					resource.TestMatchResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.day", "monday"),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.start_time", "00:00"),
					resource.TestCheckResourceAttrSet("data.digitalocean_kubernetes_cluster.foobar", "maintenance_policy.0.duration"),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_utilization_threshold", "0.5"),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.scale_down_unneeded_time", "1m30s"),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "cluster_autoscaler_configuration.0.expanders.0", "priority"),
					resource.TestCheckResourceAttr("data.digitalocean_kubernetes_cluster.foobar", "routing_agent.0.enabled", "true"),
				),
			},
		},
	})
}

func testAccDigitalOceanKubernetesConfigForDataSource(version string, rName string) string {
	return fmt.Sprintf(`%s

resource "digitalocean_kubernetes_cluster" "foo" {
  name         = "%s"
  region       = "lon1"
  version      = data.digitalocean_kubernetes_versions.test.latest_version
  tags         = ["foo", "bar"]
  auto_upgrade = true

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"]
    labels = {
      priority = "high"
    }
  }
  maintenance_policy {
    day        = "monday"
    start_time = "00:00"
  }

  cluster_autoscaler_configuration {
    scale_down_utilization_threshold = 0.5
    scale_down_unneeded_time         = "1m30s"
    expanders                        = ["priority"]
  }

  routing_agent {
    enabled = true
  }
}`, version, rName)
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

func TestAccDataSourceDigitalOceanKubernetesCluster_NodeCountZero(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	clusterConfig := fmt.Sprintf(`%s
resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.digitalocean_kubernetes_versions.test.latest_version
  tags    = ["foo", "bar"]

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 0
    tags       = ["one", "two"]
  }
}
`, testClusterVersionLatest, rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: clusterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
				),
			},
		},
	})
}
