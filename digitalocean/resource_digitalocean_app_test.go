package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/internal/setutil"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_app", &resource.Sweeper{
		Name: "digitalocean_app",
		F:    testSweepApp,
	})

}

func testSweepApp(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*CombinedConfig).godoClient()

	opt := &godo.ListOptions{PerPage: 200}
	apps, _, err := client.Apps.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, app := range apps {
		if strings.HasPrefix(app.Spec.Name, testNamePrefix) {
			log.Printf("Destroying app %s", app.Spec.Name)

			if _, err := client.Apps.Delete(context.Background(), app.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

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
						"digitalocean_app.foobar", "spec.0.service.0.instance_count", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.instance_size_slug", "basic-xxs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.git.0.repo_clone_url",
						"https://github.com/digitalocean/sample-golang.git"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.git.0.branch", "main"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.health_check.0.http_path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.health_check.0.timeout_seconds", "10"),
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
						"digitalocean_app.foobar", "spec.0.static_site.0.catchall_document", "404.html"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.routes.0.path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.build_command", "bundle exec jekyll build -d ./public"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.output_dir", "/public"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.git.0.repo_clone_url",
						"https://github.com/digitalocean/sample-jekyll.git"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.git.0.branch", "main"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_Envs(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	oneEnv := `
      env {
        key   = "COMPONENT_FOO"
        value = "bar"
      }
`

	twoEnvs := `
      env {
        key   = "COMPONENT_FOO"
        value = "bar"
      }

      env {
        key   = "COMPONENT_FIZZ"
        value = "pop"
        scope = "BUILD_TIME"
      }
`

	oneEnvUpdated := `
      env {
        key   = "COMPONENT_FOO"
        value = "baz"
        scope = "RUN_TIME"
        type  = "GENERAL"
      }
`

	oneAppEnv := `
      env {
        key   = "APP_FOO"
        value = "bar"
      }
`

	twoAppEnvs := `
      env {
        key   = "APP_FOO"
        value = "bar"
      }

      env {
        key   = "APP_FIZZ"
        value = "pop"
        scope = "BUILD_TIME"
      }
`

	oneAppEnvUpdated := `
      env {
        key   = "APP_FOO"
        value = "baz"
        scope = "RUN_TIME"
        type  = "GENERAL"
      }
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Envs, appName, oneEnv, oneAppEnv),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.env.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.service.0.env.*",
						map[string]string{
							"key":   "COMPONENT_FOO",
							"value": "bar",
							"scope": "RUN_AND_BUILD_TIME",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.env.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.env.*",
						map[string]string{
							"key":   "APP_FOO",
							"value": "bar",
							"scope": "RUN_AND_BUILD_TIME",
						},
					),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Envs, appName, twoEnvs, twoAppEnvs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.env.#", "2"),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.service.0.env.*",
						map[string]string{
							"key":   "COMPONENT_FOO",
							"value": "bar",
							"scope": "RUN_AND_BUILD_TIME",
						},
					),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.service.0.env.*",
						map[string]string{
							"key":   "COMPONENT_FIZZ",
							"value": "pop",
							"scope": "BUILD_TIME",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.env.#", "2"),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.env.*",
						map[string]string{
							"key":   "APP_FOO",
							"value": "bar",
							"scope": "RUN_AND_BUILD_TIME",
						},
					),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.env.*",
						map[string]string{
							"key":   "APP_FIZZ",
							"value": "pop",
							"scope": "BUILD_TIME",
						},
					),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Envs, appName, oneEnvUpdated, oneAppEnvUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.env.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.service.0.env.*",
						map[string]string{
							"key":   "COMPONENT_FOO",
							"value": "baz",
							"scope": "RUN_TIME",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.env.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.env.*",
						map[string]string{
							"key":   "APP_FOO",
							"value": "baz",
							"scope": "RUN_TIME",
						},
					),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_Worker(t *testing.T) {
	var app godo.App
	appName := randomTestName()
	workerConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_worker, appName, "basic-xxs")
	upgradedWorkerConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_worker, appName, "professional-xs")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: workerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.worker.0.instance_count", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.worker.0.instance_size_slug", "basic-xxs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.worker.0.git.0.repo_clone_url",
						"https://github.com/digitalocean/sample-sleeper.git"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.worker.0.git.0.branch", "main"),
				),
			},
			{
				Config: upgradedWorkerConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.worker.0.instance_size_slug", "professional-xs"),
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
      instance_size_slug = "basic-xxs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }

      health_check {
        http_path       = "/"
        timeout_seconds = 10
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
      instance_size_slug = "basic-xxs"

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
      instance_size_slug = "basic-xxs"

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
      instance_size_slug = "basic-xxs"

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
      name              = "sample-jekyll"
      build_command     = "bundle exec jekyll build -d ./public"
      output_dir        = "/public"
      environment_slug  = "jekyll"
      catchall_document = "404.html"

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

var testAccCheckDigitalOceanAppConfig_Envs = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "basic-xxs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }

%s
    }

%s
  }
}`

var testAccCheckDigitalOceanAppConfig_worker = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    worker {
      name               = "go-worker"
      instance_count     = 1
      instance_size_slug = "%s"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-sleeper.git"
        branch         = "main"
      }
    }
  }
}`
