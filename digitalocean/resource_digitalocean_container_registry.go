package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const RegistryHostname = "registry.digitalocean.com"

func resourceDigitalOceanContainerRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanContainerRegistryCreate,
		ReadContext:   resourceDigitalOceanContainerRegistryRead,
		UpdateContext: resourceDigitalOceanContainerRegistryUpdate,
		DeleteContext: resourceDigitalOceanContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"subscription_tier_slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"starter",
					"basic",
					"professional",
				}, false),
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	// Build up our creation options
	opts := &godo.RegistryCreateRequest{
		Name:                 d.Get("name").(string),
		SubscriptionTierSlug: d.Get("subscription_tier_slug").(string),
	}

	log.Printf("[DEBUG] Container Registry create configuration: %#v", opts)
	reg, _, err := client.Registry.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating container registry: %s", err)
	}

	d.SetId(reg.Name)
	log.Printf("[INFO] Container Registry: %s", reg.Name)

	return resourceDigitalOceanContainerRegistryRead(ctx, d, meta)
}

func resourceDigitalOceanContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	reg, resp, err := client.Registry.Get(context.Background())
	if err != nil {
		// If the registry is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving container registry: %s", err)
	}

	d.SetId(reg.Name)
	d.Set("name", reg.Name)
	d.Set("endpoint", fmt.Sprintf("%s/%s", RegistryHostname, reg.Name))
	d.Set("server_url", RegistryHostname)

	sub, _, err := client.Registry.GetSubscription(context.Background())
	if err != nil {
		return diag.Errorf("Error retrieving container registry subscription: %s", err)
	}
	d.Set("subscription_tier_slug", sub.Tier.Slug)

	return nil
}

func resourceDigitalOceanContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()
	if d.HasChange("subscription_tier_slug") {
		req := &godo.RegistrySubscriptionUpdateRequest{
			TierSlug: d.Get("subscription_tier_slug").(string),
		}

		_, _, err := client.Registry.UpdateSubscription(ctx, req)
		if err != nil {
			return diag.Errorf("Error updating container registry subscription: %s", err)
		}
	}
	return resourceDigitalOceanContainerRegistryRead(ctx, d, meta)
}

func resourceDigitalOceanContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting container registry: %s", d.Id())
	_, err := client.Registry.Delete(context.Background())
	if err != nil {
		return diag.Errorf("Error deleting container registry: %s", err)
	}
	d.SetId("")
	return nil
}
