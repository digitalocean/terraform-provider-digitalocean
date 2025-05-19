package spaces

import (
	"context"
	"log"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanSpacesKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSpacesKeyRead,

		Schema: spacesKeyDataSourceSchema(),
	}
}

func dataSourceSpacesKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	key, _, err := client.SpacesKeys.Get(ctx, d.Get("access_key").(string))
	if err != nil {
		return diag.Errorf("Error reading Spaces key: %s", err)
	}

	if key == nil {
		log.Printf("[WARN] Key not found: %s", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(key.AccessKey)
	d.Set("name", key.Name)
	d.Set("access_key", key.AccessKey)
	d.Set("grant", flattenGrants(key.Grants))
	d.Set("created_at", key.CreatedAt)
	return nil
}
