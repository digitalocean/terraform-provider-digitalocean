package loadbalancer_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanLoadbalancer_Basic(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_basic(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "1"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_loadbalancer.foobar", "vpc_uuid"),
					resource.TestMatchResourceAttr(
						"digitalocean_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_backend_keepalive", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "disable_lets_encrypt_dns_records", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "http_idle_timeout_seconds", "90"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_loadbalancer.foobar", "project_id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_Updated(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_basic(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_backend_keepalive", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "disable_lets_encrypt_dns_records", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "http_idle_timeout_seconds", "90"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_loadbalancer.foobar", "project_id"),
				),
			},
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_updated(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "81",
							"entry_protocol":  "http",
							"target_port":     "81",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_backend_keepalive", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "disable_lets_encrypt_dns_records", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "http_idle_timeout_seconds", "120"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_loadbalancer.foobar", "project_id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_dropletTag(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_dropletTag(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_tag", "sample"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_minimal(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_minimal(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "sticky_sessions.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "sticky_sessions.0.type", "none"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_backend_keepalive", "false"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_loadbalancer.foobar", "project_id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_NonDefaultProject(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	lbName := acceptance.RandomTestName()
	projectName := acceptance.RandomTestName()

	projectConfig := `


resource "digitalocean_project" "test" {
  name = "%s"
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(projectConfig, projectName) + testAccCheckDigitalOceanLoadbalancerConfig_NonDefaultProject(projectName, lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.test", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.test", "name", lbName),
					resource.TestCheckResourceAttrPair(
						"digitalocean_loadbalancer.test", "project_id", "digitalocean_project.test", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.test", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.test", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.test", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.test",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.test", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.test", "healthcheck.0.port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.test", "healthcheck.0.protocol", "http"),
				),
			},
			{
				// The load balancer must be destroyed before the project which
				// discovers that asynchronously.
				Config: fmt.Sprintf(projectConfig, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(10 * time.Second)
							return nil
						},
					),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_minimalUDP(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_minimalUDP(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "udp",
							"target_port":     "80",
							"target_protocol": "udp",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "sticky_sessions.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "sticky_sessions.0.type", "none"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_backend_keepalive", "false"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_stickySessions(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_stickySessions(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "sticky_sessions.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "sticky_sessions.0.type", "cookies"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "sticky_sessions.0.cookie_name", "sessioncookie"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "sticky_sessions.0.cookie_ttl_seconds", "1800"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_sslTermination(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)
	certName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(
					certName, rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": certName,
							"tls_passthrough":  "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "redirect_http_to_https", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_sslCertByName(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := acceptance.GenerateTestCertMaterial(t)
	certName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(
					certName, rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": certName,
							"tls_passthrough":  "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "redirect_http_to_https", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "true"),
				),
			},
		},
	})
}

// Load balancers can only be resized once an hour. The initial create counts
// as a "resize" in this context. This test can not perform a resize, but it
// does ensure that the the PUT includes the expected content by checking for
// the failure.
func TestAccDigitalOceanLoadbalancer_resizeExpectedFailure(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	lbConfig := `resource "digitalocean_loadbalancer" "foobar" {
  name      = "loadbalancer-%d"
  region    = "nyc3"
  size_unit = %d

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(lbConfig, rInt, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size_unit", "1"),
				),
			},
			{
				Config:      fmt.Sprintf(lbConfig, rInt, 2),
				ExpectError: regexp.MustCompile("Load Balancer can only be resized once every hour, last resized at:"),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_multipleRules(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_multipleRules(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "443",
							"entry_protocol":  "https",
							"target_port":     "443",
							"target_protocol": "https",
							"tls_passthrough": "true",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
							"tls_passthrough": "false",
						},
					),
				),
			},
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_multipleRulesUDP(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "443",
							"entry_protocol":  "udp",
							"target_port":     "443",
							"target_protocol": "udp",
							"tls_passthrough": "false",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "444",
							"entry_protocol":  "udp",
							"target_port":     "444",
							"target_protocol": "udp",
							"tls_passthrough": "false",
						},
					),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_WithVPC(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	lbName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_WithVPC(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", lbName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_Firewall(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	lbName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_Firewall(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", lbName),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "firewall.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "firewall.0.deny.0", "cidr:1.2.0.0/16"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "firewall.0.deny.1", "ip:2.3.4.5"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "firewall.0.allow.0", "ip:1.2.3.4"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "firewall.0.allow.1", "cidr:2.3.4.0/24"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanLoadbalancerDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_loadbalancer" {
			continue
		}

		_, _, err := client.LoadBalancers.Get(context.Background(), rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for loadbalancer (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckDigitalOceanLoadbalancerExists(n string, loadbalancer *godo.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Loadbalancer ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		lb, _, err := client.LoadBalancers.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if lb.ID != rs.Primary.ID {
			return fmt.Errorf("Loabalancer not found")
		}

		*loadbalancer = *lb

		return nil
	}
}

func testAccCheckDigitalOceanLoadbalancerConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name   = "loadbalancer-%d"
  region = "nyc3"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  enable_proxy_protocol     = true
  enable_backend_keepalive  = true
  http_idle_timeout_seconds = 90

  droplet_ids = [digitalocean_droplet.foobar.id]
}`, acceptance.RandomTestName(), rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_updated(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name   = "loadbalancer-%d"
  region = "nyc3"

  forwarding_rule {
    entry_port     = 81
    entry_protocol = "http"

    target_port     = 81
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  enable_proxy_protocol            = false
  enable_backend_keepalive         = false
  disable_lets_encrypt_dns_records = true
  http_idle_timeout_seconds        = 120

  droplet_ids = [digitalocean_droplet.foobar.id, digitalocean_droplet.foo.id]
}`, acceptance.RandomTestName(), acceptance.RandomTestName(), rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_dropletTag(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "barbaz" {
  name = "sample"
}

resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  tags   = [digitalocean_tag.barbaz.id]
}

resource "digitalocean_loadbalancer" "foobar" {
  name   = "loadbalancer-%d"
  region = "nyc3"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  droplet_tag = digitalocean_tag.barbaz.name

  depends_on = [digitalocean_droplet.foobar]
}`, acceptance.RandomTestName(), rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_minimal(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name      = "loadbalancer-%d"
  region    = "nyc3"
  size_unit = 1

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  droplet_ids = [digitalocean_droplet.foobar.id]
}`, acceptance.RandomTestName(), rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_NonDefaultProject(projectName, lbName string) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "test" {
  name = "%s"
}

resource "digitalocean_loadbalancer" "test" {
  name       = "%s"
  region     = "nyc3"
  size       = "lb-small"
  project_id = digitalocean_project.test.id

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  droplet_tag = digitalocean_tag.test.name
}`, projectName, lbName)
}

func testAccCheckDigitalOceanLoadbalancerConfig_minimalUDP(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name   = "loadbalancer-%d"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "udp"

    target_port     = 80
    target_protocol = "udp"
  }

  droplet_ids = [digitalocean_droplet.foobar.id]
}`, acceptance.RandomTestName(), rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_stickySessions(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name   = "loadbalancer-%d"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  sticky_sessions {
    type               = "cookies"
    cookie_name        = "sessioncookie"
    cookie_ttl_seconds = 1800
  }

  droplet_ids = [digitalocean_droplet.foobar.id]
}`, acceptance.RandomTestName(), rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(certName string, rInt int, privateKeyMaterial, leafCert, certChain, certAttribute string) string {
	return fmt.Sprintf(`
resource "digitalocean_certificate" "foobar" {
  name              = "%s"
  private_key       = <<EOF
%s
EOF
  leaf_certificate  = <<EOF
%s
EOF
  certificate_chain = <<EOF
%s
EOF
}

resource "digitalocean_loadbalancer" "foobar" {
  name                   = "loadbalancer-%d"
  region                 = "nyc3"
  size                   = "lb-small"
  redirect_http_to_https = true
  enable_proxy_protocol  = true

  forwarding_rule {
    entry_port     = 443
    entry_protocol = "https"

    target_port     = 80
    target_protocol = "http"

    %s = digitalocean_certificate.foobar.id
  }
}`, certName, privateKeyMaterial, leafCert, certChain, rInt, certAttribute)
}

func testAccCheckDigitalOceanLoadbalancerConfig_multipleRules(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 443
    entry_protocol = "https"

    target_port     = 443
    target_protocol = "https"

    tls_passthrough = true
  }

  forwarding_rule {
    entry_port      = 80
    target_protocol = "http"
    entry_protocol  = "http"
    target_port     = 80
  }
}`, rName)
}

func testAccCheckDigitalOceanLoadbalancerConfig_multipleRulesUDP(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 443
    entry_protocol = "udp"

    target_port     = 443
    target_protocol = "udp"
  }

  forwarding_rule {
    entry_port      = 444
    target_protocol = "udp"
    entry_protocol  = "udp"
    target_port     = 444
  }
}`, rName)
}

func testAccCheckDigitalOceanLoadbalancerConfig_WithVPC(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  name     = "%s"
  size     = "s-1vcpu-1gb"
  image    = "ubuntu-22-04-x64"
  region   = "nyc3"
  vpc_uuid = digitalocean_vpc.foobar.id
}

resource "digitalocean_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  vpc_uuid    = digitalocean_vpc.foobar.id
  droplet_ids = [digitalocean_droplet.foobar.id]
}`, acceptance.RandomTestName(), acceptance.RandomTestName(), name)
}

func testAccCheckDigitalOceanLoadbalancerConfig_Firewall(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "lb-small"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  firewall {
    deny  = ["cidr:1.2.0.0/16", "ip:2.3.4.5"]
    allow = ["ip:1.2.3.4", "cidr:2.3.4.0/24"]
  }

  droplet_ids = [digitalocean_droplet.foobar.id]
}`, acceptance.RandomTestName(), name)
}
