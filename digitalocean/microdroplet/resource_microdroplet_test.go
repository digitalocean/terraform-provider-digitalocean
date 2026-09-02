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
	config := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletOCIRef)

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
					resource.TestCheckResourceAttr(resourceName, "size.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "size.0.memory", "4096"),
					resource.TestCheckResourceAttr(resourceName, "source.0.oci_ref", testMicroDropletOCIRef),
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
	config := fmt.Sprintf(testAccMicroDropletConfig_Full, name, testMicroDropletOCIRef)

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
					resource.TestCheckResourceAttr(resourceName, "ports.#", "2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMicroDroplet_Pause(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet.foobar"

	running := fmt.Sprintf(testAccMicroDropletConfig_State, name, string(godo.MicroDropletStateRunning), testMicroDropletOCIRef)
	paused := fmt.Sprintf(testAccMicroDropletConfig_State, name, string(godo.MicroDropletStatePaused), testMicroDropletOCIRef)

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

	paused := fmt.Sprintf(testAccMicroDropletConfig_State, name, string(godo.MicroDropletStatePaused), testMicroDropletOCIRef)
	running := fmt.Sprintf(testAccMicroDropletConfig_State, name, string(godo.MicroDropletStateRunning), testMicroDropletOCIRef)

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

func TestAccDigitalOceanMicroDroplet_RecreateOnAutoPauseChange(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microdroplet.foobar"

	without := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletOCIRef)
	withFive := fmt.Sprintf(testAccMicroDropletConfig_AutoPause, name, testMicroDropletOCIRef, "5m")
	withTen := fmt.Sprintf(testAccMicroDropletConfig_AutoPause, name, testMicroDropletOCIRef, "10m")
	withDisabled := fmt.Sprintf(testAccMicroDropletConfig_AutoPauseDisabled, name, testMicroDropletOCIRef)

	var firstID, secondID, thirdID, fourthID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: without,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroDropletExists(resourceName),
					captureMicroDropletID(resourceName, &firstID),
				),
			},
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

	first := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletOCIRef)
	second := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletOCIRefAlt)

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
					resource.TestCheckResourceAttr(resourceName, "source.0.oci_ref", testMicroDropletOCIRefAlt),
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
	testMicroDropletRegion = "nyc3"

	testMicroDropletOCIRef = "docker.io/library/nginx:latest"

	testMicroDropletOCIRefAlt = "docker.io/library/httpd:latest"
)

const testAccMicroDropletConfig_Basic = `
resource "digitalocean_microdroplet" "foobar" {
  name   = "%s"
  region = "nyc3"

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "%s"
  }
}
`

const testAccMicroDropletConfig_State = `
resource "digitalocean_microdroplet" "foobar" {
  name   = "%s"
  region = "nyc3"
  state  = "%s"

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "%s"
  }
}
`

const testAccMicroDropletConfig_AutoPause = `
resource "digitalocean_microdroplet" "foobar" {
  name   = "%s"
  region = "nyc3"

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "%s"
  }

  auto_pause {
    enabled      = true
    idle_timeout = "%s"
  }
}
`

const testAccMicroDropletConfig_AutoPauseDisabled = `
resource "digitalocean_microdroplet" "foobar" {
  name   = "%s"
  region = "nyc3"

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "%s"
  }

  auto_pause {
    enabled = false
  }
}
`

const testAccMicroDropletConfig_Full = `
resource "digitalocean_microdroplet" "foobar" {
  name          = "%s"
  region        = "nyc3"
  http_port     = 8080
  http_protocol = "http"
  auto_resume   = true
  ports         = [80, 8080]

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "%s"
  }

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
