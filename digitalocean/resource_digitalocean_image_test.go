package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanImageFull(t *testing.T) {
	rString := acctest.RandString(5)
	name := "digitalocean_image.flatcar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanImageConfig(rString),
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
		},
	})
}

func testAccCheckDigitalOceanImageConfig(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_image" "flatcar" {
	name = "%s-name"
	url = "https://stable.release.flatcar-linux.net/amd64-usr/2605.7.0/flatcar_production_digitalocean_image.bin.bz2"
	regions = ["nyc3"]
	description = "%s-description"
	tags = [
		"flatcar"
	]
}
`, name, name)
}

func testAccCheckDigitalOceanImageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_image" {
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
