package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanFunctionsNamespace_ByLabelAndID(t *testing.T) {
	label := acceptance.RandomTestName()
	resourceConfig := testAccCheckDigitalOceanFunctionsNamespaceConfig(label, "nyc1")
	dataSourceConfig := `
data "digitalocean_functions_namespace" "by_label" {
  label = digitalocean_functions_namespace.test.label
}

data "digitalocean_functions_namespace" "by_id" {
  namespace_id = digitalocean_functions_namespace.test.id
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFunctionsNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.digitalocean_functions_namespace.by_label", "namespace_id", "digitalocean_functions_namespace.test", "id"),
					resource.TestCheckResourceAttrPair("data.digitalocean_functions_namespace.by_label", "api_host", "digitalocean_functions_namespace.test", "api_host"),
					resource.TestCheckResourceAttrPair("data.digitalocean_functions_namespace.by_id", "label", "digitalocean_functions_namespace.test", "label"),
					resource.TestCheckResourceAttr("data.digitalocean_functions_namespace.by_id", "region", "nyc1"),
				),
			},
			{
				Config: resourceConfig + fmt.Sprintf(`
data "digitalocean_functions_namespace" "missing" {
  label = "%s-missing"
}
`, label),
				ExpectError: regexp.MustCompile("no Functions namespace found with label"),
			},
		},
	})
}
