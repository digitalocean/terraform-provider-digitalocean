package microdroplet_test

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

func TestAccDigitalOceanMicroDropletImage_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet_image.foobar"
	config := fmt.Sprintf(testAccMicroDropletImageConfig_Basic, name, testMicroDropletImage)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "source", testMicroDropletImage),
					resource.TestCheckResourceAttr(resourceName, "status", string(godo.MicroDropletImageStatusAvailable)),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
				),
			},
		},
	})
}

func testAccCheckMicroDropletImageExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID set for MicroDroplet image resource: %s", name)
		}

		img, _, err := client.MicroDropletImages.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if img.ID != rs.Primary.ID {
			return fmt.Errorf("MicroDroplet image not found: %s / %s", name, rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckMicroDropletImageDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_microdroplet_image" {
			continue
		}
		img, _, err := client.MicroDropletImages.Get(context.Background(), rs.Primary.ID)
		if err == nil && img.Status != godo.MicroDropletImageStatusDeleted {
			return fmt.Errorf("MicroDroplet image %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const testAccMicroDropletImageConfig_Basic = `
resource "digitalocean_microdroplet_image" "foobar" {
  name   = "%s"
  source = "%s"
}
`
