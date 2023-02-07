package tag

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanTags() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"total_resource_count": {
				Type: schema.TypeInt,
			},
			"droplets_count": {
				Type: schema.TypeInt,
			},
			"images_count": {
				Type: schema.TypeInt,
			},
			"volumes_count": {
				Type: schema.TypeInt,
			},
			"volume_snapshots_count": {
				Type: schema.TypeInt,
			},
			"databases_count": {
				Type: schema.TypeInt,
			},
		},
		ResultAttributeName: "tags",
		FlattenRecord:       flattenDigitalOceanTag,
		GetRecords:          getDigitalOceanTags,
	}

	return datalist.NewResource(dataListConfig)
}

func getDigitalOceanTags(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	tagsList := []interface{}{}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		tags, resp, err := client.Tags.List(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving tags: %s", err)
		}

		for _, tag := range tags {
			tagsList = append(tagsList, tag)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving tags: %s", err)
		}

		opts.Page = page + 1
	}

	return tagsList, nil
}

func flattenDigitalOceanTag(tag, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	t := tag.(godo.Tag)

	flattenedTag := map[string]interface{}{}
	flattenedTag["name"] = t.Name
	flattenedTag["total_resource_count"] = t.Resources.Count
	flattenedTag["droplets_count"] = t.Resources.Droplets.Count
	flattenedTag["images_count"] = t.Resources.Images.Count
	flattenedTag["volumes_count"] = t.Resources.Volumes.Count
	flattenedTag["volume_snapshots_count"] = t.Resources.VolumeSnapshots.Count
	flattenedTag["databases_count"] = t.Resources.Databases.Count

	return flattenedTag, nil
}
