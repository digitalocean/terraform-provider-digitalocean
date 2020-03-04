package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDigitalOceanProject_DefaultProject(t *testing.T) {
	config := `
data "digitalocean_project" "default" {
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
	nonDefaultProjectName := randomName("tf-acc-project-", 6)
	config := fmt.Sprintf(`
resource "digitalocean_project" "foo" {
	name = "%s"
}

data "digitalocean_project" "bar" {
  	id = digitalocean_project.foo.id
}

data "digitalocean_project" "barfoo" {
    name = "%s"
}
`, nonDefaultProjectName, nonDefaultProjectName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
