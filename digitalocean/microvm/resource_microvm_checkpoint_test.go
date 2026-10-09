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

func TestAccDigitalOceanMicroVMCheckpoint_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm_checkpoint.foobar"
	cfg := fmt.Sprintf(testAccMicroVMCheckpointConfig_Basic, name, testMicroVMOCIRef, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMCheckpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "microvm_id", "digitalocean_microvm.foobar", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "microvm_name", name),
					resource.TestCheckResourceAttr(resourceName, "region", testMicroVMRegion),
					resource.TestCheckResourceAttr(resourceName, "status", string(godo.MicroVMCheckpointStatusAvailable)),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDigitalOceanMicroVMCheckpoint_Restore(t *testing.T) {
	name := acceptance.RandomTestName()
	restoredName := acceptance.RandomTestName("restored")
	cfg := fmt.Sprintf(testAccMicroVMCheckpointConfig_Basic, name, testMicroVMOCIRef, name) +
		fmt.Sprintf(testAccMicroVMConfig_FromCheckpoint, restoredName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckMicroVMDestroy,
			testAccCheckMicroVMCheckpointDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists("digitalocean_microvm.restored"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_microvm.restored", "source.0.checkpoint_id",
						"digitalocean_microvm_checkpoint.foobar", "id",
					),
					resource.TestCheckResourceAttr("digitalocean_microvm.restored", "region", testMicroVMRegion),
					resource.TestCheckResourceAttr("digitalocean_microvm.restored", "current_state", string(godo.MicroVMStateRunning)),
				),
			},
		},
	})
}

func testAccCheckMicroVMCheckpointDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_microvm_checkpoint" {
			continue
		}
		cp, resp, err := client.MicroVMs.GetCheckpoint(context.Background(), rs.Primary.ID)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				continue
			}
			return err
		}
		if cp.Status != godo.MicroVMCheckpointStatusDeleted {
			return fmt.Errorf("MicroVM checkpoint %s still exists (status %s)", rs.Primary.ID, cp.Status)
		}
	}
	return nil
}

const testAccMicroVMCheckpointConfig_Basic = `
resource "digitalocean_microvm" "foobar" {
  name   = "%s"
  region = "nyc3"

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "%s"
  }
}

resource "digitalocean_microvm_checkpoint" "foobar" {
  microvm_id = digitalocean_microvm.foobar.id
  name       = "%s"
}
`

const testAccMicroVMConfig_FromCheckpoint = `
resource "digitalocean_microvm" "restored" {
  name = "%s"

  source {
    checkpoint_id = digitalocean_microvm_checkpoint.foobar.id
  }
}
`
