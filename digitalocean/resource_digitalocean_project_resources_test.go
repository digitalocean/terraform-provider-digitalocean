package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanProjectResources_Basic(t *testing.T) {
	projectName := generateProjectName()
	dropletName := generateDropletName()

	config := fmt.Sprintf(`
resource "digitalocean_project" "foo" {
  name = "%s"
}

resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
  user_data = "foobar"
}

resource "digitalocean_project_resource" "barfoo" {
  project = digitalocean_project.foo.id
  resource = digitalocean_droplet.foobar.urn
}
`, projectName, dropletName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanProjectResourcesDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testResourceInstanceState("digitalocean_project_resource.barfoo", testAccCheckDigitalOceanProjectResourcesExists),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanProjectResourcesExists(is *terraform.InstanceState) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	projectId, ok := is.Attributes["project"]
	if !ok {
		return fmt.Errorf("project attribute not set")
	}

	urn, ok := is.Attributes["resource"]
	if !ok {
		return fmt.Errorf("resource attribute not set")
	}

	projectResources, _, err := client.Projects.ListResources(context.Background(), projectId, nil)
	if err != nil {
		return fmt.Errorf("Error Retrieving project resources to confrim.")
	}

	for _, v := range projectResources {
		if v.URN == urn {
			return nil
		}
	}

	return fmt.Errorf("URN %s was not assigned to project", urn)
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
