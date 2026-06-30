package nfs_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/nfs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanNfsAccessPoint_Basic(t *testing.T) {
	resourceName := "digitalocean_nfs_access_point.foobar"
	nfsName := acceptance.RandomTestName()
	vpcName := acceptance.RandomTestName()
	apName := acceptance.RandomTestName("ap")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsAccessPointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanNfsAccessPointConfig_basic(nfsName, vpcName, apName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive("digitalocean_nfs.foo"),
					testAccCheckDigitalOceanNfsAccessPointExists(resourceName),
					testAccCheckDigitalOceanNfsAccessPointIsActive(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", apName),
					resource.TestCheckResourceAttr(resourceName, "path", "/data"),
					resource.TestCheckResourceAttrPair(resourceName, "share_id", "digitalocean_nfs.foo", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "digitalocean_vpc.foobar", "id"),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.squash_config", "ROOT_SQUASH"),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.protocols.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.anonuid", "65534"),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.anongid", "65534"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDigitalOceanNfsAccessPointConfig_basic(nfsName, vpcName, apName string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_nfs" "foo" {
  region           = "atl1"
  name             = "%s"
  size             = 50
  vpc_id           = digitalocean_vpc.foobar.id
  performance_tier = "standard"
}

resource "digitalocean_nfs_access_point" "foobar" {
  name     = "%s"
  share_id = digitalocean_nfs.foo.id
  path     = "/data"
  vpc_id   = digitalocean_vpc.foobar.id

  access_policy {
    protocols     = ["NFS4"]
    squash_config = "ROOT_SQUASH"
  }
}
`, vpcName, nfsName, apName)
}

func testAccCheckDigitalOceanNfsAccessPointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No NFS access point ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		accessPoint, _, err := nfs.GetNfsAccessPoint(context.Background(), client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if accessPoint.ID != rs.Primary.ID {
			return fmt.Errorf("NFS access point not found")
		}

		return nil
	}
}

func testAccCheckDigitalOceanNfsAccessPointIsActive(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No NFS access point ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		for i := 0; i < 20; i++ {
			accessPoint, _, err := nfs.GetNfsAccessPoint(context.Background(), client, rs.Primary.ID)
			if err != nil {
				return err
			}
			if accessPoint.Status == godo.NfsAccessPointActive {
				time.Sleep(5 * time.Second)
				return nil
			}
			time.Sleep(10 * time.Second)
		}
		return fmt.Errorf("NFS access point did not become ACTIVE in time")
	}
}

func testAccCheckDigitalOceanNfsAccessPointDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_nfs_access_point" {
			continue
		}

		for i := 0; i < 30; i++ {
			accessPoint, resp, err := nfs.GetNfsAccessPoint(context.Background(), client, rs.Primary.ID)
			if err != nil && resp != nil && resp.StatusCode == 404 {
				return nil
			}
			if err == nil && (accessPoint == nil || accessPoint.Status == godo.NfsAccessPointDeleted) {
				return nil
			}
			time.Sleep(2 * time.Second)
		}
		return fmt.Errorf("NFS access point still exists")
	}
	return nil
}
