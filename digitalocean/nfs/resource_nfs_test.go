package nfs_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccDigitalOceanNfsConfigBasic = `
resource "digitalocean_nfs" "foobar" {
  region = "atl1"
  name   = "%s"
  size   = 60
  vpc_id = digitalocean_vpc.foobar.id
}
resource "digitalocean_vpc" "foobar" {
  name   = "%s-vpc"
  region = "atl1"
}
`

func testAccCheckDigitalOceanNfsConfig_resize(name string, size int) string {
	return fmt.Sprintf(`
resource "digitalocean_nfs" "foobar" {
  region = "atl1"
  name   = "%s"
  size   = %d
  vpc_id = digitalocean_vpc.foobar.id
}

resource "digitalocean_vpc" "foobar" {
  name   = "%s-vpc"
  region = "atl1"
}`, name, size, name)
}

func testAccCheckDigitalOceanNfsConfig_performanceTier(name, tier string) string {
	return fmt.Sprintf(`
resource "digitalocean_nfs" "foobar" {
  region            = "atl1"
  name              = "%s"
  size              = 60
  vpc_id            = digitalocean_vpc.foobar.id
  performance_tier  = "%s"
}

resource "digitalocean_vpc" "foobar" {
  name   = "%s-vpc"
  region = "atl1"
}`, name, tier, name)
}

func TestAccDigitalOceanNfs_Basic(t *testing.T) {
	name := acceptance.RandomTestName("nfs")
	var share godo.Nfs
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDigitalOceanNfsConfigBasic, name, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanNfsExists("digitalocean_nfs.foobar", &share),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "name", name),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "region", "atl1"),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "size", "60"),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "performance_tier", "standard"),
					resource.TestCheckResourceAttrSet("digitalocean_nfs.foobar", "host"),
					resource.TestCheckResourceAttrSet("digitalocean_nfs.foobar", "mount_path"),
				),
			},
		},
	})
}

func TestAccDigitalOceanNfs_Resize(t *testing.T) {
	resourceName := "digitalocean_nfs.foobar"
	name := acceptance.RandomTestName("nfs")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanNfsConfig_resize(name, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive(resourceName),
					testAccCheckDigitalOceanNfsSize(resourceName, 50),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "size", "50"),
				),
			},
			{
				Config: testAccCheckDigitalOceanNfsConfig_resize(name, 60),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive(resourceName),
					testAccCheckDigitalOceanNfsSize(resourceName, 60),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "size", "60"),
				),
			},
		},
	})
}

func TestAccDigitalOceanNfs_ShrinkError(t *testing.T) {
	resourceName := "digitalocean_nfs.foobar"
	name := acceptance.RandomTestName("nfs")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanNfsConfig_resize(name, 60),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive(resourceName),
					testAccCheckDigitalOceanNfsSize(resourceName, 60),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "size", "60"),
				),
			},
			{
				Config:      testAccCheckDigitalOceanNfsConfig_resize(name, 50),
				ExpectError: regexp.MustCompile(`share ` + "`size`" + ` can only be expanded and not shrunk`),
			},
		},
	})
}

func TestAccDigitalOceanNfs_PerformanceTier(t *testing.T) {
	resourceName := "digitalocean_nfs.foobar"
	name := acceptance.RandomTestName("nfs")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanNfsConfig_performanceTier(name, "standard"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive(resourceName),
					testAccCheckDigitalOceanNfsPerformanceTier(resourceName, "standard"),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "performance_tier", "standard"),
				),
			},
			{
				Config: testAccCheckDigitalOceanNfsConfig_performanceTier(name, "high"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive(resourceName),
					testAccCheckDigitalOceanNfsPerformanceTier(resourceName, "high"),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "performance_tier", "high"),
				),
			},
			{
				Config: testAccCheckDigitalOceanNfsConfig_performanceTier(name, "standard"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive(resourceName),
					testAccCheckDigitalOceanNfsPerformanceTier(resourceName, "standard"),
					resource.TestCheckResourceAttr("digitalocean_nfs.foobar", "performance_tier", "standard"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanNfsExists(rn string, share *godo.Nfs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok || rs.Primary.ID == "" {
			return fmt.Errorf("share not found")
		}
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		got, _, err := client.Nfs.Get(context.Background(), rs.Primary.ID, rs.Primary.Attributes["region"])
		if err != nil {
			return err
		}
		*share = *got
		return nil
	}
}
func testAccCheckDigitalOceanNfsSize(resourceName string, expectedSize int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		region := rs.Primary.Attributes["region"]

		// Poll until API matches (5 minutes timeout)
		for i := 0; i < 60; i++ {
			share, _, err := client.Nfs.Get(context.Background(), rs.Primary.ID, region)
			if err != nil {
				return err
			}

			if share.SizeGib == expectedSize && share.Status == "ACTIVE" {
				// Give Terraform a moment to refresh state after API is ready
				time.Sleep(2 * time.Second)
				return nil
			}

			time.Sleep(5 * time.Second)
		}

		return fmt.Errorf("NFS share did not reach expected size %d within timeout", expectedSize)
	}
}

func testAccCheckDigitalOceanNfsPerformanceTier(resourceName string, expectedTier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		region := rs.Primary.Attributes["region"]

		// Poll until API matches (5 minutes timeout)
		for i := 0; i < 60; i++ {
			share, _, err := client.Nfs.Get(context.Background(), rs.Primary.ID, region)
			if err != nil {
				return err
			}

			if share.PerformanceTier == expectedTier && share.Status == "ACTIVE" {
				// Give Terraform a moment to refresh state after API is ready
				time.Sleep(2 * time.Second)
				return nil
			}

			time.Sleep(5 * time.Second)
		}

		return fmt.Errorf("NFS share did not reach expected performance tier %s within timeout", expectedTier)
	}
}
