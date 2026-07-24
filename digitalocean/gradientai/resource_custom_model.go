package gradientai

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// customModelAlwaysMutableMetadataFields are accepted by the metadata-update
// API for every source type.
var customModelAlwaysMutableMetadataFields = []string{
	"description",
	"tags",
}

// customModelSpacesOnlyMetadataFields are accepted by the metadata-update API
// only when the model was imported via SOURCE_TYPE_SPACES_BUCKET. For any
// other source type the PATCH is silently ignored while GET keeps returning
// the importer-reported values, which would cause a permanent terraform plan
// diff — see resourceDigitalOceanCustomModelCustomizeDiff.
var customModelSpacesOnlyMetadataFields = []string{
	"license",
	"parameters",
	"input_modalities",
	"output_modalities",
}

// Valid source_type and source_ref.access_type values accepted by the API.
// The provider passes these strings through to the API as-is.
var (
	customModelSourceTypes = []string{
		string(godo.CustomModelSourceTypeHuggingFace),
		string(godo.CustomModelSourceTypeSpacesBucket),
		string(godo.CustomModelSourceTypeSDKUpload),
		string(godo.CustomModelSourceTypeFineTuning),
	}
	customModelAccessTypes = []string{
		string(godo.CustomModelSourceRefAccessTypePublic),
		string(godo.CustomModelSourceRefAccessTypePrivate),
		string(godo.CustomModelSourceRefAccessTypeGated),
	}
)

func ResourceDigitalOceanCustomModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanCustomModelCreate,
		ReadContext:   resourceDigitalOceanCustomModelRead,
		UpdateContext: resourceDigitalOceanCustomModelUpdate,
		DeleteContext: resourceDigitalOceanCustomModelDelete,
		CustomizeDiff: resourceDigitalOceanCustomModelCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A human-readable name for the custom model.",
				ValidateFunc: validation.NoZeroValues,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the custom model.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "User-defined tags associated with the custom model.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"source_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Source of the model to import. One of SOURCE_TYPE_HUGGINGFACE, SOURCE_TYPE_SPACES_BUCKET, SOURCE_TYPE_SDK_UPLOAD, SOURCE_TYPE_FINE_TUNING.",
				ValidateFunc: validation.StringInSlice(customModelSourceTypes, false),
			},
			"source_ref": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Reference to the source from which to import the custom model.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Repository identifier (e.g. the HuggingFace repo). Required for SOURCE_TYPE_HUGGINGFACE sources.",
						},
						"commit_sha": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Commit SHA to pin for the import. If omitted, the API resolves and returns the SHA actually imported.",
						},
						"access_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Access type for the source repository. One of ACCESS_TYPE_PUBLIC, ACCESS_TYPE_PRIVATE, ACCESS_TYPE_GATED.",
							ValidateFunc: validation.StringInSlice(customModelAccessTypes, false),
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Spaces bucket name for SOURCE_TYPE_SPACES_BUCKET sources.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region of the source bucket.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key prefix inside the source bucket.",
						},
						"hf_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "HuggingFace token used to access ACCESS_TYPE_PRIVATE or ACCESS_TYPE_GATED repositories. Write-only.",
						},
					},
				},
			},
			"preferred_gpu_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preferred GPU region where the model artifacts should be staged.",
			},
			"accept_terms_and_conditions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the caller accepts the model provider's terms and conditions. Write-only.",
			},

			// Computed fields populated by Read.
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the custom model.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the custom model.",
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Error message if the custom model import failed.",
			},
			"architecture": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Model architecture as reported by the importer.",
			},
			"total_size_bytes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Total size of the imported model artifacts in bytes.",
			},
			"file_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of files that make up the imported model.",
			},
			"license": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "License of the model. Defaults to the value reported by the importer. Caller-supplied overrides are honored only for SOURCE_TYPE_SPACES_BUCKET imports.",
			},
			"context_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum context length supported by the model.",
			},
			"cost_estimate_per_month": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Estimated monthly cost of running the custom model.",
			},
			"input_modalities": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Input modalities supported by the model. Defaults to the values reported by the importer. Caller-supplied overrides are honored only for SOURCE_TYPE_SPACES_BUCKET imports.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"output_modalities": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Output modalities produced by the model. Defaults to the values reported by the importer. Caller-supplied overrides are honored only for SOURCE_TYPE_SPACES_BUCKET imports.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"parameters": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Parameter-count summary for the model. Defaults to the value reported by the importer. Caller-supplied overrides are honored only for SOURCE_TYPE_SPACES_BUCKET imports.",
			},
			"team_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the team that owns the custom model.",
			},
			"storage_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region where the custom model artifacts are stored.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the custom model was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the custom model was last updated.",
			},
			"active_deployments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Active dedicated inference deployments referencing this custom model.",
				Elem:        customModelActiveDeploymentSchemaRead(),
			},
		},
	}
}

func resourceDigitalOceanCustomModelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := buildCustomModelImportRequest(d)

	importResp, _, err := client.GradientAI.ImportCustomModel(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error importing custom model: %w", err))
	}
	if importResp == nil || importResp.Model == nil || importResp.Model.Uuid == "" {
		return diag.Errorf("custom model import returned no model UUID")
	}

	d.SetId(importResp.Model.Uuid)

	if err := waitForCustomModelReady(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(fmt.Errorf("custom model (%s) did not become ready: %w", d.Id(), err))
	}

	return resourceDigitalOceanCustomModelRead(ctx, d, meta)
}

func resourceDigitalOceanCustomModelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	model, resp, err := client.GradientAI.GetCustomModel(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading custom model (%s): %w", d.Id(), err))
	}
	if model == nil {
		d.SetId("")
		return nil
	}

	flat, err := FlattenDigitalOceanCustomModel(model)
	if err != nil {
		return diag.FromErr(err)
	}

	// Preserve write-only source_ref fields (access_type, hf_token, bucket,
	// region, prefix) that the API never echoes back, so subsequent plans
	// don't show spurious drift.
	flat["source_ref"] = normalizeCustomModelSourceRefForState(flat["source_ref"], d.Get("source_ref"))

	for key, value := range flat {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(fmt.Errorf("error setting %s: %w", key, err))
		}
	}

	return nil
}

// normalizeCustomModelSourceRefForState prepares the API-derived source_ref
// for writing to state. For write-only inputs that the API never echoes back
// (access_type, hf_token, bucket, region, prefix), it preserves whatever the
// caller previously supplied so subsequent plans do not show spurious drift.
func normalizeCustomModelSourceRefForState(apiVal, existingVal interface{}) []interface{} {
	apiList, _ := apiVal.([]interface{})
	if len(apiList) == 0 {
		// API returned no source_ref; fall back entirely to existing state so
		// we don't lose user-supplied write-only fields.
		if existingList, ok := existingVal.([]interface{}); ok {
			return existingList
		}
		return []interface{}{}
	}
	row, ok := apiList[0].(map[string]interface{})
	if !ok {
		return apiList
	}

	// Fields below are accepted on import but not returned by GetCustomModel.
	// Carry them forward from the existing resource state.
	writeOnlyFields := []string{"access_type", "hf_token", "bucket", "region", "prefix"}
	var existingRow map[string]interface{}
	if existingList, ok := existingVal.([]interface{}); ok && len(existingList) > 0 {
		existingRow, _ = existingList[0].(map[string]interface{})
	}
	for _, key := range writeOnlyFields {
		cur, _ := row[key].(string)
		if cur != "" {
			continue
		}
		if existingRow == nil {
			continue
		}
		if v, ok := existingRow[key].(string); ok {
			row[key] = v
		}
	}
	return []interface{}{row}
}

// resourceDigitalOceanCustomModelCustomizeDiff rejects plans that set any of
// the Spaces-only metadata fields on a non-SOURCE_TYPE_SPACES_BUCKET resource.
// The API silently ignores those fields on PATCH for other source types, which
// would otherwise produce a permanent plan diff that no apply can resolve.
//
// HasChange + a non-empty value distinguishes caller-set values from
// API-reflected state values (e.g. license read back after `terraform import`
// of an HF model), so the latter does not trigger a spurious error.
func resourceDigitalOceanCustomModelCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	sourceType := d.Get("source_type").(string)
	if sourceType == string(godo.CustomModelSourceTypeSpacesBucket) {
		return nil
	}

	var conflicts []string
	for _, field := range customModelSpacesOnlyMetadataFields {
		if !d.HasChange(field) {
			continue
		}
		switch v := d.Get(field).(type) {
		case string:
			if v != "" {
				conflicts = append(conflicts, field)
			}
		case []interface{}:
			if len(v) > 0 {
				conflicts = append(conflicts, field)
			}
		}
	}
	if len(conflicts) > 0 {
		return fmt.Errorf(
			"%s can only be set when source_type = %q (got %q). "+
				"Either set source_type to %q or remove these attributes from your configuration.",
			strings.Join(conflicts, ", "),
			godo.CustomModelSourceTypeSpacesBucket,
			sourceType,
			godo.CustomModelSourceTypeSpacesBucket,
		)
	}
	return nil
}

func resourceDigitalOceanCustomModelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	isSpaces := d.Get("source_type").(string) == string(godo.CustomModelSourceTypeSpacesBucket)

	mutableFields := append([]string{}, customModelAlwaysMutableMetadataFields...)
	if isSpaces {
		mutableFields = append(mutableFields, customModelSpacesOnlyMetadataFields...)
	}
	if !d.HasChanges(mutableFields...) {
		return resourceDigitalOceanCustomModelRead(ctx, d, meta)
	}

	// Name is intentionally omitted from the payload because the API does not
	// support renaming a custom model.
	updateReq := &godo.CustomModelMetadataUpdateRequest{
		Description: d.Get("description").(string),
		Tags:        expandCustomModelTags(d.Get("tags")),
	}
	if isSpaces {
		// The Spaces-only metadata fields are gated on source_type both here
		// and in CustomizeDiff so non-Spaces PATCHes never carry attributes
		// the API would silently discard.
		updateReq.License = d.Get("license").(string)
		updateReq.Parameters = d.Get("parameters").(string)
		updateReq.InputModalities = expandCustomModelStringList(d.Get("input_modalities"))
		updateReq.OutputModalities = expandCustomModelStringList(d.Get("output_modalities"))
	}

	if _, _, err := client.GradientAI.UpdateCustomModelMetadata(ctx, d.Id(), updateReq); err != nil {
		return diag.FromErr(fmt.Errorf("error updating custom model (%s) metadata: %w", d.Id(), err))
	}

	return resourceDigitalOceanCustomModelRead(ctx, d, meta)
}

func resourceDigitalOceanCustomModelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	delResp, resp, err := client.GradientAI.DeleteCustomModel(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting custom model (%s): %w", d.Id(), err))
	}
	if delResp != nil && delResp.Status == godo.DeleteCustomModelStatusFail {
		return diag.Errorf("custom model (%s) delete failed: %s", d.Id(), delResp.Error)
	}

	d.SetId("")
	return nil
}

func waitForCustomModelReady(ctx context.Context, client *godo.Client, id string, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		model, _, err := client.GradientAI.GetCustomModel(ctx, id)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error polling custom model (%s): %w", id, err))
		}
		if model == nil {
			return retry.RetryableError(fmt.Errorf("custom model (%s) not yet visible, waiting", id))
		}
		switch model.Status {
		case godo.CustomModelStatusReady:
			return nil
		case godo.CustomModelStatusFailed:
			if model.ErrorMessage != "" {
				return retry.NonRetryableError(fmt.Errorf("custom model (%s) entered failed state: %s", id, model.ErrorMessage))
			}
			return retry.NonRetryableError(fmt.Errorf("custom model (%s) entered failed state", id))
		case godo.CustomModelStatusDeleted:
			return retry.NonRetryableError(fmt.Errorf("custom model (%s) was deleted while waiting for ready", id))
		default:
			return retry.RetryableError(fmt.Errorf("custom model (%s) is %s, waiting for ready", id, model.Status))
		}
	})
}

func buildCustomModelImportRequest(d *schema.ResourceData) *godo.CustomModelImportRequest {
	req := &godo.CustomModelImportRequest{
		Name:                     d.Get("name").(string),
		SourceType:               godo.CustomModelSourceType(d.Get("source_type").(string)),
		Description:              d.Get("description").(string),
		PreferredGpuRegion:       d.Get("preferred_gpu_region").(string),
		AcceptTermsAndConditions: d.Get("accept_terms_and_conditions").(bool),
		SourceRef:                expandCustomModelSourceRef(d.Get("source_ref").([]interface{})),
		Tags:                     expandCustomModelTags(d.Get("tags")),
	}
	return req
}

func expandCustomModelSourceRef(raw []interface{}) *godo.CustomModelSourceRef {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	return &godo.CustomModelSourceRef{
		RepoId:     m["repo_id"].(string),
		CommitSha:  m["commit_sha"].(string),
		AccessType: godo.CustomModelSourceRefAccessType(m["access_type"].(string)),
		Bucket:     m["bucket"].(string),
		Region:     m["region"].(string),
		Prefix:     m["prefix"].(string),
		HfToken:    m["hf_token"].(string),
	}
}

func expandCustomModelTags(raw interface{}) *godo.CustomModelTags {
	set, ok := raw.(*schema.Set)
	if !ok {
		return nil
	}
	tags := make([]string, 0, set.Len())
	for _, v := range set.List() {
		if s, ok := v.(string); ok && s != "" {
			tags = append(tags, s)
		}
	}
	return &godo.CustomModelTags{Tags: tags}
}

// expandCustomModelStringList converts a TypeList schema value (e.g.
// input_modalities, output_modalities) into a plain []string suitable for the
// metadata update request. Empty/missing input yields nil so the field is
// omitted from the JSON body.
func expandCustomModelStringList(raw interface{}) []string {
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, v := range list {
		if s, ok := v.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}
