package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceDigitalOceanDomain_Basic(t *testing.T) {
	var domain godo.Domain
	domainName := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanDomainConfig_basic, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDomainExists("data.digitalocean_domain.foobar", &domain),
					testAccCheckDataSourceDigitalOceanDomainAttributes(&domain, domainName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_domain.foobar", "name", domainName),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanDomainAttributes(domain *godo.Domain, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if domain.Name != name {
			return fmt.Errorf("Bad name: %s", domain.Name)
		}

		return nil
	}
}

func testAccCheckDataSourceDigitalOceanDomainExists(n string, domain *godo.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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

const testAccCheckDataSourceDigitalOceanDomainConfig_basic = `
resource "digitalocean_domain" "foo" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

data "digitalocean_domain" "foobar" {
  name       = "${digitalocean_domain.foo.name}"
}`
