package microvm_test

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

func TestAccDigitalOceanMicroVM_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm.foobar"
	config := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRef)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "region", testMicroVMRegion),
					resource.TestCheckResourceAttr(resourceName, "size.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "size.0.memory", "4096"),
					resource.TestCheckResourceAttr(resourceName, "source.0.oci_ref", testMicroVMOCIRef),
					resource.TestCheckResourceAttr(resourceName, "state", string(godo.MicroVMStateRunning)),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroVMStateRunning)),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMicroVM_Full(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm.foobar"
	config := fmt.Sprintf(testAccMicroVMConfig_Full, name, testMicroVMOCIRef)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
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

func TestAccDigitalOceanMicroVM_Pause(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm.foobar"

	running := fmt.Sprintf(testAccMicroVMConfig_State, name, string(godo.MicroVMStateRunning), testMicroVMOCIRef)
	paused := fmt.Sprintf(testAccMicroVMConfig_State, name, string(godo.MicroVMStatePaused), testMicroVMOCIRef)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: running,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", string(godo.MicroVMStateRunning)),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroVMStateRunning)),
				),
			},
			{
				Config: paused,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", string(godo.MicroVMStatePaused)),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroVMStatePaused)),
				),
			},
		},
	})
}

func TestAccDigitalOceanMicroVM_ResumeAfterPause(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm.foobar"

	paused := fmt.Sprintf(testAccMicroVMConfig_State, name, string(godo.MicroVMStatePaused), testMicroVMOCIRef)
	running := fmt.Sprintf(testAccMicroVMConfig_State, name, string(godo.MicroVMStateRunning), testMicroVMOCIRef)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: paused,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroVMStatePaused)),
				),
			},
			{
				Config: running,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "current_state", string(godo.MicroVMStateRunning)),
				),
			},
		},
	})
}

func TestAccDigitalOceanMicroVM_RecreateOnAutoPauseChange(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm.foobar"

	without := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRef)
	withFive := fmt.Sprintf(testAccMicroVMConfig_AutoPause, name, testMicroVMOCIRef, "5m")
	withTen := fmt.Sprintf(testAccMicroVMConfig_AutoPause, name, testMicroVMOCIRef, "10m")
	withDisabled := fmt.Sprintf(testAccMicroVMConfig_AutoPauseDisabled, name, testMicroVMOCIRef)

	var firstID, secondID, thirdID, fourthID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: without,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					captureMicroVMID(resourceName, &firstID),
				),
			},
			{
				Config: withFive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.idle_timeout", "5m"),
					captureMicroVMID(resourceName, &secondID),
					assertIDChanged(&firstID, &secondID),
				),
			},
			{
				Config: withTen,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.idle_timeout", "10m"),
					captureMicroVMID(resourceName, &thirdID),
					assertIDChanged(&secondID, &thirdID),
				),
			},
			{
				Config: withDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_pause.0.enabled", "false"),
					captureMicroVMID(resourceName, &fourthID),
					assertIDChanged(&thirdID, &fourthID),
				),
			},
		},
	})
}

func captureMicroVMID(name string, out *string) resource.TestCheckFunc {
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
			return fmt.Errorf("expected MicroVM to be recreated (ForceNew), but ID stayed %q", *before)
		}
		return nil
	}
}

func TestAccDigitalOceanMicroVM_ImmutableFields(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_microvm.foobar"

	first := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRef)
	second := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRefAlt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMicroVMDestroy,
		Steps: []resource.TestStep{
			{Config: first},
			{
				Config: second,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicroVMExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "source.0.oci_ref", testMicroVMOCIRefAlt),
				),
			},
		},
	})
}

func testAccCheckMicroVMExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID set for MicroVM resource: %s", name)
		}

		m, _, err := client.MicroVMs.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if m.ID != rs.Primary.ID {
			return fmt.Errorf("MicroVM not found: %s / %s", name, rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckMicroVMDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_microvm" {
			continue
		}
		_, _, err := client.MicroVMs.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("MicroVM %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const (
	testMicroVMRegion = "nyc3"

	testMicroVMOCIRef = "docker.io/library/nginx:latest"

	testMicroVMOCIRefAlt = "docker.io/library/httpd:latest"
)

const testAccMicroVMConfig_Basic = `
resource "digitalocean_microvm" "foobar" {
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

const testAccMicroVMConfig_State = `
resource "digitalocean_microvm" "foobar" {
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

const testAccMicroVMConfig_AutoPause = `
resource "digitalocean_microvm" "foobar" {
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

const testAccMicroVMConfig_AutoPauseDisabled = `
resource "digitalocean_microvm" "foobar" {
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

const testAccMicroVMConfig_Full = `
resource "digitalocean_microvm" "foobar" {
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

  tags = ["tf-acc-test-microvm"]
}
`
