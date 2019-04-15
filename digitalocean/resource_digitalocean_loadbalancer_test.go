package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
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
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "1"),
					resource.TestMatchResourceAttr("digitalocean_loadbalancer.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_Updated(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
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
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "1"),
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
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_port", "81"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_port", "81"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.port", "22"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "healthcheck.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "droplet_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_dropletTag(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
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
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_protocol", "http"),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
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
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_protocol", "http"),
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
				),
			},
		},
	})
}

func TestAccDigitalOceanLoadbalancer_stickySessions(t *testing.T) {
	var loadbalancer godo.LoadBalancer
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
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
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_protocol", "http"),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanLoadbalancerExists("digitalocean_loadbalancer.foobar", &loadbalancer),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "name", fmt.Sprintf("loadbalancer-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_port", "443"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.entry_protocol", "https"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_port", "80"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "forwarding_rule.0.target_protocol", "http"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "redirect_http_to_https", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_loadbalancer.foobar", "enable_proxy_protocol", "true"),
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
  size      = "512mb"
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

  droplet_ids = ["${digitalocean_droplet.foobar.id}"]
}`, rInt, rInt)
}

func testAccCheckDigitalOceanLoadbalancerConfig_updated(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
}

resource "digitalocean_droplet" "foo" {
  name      = "foo-%d"
  size      = "512mb"
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

  droplet_ids = ["${digitalocean_droplet.foobar.id}","${digitalocean_droplet.foo.id}"]
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

func testAccCheckDigitalOceanLoadbalancerConfig_sslTermination(rInt int, privateKeyMaterial, leafCert, certChain string) string {
	return fmt.Sprintf(`
resource "digitalocean_certificate" "foobar" {
  name = "certificate-%d"
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
  redirect_http_to_https = true
  enable_proxy_protocol  = true

  forwarding_rule {
    entry_port      = 443
    entry_protocol  = "https"

    target_port     = 80
    target_protocol = "http"

    certificate_id  = "${digitalocean_certificate.foobar.id}"
  }
}`, rInt, privateKeyMaterial, leafCert, certChain, rInt)
}
