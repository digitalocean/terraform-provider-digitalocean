package domain_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanDomain_Basic(t *testing.T) {
	var domain godo.Domain
	domainName := acceptance.RandomTestName() + ".com"
	expectedURN := fmt.Sprintf("do:domain:%s", domainName)

	resourceConfig := fmt.Sprintf(`
resource "digitalocean_domain" "foo" {
  name       = "%s"
  ip_address = "192.168.0.10"
}
`, domainName)

	dataSourceConfig := `
data "digitalocean_domain" "foobar" {
  name = digitalocean_domain.foo.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDomainExists("data.digitalocean_domain.foobar", &domain),
					testAccCheckDataSourceDigitalOceanDomainAttributes(&domain, domainName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_domain.foobar", "name", domainName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_domain.foobar", "urn", expectedURN),
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
