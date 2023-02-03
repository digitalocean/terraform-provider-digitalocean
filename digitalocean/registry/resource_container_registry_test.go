package registry_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/registry"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanContainerRegistry_Basic(t *testing.T) {
	var reg godo.Registry
	name := acceptance.RandomTestName()
	starterConfig := fmt.Sprintf(testAccCheckDigitalOceanContainerRegistryConfig_basic, name, "starter", "")
	basicConfig := fmt.Sprintf(testAccCheckDigitalOceanContainerRegistryConfig_basic, name, "basic", "")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: starterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryAttributes(&reg, name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "endpoint", "registry.digitalocean.com/"+name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "server_url", "registry.digitalocean.com"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "subscription_tier_slug", "starter"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "region"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "storage_usage_bytes"),
				),
			},
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryAttributes(&reg, name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "endpoint", "registry.digitalocean.com/"+name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "server_url", "registry.digitalocean.com"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "subscription_tier_slug", "basic"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "region"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "storage_usage_bytes"),
				),
			},
		},
	})
}

func TestAccDigitalOceanContainerRegistry_CustomRegion(t *testing.T) {
	var reg godo.Registry
	name := acceptance.RandomTestName()
	starterConfig := fmt.Sprintf(testAccCheckDigitalOceanContainerRegistryConfig_basic, name, "starter", `  region = "sfo3"`)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: starterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanContainerRegistryExists("digitalocean_container_registry.foobar", &reg),
					testAccCheckDigitalOceanContainerRegistryAttributes(&reg, name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "endpoint", "registry.digitalocean.com/"+name),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "server_url", "registry.digitalocean.com"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "subscription_tier_slug", "starter"),
					resource.TestCheckResourceAttr(
						"digitalocean_container_registry.foobar", "region", "sfo3"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_container_registry.foobar", "storage_usage_bytes"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanContainerRegistryDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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

func testAccCheckDigitalOceanContainerRegistryAttributes(reg *godo.Registry, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if reg.Name != name {
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

var testAccCheckDigitalOceanContainerRegistryConfig_basic = `
resource "digitalocean_container_registry" "foobar" {
  name                   = "%s"
  subscription_tier_slug = "%s"
  %s
}`

func TestRevokeOAuthToken(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	token, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("error creating fake token: %s", err.Error())
	}

	mux.HandleFunc("/revoke", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost != r.Method {
			t.Errorf("method = %v, expected %v", r.Method, http.MethodPost)
		}

		authHeader := r.Header.Get("Authorization")
		expectedAuth := fmt.Sprintf("Bearer %s", token)
		if authHeader != expectedAuth {
			t.Errorf("auth header  = %v, expected %v", authHeader, expectedAuth)
		}

		r.ParseForm()
		bodyToken := r.Form.Get("token")
		if token != bodyToken {
			t.Errorf("token  = %v, expected %v", bodyToken, token)
		}

		w.WriteHeader(http.StatusOK)
	})

	err = registry.RevokeOAuthToken(token, server.URL+"/revoke")
	if err != nil {
		t.Errorf("error revoking token: %s", err.Error())
	}
}
