package byoip

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanBYOIPAddresses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanBYOIPAddressesRead,
		Schema: map[string]*schema.Schema{
			"byoip_prefix_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of the BYOIP prefix to list addresses from",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of IP addresses allocated from the BYOIP prefix",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the IP address allocation",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region where the IP is allocated",
						},
						"assigned_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp when the IP was assigned",
						},
					},
				},
			},
		},
	}
}

func dataSourceDigitalOceanBYOIPAddressesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := getBYOIPService(meta)
	uuid := d.Get("byoip_prefix_uuid").(string)

	// List all addresses from the BYOIP prefix
	addresses, _, err := service.GetResources(context.Background(), uuid, nil)
	if err != nil {
		return diag.Errorf("Error retrieving BYOIP addresses: %s", err)
	}

	d.SetId(uuid)

	if err := d.Set("addresses", flattenBYOIPAddresses(addresses)); err != nil {
		return diag.Errorf("Error setting addresses: %s", err)
	}

	return nil
}

func flattenBYOIPAddresses(addresses []godo.BYOIPPrefixResource) []interface{} {
	if addresses == nil {
		return nil
	}

	flattenedAddresses := make([]interface{}, len(addresses))
	for i, addr := range addresses {
		rawAddress := map[string]interface{}{
			"id":          int(addr.ID),
			"ip_address":  addr.Resource,
			"region":      addr.Region,
			"assigned_at": addr.AssignedAt.UTC().String(),
		}
		flattenedAddresses[i] = rawAddress
	}

	return flattenedAddresses
}
