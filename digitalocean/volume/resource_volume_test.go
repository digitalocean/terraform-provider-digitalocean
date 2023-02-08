package volume_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanVolume_Basic(t *testing.T) {
	name := fmt.Sprintf("volume-%s", acctest.RandString(10))

	expectedURNRegEx, _ := regexp.Compile(`do:volume:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	volume := godo.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeConfig_basic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "size", "100"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "description", "peace makes plenty"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "tags.#", "2"),
					resource.TestMatchResourceAttr("digitalocean_volume.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanVolumeConfig_basic = `
resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
  tags        = ["foo", "bar"]
}`

func testAccCheckDigitalOceanVolumeExists(rn string, volume *godo.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no volume ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		got, _, err := client.Storage.GetVolume(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if got.Name != volume.Name {
			return fmt.Errorf("wrong volume found, want %q got %q", volume.Name, got.Name)
		}
		// get the computed volume details
		*volume = *got
		return nil
	}
}

func testAccCheckDigitalOceanVolumeDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_volume" {
			continue
		}

		// Try to find the volume
		_, _, err := client.Storage.GetVolume(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func TestAccDigitalOceanVolume_Droplet(t *testing.T) {
	var (
		volume  = godo.Volume{Name: fmt.Sprintf("volume-%s", acctest.RandString(10))}
		droplet godo.Droplet
	)
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanVolumeConfig_droplet(rInt, volume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					// the droplet should see an attached volume
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "volume_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVolumeConfig_droplet(rInt int, vName string) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_droplet" "foobar" {
  name               = "baz-%d"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc1"
  ipv6               = true
  private_networking = true
  volume_ids         = ["${digitalocean_volume.foobar.id}"]
}`, vName, rInt)
}

func TestAccDigitalOceanVolume_LegacyFilesystemType(t *testing.T) {
	name := fmt.Sprintf("volume-%s", acctest.RandString(10))

	volume := godo.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeConfig_legacy_filesystem_type, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "size", "100"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "description", "peace makes plenty"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "filesystem_type", "xfs"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanVolumeConfig_legacy_filesystem_type = `
resource "digitalocean_volume" "foobar" {
  region          = "nyc1"
  name            = "%s"
  size            = 100
  description     = "peace makes plenty"
  filesystem_type = "xfs"
}`

func TestAccDigitalOceanVolume_FilesystemType(t *testing.T) {
	name := fmt.Sprintf("volume-%s", acctest.RandString(10))

	volume := godo.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeConfig_filesystem_type, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "size", "100"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "description", "peace makes plenty"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "initial_filesystem_type", "xfs"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "initial_filesystem_label", "label"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "filesystem_type", "xfs"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "filesystem_label", "label"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanVolumeConfig_filesystem_type = `
resource "digitalocean_volume" "foobar" {
  region                   = "nyc1"
  name                     = "%s"
  size                     = 100
  description              = "peace makes plenty"
  initial_filesystem_type  = "xfs"
  initial_filesystem_label = "label"
}`

func TestAccDigitalOceanVolume_Resize(t *testing.T) {
	var (
		volume  = godo.Volume{Name: fmt.Sprintf("volume-%s", acctest.RandString(10))}
		droplet godo.Droplet
	)
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanVolumeConfig_resize(rInt, volume.Name, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					// the droplet should see an attached volume
					resource.TestCheckResourceAttr("digitalocean_droplet.foobar", "volume_ids.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_volume.foobar", "size", "20"),
				),
			},
			{
				Config: testAccCheckDigitalOceanVolumeConfig_resize(rInt, volume.Name, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					// the droplet should see an attached volume
					resource.TestCheckResourceAttr("digitalocean_droplet.foobar", "volume_ids.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_volume.foobar", "size", "50"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVolumeConfig_resize(rInt int, vName string, vSize int) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = %d
  description = "peace makes plenty"
}

resource "digitalocean_droplet" "foobar" {
  name               = "baz-%d"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc1"
  ipv6               = true
  private_networking = true
  volume_ids         = ["${digitalocean_volume.foobar.id}"]
}`, vName, vSize, rInt)
}

func TestAccDigitalOceanVolume_CreateFromSnapshot(t *testing.T) {
	rInt := acctest.RandInt()

	volume := godo.Volume{
		Name: fmt.Sprintf("volume-snap-%d", rInt),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanVolumeConfig_create_from_snapshot(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					// the droplet should see an attached volume
					resource.TestCheckResourceAttr("digitalocean_volume.foobar", "region", "nyc1"),
					resource.TestCheckResourceAttr("digitalocean_volume.foobar", "size", "100"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVolumeConfig_create_from_snapshot(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foo" {
  region      = "nyc1"
  name        = "volume-%d"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foo" {
  name      = "snapshot-%d"
  volume_id = "${digitalocean_volume.foo.id}"
}

resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "volume-snap-%d"
  size        = "${digitalocean_volume_snapshot.foo.min_disk_size}"
  snapshot_id = "${digitalocean_volume_snapshot.foo.id}"
}`, rInt, rInt, rInt)
}

func TestAccDigitalOceanVolume_UpdateTags(t *testing.T) {
	name := fmt.Sprintf("volume-%s", acctest.RandString(10))

	volume := godo.Volume{
		Name: name,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeConfig_basic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "tags.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeConfig_basic_tag_update, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"digitalocean_volume.foobar", "tags.#", "3"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanVolumeConfig_basic_tag_update = `
resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
  tags        = ["foo", "bar", "baz"]
}`
