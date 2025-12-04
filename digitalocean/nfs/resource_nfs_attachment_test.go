package nfs_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanNfsAttachment_Basic(t *testing.T) {
	var share godo.Nfs
	shareName := acceptance.RandomTestName("nfs")
	vpcName := acceptance.RandomTestName("vpc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanNfsAttachmentConfig_basic(shareName, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanNfsExists("digitalocean_nfs.foobar", &share),
					testAccCheckDigitalOceanNfsAttachmentExists("digitalocean_nfs_attachment.foobar"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_nfs_attachment.foobar", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_nfs_attachment.foobar", "vpc_id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_nfs_attachment.foobar", "share_id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanNfsAttachment_MoveToNewVPC(t *testing.T) {
	var share godo.Nfs
	shareName := acceptance.RandomTestName("nfs")
	initialVpcName := acceptance.RandomTestName("vpc-initial")
	attachVpcName := acceptance.RandomTestName("vpc-attach")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanNfsAttachmentConfig_moveToNewVPC_initial(shareName, initialVpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanNfsExists("digitalocean_nfs.foobar", &share),
					resource.TestCheckResourceAttrPair(
						"digitalocean_nfs.foobar", "vpc_id", "digitalocean_vpc.initial", "id"),
				),
			},
			{
				Config: testAccCheckDigitalOceanNfsAttachmentConfig_moveToNewVPC_attach(shareName, initialVpcName, attachVpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanNfsExists("digitalocean_nfs.foobar", &share),
					testAccCheckDigitalOceanNfsAttachmentExists("digitalocean_nfs_attachment.foobar"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_nfs_attachment.foobar", "vpc_id", "digitalocean_vpc.attach", "id"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_nfs_attachment.foobar", "share_id", "digitalocean_nfs.foobar", "id"),
				),
			},
		},
	})
}

func TestAccDigitalOceanNfsAttachment_VPCAlreadyHasShare(t *testing.T) {
	shareName := acceptance.RandomTestName("nfs")
	secondShareName := acceptance.RandomTestName("nfs-second")
	firstVpcName := acceptance.RandomTestName("vpc-first")
	secondVpcName := acceptance.RandomTestName("vpc-second")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDigitalOceanNfsAttachmentConfig_vpcAlreadyHasShare(shareName, secondShareName, firstVpcName, secondVpcName),
				ExpectError: regexp.MustCompile("Error attaching share|already attached|conflict|only one share"),
			},
		},
	})
}

func testAccCheckDigitalOceanNfsAttachmentExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no NFS attachment ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		shareId := rs.Primary.Attributes["share_id"]
		vpcId := rs.Primary.Attributes["vpc_id"]
		region := rs.Primary.Attributes["region"]

		share, _, err := client.Nfs.Get(context.Background(), shareId, region)
		if err != nil {
			return err
		}

		if len(share.VpcIDs) == 0 || share.VpcIDs[0] != vpcId {
			return fmt.Errorf("wrong NFS attachment found for share %s, got %q wanted %q", shareId, share.VpcIDs, vpcId)
		}

		return nil
	}
}

func testAccCheckDigitalOceanNfsAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_nfs_attachment" {
			continue
		}
	}

	return nil
}
func testAccCheckDigitalOceanNfsAttachmentConfig_basic(shareName, vpcName string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_nfs" "foobar" {
  name   = "%s"
  region = "atl1"
  size   = 50
  vpc_id = digitalocean_vpc.foobar.id
}

resource "digitalocean_nfs_attachment" "foobar" {
  vpc_id   = digitalocean_vpc.foobar.id
  share_id = digitalocean_nfs.foobar.id
  region   = "atl1"
}`, vpcName, shareName)
}

func testAccCheckDigitalOceanNfsAttachmentConfig_moveToNewVPC_initial(shareName, initialVpcName string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "initial" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_nfs" "foobar" {
  name   = "%s"
  region = "atl1"
  size   = 50
  vpc_id = digitalocean_vpc.initial.id
}`, initialVpcName, shareName)
}

func testAccCheckDigitalOceanNfsAttachmentConfig_moveToNewVPC_attach(shareName, initialVpcName, attachVpcName string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "initial" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_vpc" "attach" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_nfs" "foobar" {
  name   = "%s"
  region = "atl1"
  size   = 50
  vpc_id = digitalocean_vpc.initial.id

  lifecycle {
    ignore_changes = [vpc_id]
  }
}

resource "digitalocean_nfs_attachment" "foobar" {
  vpc_id   = digitalocean_vpc.attach.id
  share_id = digitalocean_nfs.foobar.id
  region   = "atl1"
}`, initialVpcName, attachVpcName, shareName)
}

func testAccCheckDigitalOceanNfsAttachmentConfig_vpcAlreadyHasShare(firstShareName, secondShareName, firstVpcName, secondVpcName string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "first" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_vpc" "second" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_nfs" "first" {
  name   = "%s"
  region = "atl1"
  size   = 50
  vpc_id = digitalocean_vpc.first.id
}

resource "digitalocean_nfs" "second" {
  name   = "%s"
  region = "atl1"
  size   = 50
  vpc_id = digitalocean_vpc.second.id
}

resource "digitalocean_nfs_attachment" "first" {
  vpc_id   = digitalocean_vpc.first.id
  share_id = digitalocean_nfs.first.id
  region   = "atl1"
}

resource "digitalocean_nfs_attachment" "second" {
  vpc_id   = digitalocean_vpc.first.id
  share_id = digitalocean_nfs.second.id
  region   = "atl1"
}`, firstVpcName, secondVpcName, firstShareName, secondShareName)
}
