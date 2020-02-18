package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDigitalOceanProjects_Basic(t *testing.T) {
	config := `
resource "digitalocean_project" "prod" {
	name = "A production project"
	environment = "Production"
}

resource "digitalocean_project" "staging" {
	name = "A staging project"
	environment = "Staging"
}

data "digitalocean_projects" "prod" {
	filter {
      key = "environment"
      values = ["Production"]
    }
	depends_on = [digitalocean_project.prod, digitalocean_project.staging]
}

data "digitalocean_projects" "staging" {
	filter {
      key = "name"
      values = ["A staging project"]
    }
	depends_on = [digitalocean_project.prod, digitalocean_project.staging]
}

data "digitalocean_projects" "both" {
	filter {
      key = "environment"
      values = ["Production"]
    }
	filter {
      key = "name"
      values = ["A staging project"]
    }
	depends_on = [digitalocean_project.prod, digitalocean_project.staging]
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_projects.prod", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.prod", "projects.0.name", "A production project"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.prod", "projects.0.environment", "Production"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.staging", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.staging", "projects.0.name", "A staging project"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.staging", "projects.0.environment", "Staging"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.staging", "projects.#", "0"),
				),
			},
		},
	})
}
