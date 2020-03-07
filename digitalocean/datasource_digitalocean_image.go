package digitalocean

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceDigitalOceanImage() *schema.Resource {
	recordSchema := imageSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["id"].Optional = true
	recordSchema["id"].ValidateFunc = validation.NoZeroValues
	recordSchema["id"].ExactlyOneOf = []string{"id", "slug", "name"}

	recordSchema["name"].Optional = true
	recordSchema["name"].ValidateFunc = validation.StringIsNotEmpty
	recordSchema["name"].ExactlyOneOf = []string{"id", "slug", "name"}

	recordSchema["slug"].Optional = true
	recordSchema["slug"].ValidateFunc = validation.StringIsNotEmpty
	recordSchema["slug"].ExactlyOneOf = []string{"id", "slug", "name"}

	recordSchema["private"].Optional = true
	recordSchema["private"].ConflictsWith = []string{"id", "slug"}

	return &schema.Resource{
		Read:   dataSourceDigitalOceanImageRead,
		Schema: recordSchema,
	}
}

func dataSourceDigitalOceanImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	var foundImage *godo.Image

	if id, ok := d.GetOk("id"); ok {
		image, resp, err := client.Images.GetByID(context.Background(), id.(int))
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return fmt.Errorf("image ID %d not found: %s", id.(int), err)
			}
			return fmt.Errorf("Error retrieving image: %s", err)
		}
		foundImage = image
	} else if slug, ok := d.GetOk("slug"); ok {
		image, resp, err := client.Images.GetBySlug(context.Background(), slug.(string))
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return fmt.Errorf("image not found: %s", err)
			}
			return fmt.Errorf("Error retrieving image: %s", err)
		}
		foundImage = image
	} else if name, ok := d.GetOk("name"); ok {
		private := d.Get("private")

		var allImages []interface{}
		if private.(bool) {
			images, err := listDigitalOceanImages(client.Images.ListUser)
			if err != nil {
				return err
			}
			allImages = images
		} else {
			images, err := listDigitalOceanImages(client.Images.List)
			if err != nil {
				return err
			}
			allImages = images
		}

		var results []interface{}

		for _, image := range allImages {
			if image.(godo.Image).Name == name {
				results = append(results, image)
			}
		}

		if len(results) == 0 {
			return fmt.Errorf("no image found with name %s", name)
		} else if len(results) > 1 {
			return fmt.Errorf("too many images found with name %s (found %d, expected 1)", name, len(results))
		}

		result := results[0].(godo.Image)
		foundImage = &result
	} else {
		return fmt.Errorf("Illegal state: one of id, name, or slug must be set")
	}

	flattenedImage, err := flattenDigitalOceanImage(*foundImage, meta)
	if err != nil {
		return err
	}

	if err := setResourceDataFromMap(d, flattenedImage); err != nil {
		return err
	}

	d.SetId(strconv.Itoa(foundImage.ID))

	return nil
}
