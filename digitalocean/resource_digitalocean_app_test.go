package digitalocean

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

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

func TestAccDigitalOceanApp_Image(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_addImage, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.image.0.registry_type", "DOCKER_HUB"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.image.0.registry", "caddy"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.image.0.repository", "caddy"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.image.0.tag", "2.2.1-alpine"),
				),
			},
		},
	})
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
						"digitalocean_app.foobar", "spec.0.alert.0.rule", "DEPLOYMENT_FAILED"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.instance_count", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.instance_size_slug", "basic-xxs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.git.0.repo_clone_url",
						"https://github.com/digitalocean/sample-golang.git"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.git.0.branch", "main"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.health_check.0.http_path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.health_check.0.timeout_seconds", "10"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.alert.0.value", "75"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.alert.0.operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.alert.0.window", "TEN_MINUTES"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.alert.0.rule", "CPU_UTILIZATION"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.log_destination.0.name", "ServiceLogs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.log_destination.0.papertrail.0.endpoint", "syslog+tls://example.com:12345"),
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
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.1.name", "python-service"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.1.routes.0.path", "/python"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.1.routes.0.preserve_path_prefix", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.alert.0.value", "85"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.alert.0.operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.alert.0.window", "FIVE_MINUTES"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.alert.0.rule", "CPU_UTILIZATION"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.log_destination.0.name", "ServiceLogs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.log_destination.0.papertrail.0.endpoint", "syslog+tls://example.com:12345"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_addDatabase, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.routes.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.database.0.name", "test-db"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.database.0.engine", "PG"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.log_destination.0.name", "ServiceLogs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.log_destination.0.papertrail.0.endpoint", "syslog+tls://example.com:12345"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_Job(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_addJob, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.0.name", "example-pre-job"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.0.kind", "PRE_DEPLOY"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.0.run_command", "echo 'This is a pre-deploy job.'"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.1.name", "example-post-job"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.1.kind", "POST_DEPLOY"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.1.run_command", "echo 'This is a post-deploy job.'"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.1.log_destination.0.name", "JobLogs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.1.log_destination.0.datadog.0.endpoint", "https://example.com"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.1.log_destination.0.datadog.0.api_key", "test-api-key"),
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
						"digitalocean_app.foobar", "spec.0.static_site.0.routes.0.preserve_path_prefix", "false"),
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

