package database

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanDatabaseMetricsCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDatabaseMetricsCredentialsRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceDigitalOceanDatabaseMetricsCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	creds, _, err := client.Databases.GetMetricsCredentials(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving database metrics credentials: %s", err)
	}

	d.SetId("metrics-credentials")
	d.Set("username", creds.BasicAuthUsername)
	d.Set("password", creds.BasicAuthPassword)

	return nil
}
