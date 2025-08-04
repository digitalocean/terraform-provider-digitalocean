package partnernetworkconnect

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanPartnerAttachmentServiceKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanPartnerAttachmentServiceKeyRead,
		Schema: map[string]*schema.Schema{
			"attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Partner Attachment for which to retrieve the service key",
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanPartnerAttachmentServiceKeyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	attachmentID := d.Get("attachment_id").(string)
	serviceKey, _, err := client.PartnerAttachment.GetServiceKey(ctx, attachmentID)
	if err != nil {
		return diag.Errorf("error retrieving service key for partner attachment %q: %s", attachmentID, err)
	}

	d.SetId(attachmentID)
	d.Set("value", serviceKey.Value)
	d.Set("state", serviceKey.State)
	d.Set("created_at", serviceKey.CreatedAt.UTC().String())

	return nil
}
