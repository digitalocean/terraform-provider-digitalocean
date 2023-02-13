package droplet_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/droplet"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	defaultSize  = "s-1vcpu-1gb"
	defaultImage = "ubuntu-22-04-x64"
)

func TestAccDigitalOceanDroplet_Basic(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", defaultSize),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "price_hourly", "0.00893"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "price_monthly", "6"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "image", defaultImage),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "user_data", util.HashString("foobar")),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet.foobar", "ipv4_address_private"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "ipv6_address", ""),
					resource.TestCheckResourceAttrSet("digitalocean_droplet.foobar", "urn"),
					resource.TestCheckResourceAttrSet("digitalocean_droplet.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_WithID(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_withID(name, defaultImage),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_withSSH(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("digitalocean@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_withSSH(name, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", defaultSize),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "image", defaultImage),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "user_data", util.HashString("foobar")),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_Update(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()
	newName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_RenameAndResize(newName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletRenamedAndResized(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", newName),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "50"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_ResizeWithOutDisk(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_resize_without_disk(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletResizeWithOutDisk(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "25"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_ResizeSmaller(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
				),
			},
			// Test moving to larger plan with resize_disk = false only increases RAM, not disk.
			{
				Config: testAccCheckDigitalOceanDropletConfig_resize_without_disk(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletResizeWithOutDisk(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "25"),
				),
			},
			// Test that we can downgrade a Droplet plan as long as the disk remains the same
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "s-1vcpu-1gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "25"),
				),
			},
			// Test that resizing resize_disk = true increases the disk
			{
				Config: testAccCheckDigitalOceanDropletConfig_resize(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletResizeSmaller(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "disk", "50"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_UpdateUserData(t *testing.T) {
	var afterCreate, afterUpdate godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterCreate),
					testAccCheckDigitalOceanDropletAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_userdata_update(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar",
						"user_data",
						util.HashString("foobar foobar")),
					testAccCheckDigitalOceanDropletRecreated(
						t, &afterCreate, &afterUpdate),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_UpdateTags(t *testing.T) {
	var afterCreate, afterUpdate godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterCreate),
					testAccCheckDigitalOceanDropletAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_tag_update(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar",
						"tags.#",
						"1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_VPCAndIpv6(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_VPCAndIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes_PrivateNetworkingIpv6(&droplet),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "ipv6", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_UpdatePrivateNetworkingIpv6(t *testing.T) {
	var afterCreate, afterUpdate godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterCreate),
					testAccCheckDigitalOceanDropletAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
				),
			},
			// For "private_networking," this is now a effectively a no-opt only updating state.
			// All Droplets are assigned to a VPC by default. The API should still respond successfully.
			{
				Config: testAccCheckDigitalOceanDropletConfig_PrivateNetworkingIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "private_networking", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "ipv6", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_Monitoring(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_Monitoring(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "monitoring", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_conditionalVolumes(t *testing.T) {
	var firstDroplet godo.Droplet
	var secondDroplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_conditionalVolumes(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar.0", &firstDroplet),
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar.1", &secondDroplet),
					resource.TestCheckResourceAttr("digitalocean_droplet.foobar.0", "volume_ids.#", "1"),

					// This could be improved in core/HCL to make it less confusing
					// but it's the only way to use conditionals in this context for now and "it works"
					resource.TestCheckResourceAttr("digitalocean_droplet.foobar.1", "volume_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_EnableAndDisableBackups(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "backups", "false"),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_EnableBackups(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "backups", "true"),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_DisableBackups(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "backups", "false"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_EnableAndDisableGracefulShutdown(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					testAccCheckDigitalOceanDropletAttributes(&droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "graceful_shutdown", "false"),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_EnableGracefulShutdown(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "graceful_shutdown", "true"),
				),
			},

			{
				Config: testAccCheckDigitalOceanDropletConfig_DisableGracefulShutdown(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "graceful_shutdown", "false"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDroplet_withDropletAgentSetTrue tests that no error is returned
// from the API when creating a Droplet using an OS that supports the agent
// if the `droplet_agent` field is explicitly set to true.
func TestAccDigitalOceanDroplet_withDropletAgentSetTrue(t *testing.T) {
	var droplet godo.Droplet
	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("digitalocean@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	dropletName := acceptance.RandomTestName()
	agent := "droplet_agent = true"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_DropletAgent(keyName, publicKeyMaterial, dropletName, "ubuntu-20-04-x64", agent),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", dropletName),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "droplet_agent", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "image", "ubuntu-20-04-x64"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDroplet_withDropletAgentSetFalse tests that no error is returned
// from the API when creating a Droplet using an OS that does not support the agent
// if the `droplet_agent` field is explicitly set to false.
func TestAccDigitalOceanDroplet_withDropletAgentSetFalse(t *testing.T) {
	var droplet godo.Droplet
	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("digitalocean@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	dropletName := acceptance.RandomTestName()
	agent := "droplet_agent = false"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_DropletAgent(keyName, publicKeyMaterial, dropletName, "rancheros", agent),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", dropletName),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "droplet_agent", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "image", "rancheros"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDroplet_withDropletAgentNotSet tests that no error is returned
// from the API when creating a Droplet using an OS that does not support the agent
// if the `droplet_agent` field is not explicitly set.
func TestAccDigitalOceanDroplet_withDropletAgentNotSet(t *testing.T) {
	var droplet godo.Droplet
	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("digitalocean@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	dropletName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanDropletConfig_DropletAgent(keyName, publicKeyMaterial, dropletName, "rancheros", ""),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", dropletName),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "s-1vcpu-1gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "image", "rancheros"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "region", "nyc3"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDroplet_withDropletAgentExpectError tests that an error is returned
// from the API when creating a Droplet using an OS that does not support the agent
// if the `droplet_agent` field is explicitly set to true.
func TestAccDigitalOceanDroplet_withDropletAgentExpectError(t *testing.T) {
	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("digitalocean@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	dropletName := acceptance.RandomTestName()
	agent := "droplet_agent = true"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDigitalOceanDropletConfig_DropletAgent(keyName, publicKeyMaterial, dropletName, "rancheros", agent),
				ExpectError: regexp.MustCompile(`is not supported`),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_withTimeout(t *testing.T) {
	dropletName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  timeouts {
    create = "5s"
  }
}`, dropletName),
				ExpectError: regexp.MustCompile(`timeout while waiting for state`),
			},
		},
	})
}

func TestAccDigitalOceanDroplet_Regionless(t *testing.T) {
	var droplet godo.Droplet
	dropletName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckDigitalOceanDropletDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name  = "%s"
  size  = "s-1vcpu-1gb"
  image = "ubuntu-22-04-x64"
}`, dropletName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "name", dropletName),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "size", "s-1vcpu-1gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet.foobar", "image", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet.foobar", "region"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDropletAttributes(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.URN() != fmt.Sprintf("do:droplet:%d", droplet.ID) {
			return fmt.Errorf("Bad URN: %s", droplet.URN())
		}

		if droplet.Image.Slug != "ubuntu-22-04-x64" {
			return fmt.Errorf("Bad image_slug: %s", droplet.Image.Slug)
		}

		if droplet.Size.Slug != "s-1vcpu-1gb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.Size.Slug)
		}

		if droplet.Size.PriceHourly != 0.00893 {
			return fmt.Errorf("Bad price_hourly: %v", droplet.Size.PriceHourly)
		}

		if droplet.Size.PriceMonthly != 6.0 {
			return fmt.Errorf("Bad price_monthly: %v", droplet.Size.PriceMonthly)
		}

		if droplet.Region.Slug != "nyc3" {
			return fmt.Errorf("Bad region_slug: %s", droplet.Region.Slug)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletRenamedAndResized(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.Size.Slug != "s-1vcpu-2gb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.SizeSlug)
		}

		if droplet.Disk != 50 {
			return fmt.Errorf("Bad disk: %d", droplet.Disk)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletResizeWithOutDisk(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.Size.Slug != "s-1vcpu-2gb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.SizeSlug)
		}

		if droplet.Disk != 25 {
			return fmt.Errorf("Bad disk: %d", droplet.Disk)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletResizeSmaller(droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if droplet.Size.Slug != "s-1vcpu-2gb" {
			return fmt.Errorf("Bad size_slug: %s", droplet.SizeSlug)
		}

		if droplet.Disk != 50 {
			return fmt.Errorf("Bad disk: %d", droplet.Disk)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletAttributes_PrivateNetworkingIpv6(d *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if d.Image.Slug != "ubuntu-22-04-x64" {
			return fmt.Errorf("Bad image_slug: %s", d.Image.Slug)
		}

		if d.Size.Slug != "s-1vcpu-1gb" {
			return fmt.Errorf("Bad size_slug: %s", d.Size.Slug)
		}

		if d.Region.Slug != "nyc3" {
			return fmt.Errorf("Bad region_slug: %s", d.Region.Slug)
		}

		if droplet.FindIPv4AddrByType(d, "private") == "" {
			return fmt.Errorf("No ipv4 private: %s", droplet.FindIPv4AddrByType(d, "private"))
		}

		if droplet.FindIPv4AddrByType(d, "public") == "" {
			return fmt.Errorf("No ipv4 public: %s", droplet.FindIPv4AddrByType(d, "public"))
		}

		if droplet.FindIPv6AddrByType(d, "public") == "" {
			return fmt.Errorf("No ipv6 public: %s", droplet.FindIPv6AddrByType(d, "public"))
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "digitalocean_droplet" {
				continue
			}
			if rs.Primary.Attributes["ipv6_address"] != strings.ToLower(droplet.FindIPv6AddrByType(d, "public")) {
				return fmt.Errorf("IPV6 Address should be lowercase")
			}

		}

		return nil
	}
}

func testAccCheckDigitalOceanDropletRecreated(t *testing.T,
	before, after *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID == after.ID {
			t.Fatalf("Expected change of droplet IDs, but both were %v", before.ID)
		}
		return nil
	}
}

func testAccCheckDigitalOceanDropletConfig_withID(name string, slug string) string {
	return fmt.Sprintf(`
data "digitalocean_image" "foobar" {
  slug = "%s"
}

resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = data.digitalocean_image.foobar.id
  region    = "nyc3"
  user_data = "foobar"
}`, slug, name, defaultSize)
}

func testAccCheckDigitalOceanDropletConfig_withSSH(name string, testAccValidPublicKey string) string {
	return fmt.Sprintf(`
resource "digitalocean_ssh_key" "foobar" {
  name       = "%s-key"
  public_key = "%s"
}

resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  ssh_keys  = [digitalocean_ssh_key.foobar.id]
}`, name, testAccValidPublicKey, name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_tag_update(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "barbaz" {
  name = "barbaz"
}

resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  tags      = [digitalocean_tag.barbaz.id]
}
`, name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_userdata_update(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar foobar"
}
`, name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_RenameAndResize(newName string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-2gb"
  image  = "%s"
  region = "nyc3"
}
`, newName, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_resize_without_disk(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name        = "%s"
  size        = "s-1vcpu-2gb"
  image       = "%s"
  region      = "nyc3"
  user_data   = "foobar"
  resize_disk = false
}
`, name, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_resize(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name        = "%s"
  size        = "s-1vcpu-2gb"
  image       = "%s"
  region      = "nyc3"
  user_data   = "foobar"
  resize_disk = true
}
`, name, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_PrivateNetworkingIpv6(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name               = "%s"
  size               = "%s"
  image              = "%s"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}
`, name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_VPCAndIpv6(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "nyc3"
  ipv6     = true
  vpc_uuid = digitalocean_vpc.foobar.id
}
`, acceptance.RandomTestName(), name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_Monitoring(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name       = "%s"
  size       = "%s"
  image      = "%s"
  region     = "nyc3"
  monitoring = true
}
 `, name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_conditionalVolumes(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "myvol-01" {
  region      = "sfo3"
  name        = "%s-01"
  size        = 1
  description = "an example volume"
}

resource "digitalocean_volume" "myvol-02" {
  region      = "sfo3"
  name        = "%s-02"
  size        = 1
  description = "an example volume"
}

resource "digitalocean_droplet" "foobar" {
  count      = 2
  name       = "%s-${count.index}"
  region     = "sfo3"
  image      = "%s"
  size       = "%s"
  volume_ids = [count.index == 0 ? digitalocean_volume.myvol-01.id : digitalocean_volume.myvol-02.id]
}
`, name, name, name, defaultImage, defaultSize)
}

func testAccCheckDigitalOceanDropletConfig_EnableBackups(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  backups   = true
}`, name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_DisableBackups(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  backups   = false
}`, name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_DropletAgent(keyName, testAccValidPublicKey, dropletName, image, agent string) string {
	return fmt.Sprintf(`
resource "digitalocean_ssh_key" "foobar" {
  name       = "%s"
  public_key = "%s"
}

resource "digitalocean_droplet" "foobar" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "nyc3"
  ssh_keys = [digitalocean_ssh_key.foobar.id]
  %s
}`, keyName, testAccValidPublicKey, dropletName, defaultSize, image, agent)
}

func testAccCheckDigitalOceanDropletConfig_EnableGracefulShutdown(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name              = "%s"
  size              = "%s"
  image             = "%s"
  region            = "nyc3"
  user_data         = "foobar"
  graceful_shutdown = true
}`, name, defaultSize, defaultImage)
}

func testAccCheckDigitalOceanDropletConfig_DisableGracefulShutdown(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name              = "%s"
  size              = "%s"
  image             = "%s"
  region            = "nyc3"
  user_data         = "foobar"
  graceful_shutdown = false
}`, name, defaultSize, defaultImage)
}
