package kubernetes_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanKubernetesClusterCredentials_Ephemeral(t *testing.T) {
	rName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acceptance.TestAccProtoV5ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"kubernetes": {
				Source:            "hashicorp/kubernetes",
				VersionConstraint: "2.38.0",
			},
		},
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`%s

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

ephemeral "digitalocean_kubernetes_cluster_credentials" "foobar" {
  cluster_id = digitalocean_kubernetes_cluster.foobar.id
}

provider "kubernetes" {
  host                   = ephemeral.digitalocean_kubernetes_cluster_credentials.foobar.host
  cluster_ca_certificate = base64decode(ephemeral.digitalocean_kubernetes_cluster_credentials.foobar.cluster_ca_certificate)
  token                  = ephemeral.digitalocean_kubernetes_cluster_credentials.foobar.token
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "example-namespace"
  }
}
`, testClusterVersionLatest, rName),
				Check: resource.TestCheckResourceAttr("kubernetes_namespace.example", "metadata.0.name", "example-namespace"),
			},
		},
	})
}
