package digitalocean

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDigitalOceanProject_CreateWithDefaults(t *testing.T) {

	expectedName := generateProjectName()
	createConfig := fixtureCreateWithDefaults(expectedName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "description", ""),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "purpose", "Web Application"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "environment", "Development"),
					resource.TestCheckResourceAttrSet("digitalocean_project.myproj", "id"),
					resource.TestCheckResourceAttrSet("digitalocean_project.myproj", "owner_uuid"),
					resource.TestCheckResourceAttrSet("digitalocean_project.myproj", "owner_id"),
					resource.TestCheckResourceAttrSet("digitalocean_project.myproj", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_project.myproj", "updated_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanProject_CreateWithInitialValues(t *testing.T) {

	expectedName := generateProjectName()
	expectedDescription := "A simple project for a web app."
	expectedPurpose := "My Basic Web App"
	expectedEnvironment := "Production"

	createConfig := fixtureCreateWithInitialValues(expectedName, expectedDescription,
		expectedPurpose, expectedEnvironment)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "description", expectedDescription),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "purpose", expectedPurpose),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "environment", expectedEnvironment),
				),
			},
		},
	})
}

func TestAccDigitalOceanProject_UpdateWithInitialValues(t *testing.T) {

	expectedName := generateProjectName()
	expectedDesc := "A simple project for a web app."
	expectedPurpose := "My Basic Web App"
	expectedEnv := "Production"

	createConfig := fixtureCreateWithInitialValues(expectedName, expectedDesc,
		expectedPurpose, expectedEnv)

	expectedUpdateName := generateProjectName()
	expectedUpdateDesc := "A simple project for Beta testing."
	expectedUpdatePurpose := "MyWeb App, (Beta)"
	expectedUpdateEnv := "Staging"

	updateConfig := fixtureUpdateWithValues(expectedUpdateName, expectedUpdateDesc,
		expectedUpdatePurpose, expectedUpdateEnv)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "description", expectedDesc),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "purpose", expectedPurpose),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "environment", expectedEnv),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedUpdateName),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "description", expectedUpdateDesc),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "purpose", expectedUpdatePurpose),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "environment", expectedUpdateEnv),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "digitalocean_project" {
			continue
		}

		_, _, err := client.Projects.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Project resource still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanProjectExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundProject, _, err := client.Projects.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundProject.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		return nil
	}
}

func generateProjectName() string {
	return fmt.Sprintf("tf-proj-test-%d", acctest.RandInt())
}

func fixtureCreateWithDefaults(name string) string {
	return fmt.Sprintf(`
		resource "digitalocean_project" "myproj" {
			name = "%s"
		}`, name)
}

func fixtureUpdateWithValues(name, description, purpose, environment string) string {
	return fixtureCreateWithInitialValues(name, description, purpose, environment)
}

func fixtureCreateWithInitialValues(name, description, purpose, environment string) string {
	return fmt.Sprintf(`
		resource "digitalocean_project" "myproj" {
			name = "%s"
			description = "%s"
			purpose = "%s"
			environment = "%s"
		}`, name, description, purpose, environment)
}
