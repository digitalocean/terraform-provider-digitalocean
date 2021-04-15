package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanCustomImageFull(t *testing.T) {
	rString := randomTestName()
	name := fmt.Sprintf("digitalocean_custom_image.%s", rString)
	regions := `["nyc3"]`
	updatedString := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCustomImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanCustomImageConfig(rString, rString, regions, "Unknown"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", fmt.Sprintf("%s-name", rString)),
					resource.TestCheckResourceAttr(name, "description", fmt.Sprintf("%s-description", rString)),
					resource.TestCheckResourceAttr(name, "distribution", "Unknown"),
					resource.TestCheckResourceAttr(name, "public", "false"),
					resource.TestCheckResourceAttr(name, "regions.0", "nyc3"),
					resource.TestCheckResourceAttr(name, "status", "available"),
					resource.TestCheckResourceAttr(name, "tags.0", "flatcar"),
					resource.TestCheckResourceAttr(name, "type", "custom"),
					resource.TestCheckResourceAttr(name, "slug", ""),
					resource.TestCheckResourceAttrSet(name, "created_at"),
					resource.TestCheckResourceAttrSet(name, "image_id"),
					resource.TestCheckResourceAttrSet(name, "min_disk_size"),
					resource.TestCheckResourceAttrSet(name, "size_gigabytes"),
				),
			},
			{
				Config: testAccCheckDigitalOceanCustomImageConfig(rString, updatedString, regions, "CoreOS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", fmt.Sprintf("%s-name", updatedString)),
					resource.TestCheckResourceAttr(name, "description", fmt.Sprintf("%s-description", updatedString)),
					resource.TestCheckResourceAttr(name, "distribution", "CoreOS"),
				),
			},
		},
	})
}

func TestAccDigitalOceanCustomImageMultiRegion(t *testing.T) {
	rString := randomTestName()
	name := fmt.Sprintf("digitalocean_custom_image.%s", rString)
	regions := `["nyc3", "nyc2"]`
	regionsUpdated := `["nyc3", "nyc2", "tor1"]`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCustomImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanCustomImageConfig(rString, rString, regions, "Unknown"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", fmt.Sprintf("%s-name", rString)),
					resource.TestCheckResourceAttr(name, "regions.0", "nyc2"),
					resource.TestCheckResourceAttr(name, "regions.1", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanCustomImageConfig(rString, rString, regionsUpdated, "Unknown"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", fmt.Sprintf("%s-name", rString)),
					resource.TestCheckResourceAttr(name, "regions.0", "nyc2"),
					resource.TestCheckResourceAttr(name, "regions.1", "nyc3"),
					resource.TestCheckResourceAttr(name, "regions.2", "tor1"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanCustomImageConfig(rName string, name string, regions string, distro string) string {
	return fmt.Sprintf(`
resource "digitalocean_custom_image" "%s" {
	name = "%s-name"
	url  = "https://stable.release.flatcar-linux.net/amd64-usr/2605.7.0/flatcar_production_digitalocean_image.bin.bz2"
	regions      = %v
	description  = "%s-description"
	distribution = "%s"
	tags = [
		"flatcar"
	]
}
`, rName, name, regions, name, distro)
}

func testAccCheckDigitalOceanCustomImageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_custom_image" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Image by ID
		resp, _, err := client.Images.GetByID(context.Background(), id)

		if resp.Status != imageDeletedStatus {
			return fmt.Errorf("Image %d not destroyed", id)
		}
	}

	return nil
}
