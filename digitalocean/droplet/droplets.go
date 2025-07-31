package droplet

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dropletSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Description: "id of the Droplet",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "name of the Droplet",
		},
		"project_id": {
			Type:        schema.TypeString,
			Description: "ID of the project to which the Droplet belongs",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "the creation date for the Droplet",
		},
		"urn": {
			Type:        schema.TypeString,
			Description: "the uniform resource name for the Droplet",
		},
		"region": {
			Type:        schema.TypeString,
			Description: "the region that the Droplet instance is deployed in",
		},
		"image": {
			Type:        schema.TypeString,
			Description: "the image id or slug of the Droplet",
		},
		"size": {
			Type:        schema.TypeString,
			Description: "the current size of the Droplet",
		},
		"disk": {
			Type:        schema.TypeInt,
			Description: "the size of the Droplets disk in gigabytes",
		},
		"vcpus": {
			Type:        schema.TypeInt,
			Description: "the number of virtual cpus",
		},
		"memory": {
			Type:        schema.TypeInt,
			Description: "memory of the Droplet in megabytes",
		},
		"price_hourly": {
			Type:        schema.TypeFloat,
			Description: "the Droplets hourly price",
		},
		"price_monthly": {
			Type:        schema.TypeFloat,
			Description: "the Droplets monthly price",
		},
		"status": {
			Type:        schema.TypeString,
			Description: "state of the Droplet instance",
		},
		"locked": {
			Type:        schema.TypeBool,
			Description: "whether the Droplet has been locked",
		},
		"ipv4_address": {
			Type:        schema.TypeString,
			Description: "the Droplets public ipv4 address",
		},
		"ipv4_address_private": {
			Type:        schema.TypeString,
			Description: "the Droplets private ipv4 address",
		},
		"ipv6_address": {
			Type:        schema.TypeString,
			Description: "the Droplets public ipv6 address",
		},
		"ipv6_address_private": {
			Type:        schema.TypeString,
			Description: "the Droplets private ipv4 address",
		},
		"backups": {
			Type:        schema.TypeBool,
			Description: "whether the Droplet has backups enabled",
		},
		"ipv6": {
			Type:        schema.TypeBool,
			Description: "whether the Droplet has ipv6 enabled",
		},
		"private_networking": {
			Type:        schema.TypeBool,
			Description: "whether the Droplet has private networking enabled",
		},
		"monitoring": {
			Type:        schema.TypeBool,
			Description: "whether the Droplet has monitoring enabled",
		},
		"volume_ids": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "list of volumes attached to the Droplet",
		},
		"tags": tag.TagsDataSourceSchema(),
		"vpc_uuid": {
			Type:        schema.TypeString,
			Description: "UUID of the VPC in which the Droplet is located",
		},
	}
}

func getDigitalOceanDroplets(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	gpus, _ := extra["gpus"].(bool)

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var dropletList []interface{}

	for {
		var (
			droplets []godo.Droplet
			resp     *godo.Response
			err      error
		)
		if gpus {
			droplets, resp, err = client.Droplets.ListWithGPUs(context.Background(), opts)
		} else {
			droplets, resp, err = client.Droplets.List(context.Background(), opts)
		}

		if err != nil {
			return nil, fmt.Errorf("Error retrieving droplets: %s", err)
		}

		for _, droplet := range droplets {
			dropletList = append(dropletList, droplet)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving droplets: %s", err)
		}

		opts.Page = page + 1
	}

	return dropletList, nil
}

func flattenDigitalOceanDroplet(rawDroplet, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	droplet := rawDroplet.(godo.Droplet)

	flattenedDroplet := map[string]interface{}{
		"id":            droplet.ID,
		"project_id":    droplet.ProjectID,
		"name":          droplet.Name,
		"urn":           droplet.URN(),
		"region":        droplet.Region.Slug,
		"size":          droplet.Size.Slug,
		"price_hourly":  droplet.Size.PriceHourly,
		"price_monthly": droplet.Size.PriceMonthly,
		"disk":          droplet.Disk,
		"vcpus":         droplet.Vcpus,
		"memory":        droplet.Memory,
		"status":        droplet.Status,
		"locked":        droplet.Locked,
		"created_at":    droplet.Created,
		"vpc_uuid":      droplet.VPCUUID,
	}

	if droplet.Image.Slug == "" {
		flattenedDroplet["image"] = strconv.Itoa(droplet.Image.ID)
	} else {
		flattenedDroplet["image"] = droplet.Image.Slug
	}

	if publicIPv4 := FindIPv4AddrByType(&droplet, "public"); publicIPv4 != "" {
		flattenedDroplet["ipv4_address"] = publicIPv4
	}

	if privateIPv4 := FindIPv4AddrByType(&droplet, "private"); privateIPv4 != "" {
		flattenedDroplet["ipv4_address_private"] = privateIPv4
	}

	if publicIPv6 := FindIPv6AddrByType(&droplet, "public"); publicIPv6 != "" {
		flattenedDroplet["ipv6_address"] = strings.ToLower(publicIPv6)
	}

	if privateIPv6 := FindIPv6AddrByType(&droplet, "private"); privateIPv6 != "" {
		flattenedDroplet["ipv6_address_private"] = strings.ToLower(privateIPv6)
	}

	if features := droplet.Features; features != nil {
		flattenedDroplet["backups"] = containsDigitalOceanDropletFeature(features, "backups")
		flattenedDroplet["ipv6"] = containsDigitalOceanDropletFeature(features, "ipv6")
		flattenedDroplet["private_networking"] = containsDigitalOceanDropletFeature(features, "private_networking")
		flattenedDroplet["monitoring"] = containsDigitalOceanDropletFeature(features, "monitoring")
	}

	flattenedDroplet["volume_ids"] = flattenDigitalOceanDropletVolumeIds(droplet.VolumeIDs)

	flattenedDroplet["tags"] = tag.FlattenTags(droplet.Tags)

	return flattenedDroplet, nil
}
