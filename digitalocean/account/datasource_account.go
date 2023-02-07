package account

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanAccountRead,
		Schema: map[string]*schema.Schema{
			"droplet_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of Droplets current user or team may have active at one time.",
			},
			"floating_ip_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of Floating IPs the current user or team may have.",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address used by the current user to register for DigitalOcean.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique universal identifier for the current user.",
			},
			"email_verified": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, the user has verified their account via email. False otherwise.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This value is one of \"active\", \"warning\" or \"locked\".",
			},
			"status_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-readable message giving more details about the status of the account.",
			},
		},
	}
}

func dataSourceDigitalOceanAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	account, _, err := client.Account.Get(context.Background())
	if err != nil {
		return diag.Errorf("Error retrieving account: %s", err)
	}

	d.SetId(account.UUID)
	d.Set("droplet_limit", account.DropletLimit)
	d.Set("floating_ip_limit", account.FloatingIPLimit)
	d.Set("email", account.Email)
	d.Set("uuid", account.UUID)
	d.Set("email_verified", account.EmailVerified)
	d.Set("status", account.Status)
	d.Set("status_message", account.StatusMessage)

	return nil
}
