package byoip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanBYOIPPrefix() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanBYOIPPrefixRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of the BYOIP prefix",
				ValidateFunc: validation.NoZeroValues,
			},
			"prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIDR notation of the prefix",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region where the prefix is deployed",
			},
			"advertised": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the prefix is advertised",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the BYOIP prefix",
			},
			"failure_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason for failure if status is failed",
			},
		},
	}
}

func dataSourceDigitalOceanBYOIPPrefixRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := getBYOIPService(meta)
	uuid := d.Get("uuid").(string)

	prefix, _, err := service.Get(context.Background(), uuid)
	if err != nil {
		return diag.Errorf("Error retrieving BYOIP prefix: %s", err)
	}

	d.SetId(prefix.UUID)
	d.Set("uuid", prefix.UUID)
	d.Set("prefix", prefix.Prefix)
	d.Set("region", prefix.Region)
	d.Set("status", prefix.Status)
	d.Set("advertised", prefix.Advertised)
	d.Set("failure_reason", prefix.FailureReason)

	return nil
}
