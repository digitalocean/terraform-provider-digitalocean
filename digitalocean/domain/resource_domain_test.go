package domain_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_domain", &resource.Sweeper{
		Name: "digitalocean_domain",
		F:    testSweepDomain,
	})

}

func testSweepDomain(region string) error {
	meta, err := acceptance.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	domains, _, err := client.Domains.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, d := range domains {
		if strings.HasPrefix(d.Name, acceptance.TestNamePrefix) {
			log.Printf("Destroying domain %s", d.Name)

			if _, err := client.Domains.Delete(context.Background(), d.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanDomain_Basic(t *testing.T) {
	var domain godo.Domain
	domainName := acceptance.RandomTestName() + ".com"

	expectedURN := fmt.Sprintf("do:domain:%s", domainName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDomainConfig_basic, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDomainExists("digitalocean_domain.foobar", &domain),
					testAccCheckDigitalOceanDomainAttributes(&domain, domainName),
					resource.TestCheckResourceAttr(
						"digitalocean_domain.foobar", "name", domainName),
					resource.TestCheckResourceAttr(
						"digitalocean_domain.foobar", "ip_address", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"digitalocean_domain.foobar", "urn", expectedURN),
				),
			},
		},
	})
}

func TestAccDigitalOceanDomain_WithoutIp(t *testing.T) {
	var domain godo.Domain
	domainName := acceptance.RandomTestName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDomainConfig_withoutIp, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDomainExists("digitalocean_domain.foobar", &domain),
					testAccCheckDigitalOceanDomainAttributes(&domain, domainName),
					resource.TestCheckResourceAttr(
						"digitalocean_domain.foobar", "name", domainName),
					resource.TestCheckNoResourceAttr(
						"digitalocean_domain.foobar", "ip_address"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDomainDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_domain" {
			continue
		}

		// Try to find the domain
		_, _, err := client.Domains.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Domain still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanDomainAttributes(domain *godo.Domain, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if domain.Name != name {
			return fmt.Errorf("Bad name: %s", domain.Name)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDomainExists(n string, domain *godo.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundDomain, _, err := client.Domains.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDomain.Name != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*domain = *foundDomain

		return nil
	}
}

const testAccCheckDigitalOceanDomainConfig_basic = `
resource "digitalocean_domain" "foobar" {
	name       = "%s"
	ip_address = "192.168.0.10"
}`

const testAccCheckDigitalOceanDomainConfig_withoutIp = `
resource "digitalocean_domain" "foobar" {
	name       = "%s"
}`
