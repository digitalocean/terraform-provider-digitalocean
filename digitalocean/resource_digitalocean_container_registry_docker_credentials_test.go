package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanContainerRegistryDockerCredentials_Basic(t *testing.T) {
	var reg godo.Registry

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanContainerRegistryDockerCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanContainerRegistryDockerCredentialsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryDockerCredentialsExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryDockerCredentialsAttributes(&reg),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry_docker_credentials.foobar", "registry_name", "foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry_docker_credentials.foobar", "write", "true"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry_docker_credentials.foobar", "docker_credentials"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry_docker_credentials.foobar", "credential_expiration_time"),
				),
			},
		},
	})
}

func TestAccDigitalOceanContainerRegistryDockerCredentials_withExpiry(t *testing.T) {
	var reg godo.Registry

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanContainerRegistryDockerCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanContainerRegistryDockerCredentialsConfig_withExpiry,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryDockerCredentialsExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryDockerCredentialsAttributes(&reg),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry_docker_credentials.foobar", "registry_name", "foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry_docker_credentials.foobar", "write", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry_docker_credentials.foobar", "expiry_seconds", "3600"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry_docker_credentials.foobar", "docker_credentials"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry_docker_credentials.foobar", "credential_expiration_time"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanContainerRegistryDockerCredentialsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_container_registry_docker_credentials" {
			continue
		}

		// Try to find the key
		_, _, err := client.Registry.Get(context.Background())

		if err == nil {
			return fmt.Errorf("Container Registry still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanContainerRegistryDockerCredentialsAttributes(reg *godo.Registry) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if reg.Name != "foobar" {
			return fmt.Errorf("Bad name: %s", reg.Name)
		}

		return nil
	}
}

func testAccCheckDigitalOceanContainerRegistryDockerCredentialsExists(n string, reg *godo.Registry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		// Try to find the registry
		foundReg, _, err := client.Registry.Get(context.Background())

		if err != nil {
			return err
		}

		*reg = *foundReg

		return nil
	}
}

var testAccCheckDigitalOceanContainerRegistryDockerCredentialsConfig_basic = `
resource "digitalocean_container_registry" "foobar" {
	name = "foobar"
}

resource "digitalocean_container_registry_docker_credentials" "foobar" {
	registry_name = digitalocean_container_registry.foobar.name
	write = true
}`

var testAccCheckDigitalOceanContainerRegistryDockerCredentialsConfig_withExpiry = `
resource "digitalocean_container_registry" "foobar" {
	name = "foobar"
}

resource "digitalocean_container_registry_docker_credentials" "foobar" {
	registry_name = digitalocean_container_registry.foobar.name
	write = true
	expiry_seconds = 3600
}`
