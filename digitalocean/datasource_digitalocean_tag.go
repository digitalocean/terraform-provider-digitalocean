package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDigitalOceanTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanTagRead,
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the tag",
				ValidateFunc: validateTag,
			},
		},
	}
}

func dataSourceDigitalOceanTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	name := d.Get("name").(string)

	tag, resp, err := client.Tags.Get(context.Background(), name)
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("tag not found: %s", err)
		}
		return fmt.Errorf("Error retrieving tag: %s", err)
	}

	d.SetId(tag.Name)
	d.Set("name", tag.Name)

	return nil
}
