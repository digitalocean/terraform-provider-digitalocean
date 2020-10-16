package digitalocean

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanApp_Basic(t *testing.T) {
	var app godo.App
	appName := randomTestName()
	appCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_basic, appName)
	appDataConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanAppConfig, appCreateConfig)

	updatedAppCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanAppConfig_addService, appName)
	updatedAppDataConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanAppConfig, updatedAppCreateConfig)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: appCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAppExists("digitalocean_app.foobar", &app),
				),
			},
			{
				Config: appDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.name", appName),
					resource.TestCheckResourceAttrPair("digitalocean_app.foobar", "default_ingress",
						"data.digitalocean_app.foobar", "default_ingress"),
					resource.TestCheckResourceAttrPair("digitalocean_app.foobar", "live_url",
						"data.digitalocean_app.foobar", "live_url"),
					resource.TestCheckResourceAttrPair("digitalocean_app.foobar", "active_deployment_id",
						"data.digitalocean_app.foobar", "active_deployment_id"),
					resource.TestCheckResourceAttrPair("digitalocean_app.foobar", "updated_at",
						"data.digitalocean_app.foobar", "updated_at"),
					resource.TestCheckResourceAttrPair("digitalocean_app.foobar", "created_at",
						"data.digitalocean_app.foobar", "created_at"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.0.instance_count", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.0.instance_size_slug", "basic-xxs"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.0.git.0.repo_clone_url",
						"https://github.com/digitalocean/sample-golang.git"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.0.git.0.branch", "main"),
				),
			},
			{
				Config: updatedAppDataConfig,
			},
			{
				Config: updatedAppDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.0.name", "go-service"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.0.routes.0.path", "/go"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.1.name", "python-service"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_app.foobar", "spec.0.service.1.routes.0.path", "/python"),
				),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanAppConfig = `
%s

data "digitalocean_app" "foobar" {
  app_id = digitalocean_app.foobar.id
}`
