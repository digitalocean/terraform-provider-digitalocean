package digitalocean

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const RegistryHostname = "registry.digitalocean.com"

func dataSourceDigitalOceanContainerRegistry() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanContainerRegistryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the container registry",
				ValidateFunc: validation.NoZeroValues,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	reg, response, err := client.Registry.Get(context.Background())

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return diag.Errorf("registry not found: %s", err)
		}
		return diag.Errorf("Error retrieving registry: %s", err)
	}

	d.SetId(reg.Name)
	d.Set("name", reg.Name)
	d.Set("endpoint", fmt.Sprintf("%s/%s", RegistryHostname, reg.Name))
	d.Set("server_url", RegistryHostname)
	return nil
}
