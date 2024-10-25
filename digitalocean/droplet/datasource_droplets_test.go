package droplet_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDroplets_Basic(t *testing.T) {
	name1 := acceptance.RandomTestName("01")
	name2 := acceptance.RandomTestName("02")

	resourcesConfig := fmt.Sprintf(`
resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
}

resource "digitalocean_droplet" "bar" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
}
`, name1, defaultSize, defaultImage, name2, defaultSize, defaultImage)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_droplets" "result" {
  filter {
    key    = "name"
    values = ["%s"]
  }
}
`, name1)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_droplets.result", "droplets.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_droplets.result", "droplets.0.name", name1),
					resource.TestCheckResourceAttrPair("data.digitalocean_droplets.result", "droplets.0.id", "digitalocean_droplet.foo", "id"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDroplets_GPUDroplet(t *testing.T) {
	runGPU := os.Getenv(runGPUEnvVar)
	if runGPU == "" {
		t.Skip("'DO_RUN_GPU_TESTS' env var not set; Skipping tests that requires a GPU Droplet")
	}

	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("digitalocean@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	name1 := acceptance.RandomTestName("gpu")
	name2 := acceptance.RandomTestName("regular")

	resourcesConfig := fmt.Sprintf(`
resource "digitalocean_ssh_key" "foobar" {
  name       = "%s"
  public_key = "%s"
}

resource "digitalocean_droplet" "gpu" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "nyc2"
  ssh_keys = [digitalocean_ssh_key.foobar.id]
}

resource "digitalocean_droplet" "regular" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc2"
}
`, keyName, publicKeyMaterial, name1, gpuSize, gpuImage, name2, defaultSize, defaultImage)

	datasourceConfig := `
data "digitalocean_droplets" "result" {
  gpus = true
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_droplets.result", "droplets.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_droplets.result", "droplets.0.name", name1),
					resource.TestCheckResourceAttrPair("data.digitalocean_droplets.result", "droplets.0.id", "digitalocean_droplet.gpu", "id"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