func TestAccDigitalOceanApp_InternalPort(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_addInternalPort, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.instance_count", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.instance_size_slug", "basic-xxs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.git.0.repo_clone_url",
						"https://github.com/digitalocean/sample-golang.git"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.git.0.branch", "main"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.internal_ports.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.internal_ports.0", "5000"),
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
					resource.TestCheckTypeSetElemNestedAttrs(
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
					resource.TestCheckTypeSetElemNestedAttrs(
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
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.service.0.env.*",
						map[string]string{
							"key":   "COMPONENT_FOO",
							"value": "bar",
							"scope": "RUN_AND_BUILD_TIME",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
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
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_app.foobar",
						"spec.0.env.*",
						map[string]string{
							"key":   "APP_FOO",
							"value": "bar",
							"scope": "RUN_AND_BUILD_TIME",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
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
					resource.TestCheckTypeSetElemNestedAttrs(
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
					resource.TestCheckTypeSetElemNestedAttrs(
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
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.worker.0.log_destination.0.name", "WorkerLogs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.worker.0.log_destination.0.logtail.0.token", "test-api-token"),
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

func TestAccDigitalOceanApp_Function(t *testing.T) {
	var app godo.App
	appName := randomTestName()
	fnConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_function, appName, "")

	corsConfig := `
       cors {
         allow_origins {
           prefix = "https://example.com"
         }
         allow_methods     = ["GET"]
         allow_headers     = ["X-Custom-Header"]
         expose_headers    = ["Content-Encoding", "ETag"]
         max_age           = "1h"
       }
`
	updatedFnConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_function, appName, corsConfig)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.function.0.source_dir", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.function.0.routes.0.path", "/api"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.function.0.git.0.repo_clone_url",
						"https://github.com/digitalocean/sample-functions-nodejs-helloworld.git"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.function.0.git.0.branch", "master"),
				),
			},
			{
				Config: updatedFnConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.function.0.cors.0.allow_origins.0.prefix", "https://example.com"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.function.0.cors.0.allow_methods.*", "GET"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.function.0.cors.0.allow_headers.*", "X-Custom-Header"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.function.0.cors.0.expose_headers.*", "Content-Encoding"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.function.0.cors.0.expose_headers.*", "ETag"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.function.0.cors.0.max_age", "1h"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_Domain(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	domain := fmt.Sprintf(`
       domain {
         name     = "%s.com"
         wildcard = true
       }
`, appName)

	updatedDomain := fmt.Sprintf(`
       domain {
         name     = "%s.net"
         wildcard = true
       }
`, appName)

	domainsConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Domains, appName, domain)
	updatedDomainConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Domains, appName, updatedDomain)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: domainsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.domain.0.name", appName+".com"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.domain.0.wildcard", "true"),
				),
			},
			{
				Config: updatedDomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.domain.0.name", appName+".net"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.domain.0.wildcard", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_DomainsDeprecation(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	deprecatedStyleDomain := fmt.Sprintf(`
       domains = ["%s.com"]
`, appName)

	updatedDeprecatedStyleDomain := fmt.Sprintf(`
       domains = ["%s.net"]
`, appName)

	newStyleDomain := fmt.Sprintf(`
       domain {
         name     = "%s.com"
         wildcard = true
       }
`, appName)

	domainsConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Domains, appName, deprecatedStyleDomain)
	updateDomainsConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Domains, appName, updatedDeprecatedStyleDomain)
	replaceDomainsConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Domains, appName, newStyleDomain)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: domainsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.domains.0", appName+".com"),
				),
			},
			{
				Config: updateDomainsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.domains.0", appName+".net"),
				),
			},
			{
				Config: replaceDomainsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.domain.0.name", appName+".com"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.domain.0.wildcard", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_CORS(t *testing.T) {
	var app godo.App
	appName := randomTestName()

	allowedOrginExact := `
       cors {
         allow_origins {
           exact = "https://example.com"
         }
       }
`

	allowedOrginPrefix := `
       cors {
         allow_origins {
           prefix = "https://example.com"
         }
       }
`

	allowedOrginRegex := `
       cors {
         allow_origins {
           regex = "https://[0-9a-z]*.digitalocean.com"
         }
       }
`

	fullConfig := `
       cors {
         allow_origins {
           prefix = "https://example.com"
         }
         allow_methods     = ["GET", "PUT"]
         allow_headers     = ["X-Custom-Header", "Upgrade-Insecure-Requests"]
         expose_headers    = ["Content-Encoding", "ETag"]
         max_age           = "1h"
         allow_credentials = true
       }
`

	allowedOrginExactConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_CORS,
		appName, allowedOrginExact,
	)
	allowedOrginPrefixConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_CORS,
		appName, allowedOrginPrefix,
	)
	allowedOrginRegexConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_CORS,
		appName, allowedOrginRegex,
	)
	updatedConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_CORS,
		appName, fullConfig,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: allowedOrginExactConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_origins.0.exact", "https://example.com"),
				),
			},
			{
				Config: allowedOrginPrefixConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_origins.0.prefix", "https://example.com"),
				),
			},
			{
				Config: allowedOrginRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_origins.0.regex", "https://[0-9a-z]*.digitalocean.com"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_origins.0.prefix", "https://example.com"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_methods.*", "GET"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_methods.*", "PUT"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_headers.*", "X-Custom-Header"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_headers.*", "Upgrade-Insecure-Requests"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.expose_headers.*", "Content-Encoding"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.expose_headers.*", "ETag"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.max_age", "1h"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.cors.0.allow_credentials", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_TimeoutConfig(t *testing.T) {
	appName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanAppConfig_withTimeout, appName),
				ExpectError: regexp.MustCompile("timeout waiting for app"),
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

    alert {
      rule = "DEPLOYMENT_FAILED"
    }

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

      alert {
        value    = 75
        operator = "GREATER_THAN"
        window   = "TEN_MINUTES"
        rule     = "CPU_UTILIZATION"
      }

      log_destination {
        name = "ServiceLogs"
        papertrail {
          endpoint = "syslog+tls://example.com:12345"
        }
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_withTimeout = `
resource "digitalocean_app" "foobar" {
  timeouts {
    create = "10s"
  }

  spec {
    name = "%s"
    region = "ams"

    service {
      name               = "go-service-with-timeout"
      instance_size_slug = "basic-xxs"

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

    alert {
      rule = "DEPLOYMENT_FAILED"
    }

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

      alert {
        value    = 85
        operator = "GREATER_THAN"
        window   = "FIVE_MINUTES"
        rule     = "CPU_UTILIZATION"
      }

      log_destination {
        name = "ServiceLogs"
        papertrail {
          endpoint = "syslog+tls://example.com:12345"
        }
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
        path                 = "/python"
        preserve_path_prefix = true
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_addImage = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    service {
      name               = "image-service"
      instance_count     = 1
      instance_size_slug = "basic-xxs"

      image {
        registry_type = "DOCKER_HUB"
        registry      = "caddy"
        repository    = "caddy"
        tag           = "2.2.1-alpine"
      }

      http_port = 80
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_addInternalPort = `
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

	  internal_ports = [ 5000 ]
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_addDatabase = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    alert {
      rule = "DEPLOYMENT_FAILED"
    }

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

      alert {
        value    = 85
        operator = "GREATER_THAN"
        window   = "FIVE_MINUTES"
        rule     = "CPU_UTILIZATION"
      }

      log_destination {
        name = "ServiceLogs"
        papertrail {
          endpoint = "syslog+tls://example.com:12345"
        }
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

      routes {
        path = "/foo"
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_function = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "nyc"

    function {
      name              = "example"
	  source_dir        = "/"
      git {
        repo_clone_url = "https://github.com/digitalocean/sample-functions-nodejs-helloworld.git"
        branch         = "master"
      }
      routes {
        path = "/api"
      }

%s

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

      log_destination {
        name = "WorkerLogs"
        logtail {
          token = "test-api-token"
        }
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_addJob = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    job {
      name               = "example-pre-job"
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      kind = "PRE_DEPLOY"
      run_command = "echo 'This is a pre-deploy job.'"

      image {
        registry_type = "DOCKER_HUB"
        registry      = "frolvlad"
        repository    = "alpine-bash"
        tag           = "latest"
      }
    }

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

    job {
      name               = "example-post-job"
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      kind = "POST_DEPLOY"
      run_command = "echo 'This is a post-deploy job.'"

      image {
        registry_type = "DOCKER_HUB"
        registry      = "frolvlad"
        repository    = "alpine-bash"
        tag           = "latest"
      }

      log_destination {
        name = "JobLogs"
        datadog {
          endpoint = "https://example.com"
          api_key = "test-api-key"
        }
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_Domains = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "ams"

    %s

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "basic-xxs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_CORS = `
resource "digitalocean_app" "foobar" {
  spec {
    name = "%s"
    region = "nyc"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "basic-xxs"

      %s

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }
    }
  }
}`
