package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanRecord_Basic(t *testing.T) {
	var record godo.DomainRecord
	recordDomain := fmt.Sprintf("%s.com", randomTestName())
	recordName := randomTestName("record")
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_domain" "foo" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foo" {
  domain = digitalocean_domain.foo.name
  type   = "A"
  name   = "%s"
  value  = "192.168.0.10"
}`, recordDomain, recordName)
	dataSourceConfig := `
data "digitalocean_record" "foobar" {
  name      = digitalocean_record.foo.name
  domain    = digitalocean_domain.foo.name
}`

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
					testAccCheckDataSourceDigitalOceanRecordExists("data.digitalocean_record.foobar", &record),
					testAccCheckDataSourceDigitalOceanRecordAttributes(&record, recordName, "A"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_record.foobar", "name", recordName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_record.foobar", "type", "A"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanRecordAttributes(record *godo.DomainRecord, name string, r_type string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Name != name {
			return fmt.Errorf("Bad name: %s", record.Name)
		}

		if record.Type != r_type {
			return fmt.Errorf("Bad type: %s", record.Type)
		}

		return nil
	}
}

func testAccCheckDataSourceDigitalOceanRecordExists(n string, record *godo.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		domain := rs.Primary.Attributes["domain"]
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		foundRecord, _, err := client.Domains.Record(context.Background(), domain, id)
		if err != nil {
			return err
		}

		if foundRecord.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Record not found")
		}

		*record = *foundRecord

		return nil
	}
}
