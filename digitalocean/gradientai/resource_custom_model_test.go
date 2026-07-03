package gradientai_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/gradientai"
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

	// HF imports only support description/tags updates. The Spaces-only
	// metadata fields (license, parameters, input/output modalities) are
	// exercised separately by TestAccDigitalOceanCustomModel_UpdateSpacesMetadata.
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

// TestAccDigitalOceanCustomModel_UpdateSpacesMetadata exercises the Spaces-only
// editable metadata fields (license, parameters, input_modalities,
// output_modalities) which the API only honors for SOURCE_TYPE_SPACES_BUCKET
// imports.
//
// The test is skipped unless the caller exports a Spaces bucket / region /
// prefix that already contains importable model artifacts:
//
//	DO_CUSTOM_MODEL_SPACES_BUCKET
//	DO_CUSTOM_MODEL_SPACES_REGION
//	DO_CUSTOM_MODEL_SPACES_PREFIX
func TestAccDigitalOceanCustomModel_UpdateSpacesMetadata(t *testing.T) {
	bucket := os.Getenv("DO_CUSTOM_MODEL_SPACES_BUCKET")
	region := os.Getenv("DO_CUSTOM_MODEL_SPACES_REGION")
	prefix := os.Getenv("DO_CUSTOM_MODEL_SPACES_PREFIX")
	if bucket == "" || region == "" || prefix == "" {
		t.Skip("DO_CUSTOM_MODEL_SPACES_BUCKET / DO_CUSTOM_MODEL_SPACES_REGION / DO_CUSTOM_MODEL_SPACES_PREFIX must be set for the Spaces metadata update test")
	}

	var model godo.CustomModel
	name := acceptance.RandomTestName() + "-cm"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCustomModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_spacesBasic(name, bucket, region, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCustomModelExists("digitalocean_gradientai_custom_model.test", &model),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "name", name),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "source_type", "SOURCE_TYPE_SPACES_BUCKET"),
				),
			},
			{
				Config: testAccCustomModelConfig_spacesUpdated(name, bucket, region, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCustomModelExists("digitalocean_gradientai_custom_model.test", &model),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "name", name),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "description", "updated by terraform acceptance test"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "license", "apache-2.0"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "parameters", "8B"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "input_modalities.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "input_modalities.0", "text"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "output_modalities.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "output_modalities.0", "text"),
				),
			},
		},
	})
}

// TestAccDigitalOceanCustomModel_RejectSpacesOnlyFieldsOnHF asserts that
// setting any of the Spaces-only metadata fields on a non-Spaces source is
// rejected at plan time by CustomizeDiff, so users do not get into a permanent
// plan-diff state.
func TestAccDigitalOceanCustomModel_RejectSpacesOnlyFieldsOnHF(t *testing.T) {
	name := acceptance.RandomTestName() + "-cm"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCustomModelConfig_hfWithSpacesOnlyFields(name),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`can only be set when source_type = "SOURCE_TYPE_SPACES_BUCKET"`),
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
	if t := testCustomModelHfToken(); t != "" {
		return fmt.Sprintf(`
resource "digitalocean_gradientai_custom_model" "test" {
  name                        = %q
  source_type                 = "SOURCE_TYPE_HUGGINGFACE"
  preferred_gpu_region        = %q
  accept_terms_and_conditions = true

  source_ref {
    repo_id     = %q
    access_type = %q
    hf_token    = %q
  }
}
`, name, testCustomModelGpuRegion(), testCustomModelRepo(), testCustomModelAccessType(), t)
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
  }
}
`, name, testCustomModelGpuRegion(), testCustomModelRepo(), testCustomModelAccessType())
}

func testAccCustomModelConfig_updated(name string) string {
	if t := testCustomModelHfToken(); t != "" {
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
    hf_token    = %q
  }
}
`, name, testCustomModelGpuRegion(), testCustomModelRepo(), testCustomModelAccessType(), t)
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
  }
}
`, name, testCustomModelGpuRegion(), testCustomModelRepo(), testCustomModelAccessType())
}

func testAccCustomModelConfig_spacesBasic(name, bucket, region, prefix string) string {
	return fmt.Sprintf(`
resource "digitalocean_gradientai_custom_model" "test" {
  name                        = %q
  source_type                 = "SOURCE_TYPE_SPACES_BUCKET"
  preferred_gpu_region        = %q
  accept_terms_and_conditions = true

  source_ref {
    bucket = %q
    region = %q
    prefix = %q
  }
}
`, name, testCustomModelGpuRegion(), bucket, region, prefix)
}

func testAccCustomModelConfig_spacesUpdated(name, bucket, region, prefix string) string {
	return fmt.Sprintf(`
resource "digitalocean_gradientai_custom_model" "test" {
  name                        = %q
  description                 = "updated by terraform acceptance test"
  source_type                 = "SOURCE_TYPE_SPACES_BUCKET"
  preferred_gpu_region        = %q
  accept_terms_and_conditions = true
  tags                        = ["terraform-test", "acceptance"]

  license           = "apache-2.0"
  parameters        = "8B"
  input_modalities  = ["text"]
  output_modalities = ["text"]

  source_ref {
    bucket = %q
    region = %q
    prefix = %q
  }
}
`, name, testCustomModelGpuRegion(), bucket, region, prefix)
}

// testAccCustomModelConfig_hfWithSpacesOnlyFields builds a config that combines
// a HuggingFace source with the Spaces-only metadata fields. CustomizeDiff
// should reject this at plan time before any API call is made, so the config
// uses a clearly-fake repo_id to make the intent obvious and to make sure no
// real import is ever issued.
func testAccCustomModelConfig_hfWithSpacesOnlyFields(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_gradientai_custom_model" "test" {
  name                        = %q
  source_type                 = "SOURCE_TYPE_HUGGINGFACE"
  preferred_gpu_region        = %q
  accept_terms_and_conditions = true

  license          = "apache-2.0"
  input_modalities = ["text"]

  source_ref {
    repo_id     = "example/should-never-import"
    access_type = "ACCESS_TYPE_PUBLIC"
  }
}
`, name, testCustomModelGpuRegion())
}

