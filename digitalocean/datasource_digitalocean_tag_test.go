package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceDigitalOceanTag_Basic(t *testing.T) {
	var tag godo.Tag
	tagName := fmt.Sprintf("foo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanTagConfig_basic, tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanTagExists("data.digitalocean_tag.foobar", &tag),
					resource.TestCheckResourceAttr(
						"data.digitalocean_tag.foobar", "name", tagName),
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

		client := testAccProvider.Meta().(*godo.Client)

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
