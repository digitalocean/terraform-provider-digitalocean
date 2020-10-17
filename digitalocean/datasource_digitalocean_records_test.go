package digitalocean

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceDigitalOceanRecords_Basic(t *testing.T) {
	name1 := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resourcesConfig := fmt.Sprintf(`
resource "digitalocean_domain" "foo" {
  name = "%s"
}

resource "digitalocean_record" "mail" {
  name = "mail"
  domain = digitalocean_domain.foo.name
  type = "MX"
  priority = 10
  value = "mail.example.com."
}

resource "digitalocean_record" "www" {
  name = "www"
  domain = digitalocean_domain.foo.name
  type = "A"
  value = "192.168.1.1"
}
`, name1)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_records" "result" {
  domain = "%s"
  filter {
    key = "type"
    values = ["A"]
  }
}
`, name1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_records.result", "records.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_records.result", "records.0.domain", name1),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.id", "digitalocean_record.www", "id"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.name", "digitalocean_record.www", "name"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.type", "digitalocean_record.www", "type"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.value", "digitalocean_record.www", "value"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.priority", "digitalocean_record.www", "priority"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.port", "digitalocean_record.www", "port"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.ttl", "digitalocean_record.www", "ttl"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.weight", "digitalocean_record.www", "weight"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.flags", "digitalocean_record.www", "flags"),
					resource.TestCheckResourceAttrPair("data.digitalocean_records.result", "records.0.tag", "digitalocean_record.www", "tag"),
				),
			},
		},
	})
}
