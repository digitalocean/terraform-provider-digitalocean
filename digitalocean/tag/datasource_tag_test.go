package tag_test

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

func TestAccDataSourceDigitalOceanTag_Basic(t *testing.T) {
	var tag godo.Tag
	tagName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_tag" "foo" {
  name = "%s"
}`, tagName)
	dataSourceConfig := `
data "digitalocean_tag" "foobar" {
  name = "${digitalocean_tag.foo.name}"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanTagExists("data.digitalocean_tag.foobar", &tag),
					resource.TestCheckResourceAttr(
						"data.digitalocean_tag.foobar", "name", tagName),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tag.foobar", "total_resource_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tag.foobar", "droplets_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tag.foobar", "images_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tag.foobar", "volumes_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tag.foobar", "volume_snapshots_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tag.foobar", "databases_count"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanTagExists(n string, tag *godo.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No tag ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundTag, _, err := client.Tags.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundTag.Name != rs.Primary.ID {
			return fmt.Errorf("Tag not found")
		}

		*tag = *foundTag

		return nil
	}
}
