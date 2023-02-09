package project_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanProject_CreateWithDefaults(t *testing.T) {

	expectedName := generateProjectName()
	createConfig := fixtureCreateWithDefaults(expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
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
						"digitalocean_project.myproj", "environment", ""),
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

func TestAccDigitalOceanProject_CreateWithIsDefault(t *testing.T) {
	expectedName := generateProjectName()
	expectedIsDefault := "true"
	createConfig := fixtureCreateWithIsDefault(expectedName, expectedIsDefault)

	var (
		originalDefaultProject = &godo.Project{}
		client                 = &godo.Client{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)

			// Get an store original default project ID
			client = acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
			defaultProject, _, defaultProjErr := client.Projects.GetDefault(context.Background())
			if defaultProjErr != nil {
				t.Errorf("Error locating default project %s", defaultProjErr)
			}
			originalDefaultProject = defaultProject
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config:             createConfig,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					// Restore original default project. This must happen here
					// to ensure it runs even if the tests fails.
					func(*terraform.State) error {
						t.Logf("Restoring original default project: %s (%s)", originalDefaultProject.Name, originalDefaultProject.ID)
						originalDefaultProject.IsDefault = true
						updateReq := &godo.UpdateProjectRequest{
							Name:        originalDefaultProject.Name,
							Description: originalDefaultProject.Description,
							Purpose:     originalDefaultProject.Purpose,
							Environment: originalDefaultProject.Environment,
							IsDefault:   true,
						}
						_, _, err := client.Projects.Update(context.Background(), originalDefaultProject.ID, updateReq)
						if err != nil {
							return fmt.Errorf("Error restoring default project %s", err)
						}
						return nil
					},
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "description", ""),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "purpose", "Web Application"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "environment", ""),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "is_default", expectedIsDefault),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
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

func TestAccDigitalOceanProject_CreateWithDropletResource(t *testing.T) {

	expectedName := generateProjectName()
	expectedDropletName := generateDropletName()

	createConfig := fixtureCreateWithDropletResource(expectedDropletName, expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanProject_UpdateWithDropletResource(t *testing.T) {

	expectedName := generateProjectName()
	expectedDropletName := generateDropletName()

	createConfig := fixtureCreateWithDropletResource(expectedDropletName, expectedName)

	updateConfig := fixtureCreateWithDefaults(expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "1"),
				),
			},
			{
				Config: updateConfig,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "0"),
				),
			},
		},
	})
}

func TestAccDigitalOceanProject_UpdateFromDropletToSpacesResource(t *testing.T) {

	expectedName := generateProjectName()
	expectedDropletName := generateDropletName()
	expectedSpacesName := generateSpacesName()

	createConfig := fixtureCreateWithDropletResource(expectedDropletName, expectedName)

	updateConfig := fixtureCreateWithSpacesResource(expectedSpacesName, expectedName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "1"),
					resource.TestCheckResourceAttrSet("digitalocean_droplet.foobar", "urn"),
				),
			},
			{
				Config: updateConfig,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					testAccCheckDigitalOceanProjectResourceURNIsPresent("digitalocean_project.myproj", "do:spaces:"+generateSpacesName()),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", expectedName),
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "1"),
					resource.TestCheckResourceAttrSet("digitalocean_spaces_bucket.foobar", "urn"),
				),
			},
		},
	})
}

func TestAccDigitalOceanProject_WithManyResources(t *testing.T) {
	projectName := generateProjectName()
	domainBase := acceptance.RandomTestName("project")

	createConfig := fixtureCreateDomainResources(domainBase)
	updateConfig := fixtureWithManyResources(domainBase, projectName)
	destroyConfig := fixtureCreateWithDefaults(projectName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", projectName),
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "30"),
				),
			},
			{
				Config: destroyConfig,
			},
			{
				Config: destroyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanProjectExists("digitalocean_project.myproj"),
					resource.TestCheckResourceAttr(
						"digitalocean_project.myproj", "name", projectName),
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "0"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanProjectResourceURNIsPresent(resource, expectedURN string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		projectResources, _, err := client.Projects.ListResources(context.Background(), rs.Primary.ID, nil)
		if err != nil {
			return fmt.Errorf("Error Retrieving project resources to confrim.")
		}

		for _, v := range projectResources {

			if v.URN == expectedURN {
				return nil
			}

		}

		return nil
	}
}

func testAccCheckDigitalOceanProjectDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
	return acceptance.RandomTestName("project")
}

func generateDropletName() string {
	return acceptance.RandomTestName("droplet")
}

func generateSpacesName() string {
	return acceptance.RandomTestName("space")
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
  name        = "%s"
  description = "%s"
  purpose     = "%s"
  environment = "%s"
}`, name, description, purpose, environment)
}

func fixtureCreateWithDropletResource(dropletName, name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-22-04-x64"
  region    = "nyc3"
  user_data = "foobar"
}

resource "digitalocean_project" "myproj" {
  name      = "%s"
  resources = [digitalocean_droplet.foobar.urn]
}`, dropletName, name)

}

func fixtureCreateWithSpacesResource(spacesBucketName, name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "foobar" {
  name   = "%s"
  acl    = "public-read"
  region = "ams3"
}

resource "digitalocean_project" "myproj" {
  name      = "%s"
  resources = [digitalocean_spaces_bucket.foobar.urn]
}`, spacesBucketName, name)

}

func fixtureCreateDomainResources(domainBase string) string {
	return fmt.Sprintf(`
resource "digitalocean_domain" "foobar" {
  count = 30
  name  = "%s-${count.index}.com"
}`, domainBase)
}

func fixtureWithManyResources(domainBase string, name string) string {
	return fmt.Sprintf(`
resource "digitalocean_domain" "foobar" {
  count = 30
  name  = "%s-${count.index}.com"
}

resource "digitalocean_project" "myproj" {
  name      = "%s"
  resources = digitalocean_domain.foobar[*].urn
}`, domainBase, name)
}

func fixtureCreateWithIsDefault(name string, is_default string) string {
	return fmt.Sprintf(`
resource "digitalocean_project" "myproj" {
  name       = "%s"
  is_default = "%s"
}`, name, is_default)
}
