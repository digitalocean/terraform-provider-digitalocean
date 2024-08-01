package app_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
)

func TestAccDigitalOceanApp_Image(t *testing.T) {
	var app godo.App
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_addImage, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.image.0.registry_type", "DOCKER_HUB"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.image.0.registry", "caddy"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.image.0.repository", "caddy"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.image.0.tag", "2.2.1-alpine"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.name", "image-service"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.preserve_path_prefix", "false"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_Basic(t *testing.T) {
	var app godo.App
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_basic, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_app.foobar", "project_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "default_ingress"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "live_url"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "urn"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.alert.0.rule", "DEPLOYMENT_FAILED"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.instance_count", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.instance_size_slug", "basic-xxs"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.name", "go-service"),
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
						"digitalocean_app.foobar", "spec.0.service.0.health_check.0.port", "1234"),
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
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/go"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.name", "go-service"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.1.name", "python-service"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.1.match.0.path.0.prefix", "/python"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.1.component.0.preserve_path_prefix", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.1.component.0.name", "python-service"),
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
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.name", "go-service"),
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
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.2.name", "example-failed-job"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.job.2.kind", "FAILED_DEPLOY"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_StaticSite(t *testing.T) {
	var app godo.App
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.preserve_path_prefix", "false"),
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

func TestAccDigitalOceanApp_Egress(t *testing.T) {
	var app godo.App
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_Egress, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "default_ingress"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "dedicated_ips.0.ip"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "live_url"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.static_site.0.catchall_document", "404.html"),
					resource.TestCheckResourceAttr("digitalocean_app.foobar", "spec.0.egress.0.type", "DEDICATED_IP"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.preserve_path_prefix", "false"),
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
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
	appName := acceptance.RandomTestName()

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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
	appName := acceptance.RandomTestName()
	workerConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_worker, appName, "basic-xxs")
	upgradedWorkerConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_worker, appName, "professional-xs")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
	appName := acceptance.RandomTestName()
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/api"),
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
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_origins.0.prefix", "https://example.com"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_methods.*", "GET"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_headers.*", "X-Custom-Header"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.expose_headers.*", "Content-Encoding"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.expose_headers.*", "ETag"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.max_age", "1h"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_Domain(t *testing.T) {
	var app godo.App
	appName := acceptance.RandomTestName()

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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
	appName := acceptance.RandomTestName()

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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
	appName := acceptance.RandomTestName()

	allowedOriginExact := `
       cors {
         allow_origins {
           exact = "https://example.com"
         }
       }
`

	allowedOriginRegex := `
       cors {
         allow_origins {
           regex = "https://[0-9a-z]*.digitalocean.com"
         }
       }
`

	noAllowedOrigins := `
       cors {
         allow_methods     = ["GET", "PUT"]
         allow_headers     = ["X-Custom-Header", "Upgrade-Insecure-Requests"]
       }
`

	fullConfig := `
       cors {
         allow_origins {
           exact = "https://example.com"
         }
         allow_methods     = ["GET", "PUT"]
         allow_headers     = ["X-Custom-Header", "Upgrade-Insecure-Requests"]
         expose_headers    = ["Content-Encoding", "ETag"]
         max_age           = "1h"
         allow_credentials = true
       }
`

	allowedOriginExactConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_CORS,
		appName, allowedOriginExact,
	)
	allowedOriginRegexConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_CORS,
		appName, allowedOriginRegex,
	)
	noAllowedOriginsConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_CORS,
		appName, noAllowedOrigins,
	)
	updatedConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_CORS,
		appName, fullConfig,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: allowedOriginExactConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_origins.0.exact", "https://example.com"),
				),
			},
			{
				Config: allowedOriginRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_origins.0.regex", "https://[0-9a-z]*.digitalocean.com"),
				),
			},
			{
				Config: noAllowedOriginsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_methods.*", "GET"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_methods.*", "PUT"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_headers.*", "X-Custom-Header"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_headers.*", "Upgrade-Insecure-Requests"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_origins.0.exact", "https://example.com"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_methods.*", "GET"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_methods.*", "PUT"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_headers.*", "X-Custom-Header"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_headers.*", "Upgrade-Insecure-Requests"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.expose_headers.*", "Content-Encoding"),
					resource.TestCheckTypeSetElemAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.expose_headers.*", "ETag"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.max_age", "1h"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.cors.0.allow_credentials", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_TimeoutConfig(t *testing.T) {
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
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
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundApp, _, err := client.Apps.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		*app = *foundApp

		return nil
	}
}

