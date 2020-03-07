package digitalocean

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type imageListFunc func(ctx context.Context, opt *godo.ListOptions) ([]godo.Image, *godo.Response, error)

func imageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Description: "id of the image",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "name of the image",
		},
		"type": {
			Type:        schema.TypeString,
			Description: "type of the image",
		},
		"distribution": {
			Type:        schema.TypeString,
			Description: "distribution of the OS of the image",
		},
		"slug": {
			Type:        schema.TypeString,
			Description: "slug of the image",
		},
		"image": {
			Type:        schema.TypeString,
			Description: "slug or id of the image",
		},
		"private": {
			Type:        schema.TypeBool,
			Description: "Is the image private or non-private",
		},
		"min_disk_size": {
			Type:        schema.TypeInt,
			Description: "minimum disk size required by the image",
		},
		"size_gigabytes": {
			Type:        schema.TypeFloat,
			Description: "size in GB of the image",
		},
		"regions": {
			Type:        schema.TypeSet,
			Description: "list of the regions that the image is available in",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"created": {
			Type: schema.TypeString,
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "tags applied to the image",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"status": {
			Type:        schema.TypeString,
			Description: "status of the image",
		},
		"error_message": {
			Type:        schema.TypeString,
			Description: "error message associated with the image",
		},
	}
}

func getDigitalOceanImages(meta interface{}) ([]interface{}, error) {
	client := meta.(*CombinedConfig).godoClient()
	return listDigitalOceanImages(client.Images.List)
}

func listDigitalOceanImages(listImages imageListFunc) ([]interface{}, error) {
	var allImages []interface{}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		images, resp, err := listImages(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving images: %s", err)
		}

		for _, image := range images {
			allImages = append(allImages, image)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving images: %s", err)
		}

		opts.Page = page + 1
	}

	return allImages, nil
}

func flattenDigitalOceanImage(rawImage interface{}, meta interface{}) (map[string]interface{}, error) {
	image, ok := rawImage.(godo.Image)
	if !ok {
		return nil, fmt.Errorf("Unable to convert to godo.Image")
	}

	flattenedRegions := schema.NewSet(schema.HashString, []interface{}{})
	for _, region := range image.Regions {
		flattenedRegions.Add(region)
	}

	flattenedTags := schema.NewSet(schema.HashString, []interface{}{})
	for _, tag := range image.Tags {
		flattenedTags.Add(tag)
	}

	flattenedImage := map[string]interface{}{
		"id":             image.ID,
		"name":           image.Name,
		"type":           image.Type,
		"distribution":   image.Distribution,
		"slug":           image.Slug,
		"private":        !image.Public,
		"min_disk_size":  image.MinDiskSize,
		"size_gigabytes": image.SizeGigaBytes,
		"created":        image.Created,
		"regions":        flattenedRegions,
		"tags":           flattenedTags,
		"status":         image.Status,
		"error_message":  image.ErrorMessage,

		// Legacy attributes
		"image": strconv.Itoa(image.ID),
	}

	return flattenedImage, nil
}
