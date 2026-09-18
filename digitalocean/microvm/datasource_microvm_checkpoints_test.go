package microvm_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceDigitalOceanMicroVMCheckpoints_Basic pauses a fresh
// MicroVM (which triggers checkpoint creation platform-side) and then
// reads the digitalocean_microvm_checkpoints data source. Because
// checkpoint creation is asynchronous, we only assert that the datasource
// returns a well-formed list and that `checkpoints.#` is populated. Callers
// pausing on production traffic will get a non-zero count once the platform
// has captured the checkpoint.
func TestAccDataSourceDigitalOceanMicroVMCheckpoints_Basic(t *testing.T) {
	name := acceptance.RandomTestName()

	pausedConfig := fmt.Sprintf(testAccMicroVMConfig_State,
		name, string(godo.MicroVMStatePaused), testMicroVMOCIRef)

	dsConfig := `
data "digitalocean_microvm_checkpoints" "by_id" {
  microvm_id = digitalocean_microvm.foobar.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMDestroy,
		Steps: []resource.TestStep{
			{Config: pausedConfig},
			{
				Config: pausedConfig + dsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_microvm_checkpoints.by_id",
						"checkpoints.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_microvm_checkpoints.by_id",
						"microvm_id",
					),
				),
			},
		},
	})
}
