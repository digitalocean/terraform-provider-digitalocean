package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceDigitalOceanImages_Basic(t *testing.T) {
	config := `
data "digitalocean_images" "ubuntu" {
  filter {
    key = "distribution"
    values = ["Ubuntu"]
  }
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_images.ubuntu", "images.#"),
					testResourceInstanceState("data.digitalocean_images.ubuntu", testAccDataSourceDigitalOceanImages_VerifyImageData),
				),
			},
		},
	})
}

func testAccDataSourceDigitalOceanImages_VerifyImageData(is *terraform.InstanceState) error {
	ns, ok := is.Attributes["images.#"]
	if !ok {
		return fmt.Errorf("images.# attribute not found")
	}

	n, err := strconv.Atoi(ns)
	if err != nil {
		return fmt.Errorf("images.# attribute was not convertible to an integer: %s", err)
	}

	if n == 0 {
		return fmt.Errorf("Expected to find Ubuntu images")
	}

	// Verify the first image to ensure that it matches what the API returned.
	slug, ok := is.Attributes["images.0.slug"]
	if !ok {
		return fmt.Errorf("images.0.slug attribute not found")
	}

	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	image, _, err := client.Images.GetBySlug(context.Background(), slug)
	if err != nil {
		return err
	}
	log.Printf("image=%+v", image)

	if image.Name != is.Attributes["images.0.name"] {
		return fmt.Errorf("mismatch on `name`: expected=%s, actual=%s",
			image.Name, is.Attributes["images.0.name"])
	}

	if image.Type != is.Attributes["images.0.type"] {
		return fmt.Errorf("mismatch on `type`: expected=%s, actual=%s",
			image.Type, is.Attributes["images.0.type"])
	}

	if image.Description != is.Attributes["images.0.description"] {
		return fmt.Errorf("mismatch on `description`: expected=%s, actual=%s",
			image.Description, is.Attributes["images.0.description"])
	}

	return nil
}
