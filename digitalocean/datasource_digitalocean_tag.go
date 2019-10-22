package digitalocean

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanTagRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the tag",
				ValidateFunc: validateTag,
			},
		},
	}
}

func dataSourceDigitalOceanTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	name := d.Get("name").(string)

	tag, resp, err := client.Tags.Get(context.Background(), name)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return fmt.Errorf("tag not found: %s", err)
		}
		return fmt.Errorf("Error retrieving tag: %s", err)
	}

	d.SetId(tag.Name)
	d.Set("name", tag.Name)

	return nil
}
