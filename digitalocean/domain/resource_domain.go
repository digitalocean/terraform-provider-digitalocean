package domain

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDomainCreate,
		ReadContext:   resourceDigitalOceanDomainRead,
		DeleteContext: resourceDigitalOceanDomainDelete,
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
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	// Build up our creation options

	opts := &godo.DomainCreateRequest{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("ip_address"); ok {
		opts.IPAddress = v.(string)
	}

	log.Printf("[DEBUG] Domain create configuration: %#v", opts)
	domain, _, err := client.Domains.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating Domain: %s", err)
	}

	d.SetId(domain.Name)
	log.Printf("[INFO] Domain Name: %s", domain.Name)

	return resourceDigitalOceanDomainRead(ctx, d, meta)
}

func resourceDigitalOceanDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	domain, resp, err := client.Domains.Get(context.Background(), d.Id())
	if err != nil {
		// If the domain is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving domain: %s", err)
	}

	d.Set("name", domain.Name)
	d.Set("urn", domain.URN())
	d.Set("ttl", domain.TTL)

	return nil
}

func resourceDigitalOceanDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Deleting Domain: %s", d.Id())
	_, err := client.Domains.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting Domain: %s", err)
	}

	d.SetId("")
	return nil
}
