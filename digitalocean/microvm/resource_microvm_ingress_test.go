package microvm_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccDigitalOceanMicroVMIngress_Basic adds an ingress URL to a MicroVM
// through the dedicated resource and verifies the URL is minted and active.
func TestAccDigitalOceanMicroVMIngress_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm_ingress.foobar"
	cfg := fmt.Sprintf(testAccMicroVMIngressConfig_Basic, name, testMicroVMOCIRef)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckMicroVMDestroy,
			testAccCheckMicroVMIngressDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists("digitalocean_microvm.foobar"),
					resource.TestCheckResourceAttrPair(resourceName, "microvm_id", "digitalocean_microvm.foobar", "id"),
					resource.TestCheckResourceAttr(resourceName, "port", "9090"),
					resource.TestCheckResourceAttr(resourceName, "status", string(godo.MicroVMURLStatusActive)),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "url_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       testAccMicroVMIngressImportID(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// TestAccDigitalOceanMicroVM_IngressDisabled creates a MicroVM with
// ingress_enabled = false and verifies no default URL is minted.
func TestAccDigitalOceanMicroVM_IngressDisabled(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm.foobar"
	cfg := fmt.Sprintf(testAccMicroVMConfig_IngressDisabled, name, testMicroVMOCIRef)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ingress_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "urls.#", "0"),
				),
			},
		},
	})
}

func testAccMicroVMIngressImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["microvm_id"], rs.Primary.ID), nil
	}
}

func testAccCheckMicroVMIngressDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_microvm_ingress" {
			continue
		}
		microVMID := rs.Primary.Attributes["microvm_id"]
		m, resp, err := client.MicroVMs.Get(context.Background(), microVMID)
		if err != nil {
			// MicroVM gone means the URL is gone too.
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				continue
			}
			return err
		}
		for _, u := range m.URLs {
			if u.ID == rs.Primary.ID {
				return fmt.Errorf("MicroVM ingress URL %s still exists on MicroVM %s", rs.Primary.ID, microVMID)
			}
		}
	}
	return nil
}

const testAccMicroVMIngressConfig_Basic = `
resource "digitalocean_microvm" "foobar" {
  name      = "%s"
  region    = "nyc3"
  http_port = 8080
  ports     = [8080, 9090]

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "%s"
  }
}

resource "digitalocean_microvm_ingress" "foobar" {
  microvm_id = digitalocean_microvm.foobar.id
  port       = 9090
}
`

const testAccMicroVMConfig_IngressDisabled = `
resource "digitalocean_microvm" "foobar" {
  name            = "%s"
  region          = "nyc3"
  ingress_enabled = false

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "%s"
  }
}
`
