package microdroplet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanMicroDroplet_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet.foobar"
	config := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletImage)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "region", testMicroDropletRegion),
					resource.TestCheckResourceAttr(resourceName, "size", testMicroDropletSize),
					resource.TestCheckResourceAttr(resourceName, "state", string(godo.MicroDropletStateRunning)),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroDropletStateRunning)),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMicroDroplet_Full(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet.foobar"
	config := fmt.Sprintf(testAccMicroDropletConfig_Full, name, testMicroDropletImage)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "http_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "http_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.idle_timeout", "5m"),
					resource.TestCheckResourceAttr(resourceName, "auto_resume", "true"),
					resource.TestCheckResourceAttr(resourceName, "environment.FOO", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMicroDroplet_Pause(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet.foobar"

	running := fmt.Sprintf(testAccMicroDropletConfig_State, name, testMicroDropletImage, string(godo.MicroDropletStateRunning))
	paused := fmt.Sprintf(testAccMicroDropletConfig_State, name, testMicroDropletImage, string(godo.MicroDropletStatePaused))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: running,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", string(godo.MicroDropletStateRunning)),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroDropletStateRunning)),
				),
			},
			{
				Config: paused,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", string(godo.MicroDropletStatePaused)),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroDropletStatePaused)),
				),
			},
		},
	})
}

func TestAccDigitalOceanMicroDroplet_ResumeAfterPause(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet.foobar"

	paused := fmt.Sprintf(testAccMicroDropletConfig_State, name, testMicroDropletImage, string(godo.MicroDropletStatePaused))
	running := fmt.Sprintf(testAccMicroDropletConfig_State, name, testMicroDropletImage, string(godo.MicroDropletStateRunning))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: paused,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroDropletStatePaused)),
				),
			},
			{
				Config: running,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroDropletStateRunning)),
				),
			},
		},
	})
}

// TestAccDigitalOceanMicroDroplet_RecreateOnAutoPauseChange verifies that
// auto_pause is genuinely ForceNew at every level users can edit: adding
// the block, changing idle_timeout inside an existing block, and flipping
// enabled true -> false all recreate the resource (new ID) rather than
// silently no-op'ing on the Update path.
//
// The MicroDroplets API has no endpoint to mutate auto_pause in place, so
// recreate is the only way changes actually take effect. Without ForceNew
// on the nested `enabled` and `idle_timeout` schemas, edits inside an
// existing block produced RequiresNew=false plans that Update swallowed —
// leaving a perpetual diff for users tuning or disabling auto-pause.
func TestAccDigitalOceanMicroDroplet_RecreateOnAutoPauseChange(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet.foobar"

	without := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletImage)
	withFive := fmt.Sprintf(testAccMicroDropletConfig_AutoPause, name, testMicroDropletImage, "5m")
	withTen := fmt.Sprintf(testAccMicroDropletConfig_AutoPause, name, testMicroDropletImage, "10m")
	withDisabled := fmt.Sprintf(testAccMicroDropletConfig_AutoPauseDisabled, name, testMicroDropletImage)

	var firstID, secondID, thirdID, fourthID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			// Baseline: no auto_pause block.
			{
				Config: without,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					captureMicroDropletID(resourceName, &firstID),
				),
			},
			// Add the block (0 -> 1). Existing coverage.
			{
				Config: withFive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.idle_timeout", "5m"),
					captureMicroDropletID(resourceName, &secondID),
					assertIDChanged(&firstID, &secondID),
				),
			},
			// Edit idle_timeout inside the existing block (5m -> 10m). Guarded
			// by ForceNew on the nested `idle_timeout` field.
			{
				Config: withTen,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.idle_timeout", "10m"),
					captureMicroDropletID(resourceName, &thirdID),
					assertIDChanged(&secondID, &thirdID),
				),
			},
			// Flip enabled true -> false inside the existing block. Guarded
			// by ForceNew on the nested `enabled` field.
			{
				Config: withDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.enabled", "false"),
					captureMicroDropletID(resourceName, &fourthID),
					assertIDChanged(&thirdID, &fourthID),
				),
			},
		},
	})
}

