package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDigitalOceanProjects_Basic(t *testing.T) {
	prodProjectName := randomName("tf-acc-project-", 6)
	stagingProjectName := randomName("tf-acc-project-", 6)

	config := fmt.Sprintf(`
resource "digitalocean_project" "prod" {
	name = "%s"
	environment = "Production"
}

resource "digitalocean_project" "staging" {
	name = "%s"
	environment = "Staging"
}

data "digitalocean_projects" "prod" {
	filter {
      key = "environment"
      values = ["Production"]
    }
    filter {
      key = "is_default"
      values = ["false"]
    }
	depends_on = [digitalocean_project.prod, digitalocean_project.staging]
}

data "digitalocean_projects" "staging" {
	filter {
      key = "name"
      values = ["%s"]
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
      values = ["%s"]
    }
	depends_on = [digitalocean_project.prod, digitalocean_project.staging]
}
`, prodProjectName, stagingProjectName, stagingProjectName, stagingProjectName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_projects.prod", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.prod", "projects.0.name", prodProjectName),
					resource.TestCheckResourceAttr("data.digitalocean_projects.prod", "projects.0.environment", "Production"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.staging", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.staging", "projects.0.name", stagingProjectName),
					resource.TestCheckResourceAttr("data.digitalocean_projects.staging", "projects.0.environment", "Staging"),
					resource.TestCheckResourceAttr("data.digitalocean_projects.both", "projects.#", "0"),
				),
			},
		},
	})
}
