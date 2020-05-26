package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanContainerRegistry_Basic(t *testing.T) {
	var reg godo.Registry

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanContainerRegistryConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryAttributes(&reg),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "name", "foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "write", "false"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "docker_credentials"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "credential_expiration_time"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "endpoint", "registry.digitalocean.com/foobar"),
				),
			},
			{
				Config: testAccCheckDigitalOceanContainerRegistryConfig_write,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryAttributes(&reg),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "name", "foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "write", "true"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "docker_credentials"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "credential_expiration_time"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "endpoint", "registry.digitalocean.com/foobar"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanContainerRegistryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_container_registry" {
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

func testAccCheckDigitalOceanContainerRegistryAttributes(reg *godo.Registry) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if reg.Name != "foobar" {
			return fmt.Errorf("Bad name: %s", reg.Name)
		}

		return nil
	}
}

func testAccCheckDigitalOceanContainerRegistryExists(n string, reg *godo.Registry) resource.TestCheckFunc {
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

var testAccCheckDigitalOceanContainerRegistryConfig_basic = `
resource "digitalocean_container_registry" "foobar" {
    name = "foobar"
}`

var testAccCheckDigitalOceanContainerRegistryConfig_write = `
resource "digitalocean_container_registry" "foobar" {
    name  = "foobar"
    write = true
}`
