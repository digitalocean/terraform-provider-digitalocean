package tag_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanTags_Basic(t *testing.T) {
	var tag godo.Tag
	tagName := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_tag" "foo" {
  name = "%s"
}`, tagName)
	dataSourceConfig := `
data "digitalocean_tags" "foobar" {
  filter {
    key    = "name"
    values = [digitalocean_tag.foo.name]
  }
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
					testAccCheckDataSourceDigitalOceanTagExists("digitalocean_tag.foo", &tag),
					resource.TestCheckResourceAttr(
						"data.digitalocean_tags.foobar", "tags.0.name", tagName),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tags.foobar", "tags.0.total_resource_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tags.foobar", "tags.0.droplets_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tags.foobar", "tags.0.images_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tags.foobar", "tags.0.volumes_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tags.foobar", "tags.0.volume_snapshots_count"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_tags.foobar", "tags.0.databases_count"),
				),
			},
		},
	})
}
