package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanProjectResources_Basic(t *testing.T) {
	projectName := generateProjectName()
	dropletName := generateDropletName()

	baseConfig := fmt.Sprintf(`
resource "digitalocean_project" "foo" {
  name = "%s"
}

resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "s-1vcpu-1gb"
  image     = "centos-7-x64"
  region    = "nyc3"
  user_data = "foobar"
}
`, projectName, dropletName)

	projectResourcesConfigEmpty := `
resource "digitalocean_project_resources" "barfoo" {
  project = digitalocean_project.foo.id
  resources = []
}
`

	projectResourcesConfigWithDroplet := `
resource "digitalocean_project_resources" "barfoo" {
  project = digitalocean_project.foo.id
  resources = [digitalocean_droplet.foobar.urn]
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectResourcesDestroy,
		Steps: []resource.TestStep{
			{
				Config: baseConfig + projectResourcesConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("digitalocean_project_resources.barfoo", "project"),
					resource.TestCheckResourceAttr("digitalocean_project_resources.barfoo", "resources.#", "0"),
					testProjectMembershipCount("digitalocean_project_resources.barfoo", 0),
				),
			},
			{
				// Add a resource to the digitalocean_project_resources.
				Config: baseConfig + projectResourcesConfigWithDroplet,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("digitalocean_project_resources.barfoo", "project"),
					resource.TestCheckResourceAttr("digitalocean_project_resources.barfoo", "resources.#", "1"),
					testProjectMembershipCount("digitalocean_project_resources.barfoo", 1),
				),
			},
			{
				// Remove the resource that was added.
				Config: baseConfig + projectResourcesConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("digitalocean_project_resources.barfoo", "project"),
					resource.TestCheckResourceAttr("digitalocean_project_resources.barfoo", "resources.#", "0"),
					testProjectMembershipCount("digitalocean_project_resources.barfoo", 0),
				),
			},
		},
	})
}

func testProjectMembershipCount(name string, expectedCount int) resource.TestCheckFunc {
	return testResourceInstanceState(name, func(is *terraform.InstanceState) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		projectId, ok := is.Attributes["project"]
		if !ok {
			return fmt.Errorf("project attribute not set")
		}

		resources, err := loadResourceURNs(client, projectId)
		if err != nil {
			return fmt.Errorf("Error retrieving project resources: %s", err)
		}

		actualCount := len(*resources)

		if actualCount != expectedCount {
			return fmt.Errorf("project membership count mismatch: expected=%d, actual=%d",
				expectedCount, actualCount)
		}

		return nil
	})
}

func testAccCheckDigitalOceanProjectResourcesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "digitalocean_project":
			_, _, err := client.Projects.Get(context.Background(), rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Project resource still exists")
			}

		case "digitalocean_droplet":
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return err
			}

			_, _, err = client.Droplets.Get(context.Background(), id)
			if err == nil {
				return fmt.Errorf("Droplet resource still exists")
			}
		}
	}

	return nil
}
