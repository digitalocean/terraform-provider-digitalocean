package digitalocean

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanDroplet() *schema.Resource {
	recordSchema := dropletSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["id"].ExactlyOneOf = []string{"id", "tag", "name"}
	recordSchema["id"].Optional = true
	recordSchema["name"].ExactlyOneOf = []string{"id", "tag", "name"}
	recordSchema["name"].Optional = true

	recordSchema["tag"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "unique tag of the Droplet",
		ValidateFunc: validation.NoZeroValues,
		ExactlyOneOf: []string{"id", "tag", "name"},
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDropletRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanDropletRead(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	var foundDroplet godo.Droplet

	if id, ok := d.GetOk("id"); ok {
		droplet, _, err := client.Droplets.Get(context.Background(), id.(int))
		if err != nil {
			return err
		}

		foundDroplet = *droplet
	} else if v, ok := d.GetOk("tag"); ok {
		dropletList, err := getDigitalOceanDroplets(meta)
		if err != nil {
			return err
		}

		droplet, err := findDropletByTag(dropletList, v.(string))
		if err != nil {
			return err
		}

		foundDroplet = *droplet
	} else if v, ok := d.GetOk("name"); ok {
		dropletList, err := getDigitalOceanDroplets(meta)
		if err != nil {
			return err
		}

		droplet, err := findDropletByName(dropletList, v.(string))

		if err != nil {
			return err
		}

		foundDroplet = *droplet
	} else {
		return fmt.Errorf("Error: specify either a name, tag, or id to use to look up the droplet")
	}

	flattenedDroplet, err := flattenDigitalOceanDroplet(foundDroplet, meta)
	if err != nil {
		return err
	}

	if err := setResourceDataFromMap(d, flattenedDroplet); err != nil {
		return err
	}

	d.SetId(strconv.Itoa(foundDroplet.ID))
	return nil
}

func findDropletByName(droplets []interface{}, name string) (*godo.Droplet, error) {
	results := make([]godo.Droplet, 0)
	for _, v := range droplets {
		droplet := v.(godo.Droplet)
		if droplet.Name == name {
			results = append(results, droplet)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no droplet found with name %s", name)
	}
	return nil, fmt.Errorf("too many droplets found with name %s (found %d, expected 1)", name, len(results))
}

func findDropletByTag(droplets []interface{}, tag string) (*godo.Droplet, error) {
	results := make([]godo.Droplet, 0)
	for _, d := range droplets {
		droplet := d.(godo.Droplet)
		for _, t := range droplet.Tags {
			if t == tag {
				results = append(results, droplet)
			}
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no droplet found with tag %s", tag)
	}
	return nil, fmt.Errorf("too many droplets found with tag %s (found %d, expected 1)", tag, len(results))
}
