package database

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanDatabaseUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDatabaseUserRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"access_cert": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"mysql_auth_plugin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     userACLSchema(),
						},
						"opensearch_acl": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     userOpenSearchACLSchema(),
						},
					},
				},
			},
		},
	}
}

func dataSourceDigitalOceanDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	user, resp, err := client.Databases.GetUser(context.Background(), clusterID, name)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.Errorf("Database user not found: %s", err)
		}
		return diag.Errorf("Error retrieving database user: %s", err)
	}

	d.SetId(makeDatabaseUserID(clusterID, name))
	d.Set("role", user.Role)
	d.Set("password", user.Password)

	if user.MySQLSettings != nil {
		d.Set("mysql_auth_plugin", user.MySQLSettings.AuthPlugin)
	}

	if user.AccessCert != "" {
		d.Set("access_cert", user.AccessCert)
	}
	if user.AccessKey != "" {
		d.Set("access_key", user.AccessKey)
	}

	if err := d.Set("settings", flattenUserSettings(user.Settings)); err != nil {
		return diag.Errorf("Error setting user settings: %#v", err)
	}
	return nil
}
