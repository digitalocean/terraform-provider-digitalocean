package digitalocean

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceDigitalOceanContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanContainerRegistryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the container registry",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func dataSourceDigitalOceanContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	reg, response, err := client.Registry.Get(context.Background())

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return fmt.Errorf("registry not found: %s", err)
		}
		return fmt.Errorf("Error retrieving registry: %s", err)
	}

	d.SetId(reg.Name)
	d.Set("name", reg.Name)
	return nil
}
