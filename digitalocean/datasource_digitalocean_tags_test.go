package digitalocean

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDigitalOceanTags_Basic(t *testing.T) {
	var tag godo.Tag
	tagName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanTagsConfig_basic, tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanTagExists("digitalocean_tag.foo", &tag),
					resource.TestCheckResourceAttr(
						"data.digitalocean_tags.foobar", "tags.0.name", tagName),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tags.foobar", "tags.0.resource_count"),
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

const testAccCheckDataSourceDigitalOceanTagsConfig_basic = `
resource "digitalocean_tag" "foo" {
  name = "%s"
}

data "digitalocean_tags" "foobar" {
  filter {
    key    = "name"
    values = [digitalocean_tag.foo.name]
  }
}`
