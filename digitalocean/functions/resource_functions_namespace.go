package functions

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanFunctionsNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanFunctionsNamespaceCreate,
		ReadContext:   resourceDigitalOceanFunctionsNamespaceRead,
		DeleteContext: resourceDigitalOceanFunctionsNamespaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "The namespace's unique name.",
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The region where the namespace is deployed, e.g. nyc1.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The namespace's UUID.",
			},
			"api_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API host used to deploy and invoke functions in the namespace.",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The namespace's access key, used together with api_host by the serverless tooling.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the namespace was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the namespace was last updated.",
			},
		},
	}
}

func resourceDigitalOceanFunctionsNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	req := &godo.FunctionsNamespaceCreateRequest{
		Label:  d.Get("label").(string),
		Region: d.Get("region").(string),
	}

	log.Printf("[DEBUG] Functions namespace create request: %#v", req)
	ns, _, err := functionsService(meta).CreateNamespace(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating Functions namespace: %s", err)
	}

	d.SetId(ns.Namespace)
	log.Printf("[INFO] Functions namespace created: %s", ns.Namespace)

	return resourceDigitalOceanFunctionsNamespaceRead(ctx, d, meta)
}

func resourceDigitalOceanFunctionsNamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ns, resp, err := functionsService(meta).GetNamespace(ctx, d.Id())
	if err != nil {
		if isNotFound(resp, err) {
			log.Printf("[WARN] Functions namespace %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading Functions namespace %s: %s", d.Id(), err)
	}

	if err := setFunctionsNamespaceAttributes(d, ns); err != nil {
		return diag.Errorf("Error setting Functions namespace attributes: %s", err)
	}

	return nil
}

func resourceDigitalOceanFunctionsNamespaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Functions namespace: %s", d.Id())
	resp, err := functionsService(meta).DeleteNamespace(ctx, d.Id())
	if err != nil && !isNotFound(resp, err) {
		return diag.Errorf("Error deleting Functions namespace %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setFunctionsNamespaceAttributes(d *schema.ResourceData, ns *godo.FunctionsNamespace) error {
	attrs := map[string]interface{}{
		"label":      ns.Label,
		"region":     ns.Region,
		"uuid":       ns.UUID,
		"api_host":   ns.ApiHost,
		"key":        ns.Key,
		"created_at": formatTime(ns.CreatedAt),
		"updated_at": formatTime(ns.UpdatedAt),
	}
	for k, v := range attrs {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}
