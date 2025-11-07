package nfs_test

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
	"time"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanNfsDataSource_Basic(t *testing.T) {
	var nfs godo.Nfs
	name := acceptance.RandomTestName("nfs")
	resourceConfig := testAccCheckDataSourceDigitalOceanNfsConfig_basic(name)

	dataSourceConfig := `
data "digitalocean_nfs" "foobar" {
  name = digitalocean_nfs.foo.name
  region = "atl1"

}`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive("digitalocean_nfs.foo"),
				),
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsExists("data.digitalocean_nfs.foobar", &nfs),
					testAccCheckDataSourceDigitalOceanNfsIsActive("data.digitalocean_nfs.foobar"),

					resource.TestCheckResourceAttr("data.digitalocean_nfs.foobar", "name", name),
					resource.TestCheckResourceAttr("data.digitalocean_nfs.foobar", "region", "atl1"),
					resource.TestCheckResourceAttr("data.digitalocean_nfs.foobar", "size", "50"),
					resource.TestCheckResourceAttr("data.digitalocean_nfs.foobar", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanNfsConfig_basic(testName string) string {
	return fmt.Sprintf(`
resource "digitalocean_nfs" "foo" {
  region = "atl1"
  name   = "%s"
  size   = 50
  vpc_id = digitalocean_vpc.foobar.id
}
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "atl1"
}
`, testName, testName)
}

func testAccCheckDataSourceDigitalOceanNfsExists(n string, nfs *godo.Nfs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NFS ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundNfs, _, err := client.Nfs.Get(context.Background(), rs.Primary.ID, "atl1")

		if err != nil {
			return err
		}

		if foundNfs.ID != rs.Primary.ID {
			return fmt.Errorf("Nfs not found")
		}

		*nfs = *foundNfs

		return nil
	}
}

func testAccCheckDataSourceDigitalOceanNfsIsActive(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No NFS ID is set")
		}
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		for i := 0; i < 10; i++ {
			nfs, _, err := client.Nfs.Get(context.Background(), rs.Primary.ID, "atl1")
			if err != nil {
				return err
			}
			if nfs.Status == "ACTIVE" {
				time.Sleep(5 * time.Second) // Extra buffer for state propagation
				return nil
			}
			time.Sleep(10 * time.Second)
		}
		return fmt.Errorf("NFS did not become ACTIVE in time")
	}
}

func testAccCheckDigitalOceanNfsDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_nfs" {
			continue
		}

		// Wait for NFS to be deleted (may take time)
		for i := 0; i < 30; i++ {
			_, resp, err := client.Nfs.Get(context.Background(), rs.Primary.ID, rs.Primary.Attributes["region"])
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil // Successfully deleted
				}
				return fmt.Errorf("Error checking if NFS is destroyed: %s", err)
			}
			time.Sleep(2 * time.Second)
		}
		return fmt.Errorf("NFS still exists")
	}
	return nil
}
