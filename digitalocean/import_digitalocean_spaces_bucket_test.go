package digitalocean

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDigitalOceanBucket_importBasic(t *testing.T) {
	resourceName := "digitalocean_spaces_bucket.bucket"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanBucketConfigImport(rInt),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s,", "sfo2"),
				ImportStateVerifyIgnore: []string{"acl", "force_destroy"},
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "sfo2,nonexistent-bucket",
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}
