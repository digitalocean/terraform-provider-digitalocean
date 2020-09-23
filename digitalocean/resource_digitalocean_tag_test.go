package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanTag_Basic(t *testing.T) {
	var tag godo.Tag

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanTagConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanTagExists("digitalocean_tag.foobar", &tag),
					testAccCheckDigitalOceanTagAttributes(&tag),
					resource.TestCheckResourceAttr(
						"digitalocean_tag.foobar", "name", "foobar"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_tag.foobar", "total_resource_count"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_tag.foobar", "droplets_count"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_tag.foobar", "images_count"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_tag.foobar", "volumes_count"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_tag.foobar", "volume_snapshots_count"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_tag.foobar", "databases_count"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanTagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_tag" {
			continue
		}

		// Try to find the key
		_, _, err := client.Tags.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Tag still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanTagAttributes(tag *godo.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if tag.Name != "foobar" {
			return fmt.Errorf("Bad name: %s", tag.Name)
		}

		return nil
	}
}

func testAccCheckDigitalOceanTagExists(n string, tag *godo.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		// Try to find the tag
		foundTag, _, err := client.Tags.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		*tag = *foundTag

		return nil
	}
}

var testAccCheckDigitalOceanTagConfig_basic = fmt.Sprintf(`
resource "digitalocean_tag" "foobar" {
    name = "foobar"
}`)