// TestCustomModelCustomizeDiff_SpacesOnlyFields is a fast unit test (no TF_ACC,
// no DO API) that exercises the CustomizeDiff guard which rejects the
// Spaces-only metadata fields when source_type is not SOURCE_TYPE_SPACES_BUCKET.
// It uses schema.Resource.Diff with a raw config map to drive the plan.
func TestCustomModelCustomizeDiff_SpacesOnlyFields(t *testing.T) {
	baseHF := func(extra map[string]interface{}) map[string]interface{} {
		raw := map[string]interface{}{
			"name":                        "tf-acc-test",
			"source_type":                 "SOURCE_TYPE_HUGGINGFACE",
			"preferred_gpu_region":        "tor1",
			"accept_terms_and_conditions": true,
			"source_ref": []interface{}{
				map[string]interface{}{
					"repo_id":     "example/repo",
					"access_type": "ACCESS_TYPE_PUBLIC",
				},
			},
		}
		for k, v := range extra {
			raw[k] = v
		}
		return raw
	}

	baseSpaces := func(extra map[string]interface{}) map[string]interface{} {
		raw := map[string]interface{}{
			"name":                        "tf-acc-test",
			"source_type":                 "SOURCE_TYPE_SPACES_BUCKET",
			"preferred_gpu_region":        "tor1",
			"accept_terms_and_conditions": true,
			"source_ref": []interface{}{
				map[string]interface{}{
					"bucket": "example-bucket",
					"region": "nyc3",
					"prefix": "models/example/",
				},
			},
		}
		for k, v := range extra {
			raw[k] = v
		}
		return raw
	}

	cases := []struct {
		name           string
		attrs          map[string]interface{}
		expectErr      bool
		expectContains string
	}{
		{
			name:      "HF without Spaces-only fields is accepted",
			attrs:     baseHF(nil),
			expectErr: false,
		},
		{
			name:           "HF with license is rejected",
			attrs:          baseHF(map[string]interface{}{"license": "apache-2.0"}),
			expectErr:      true,
			expectContains: `license can only be set when source_type = "SOURCE_TYPE_SPACES_BUCKET"`,
		},
		{
			name:           "HF with parameters is rejected",
			attrs:          baseHF(map[string]interface{}{"parameters": "8B"}),
			expectErr:      true,
			expectContains: `parameters can only be set`,
		},
		{
			name:           "HF with input_modalities is rejected",
			attrs:          baseHF(map[string]interface{}{"input_modalities": []interface{}{"text"}}),
			expectErr:      true,
			expectContains: `input_modalities can only be set`,
		},
		{
			name:           "HF with output_modalities is rejected",
			attrs:          baseHF(map[string]interface{}{"output_modalities": []interface{}{"text"}}),
			expectErr:      true,
			expectContains: `output_modalities can only be set`,
		},
		{
			name: "HF with multiple Spaces-only fields lists all in error",
			attrs: baseHF(map[string]interface{}{
				"license":          "apache-2.0",
				"parameters":       "8B",
				"input_modalities": []interface{}{"text"},
			}),
			expectErr:      true,
			expectContains: `license, parameters, input_modalities can only be set`,
		},
		{
			name: "HF with multiple Spaces-only fields error names actionable remedy",
			attrs: baseHF(map[string]interface{}{
				"license": "apache-2.0",
			}),
			expectErr:      true,
			expectContains: `Either set source_type to "SOURCE_TYPE_SPACES_BUCKET" or remove these attributes`,
		},
		{
			name: "Spaces with all Spaces-only fields is accepted",
			attrs: baseSpaces(map[string]interface{}{
				"license":           "apache-2.0",
				"parameters":        "8B",
				"input_modalities":  []interface{}{"text"},
				"output_modalities": []interface{}{"text"},
			}),
			expectErr: false,
		},
		{
			name: "HF with empty Spaces-only string field is accepted",
			attrs: baseHF(map[string]interface{}{
				"license":    "",
				"parameters": "",
			}),
			expectErr: false,
		},
		{
			name: "HF with empty Spaces-only list fields is accepted",
			attrs: baseHF(map[string]interface{}{
				"input_modalities":  []interface{}{},
				"output_modalities": []interface{}{},
			}),
			expectErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := gradientai.ResourceDigitalOceanCustomModel()
			cfg := terraform.NewResourceConfigRaw(c.attrs)
			_, err := r.Diff(context.Background(), nil, cfg, nil)

			if c.expectErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", c.expectContains)
				}
				if !strings.Contains(err.Error(), c.expectContains) {
					t.Fatalf("expected error containing %q, got %q", c.expectContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %q", err.Error())
			}
		})
	}
}

func TestFlattenDigitalOceanCustomModel_ErrorMessage(t *testing.T) {
	t.Parallel()

	flat, err := gradientai.FlattenDigitalOceanCustomModel(&godo.CustomModel{
		Uuid:         "55555555-5555-5555-5555-555555555555",
		Name:         "team/failed-model",
		Status:       godo.CustomModelStatusFailed,
		ErrorMessage: "Repository not found or access denied",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := flat["error_message"], "Repository not found or access denied"; got != want {
		t.Fatalf("error_message = %q, want %q", got, want)
	}
	if got, want := flat["status"], string(godo.CustomModelStatusFailed); got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}
