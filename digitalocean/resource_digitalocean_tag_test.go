package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDigitalOceanTag_Basic(t *testing.T) {
	var tag godo.Tag

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanTagConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanTagExists("digitalocean_tag.foobar", &tag),
					testAccCheckDigitalOceanTagAttributes(&tag),
					resource.TestCheckResourceAttr(
						"digitalocean_tag.foobar", "name", "foobar"),
				),
			},
		},
	})
}

func TestAccDigitalOceanTag_NameValidation(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "foo",
			ExpectError: false,
		},
		{
			Input:       "foo-bar",
			ExpectError: false,
		},
		{
			Input:       "foo:bar",
			ExpectError: false,
		},
		{
			Input:       "foo_bar",
			ExpectError: false,
		},
		{
			Input:       "foo-001",
			ExpectError: false,
		},
		{
			Input:       "foo/bar",
			ExpectError: true,
		},
		{
			Input:       "foo\bar",
			ExpectError: true,
		},
		{
			Input:       "foo.bar",
			ExpectError: true,
		},
		{
			Input:       "foo*",
			ExpectError: true,
		},
		{
			Input:       acctest.RandString(256),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := validateTagName(tc.Input, tc.Input)

		hasError := len(errors) > 0

		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the DigitalOcean Tag Name to trigger a validation error for '%s'", tc.Input)
		}

		if hasError && !tc.ExpectError {
			t.Fatalf("Unexpected error validating the DigitalOcean Tag Name '%s': %s", tc.Input, errors[0])
		}
	}
}

func testAccCheckDigitalOceanTagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*godo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_tag" {
			continue
		}

		// Try to find the key
		_, _, err := client.Tags.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Tag still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanTagAttributes(tag *godo.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if tag.Name != "foobar" {
			return fmt.Errorf("Bad name: %s", tag.Name)
		}

		return nil
	}
}

func testAccCheckDigitalOceanTagExists(n string, tag *godo.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*godo.Client)

		// Try to find the tag
		foundTag, _, err := client.Tags.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		*tag = *foundTag

		return nil
	}
}

var testAccCheckDigitalOceanTagConfig_basic = fmt.Sprintf(`
resource "digitalocean_tag" "foobar" {
    name = "foobar"
}`)
