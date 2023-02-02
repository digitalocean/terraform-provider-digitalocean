package registry_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanContainerRegistryDockerCredentials_Basic(t *testing.T) {
	var reg godo.Registry
	name := acceptance.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanContainerRegistryDockerCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanContainerRegistryDockerCredentialsConfig_basic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryDockerCredentialsExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryDockerCredentialsAttributes(&reg, name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry_docker_credentials.foobar", "registry_name", name),
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
	name := acceptance.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanContainerRegistryDockerCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanContainerRegistryDockerCredentialsConfig_withExpiry, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryDockerCredentialsExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryDockerCredentialsAttributes(&reg, name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry_docker_credentials.foobar", "registry_name", name),
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
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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

func testAccCheckDigitalOceanContainerRegistryDockerCredentialsAttributes(reg *godo.Registry, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if reg.Name != name {
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

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
	name                   = "%s"
	subscription_tier_slug = "basic"
}

resource "digitalocean_container_registry_docker_credentials" "foobar" {
	registry_name = digitalocean_container_registry.foobar.name
	write = true
}`

var testAccCheckDigitalOceanContainerRegistryDockerCredentialsConfig_withExpiry = `
resource "digitalocean_container_registry" "foobar" {
	name                   = "%s"
	subscription_tier_slug = "basic"
}

resource "digitalocean_container_registry_docker_credentials" "foobar" {
	registry_name = digitalocean_container_registry.foobar.name
	write = true
	expiry_seconds = 3600
}`