// captureMicroDropletID snapshots the current primary ID of the named
// resource into out, so a subsequent step can assert whether the resource
// was recreated (ID changed) or updated in place (ID retained).
func captureMicroDropletID(name string, out *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		*out = rs.Primary.ID
		return nil
	}
}

// assertIDChanged fails if the two captured IDs match, i.e. Terraform kept
// the resource in place when the test expected a ForceNew recreate.
func assertIDChanged(before, after *string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if *before == "" || *after == "" {
			return fmt.Errorf("captured IDs not populated (before=%q after=%q)", *before, *after)
		}
		if *before == *after {
			return fmt.Errorf("expected MicroDroplet to be recreated (ForceNew), but ID stayed %q", *before)
		}
		return nil
	}
}

func TestAccDigitalOceanMicroDroplet_ImmutableFields(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet.foobar"

	firstImage := testMicroDropletImage
	secondImage := testMicroDropletImageAlt

	first := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, firstImage)
	second := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, secondImage)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			{Config: first},
			{
				Config: second,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image", secondImage),
				),
			},
		},
	})
}

func testAccCheckMicroDropletExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID set for MicroDroplet resource: %s", name)
		}

		m, _, err := client.MicroDroplets.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if m.ID != rs.Primary.ID {
			return fmt.Errorf("MicroDroplet not found: %s / %s", name, rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckMicroDropletDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_microdroplet" {
			continue
		}
		_, _, err := client.MicroDroplets.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("MicroDroplet %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const (
	// testMicroDropletRegion is the region used across microdroplet acceptance
	// tests. Override with DO_MICRODROPLET_REGION if the default is not
	// available on the test account.
	testMicroDropletRegion = "nyc3"

	// testMicroDropletSize is the smallest MicroDroplet size slug.
	testMicroDropletSize = "microdroplet-1"

	// testMicroDropletImage is a MicroDroplet image reference known to the
	// test account.
	testMicroDropletImage = "docker.io/library/nginx:latest"

	// testMicroDropletImageAlt is a second image reference used to prove
	// ForceNew semantics on `image` changes.
	testMicroDropletImageAlt = "docker.io/library/httpd:latest"
)

const testAccMicroDropletConfig_Basic = `
resource "digitalocean_microdroplet" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "microdroplet-1"
  image  = "%s"
}
`

const testAccMicroDropletConfig_State = `
resource "digitalocean_microdroplet" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "microdroplet-1"
  image  = "%s"
  state  = "%s"
}
`

const testAccMicroDropletConfig_AutoPause = `
resource "digitalocean_microdroplet" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "microdroplet-1"
  image  = "%s"

  auto_pause {
    enabled      = true
    idle_timeout = "%s"
  }
}
`

// testAccMicroDropletConfig_AutoPauseDisabled exercises the enabled=false
// path so tests can assert that flipping the toggle inside an existing block
// recreates the resource (guarded by ForceNew on the nested `enabled` field).
// idle_timeout is omitted because the schema marks it Optional+Computed.
const testAccMicroDropletConfig_AutoPauseDisabled = `
resource "digitalocean_microdroplet" "foobar" {
  name   = "%s"
  region = "nyc3"
  size   = "microdroplet-1"
  image  = "%s"

  auto_pause {
    enabled = false
  }
}
`

const testAccMicroDropletConfig_Full = `
resource "digitalocean_microdroplet" "foobar" {
  name          = "%s"
  region        = "nyc3"
  size          = "microdroplet-1"
  image         = "%s"
  http_port     = 8080
  http_protocol = "http"
  auto_resume   = true

  auto_pause {
    enabled      = true
    idle_timeout = "5m"
  }

  environment = {
    FOO = "bar"
  }

  tags = ["tf-acc-test-microdroplet"]
}
`
