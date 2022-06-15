package digitalocean

import (
	"context"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanFloatingIP() *schema.Resource {
	return &schema.Resource{
		// TODO: Uncomment when dates for deprecation timeline are set.
		// DeprecationMessage: "This resource is deprecated and will be removed in a future release. Please use digitalocean_reserved_ip instead.",
		CreateContext: resourceDigitalOceanFloatingIPCreate,
		UpdateContext: resourceDigitalOceanFloatingIPUpdate,
		ReadContext:   resourceDigitalOceanFloatingIPRead,
		DeleteContext: resourceDigitalOceanReservedIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanFloatingIPImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the floating ip",
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"droplet_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDigitalOceanFloatingIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resourceDigitalOceanReservedIPCreate(ctx, d, meta)
	if err != nil {
		return err
	}
	reservedIPURNtoFloatingIPURN(d)

	return nil
}

func resourceDigitalOceanFloatingIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resourceDigitalOceanReservedIPUpdate(ctx, d, meta)
	if err != nil {
		return err
	}
	reservedIPURNtoFloatingIPURN(d)

	return nil
}

func resourceDigitalOceanFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resourceDigitalOceanReservedIPRead(ctx, d, meta)
	if err != nil {
		return err
	}
	reservedIPURNtoFloatingIPURN(d)

	return nil
}

func resourceDigitalOceanFloatingIPImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	_, err := resourceDigitalOceanReservedIPImport(ctx, d, meta)
	if err != nil {
		return nil, err
	}
	reservedIPURNtoFloatingIPURN(d)

	return []*schema.ResourceData{d}, nil
}

// reservedIPURNtoFloatingIPURN re-formats a reserved IP URN as floating IP URN.
// TODO: Remove when the projects' API changes return values.
func reservedIPURNtoFloatingIPURN(d *schema.ResourceData) {
	ip := d.Get("ip_address")
	d.Set("urn", godo.FloatingIP{IP: ip.(string)}.URN())
}
