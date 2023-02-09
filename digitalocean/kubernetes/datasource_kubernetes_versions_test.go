package kubernetes_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanKubernetesVersions_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanKubernetesVersionsConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_kubernetes_versions.foobar", "latest_version"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanKubernetesVersions_Filtered(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanKubernetesVersionsConfig_filtered),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_kubernetes_versions.foobar", "valid_versions.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanKubernetesVersions_CreateCluster(t *testing.T) {
	rName := acceptance.RandomTestName()
	var k8s godo.KubernetesCluster

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanKubernetesVersionsConfig_create, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_kubernetes_versions.foobar", "latest_version"),
					testAccCheckDigitalOceanKubernetesClusterExists(
						"digitalocean_kubernetes_cluster.foobar", &k8s),
					resource.TestCheckResourceAttr(
						"digitalocean_kubernetes_cluster.foobar", "name", rName),
				),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanKubernetesVersionsConfig_basic = `
data "digitalocean_kubernetes_versions" "foobar" {}`

const testAccCheckDataSourceDigitalOceanKubernetesVersionsConfig_filtered = `
data "digitalocean_kubernetes_versions" "foobar" {
  version_prefix = "1.12." # No longer supported, should be empty
}`

const testAccCheckDataSourceDigitalOceanKubernetesVersionsConfig_create = `
data "digitalocean_kubernetes_versions" "foobar" {
}

resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = data.digitalocean_kubernetes_versions.foobar.latest_version

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}`
