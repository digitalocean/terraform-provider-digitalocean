package registry

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanContainerRegistries() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanContainerRegistriesCreate,
		ReadContext:   resourceDigitalOceanContainerRegistriesRead,
		UpdateContext: resourceDigitalOceanContainerRegistriesUpdate,
		DeleteContext: resourceDigitalOceanContainerRegistriesDelete,
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
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_usage_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanContainerRegistriesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	// Build up our creation options
	opts := &godo.RegistryCreateRequest{
		Name:                 d.Get("name").(string),
		SubscriptionTierSlug: d.Get("subscription_tier_slug").(string),
	}

	if region, ok := d.GetOk("region"); ok {
		opts.Region = region.(string)
	}

	log.Printf("[DEBUG] Container Registries create configuration: %#v", opts)
	reg, _, err := client.Registries.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating container registries: %s", err)
	}

	d.SetId(reg.Name)
	log.Printf("[INFO] Container Registries: %s", reg.Name)

	return resourceDigitalOceanContainerRegistriesRead(ctx, d, meta)
}

func resourceDigitalOceanContainerRegistriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	reg, resp, err := client.Registries.Get(context.Background(), d.Id())
	if err != nil {
		// If the registry is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving container registries: %s", err)
	}

	d.SetId(reg.Name)
	d.Set("name", reg.Name)
	d.Set("region", reg.Region)
	d.Set("endpoint", fmt.Sprintf("%s/%s", RegistryHostname, reg.Name))
	d.Set("server_url", RegistryHostname)
	d.Set("created_at", reg.CreatedAt.UTC().String())
	d.Set("storage_usage_bytes", reg.StorageUsageBytes)

	sub, _, err := client.Registries.GetSubscription(context.Background())
	if err != nil {
		return diag.Errorf("Error retrieving container registries subscription: %s", err)
	}
	d.Set("subscription_tier_slug", sub.Tier.Slug)

	return nil
}

func resourceDigitalOceanContainerRegistriesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	if d.HasChange("subscription_tier_slug") {
		req := &godo.RegistrySubscriptionUpdateRequest{
			TierSlug: d.Get("subscription_tier_slug").(string),
		}

		_, _, err := client.Registries.UpdateSubscription(ctx, req)
		if err != nil {
			return diag.Errorf("Error updating container registries subscription: %s", err)
		}
	}
	return resourceDigitalOceanContainerRegistriesRead(ctx, d, meta)
}

func resourceDigitalOceanContainerRegistriesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Deleting container registries: %s", d.Id())
	_, err := client.Registries.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting container registries: %s", err)
	}
	d.SetId("")
	return nil
}