func TestAccDigitalOceanApp_Features(t *testing.T) {
	var app godo.App
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_withFeatures, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.features.0", "buildpack-stack=ubuntu-18"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_nonDefaultProject(t *testing.T) {
	var app godo.App
	appName := acceptance.RandomTestName()
	projectName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_NonDefaultProject, projectName, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrPair(
						"digitalocean_project.foobar", "id", "digitalocean_app.foobar", "project_id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanApp_autoScale(t *testing.T) {
	var app godo.App
	appName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanAppConfig_autoScale, appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_app.foobar", "project_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "default_ingress"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "live_url"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "urn"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrSet("digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.instance_size_slug", "apps-d-1vcpu-0.5gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.match.0.path.0.prefix", "/"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.preserve_path_prefix", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.ingress.0.rule.0.component.0.name", "go-service"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.git.0.repo_clone_url",
						"https://github.com/digitalocean/sample-golang.git"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.git.0.branch", "main"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.autoscaling.0.min_instance_count", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.autoscaling.0.max_instance_count", "4"),
					resource.TestCheckResourceAttr(
						"digitalocean_app.foobar", "spec.0.service.0.autoscaling.0.metrics.0.cpu.0.percent", "60"),
				),
			},
		},
	})
}

var testAccCheckDigitalOceanAppConfig_basic = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
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
        port            = 1234
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
    name   = "%s"
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

var testAccCheckDigitalOceanAppConfig_withFeatures = `
resource "digitalocean_app" "foobar" {
  spec {
    name     = "%s"
    region   = "ams"
    features = ["buildpack-stack=ubuntu-18"]

    service {
      name               = "go-service-with-features"
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
    name   = "%s"
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
    }

    ingress {
      rule {
        component {
          name = "go-service"
        }
        match {
          path {
            prefix = "/go"
          }
        }
      }

      rule {
        component {
          name                 = "python-service"
          preserve_path_prefix = true
        }
        match {
          path {
            prefix = "/python"
          }
        }
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_addImage = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
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
    name   = "%s"
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

      internal_ports = [5000]
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_addDatabase = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
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

    ingress {
      rule {
        component {
          name = "go-service"
        }
        match {
          path {
            prefix = "/"
          }
        }
      }
    }

    database {
      name       = "test-db"
      engine     = "PG"
      production = false
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_StaticSite = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
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
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_Egress = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
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
    }

    egress {
      type = "DEDICATED_IP"
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_function = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
    region = "nyc"

    function {
      name       = "example"
      source_dir = "/"
      git {
        repo_clone_url = "https://github.com/digitalocean/sample-functions-nodejs-helloworld.git"
        branch         = "master"
      }
    }

    ingress {
      rule {
        component {
          name = "example"
        }

        match {
          path {
            prefix = "/api"
          }
        }

        %s
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_Envs = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
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
    name   = "%s"
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
    name   = "%s"
    region = "ams"

    job {
      name               = "example-pre-job"
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      kind               = "PRE_DEPLOY"
      run_command        = "echo 'This is a pre-deploy job.'"

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
    }

    job {
      name               = "example-post-job"
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      kind               = "POST_DEPLOY"
      run_command        = "echo 'This is a post-deploy job.'"

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
          api_key  = "test-api-key"
        }
      }
    }

    job {
      name               = "example-failed-job"
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      kind               = "FAILED_DEPLOY"
      run_command        = "echo 'This is a failed deploy job.'"

      image {
        registry_type = "DOCKER_HUB"
        registry      = "frolvlad"
        repository    = "alpine-bash"
        tag           = "latest"
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_Domains = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
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
    name   = "%s"
    region = "nyc"

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

    ingress {
      rule {
        component {
          name = "go-service"
        }

        match {
          path {
            prefix = "/"
          }
        }

        %s
      }
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_NonDefaultProject = `
resource "digitalocean_project" "foobar" {
  name = "%s"
}

resource "digitalocean_app" "foobar" {
  project_id = digitalocean_project.foobar.id
  spec {
    name   = "%s"
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
    }
  }
}`

var testAccCheckDigitalOceanAppConfig_autoScale = `
resource "digitalocean_app" "foobar" {
  spec {
    name   = "%s"
    region = "nyc"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_size_slug = "apps-d-1vcpu-0.5gb"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }

      autoscaling {
        min_instance_count = 2
        max_instance_count = 4
        metrics {
          cpu {
            percent = 60
          }
        }
      }
    }
  }
}`
