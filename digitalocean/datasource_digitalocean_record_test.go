package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceDigitalOceanRecord_Basic(t *testing.T) {
	var record godo.DomainRecord
	recordName := fmt.Sprintf("foo-%s", acctest.RandString(10))
	recordDomain := fmt.Sprintf("bar-test-terraform-%s.com", acctest.RandString(10))
	recordType := "A"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanRecordConfig_basic, recordName, recordDomain, recordType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDataSourceDigitalOceanRecordAttributes(&record, recordName, recordType),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", recordName),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "type", recordType),
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
		client := testAccProvider.Meta().(*godo.Client)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		domain := rs.Primary.Attributes["domain"]
		id, err := strconv.Atoi(rs.Primary.ID)

		foundRecord, _, err := client.Domains.Record(context.Background(), domain, id)

		if err != nil {
			return err
		}

		if foundRecord.Name != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*record = *foundRecord

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanRecordConfig_basic = `
data "digitalocean_record" "foobar" {
	name      = "%s"
	domain    = "%s"
	type			= "%s"
}`
