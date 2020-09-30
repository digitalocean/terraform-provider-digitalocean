package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDigitalOceanApp_Basic(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_basic, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "default_ingress"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "live_url"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_addService, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.name", "go-service"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/go"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.1.name", "python-service"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.1.routes.0.path", "/python"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_addDatabase, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.database.0.name", "test-db"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.database.0.engine", "PG"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_StaticSite(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_StaticSite, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "default_ingress"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "live_url"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.routes.0.path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.build_command", "bundle exec jekyll build -d ./public"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.output_dir", "/public"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_app" {
			continue
		}

		_, _, err := client.Apps.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Container Registry still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanAppExists(n string, app *godo.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundApp, _, err := client.Apps.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		*app = *foundApp

		return nil
	}
}

var testAccCheckDigitalOceanAppConfig_basic = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "professional-xs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_addService = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "professional-xs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }

      routes {
        path = "/go"
      }
    }

    service {
      name               = "python-service"
      environment_slug   = "python"
      instance_count     = 1
      instance_size_slug = "professional-xs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-python.git"
        branch         = "main"
      }

      routes {
        path = "/python"
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_addDatabase = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "professional-xs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }

      routes {
        path = "/"
      }
    }

    database {
      name = "test-db"
      engine = "PG"
      production = false
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_StaticSite = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    static_site {
      name             = "sample-jekyll"
      build_command    = "bundle exec jekyll build -d ./public"
	  output_dir       = "/public"
      environment_slug = "jekyll"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-jekyll.git"
        branch         = "main"
      }

      routes {
        path = "/"
      }
    }
  }
}`
