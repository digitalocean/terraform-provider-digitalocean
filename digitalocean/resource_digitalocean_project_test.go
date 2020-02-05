package digitalocean

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func TestAccDigitalOceanProject_CreateWithDropletResource(t *testing.T) {

	expectedName := generateProjectName()
	expectedDropletName := generateDropletName()

	createConfig := fixtureCreateWithDropletResource(expectedDropletName, expectedName)

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
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "1"),
				),
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
					resource.TestCheckResourceAttr("digitalocean_project.myproj", "resources.#", "1"),
					resource.TestCheckResourceAttrSet("digitalocean_droplet.foobar", "urn"),
				),
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
	domainBase := randomTestName()

	createConfig := fixtureCreateDomainResources(domainBase)
	updateConfig := fixtureWithManyResources(domainBase, projectName)
	destroyConfig := fixtureCreateWithDefaults(projectName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanProjectDestroy,
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
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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

func generateDropletName() string {
	return fmt.Sprintf("tf-proj-test-rsrc-droplet-%d", acctest.RandInt())
}

func generateSpacesName() string {
	return fmt.Sprintf("tf-proj-test-rsrc-spaces-%d", acctest.RandInt())
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

func fixtureCreateWithDropletResource(dropletName, name string) string {
	return fmt.Sprintf(`
		resource "digitalocean_droplet" "foobar" {
		  name      = "%s"
		  size      = "512mb"
		  image     = "centos-7-x64"
		  region    = "nyc3"
		  user_data = "foobar"
		}

		resource "digitalocean_project" "myproj" {
			name = "%s"
			resources = ["${digitalocean_droplet.foobar.urn}"]
		}`, dropletName, name)

}

func fixtureCreateWithSpacesResource(spacesBucketName, name string) string {
	return fmt.Sprintf(`
		resource "digitalocean_spaces_bucket" "foobar" {
			name = "%s"
			acl = "public-read"
			region = "ams3"
		}

		resource "digitalocean_project" "myproj" {
			name = "%s"
			resources = ["${digitalocean_spaces_bucket.foobar.urn}"]
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
			name = "%s"
			resources = digitalocean_domain.foobar[*].urn
		}`, domainBase, name)
}
