package gradientai_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	// The custom model import API today only accepts Qwen2ForCausalLM and
	// Qwen3ForCausalLM architectures, so we default to a known-good public Qwen
	// model. Override via DO_CUSTOM_MODEL_HF_REPO if needed.
	defaultCustomModelRepo       = "Qwen/Qwen3-8B"
	defaultCustomModelAccessType = "ACCESS_TYPE_PUBLIC"
	defaultCustomModelGpuRegion  = "tor1"
)

func testCustomModelRepo() string {
	if v := os.Getenv("DO_CUSTOM_MODEL_HF_REPO"); v != "" {
		return v
	}
	return defaultCustomModelRepo
}

func testCustomModelHfToken() string {
	return os.Getenv("DO_CUSTOM_MODEL_HF_TOKEN")
}

func testCustomModelAccessType() string {
	if v := os.Getenv("DO_CUSTOM_MODEL_ACCESS_TYPE"); v != "" {
		return v
	}
	return defaultCustomModelAccessType
}

func testCustomModelGpuRegion() string {
	if v := os.Getenv("DO_CUSTOM_MODEL_GPU_REGION"); v != "" {
		return v
	}
	return defaultCustomModelGpuRegion
}

func TestAccDigitalOceanCustomModel_Basic(t *testing.T) {
	var model godo.CustomModel
	name := acceptance.RandomTestName() + "-cm"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCustomModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCustomModelExists("digitalocean_gradientai_custom_model.test", &model),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "name", name),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "source_type", "SOURCE_TYPE_HUGGINGFACE"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "source_ref.0.repo_id", testCustomModelRepo()),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "source_ref.0.access_type", testCustomModelAccessType()),
					resource.TestCheckResourceAttrSet("digitalocean_gradientai_custom_model.test", "status"),
					resource.TestCheckResourceAttrSet("digitalocean_gradientai_custom_model.test", "uuid"),
					resource.TestCheckResourceAttrSet("digitalocean_gradientai_custom_model.test", "created_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanCustomModel_UpdateMetadata(t *testing.T) {
	var model godo.CustomModel
	name := acceptance.RandomTestName() + "-cm"

	// The API does not support renaming a custom model, so we only exercise
	// description/tags updates here. See resourceDigitalOceanCustomModelUpdate.
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCustomModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCustomModelExists("digitalocean_gradientai_custom_model.test", &model),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "name", name),
				),
			},
			{
				Config: testAccCustomModelConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCustomModelExists("digitalocean_gradientai_custom_model.test", &model),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "name", name),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "description", "updated by terraform acceptance test"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanCustomModel_Import(t *testing.T) {
	name := acceptance.RandomTestName() + "-cm"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCustomModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_basic(name),
			},
			{
				ResourceName:      "digitalocean_gradientai_custom_model.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_terms_and_conditions",
					"source_ref.0.hf_token",
				},
			},
		},
	})
}

func testAccCheckDigitalOceanCustomModelDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_gradientai_custom_model" {
			continue
		}
		_, _, err := client.GradientAI.GetCustomModel(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("custom model still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDigitalOceanCustomModelExists(resourceName string, model *godo.CustomModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no custom model UUID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		found, _, err := client.GradientAI.GetCustomModel(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Uuid != rs.Primary.ID {
			return fmt.Errorf("custom model not found")
		}
		*model = *found
		return nil
	}
}

func testAccCustomModelConfig_basic(name string) string {
	hfTokenAttr := ""
	if t := testCustomModelHfToken(); t != "" {
		hfTokenAttr = fmt.Sprintf("    hf_token    = %q\n", t)
	}
	return fmt.Sprintf(`
resource "digitalocean_gradientai_custom_model" "test" {
  name                        = %q
  source_type                 = "SOURCE_TYPE_HUGGINGFACE"
  preferred_gpu_region        = %q
  accept_terms_and_conditions = true

  source_ref {
    repo_id     = %q
    access_type = %q
%s  }
}
`, name, testCustomModelGpuRegion(), testCustomModelRepo(), testCustomModelAccessType(), hfTokenAttr)
}

func testAccCustomModelConfig_updated(name string) string {
	hfTokenAttr := ""
	if t := testCustomModelHfToken(); t != "" {
		hfTokenAttr = fmt.Sprintf("    hf_token    = %q\n", t)
	}
	return fmt.Sprintf(`
resource "digitalocean_gradientai_custom_model" "test" {
  name                        = %q
  description                 = "updated by terraform acceptance test"
  source_type                 = "SOURCE_TYPE_HUGGINGFACE"
  preferred_gpu_region        = %q
  accept_terms_and_conditions = true
  tags                        = ["terraform-test", "acceptance"]

  source_ref {
    repo_id     = %q
    access_type = %q
%s  }
}
`, name, testCustomModelGpuRegion(), testCustomModelRepo(), testCustomModelAccessType(), hfTokenAttr)
}
