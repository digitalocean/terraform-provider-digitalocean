package digitalocean

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanDatabaseCA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDatabaseCARead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanDatabaseCARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)
	d.SetId(clusterID)

	ca, _, err := client.Databases.GetCA(context.Background(), clusterID)
	if err != nil {
		return diag.Errorf("Error retrieving database CA certificate: %s", err)
	}

	d.Set("certificate", string(ca.Certificate))

	return nil
}
