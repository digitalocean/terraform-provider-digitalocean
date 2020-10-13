package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanTag_Basic(t *testing.T) {
	var tag godo.Tag
	tagName := fmt.Sprintf("foo-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanTagConfig_basic, tagName),
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

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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

const testAccCheckDataSourceDigitalOceanTagConfig_basic = `
resource "digitalocean_tag" "foo" {
  name = "%s"
}

data "digitalocean_tag" "foobar" {
  name = "${digitalocean_tag.foo.name}"
}`
