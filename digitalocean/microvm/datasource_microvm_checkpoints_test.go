package microvm_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanMicroVMCheckpoints_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroVMCheckpointConfig_Basic, name, testMicroVMOCIRef, name)

	dsConfig := `
data "digitalocean_microvm_checkpoints" "by_id" {
  microvm_id = digitalocean_microvm.foobar.id
  depends_on = [digitalocean_microvm_checkpoint.foobar]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMCheckpointDestroy,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microvm_checkpoints.by_id", "checkpoints.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_microvm_checkpoints.by_id", "checkpoints.0.id",
						"digitalocean_microvm_checkpoint.foobar", "id",
					),
					resource.TestCheckResourceAttr("data.digitalocean_microvm_checkpoints.by_id", "checkpoints.0.status", "CHECKPOINT_AVAILABLE"),
				),
			},
		},
	})
}
