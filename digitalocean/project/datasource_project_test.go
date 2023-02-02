package project_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanProject_DefaultProject(t *testing.T) {
	config := `
data "digitalocean_project" "default" {
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_project.default", "id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_project.default", "name"),
					resource.TestCheckResourceAttr("data.digitalocean_project.default", "is_default", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanProject_NonDefaultProject(t *testing.T) {
	nonDefaultProjectName := acceptance.RandomTestName("project")
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_project" "foo" {
  name = "%s"
}`, nonDefaultProjectName)
	dataSourceConfig := `
data "digitalocean_project" "bar" {
  id = digitalocean_project.foo.id
}

data "digitalocean_project" "barfoo" {
  name = digitalocean_project.foo.name
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_project.bar", "id"),
					resource.TestCheckResourceAttr("data.digitalocean_project.bar", "is_default", "false"),
					resource.TestCheckResourceAttr("data.digitalocean_project.bar", "name", nonDefaultProjectName),
					resource.TestCheckResourceAttr("data.digitalocean_project.barfoo", "is_default", "false"),
					resource.TestCheckResourceAttr("data.digitalocean_project.barfoo", "name", nonDefaultProjectName),
				),
			},
		},
	})
}
