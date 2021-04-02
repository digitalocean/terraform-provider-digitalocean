package digitalocean

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/internal/setutil"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_loadbalancer", &resource.Sweeper{
		Name: "digitalocean_loadbalancer",
		F:    testSweepLoadbalancer,
	})

}

func testSweepLoadbalancer(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*CombinedConfig).godoClient()

	opt := &godo.ListOptions{PerPage: 200}
	lbs, _, err := client.LoadBalancers.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, l := range lbs {
		if strings.HasPrefix(l.Name, "loadbalancer-") {
			log.Printf("Destroying loadbalancer %s", l.Name)

			if _, err := client.LoadBalancers.Delete(context.Background(), l.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanLoadbalancer_Basic(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
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
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_Updated(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
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
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
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
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_dropletTag(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
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
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_stickySessions(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
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
	privateKeyMaterial, leafCertMaterial, certChainMaterial := generateTestCertMaterial(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(
					"tf-acc-test-certificate-01", rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": "tf-acc-test-certificate-01",
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
	privateKeyMaterial, leafCertMaterial, certChainMaterial := generateTestCertMaterial(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(
					"tf-acc-test-certificate-02", rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					setutil.TestCheckTypeSetElemNestedAttrs(
						"digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": "tf-acc-test-certificate-02",
							"tls_passthrough":  "false",
						},
					),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
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
		name   = "loadbalancer-%d"
		region = "nyc3"
		size   = "%s"

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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(lbConfig, rInt, "lb-small"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "size", "lb-small"),
				),
			},
			{
				Config:      fmt.Sprintf(lbConfig, rInt, "lb-large"),
				ExpectError: regexp.MustCompile("Load Balancer can only be resized once every hour, last resized at:"),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_multipleRules(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
					setutil.TestCheckTypeSetElemNestedAttrs(
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
					setutil.TestCheckTypeSetElemNestedAttrs(
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
		},
	})
}

func TestAccDigitalOceanLoadbalancer_WithVPC(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	lbName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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

func testAccCheckDigitalOceanLoadbalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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
  name      = "foo-%d"
  size      = "s-1vcpu-1gb"
  image     = "centos-7-x64"
  region    = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name = "loadbalancer-%d"
  region = "nyc3"

  forwarding_rule {
    entry_port = 80
    entry_protocol = "http"

    target_port = 80
    target_protocol = "http"
  }

  healthcheck {
    port = 22
    protocol = "tcp"
  }

  enable_proxy_protocol    = true
  enable_backend_keepalive = true

  droplet_ids = [digitalocean_droplet.foobar.id]
}`, rInt, rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_updated(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "s-1vcpu-1gb"
  image     = "centos-7-x64"
  region    = "nyc3"
}

resource "digitalocean_droplet" "foo" {
  name      = "foo-%d"
  size      = "s-1vcpu-1gb"
  image     = "centos-7-x64"
  region    = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name = "loadbalancer-%d"
  region = "nyc3"

  forwarding_rule {
    entry_port = 81
    entry_protocol = "http"

    target_port = 81
    target_protocol = "http"
  }

  healthcheck {
    port = 22
    protocol = "tcp"
  }

  enable_proxy_protocol    = false
  enable_backend_keepalive = false

  droplet_ids = [digitalocean_droplet.foobar.id, digitalocean_droplet.foo.id]
}`, rInt, rInt, rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_dropletTag(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "barbaz" {
  name = "sample"
}

resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
  tags = ["${digitalocean_tag.barbaz.id}"]
}

resource "digitalocean_loadbalancer" "foobar" {
  name = "loadbalancer-%d"
  region = "nyc3"

  forwarding_rule {
    entry_port = 80
    entry_protocol = "http"

    target_port = 80
    target_protocol = "http"
  }

  healthcheck {
    port = 22
    protocol = "tcp"
  }

  droplet_tag = "${digitalocean_tag.barbaz.name}"

  depends_on = ["digitalocean_droplet.foobar"]
}`, rInt, rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_minimal(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name = "loadbalancer-%d"
  region = "nyc3"
  size = "lb-small"

  forwarding_rule {
    entry_port = 80
    entry_protocol = "http"

    target_port = 80
    target_protocol = "http"
  }

  droplet_ids = ["${digitalocean_droplet.foobar.id}"]
}`, rInt, rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_stickySessions(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
}

resource "digitalocean_loadbalancer" "foobar" {
  name = "loadbalancer-%d"
  region = "nyc3"
  size = "lb-small"

  forwarding_rule {
    entry_port = 80
    entry_protocol = "http"

    target_port = 80
    target_protocol = "http"
  }

  sticky_sessions {
	type = "cookies"
	cookie_name = "sessioncookie"
	cookie_ttl_seconds = 1800
  }

  droplet_ids = ["${digitalocean_droplet.foobar.id}"]
}`, rInt, rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(certName string, rInt int, privateKeyMaterial, leafCert, certChain, certAttribute string) string {
	return fmt.Sprintf(`
resource "digitalocean_certificate" "foobar" {
  name = "%s"
  private_key = <<EOF
%s
EOF
  leaf_certificate = <<EOF
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
    entry_port      = 443
    entry_protocol  = "https"

    target_port     = 80
    target_protocol = "http"

    %s = digitalocean_certificate.foobar.id
  }
}`, certName, privateKeyMaterial, leafCert, certChain, rInt, certAttribute)
}

func testAccCheckDigitalOceanLoadbalancerConfig_multipleRules(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_loadbalancer" "foobar" {
  name                   = "%s"
  region                 = "nyc3"
  size                 = "lb-small"

  forwarding_rule {
    entry_port      = 443
    entry_protocol  = "https"

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

func testAccCheckDigitalOceanLoadbalancerConfig_WithVPC(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "foobar" {
  name        = "%s"
  region      = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "s-1vcpu-1gb"
  image     = "centos-7-x64"
  region   = "nyc3"
  vpc_uuid = digitalocean_vpc.foobar.id
}

resource "digitalocean_loadbalancer" "foobar" {
  name = "%s"
  region = "nyc3"
  size = "lb-small"

  forwarding_rule {
    entry_port = 80
    entry_protocol = "http"

    target_port = 80
    target_protocol = "http"
  }

  vpc_uuid = digitalocean_vpc.foobar.id
  droplet_ids = [digitalocean_droplet.foobar.id]
}`, randomTestName(), randomTestName(), name)
}
