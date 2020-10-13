package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_droplet", &resource.Sweeper{
		Name: "digitalocean_droplet",
		F:    testSweepDroplets,
	})

}

func testSweepDroplets(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*CombinedConfig).godoClient()

	opt := &godo.ListOptions{PerPage: 200}
	droplets, _, err := client.Droplets.List(context.Background(), opt)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Found %d droplets to sweep", len(droplets))

	for _, d := range droplets {
		if strings.HasPrefix(d.Name, "foo-") || strings.HasPrefix(d.Name, "bar-") || strings.HasPrefix(d.Name, "baz-") || strings.HasPrefix(d.Name, "tf-acc-test-") || strings.HasPrefix(d.Name, "foobar-") {
			log.Printf("Destroying Droplet %s", d.Name)

			if _, err := client.Droplets.Delete(context.Background(), d.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanDroplet_Basic(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "512mb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "price_hourly", "0.00744"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "price_monthly", "5"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "image", "centos-7-x64"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "user_data", HashString("foobar")),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "ipv4_address_private", ""),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "ipv6_address", ""),
					resource.TestCheckResourceAttrSet("digitalocean_droplet.foobar", "urn"),
					resource.TestCheckResourceAttrSet("digitalocean_droplet.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_WithID(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()
	slug := "centos-7-x64"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_withID(rInt, slug),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_withSSH(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("digitalocean@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_withSSH(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "512mb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "image", "centos-7-x64"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "user_data", HashString("foobar")),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_Update(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_RenameAndResize(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletRenamedAndResized(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("baz-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "1gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "30"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_ResizeWithOutDisk(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_resize_without_disk(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletResizeWithOutDisk(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "1gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "20"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_ResizeSmaller(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_resize_without_disk(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletResizeWithOutDisk(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "1gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "20"),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "512mb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "20"),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_resize(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletResizeSmaller(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "1gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "30"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_UpdateUserData(t *testing.T) {
	var afterCreate, afterUpdate godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterCreate),
					testAccCheckDigitalOceanDropletAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_userdata_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar",
						"user_data",
						HashString("foobar foobar")),
					testAccCheckDigitalOceanDropletRecreated(
						t, &afterCreate, &afterUpdate),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_UpdateTags(t *testing.T) {
	var afterCreate, afterUpdate godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterCreate),
					testAccCheckDigitalOceanDropletAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_tag_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar",
						"tags.#",
						"1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_VPCAndIpv6(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_VPCAndIpv6(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes_PrivateNetworkingIpv6(&droplet),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "ipv6", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_UpdatePrivateNetworkingIpv6(t *testing.T) {
	var afterCreate, afterUpdate godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterCreate),
					testAccCheckDigitalOceanDropletAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},
			// For "private_networking," this is now a effectively a no-opt only updating state.
			// All Droplets are assigned to a VPC by default. The API should still respond successfully.
			{
				Config: testAccCheckDigitalOceanDropletConfig_PrivateNetworkingIpv6(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "private_networking", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "ipv6", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_Monitoring(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_Monitoring(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "monitoring", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_conditionalVolumes(t *testing.T) {
	var firstDroplet godo.Droplet
	var secondDroplet godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_conditionalVolumes(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar.0", &firstDroplet),
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar.1", &secondDroplet),
					resource.TestCheckResourceAttr("digitalocean_droplet.foobar.0", "volume_ids.#", "1"),

					// This could be improved in core/HCL to make it less confusing
					// but it's the only way to use conditionals in this context for now and "it works"
					resource.TestCheckResourceAttr("digitalocean_droplet.foobar.1", "volume_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_EnableAndDisableBackups(t *testing.T) {
	var droplet godo.Droplet
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "backups", "false"),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_EnableBackups(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "backups", "true"),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_DisableBackups(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "backups", "false"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDropletDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_droplet" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Droplet
		_, _, err = client.Droplets.Get(context.Background(), id)

		// Wait

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for droplet (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckDigitalOceanDropletAttributes(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.URN() != fmt.Sprintf("do:droplet:%d", droplet.ID) {
			return fmt.Errorf("Bad URN: %s", droplet.URN())
		}

		if droplet.Image.Slug != "centos-7-x64" {
			return fmt.Errorf("Bad image_slug: %s", droplet.Image.Slug)
		}

		if droplet.Size.Slug != "512mb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.Size.Slug)
		}

		if droplet.Size.PriceHourly != 0.00744 {
			return fmt.Errorf("Bad price_hourly: %v", droplet.Size.PriceHourly)
		}

		if droplet.Size.PriceMonthly != 5.0 {
			return fmt.Errorf("Bad price_monthly: %v", droplet.Size.PriceMonthly)
		}

		if droplet.Region.Slug != "nyc3" {
			return fmt.Errorf("Bad region_slug: %s", droplet.Region.Slug)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletRenamedAndResized(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.Size.Slug != "1gb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.SizeSlug)
		}

		if droplet.Disk != 30 {
			return fmt.Errorf("Bad disk: %d", droplet.Disk)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletResizeWithOutDisk(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.Size.Slug != "1gb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.SizeSlug)
		}

		if droplet.Disk != 20 {
			return fmt.Errorf("Bad disk: %d", droplet.Disk)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletResizeSmaller(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.Size.Slug != "1gb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.SizeSlug)
		}

		if droplet.Disk != 30 {
			return fmt.Errorf("Bad disk: %d", droplet.Disk)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletAttributes_PrivateNetworkingIpv6(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.Image.Slug != "centos-7-x64" {
			return fmt.Errorf("Bad image_slug: %s", droplet.Image.Slug)
		}

		if droplet.Size.Slug != "s-1vcpu-1gb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.Size.Slug)
		}

		if droplet.Region.Slug != "nyc3" {
			return fmt.Errorf("Bad region_slug: %s", droplet.Region.Slug)
		}

		if findIPv4AddrByType(droplet, "private") == "" {
			return fmt.Errorf("No ipv4 private: %s", findIPv4AddrByType(droplet, "private"))
		}

		if findIPv4AddrByType(droplet, "public") == "" {
			return fmt.Errorf("No ipv4 public: %s", findIPv4AddrByType(droplet, "public"))
		}

		if findIPv6AddrByType(droplet, "public") == "" {
			return fmt.Errorf("No ipv6 public: %s", findIPv6AddrByType(droplet, "public"))
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "digitalocean_droplet" {
				continue
			}
			if rs.Primary.Attributes["ipv6_address"] != strings.ToLower(findIPv6AddrByType(droplet, "public")) {
				return fmt.Errorf("IPV6 Address should be lowercase")
			}

		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletExists(n string, droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Droplet ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Droplet
		retrieveDroplet, _, err := client.Droplets.Get(context.Background(), id)

		if err != nil {
			return err
		}

		if strconv.Itoa(retrieveDroplet.ID) != rs.Primary.ID {
			return fmt.Errorf("Droplet not found")
		}

		*droplet = *retrieveDroplet

		return nil
	}
}

func testAccCheckDigitalOceanDropletRecreated(t *testing.T,
	before, after *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID == after.ID {
			t.Fatalf("Expected change of droplet IDs, but both were %v", before.ID)
		}
		return nil
	}
}

func testAccCheckDigitalOceanDropletConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
  user_data = "foobar"
}`, rInt)
}

func testAccCheckDigitalOceanDropletConfig_withID(rInt int, slug string) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  slug = "%s"
}

resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "${data.digitalocean_image.foobar.id}"
  region    = "nyc3"
  user_data = "foobar"
}`, slug, rInt)
}

func testAccCheckDigitalOceanDropletConfig_withSSH(rInt int, testAccValidPublicKey string) string {
	return fmt.Sprintf(`
resource "digitalocean_ssh_key" "foobar" {
  name       = "foobar-%d"
  public_key = "%s"
}

resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
  user_data = "foobar"
  ssh_keys  = ["${digitalocean_ssh_key.foobar.id}"]
}`, rInt, testAccValidPublicKey, rInt)
}

func testAccCheckDigitalOceanDropletConfig_tag_update(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "barbaz" {
  name       = "barbaz"
}

resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
  user_data = "foobar"
  tags  = ["${digitalocean_tag.barbaz.id}"]
}
`, rInt)
}

func testAccCheckDigitalOceanDropletConfig_userdata_update(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
  user_data = "foobar foobar"
}
`, rInt)
}

func testAccCheckDigitalOceanDropletConfig_RenameAndResize(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name     = "baz-%d"
  size     = "1gb"
  image    = "centos-7-x64"
  region   = "nyc3"
}
`, rInt)
}

func testAccCheckDigitalOceanDropletConfig_resize_without_disk(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name     = "foo-%d"
  size     = "1gb"
  image    = "centos-7-x64"
  region   = "nyc3"
  user_data = "foobar"
  resize_disk = false
}
`, rInt)
}

func testAccCheckDigitalOceanDropletConfig_resize(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name     = "foo-%d"
  size     = "1gb"
  image    = "centos-7-x64"
  region   = "nyc3"
  user_data = "foobar"
  resize_disk = true
}
`, rInt)
}

func testAccCheckDigitalOceanDropletConfig_PrivateNetworkingIpv6(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name               = "foo-%d"
  size               = "s-1vcpu-1gb"
  image              = "centos-7-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}
`, rInt)
}

func testAccCheckDigitalOceanDropletConfig_VPCAndIpv6(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "foobar" {
  name        = "%s"
  region      = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  name     = "foo-%d"
  size     = "s-1vcpu-1gb"
  image    = "centos-7-x64"
  region   = "nyc3"
  ipv6     = true
  vpc_uuid = digitalocean_vpc.foobar.id
}
`, randomTestName(), rInt)
}

func testAccCheckDigitalOceanDropletConfig_Monitoring(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name       = "foo-%d"
  size       = "1gb"
  image      = "centos-7-x64"
  region     = "nyc3"
  monitoring = true
 }
 `, rInt)
}

func testAccCheckDigitalOceanDropletConfig_conditionalVolumes(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "myvol-01" {
    region      = "sfo2"
    name        = "tf-acc-test-1-%d"
    size        = 1
    description = "an example volume"
}

resource "digitalocean_volume" "myvol-02" {
    region      = "sfo2"
    name        = "tf-acc-test-2-%d"
    size        = 1
    description = "an example volume"
}

resource "digitalocean_droplet" "foobar" {
  count = 2
  name = "tf-acc-test-%d-${count.index}"
  region = "sfo2"
  image = "centos-7-x64"
  size = "512mb"
  volume_ids = ["${count.index == 0 ? digitalocean_volume.myvol-01.id : digitalocean_volume.myvol-02.id}"]
}
`, rInt, rInt, rInt)
}

func testAccCheckDigitalOceanDropletConfig_EnableBackups(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
  user_data = "foobar"
  backups   = true
}`, rInt)
}

func testAccCheckDigitalOceanDropletConfig_DisableBackups(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "foo-%d"
  size      = "512mb"
  image     = "centos-7-x64"
  region    = "nyc3"
  user_data = "foobar"
  backups   = false
}`, rInt)
}
