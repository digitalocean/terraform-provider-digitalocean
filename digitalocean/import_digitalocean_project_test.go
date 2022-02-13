package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanProject_importBasic(t *testing.T) {
	name := generateProjectName()
	resourceName := "digitalocean_project.myproj"
	createConfig := fixtureCreateWithDefaults(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
