package dedicatedinference

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDedicatedInferenceToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDedicatedInferenceTokenCreate,
		ReadContext:   resourceDigitalOceanDedicatedInferenceTokenRead,
		DeleteContext: resourceDigitalOceanDedicatedInferenceTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"dedicated_inference_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The ID of the dedicated inference endpoint this token belongs to.",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "A human-readable name for the token.",
				ValidateFunc: validation.NoZeroValues,
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The token value. Only available immediately after creation and not retrievable afterwards.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the token was created.",
			},
		},
	}
}

func resourceDigitalOceanDedicatedInferenceTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	diID := d.Get("dedicated_inference_id").(string)
	req := &godo.DedicatedInferenceTokenCreateRequest{
		Name: d.Get("name").(string),
	}

	token, _, err := client.DedicatedInference.CreateToken(ctx, diID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating token for dedicated inference endpoint (%s): %w", diID, err))
	}

	d.SetId(fmt.Sprintf("%s:%s", diID, token.ID))
	d.Set("token", token.Value)

	return resourceDigitalOceanDedicatedInferenceTokenRead(ctx, d, meta)
}

func resourceDigitalOceanDedicatedInferenceTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	diID, tokenID, err := parseDedicatedInferenceTokenID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("dedicated_inference_id", diID)

	opts := &godo.ListOptions{Page: 1, PerPage: 200}
	for {
		tokens, resp, err := client.DedicatedInference.ListTokens(ctx, diID, opts)
		if err != nil {
			if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
			return diag.FromErr(fmt.Errorf("error listing tokens for dedicated inference endpoint (%s): %w", diID, err))
		}

		for _, t := range tokens {
			if t.ID == tokenID {
				d.Set("name", t.Name)
				if !t.CreatedAt.IsZero() {
					d.Set("created_at", t.CreatedAt.UTC().String())
				}
				return nil
			}
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error listing tokens for dedicated inference endpoint (%s): %w", diID, err))
		}
		opts.Page = page + 1
	}

	// Token not found; it may have been revoked outside Terraform.
	d.SetId("")
	return nil
}

func resourceDigitalOceanDedicatedInferenceTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	diID, tokenID, err := parseDedicatedInferenceTokenID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.DedicatedInference.RevokeToken(ctx, diID, tokenID)
	if err != nil {
		if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error revoking token (%s) for dedicated inference endpoint (%s): %w", tokenID, diID, err))
	}

	d.SetId("")
	return nil
}

func parseDedicatedInferenceTokenID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid dedicated inference token ID %q: expected {dedicated_inference_id}:{token_id}", id)
	}
	return parts[0], parts[1], nil
}
