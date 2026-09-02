package microdroplet_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceDigitalOceanMicroDropletCheckpoints_Basic pauses a fresh
// MicroDroplet (which triggers checkpoint creation platform-side) and then
// reads the digitalocean_microdroplet_checkpoints data source. Because
// checkpoint creation is asynchronous, we only assert that the datasource
// returns a well-formed list and that `checkpoints.#` is populated. Callers
// pausing on production traffic will get a non-zero count once the platform
// has captured the checkpoint.
func TestAccDataSourceDigitalOceanMicroDropletCheckpoints_Basic(t *testing.T) {
	name := acceptance.RandomTestName()

	pausedConfig := fmt.Sprintf(testAccMicroDropletConfig_State,
		name, testMicroDropletImage, string(godo.MicroDropletStatePaused))

	dsConfig := `
data "digitalocean_microdroplet_checkpoints" "by_id" {
  microdroplet_id = digitalocean_microdroplet.foobar.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			{Config: pausedConfig},
			{
				Config: pausedConfig + dsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_microdroplet_checkpoints.by_id",
						"checkpoints.#",
					),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_microdroplet_checkpoints.by_id",
						"microdroplet_id",
					),
				),
			},
		},
	})
}
