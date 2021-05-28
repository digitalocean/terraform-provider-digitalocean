package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanLoadBalancer_Basic(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	testName := randomTestName()
	resourceConfig := testAccCheckDataSourceDigitalOceanLoadBalancerConfig(testName, "lb-small")
	dataSourceConfig := `
data "digitalocean_loadbalancer" "foobar" {
  name = digitalocean_loadbalancer.foo.name
}`

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanLoadBalancerExists("data.digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "droplet_ids.#", "2"),
					resource.TestMatchResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "enable_backend_keepalive", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanLoadBalancer_Large(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	testName := randomTestName()
	resourceConfig := testAccCheckDataSourceDigitalOceanLoadBalancerConfig(testName, "lb-large")
	dataSourceConfig := `
data "digitalocean_loadbalancer" "foobar" {
  name = digitalocean_loadbalancer.foo.name
}`

	expectedURNRegEx, _ := regexp.Compile(`do:loadbalancer:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanLoadBalancerExists("data.digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "size", "lb-large"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "droplet_ids.#", "2"),
					resource.TestMatchResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "urn", expectedURNRegEx),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_loadbalancer.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "enable_backend_keepalive", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanLoadBalancer_multipleRules(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	testName := randomTestName()
	resourceConfig := testAccCheckDigitalOceanLoadbalancerConfig_multipleRules(testName)
	dataSourceConfig := `
data "digitalocean_loadbalancer" "foobar" {
  name = digitalocean_loadbalancer.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "name", testName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "443",
							"entry_protocol":  "https",
							"target_port":     "443",
							"target_protocol": "https",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":      "80",
							"entry_protocol":  "http",
							"target_port":     "80",
							"target_protocol": "http",
						},
					),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanLoadBalancer_tlsCert(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	testName := randomTestName()
	rInt := acctest.RandInt()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := generateTestCertMaterial(t)
	resourceConfig := testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(
		testName+"-cert", rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial, "certificate_name",
	)
	dataSourceConfig := `
data "digitalocean_loadbalancer" "foobar" {
  name = digitalocean_loadbalancer.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("data.digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.digitalocean_loadbalancer.foobar",
						"forwarding_rule.*",
						map[string]string{
							"entry_port":       "443",
							"entry_protocol":   "https",
							"target_port":      "80",
							"target_protocol":  "http",
							"certificate_name": testName + "-cert",
							"tls_passthrough":  "false",
						},
					),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "size", "lb-small"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "redirect_http_to_https", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "true"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanLoadBalancerExists(n string, loadbalancer *godo.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Load Balancer ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundLoadbalancer, _, err := client.LoadBalancers.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundLoadbalancer.ID != rs.Primary.ID {
			return fmt.Errorf("Load Balancer not found")
		}

		*loadbalancer = *foundLoadbalancer

		return nil
	}
}

func testAccCheckDataSourceDigitalOceanLoadBalancerConfig(testName string, sizeSlug string) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "foo" {
  name = "web"
}

resource "digitalocean_droplet" "foo" {
  count              = 2
  image              = "ubuntu-18-04-x64"
  name               = "%s-${count.index}"
  region             = "nyc3"
  size               = "512mb"
  private_networking = true
  tags               = [digitalocean_tag.foo.id]
}

resource "digitalocean_loadbalancer" "foo" {
  name   = "%s"
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

  droplet_tag = digitalocean_tag.foo.id
  depends_on  = ["digitalocean_droplet.foo"]
}`, testName, testName, sizeSlug)
}
