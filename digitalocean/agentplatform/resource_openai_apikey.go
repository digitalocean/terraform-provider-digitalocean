package agentplatform

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDigitalOceanOpenAIApiKey defines the DigitalOcean
func ResourceDigitalOceanOpenAIApiKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanOpenAIApiKeyCreate,
		ReadContext:   resourceDigitalOceanOpenAIApiKeyRead,
		UpdateContext: resourceDigitalOceanOpenAIApiKeyUpdate,
		DeleteContext: resourceDigitalOceanOpenAIApiKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"api_key", "api_key_wo"},
				Description:  "The OpenAI API key. This value is stored in state; use `api_key_wo` to avoid that.",
			},
			"api_key_wo": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				WriteOnly:    true,
				ExactlyOneOf: []string{"api_key", "api_key_wo"},
				RequiredWith: []string{"api_key_wo_version"},
				Description:  "The OpenAI API key as a write-only argument that is never stored in state. Requires Terraform 1.11 or later.",
			},
			"api_key_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"api_key_wo"},
				Description:  "Version of `api_key_wo`. Change this value to update the key.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name for the API key.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the API key.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the API key was created.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Who created the API key.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the API key was last updated.",
			},
			"deleted_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the API key was deleted.",
			},
			"model": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Models associated with the OpenAI API key",
				Elem:        ModelSchema(),
			},
		},
	}
}

func resourceDigitalOceanOpenAIApiKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	key, diags := openAIApiKeyValue(d)
	if diags.HasError() {
		return diags
	}

	openAIRequest := &godo.OpenAIAPIKeyCreateRequest{
		ApiKey: key,
		Name:   d.Get("name").(string),
	}

	apiKey, _, err := client.AgentPlatform.CreateOpenAIAPIKey(ctx, openAIRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(apiKey.Uuid)
	return resourceDigitalOceanOpenAIApiKeyRead(ctx, d, meta)
}

func resourceDigitalOceanOpenAIApiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	apiKey, _, err := client.AgentPlatform.GetOpenAIAPIKey(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if apiKey == nil {
		d.SetId("")
		return nil
	}

	d.Set("uuid", apiKey.Uuid)
	d.Set("name", apiKey.Name)

	if apiKey.CreatedAt != nil {
		d.Set("created_at", apiKey.CreatedAt.UTC().String())
	} else {
		d.Set("created_at", "")
	}
	d.Set("created_by", apiKey.CreatedBy)
	if apiKey.UpdatedAt != nil {
		d.Set("updated_at", apiKey.UpdatedAt.UTC().String())
	} else {
		d.Set("updated_at", "")
	}
	if apiKey.DeletedAt != nil {
		d.Set("deleted_at", apiKey.DeletedAt.UTC().String())
	} else {
		d.Set("deleted_at", "")
	}

	// Flatten models if needed
	if apiKey.Models != nil {
		d.Set("model", flattenModel(apiKey.Models))
	} else {
		d.Set("model", []interface{}{})
	}

	d.SetId(apiKey.Uuid)

	return nil
}

func resourceDigitalOceanOpenAIApiKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	key, diags := openAIApiKeyValue(d)
	if diags.HasError() {
		return diags
	}

	openAIRequest := &godo.OpenAIAPIKeyUpdateRequest{
		Name:       d.Get("name").(string),
		ApiKey:     key,
		ApiKeyUuid: d.Get("uuid").(string),
	}
	hasChanges := false

	if d.HasChange("name") {
		openAIRequest.Name = d.Get("name").(string)
		hasChanges = true
	}
	if d.HasChanges("api_key", "api_key_wo_version") {
		hasChanges = true
	}
	if !hasChanges {
		return nil
	}

	_, _, err := client.AgentPlatform.UpdateOpenAIAPIKey(ctx, d.Id(), openAIRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDigitalOceanOpenAIApiKeyRead(ctx, d, meta)

}

// openAIApiKeyValue returns the key from api_key, or from the raw config for
// api_key_wo since write-only values are never persisted to state.
func openAIApiKeyValue(d *schema.ResourceData) (string, diag.Diagnostics) {
	if v, ok := d.GetOk("api_key"); ok {
		return v.(string), nil
	}

	wo, diags := d.GetRawConfigAt(cty.GetAttrPath("api_key_wo"))
	if diags.HasError() {
		return "", diags
	}
	if wo.IsNull() || !wo.IsKnown() {
		return "", nil
	}

	return wo.AsString(), nil
}

func resourceDigitalOceanOpenAIApiKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Id()

	_, resp, err := client.AgentPlatform.DeleteOpenAIAPIKey(ctx, id)
	if err != nil {
		if resp != nil && resp.Response != nil && resp.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting OpenAI API Key (%s): %s", id, err))
	}

	d.SetId("")
	return nil
}
