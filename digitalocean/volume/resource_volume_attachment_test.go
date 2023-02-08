package volume_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanVolumeAttachment_Basic(t *testing.T) {
	var (
		volume  = godo.Volume{Name: fmt.Sprintf("volume-%s", acctest.RandString(10))}
		droplet godo.Droplet
	)
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanVolumeAttachmentConfig_basic(rInt, volume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &volume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanVolumeAttachmentExists("digitalocean_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "droplet_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanVolumeAttachment_Update(t *testing.T) {
	var (
		firstVolume  = godo.Volume{Name: fmt.Sprintf("volume-%s", acctest.RandString(10))}
		secondVolume = godo.Volume{Name: fmt.Sprintf("volume-s-%s", acctest.RandString(10))}
		droplet      godo.Droplet
	)
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanVolumeAttachmentConfig_basic(rInt, firstVolume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &firstVolume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanVolumeAttachmentExists("digitalocean_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "droplet_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "volume_id"),
				),
			},
			{
				Config: testAccCheckDigitalOceanVolumeAttachmentConfig_basic(rInt, secondVolume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &secondVolume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanVolumeAttachmentExists("digitalocean_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "droplet_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanVolumeAttachment_UpdateToSecondVolume(t *testing.T) {
	var (
		firstVolume  = godo.Volume{Name: fmt.Sprintf("volume-%s", acctest.RandString(10))}
		secondVolume = godo.Volume{Name: fmt.Sprintf("volume-s-%s", acctest.RandString(10))}
		droplet      godo.Droplet
	)
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanVolumeAttachmentConfig_multiple_volumes(rInt, firstVolume.Name, secondVolume.Name, "foobar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &firstVolume),
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar_second", &secondVolume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanVolumeAttachmentExists("digitalocean_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "droplet_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "volume_id"),
				),
			},
			{
				Config: testAccCheckDigitalOceanVolumeAttachmentConfig_multiple_volumes(rInt, firstVolume.Name, secondVolume.Name, "foobar_second"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &firstVolume),
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar_second", &secondVolume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanVolumeAttachmentExists("digitalocean_volume_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "droplet_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanVolumeAttachment_Multiple(t *testing.T) {
	var (
		firstVolume  = godo.Volume{Name: fmt.Sprintf("volume-%s", acctest.RandString(10))}
		secondVolume = godo.Volume{Name: fmt.Sprintf("volume-s-%s", acctest.RandString(10))}
		droplet      godo.Droplet
	)
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanVolumeAttachmentConfig_multiple(rInt, firstVolume.Name, secondVolume.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.foobar", &firstVolume),
					testAccCheckDigitalOceanVolumeExists("digitalocean_volume.barfoo", &secondVolume),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanVolumeAttachmentExists("digitalocean_volume_attachment.foobar"),
					testAccCheckDigitalOceanVolumeAttachmentExists("digitalocean_volume_attachment.barfoo"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "droplet_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.foobar", "volume_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.barfoo", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.barfoo", "droplet_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_attachment.barfoo", "volume_id"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVolumeAttachmentExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no volume ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		volumeId := rs.Primary.Attributes["volume_id"]
		dropletId, err := strconv.Atoi(rs.Primary.Attributes["droplet_id"])
		if err != nil {
			return err
		}

		got, _, err := client.Storage.GetVolume(context.Background(), volumeId)
		if err != nil {
			return err
		}

		if got.DropletIDs == nil || len(got.DropletIDs) == 0 || got.DropletIDs[0] != dropletId {
			return fmt.Errorf("wrong volume attachment found for volume %s, got %q wanted %q", volumeId, got.DropletIDs[0], dropletId)
		}

		return nil
	}
}

func testAccCheckDigitalOceanVolumeAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_volume_attachment" {
			continue
		}
	}

	return nil
}

func testAccCheckDigitalOceanVolumeAttachmentConfig_basic(rInt int, vName string) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "digitalocean_droplet" "foobar" {
  name   = "baz-%d"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc1"
}

resource "digitalocean_volume_attachment" "foobar" {
  droplet_id = digitalocean_droplet.foobar.id
  volume_id  = digitalocean_volume.foobar.id
}`, vName, rInt)
}

func testAccCheckDigitalOceanVolumeAttachmentConfig_multiple(rInt int, vName, vSecondName string) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "digitalocean_volume" "barfoo" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "digitalocean_droplet" "foobar" {
  name   = "baz-%d"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc1"
}

resource "digitalocean_volume_attachment" "foobar" {
  droplet_id = digitalocean_droplet.foobar.id
  volume_id  = digitalocean_volume.foobar.id
}

resource "digitalocean_volume_attachment" "barfoo" {
  droplet_id = digitalocean_droplet.foobar.id
  volume_id  = digitalocean_volume.barfoo.id
}`, vName, vSecondName, rInt)
}

func testAccCheckDigitalOceanVolumeAttachmentConfig_multiple_volumes(rInt int, vName, vSecondName, activeVolume string) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "digitalocean_volume" "foobar_second" {
  region      = "nyc1"
  name        = "%s"
  size        = 5
  description = "peace makes plenty"
}

resource "digitalocean_droplet" "foobar" {
  name   = "baz-%d"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc1"
}

resource "digitalocean_volume_attachment" "foobar" {
  droplet_id = digitalocean_droplet.foobar.id
  volume_id  = digitalocean_volume.%s.id
}`, vName, vSecondName, rInt, activeVolume)
}
