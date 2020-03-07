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

	recordSchema["name"].Optional = true
	recordSchema["name"].ValidateFunc = validation.StringIsNotEmpty
	recordSchema["name"].ExactlyOneOf = []string{"slug", "name"}

	recordSchema["slug"].Optional = true
	recordSchema["slug"].ValidateFunc = validation.StringIsNotEmpty
	recordSchema["slug"].ExactlyOneOf = []string{"slug", "name"}

	return &schema.Resource{
		Read:   dataSourceDigitalOceanImageRead,
		Schema: recordSchema,
	}
}

func dataSourceDigitalOceanImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	name, hasName := d.GetOk("name")
	slug, hasSlug := d.GetOk("slug")

	if !hasName && !hasSlug {
		return fmt.Errorf("One of `name` or `slug` must be assigned")
	}

	var image *godo.Image

	if hasName {
		opts := &godo.ListOptions{
			Page:    1,
			PerPage: 200,
		}

		imageList := []godo.Image{}

		for {
			images, resp, err := client.Images.ListUser(context.Background(), opts)

			if err != nil {
				return fmt.Errorf("Error retrieving images: %s", err)
			}

			for _, image := range images {
				imageList = append(imageList, image)
			}

			if resp.Links == nil || resp.Links.IsLastPage() {
				break
			}

			page, err := resp.Links.CurrentPage()
			if err != nil {
				return fmt.Errorf("Error retrieving images: %s", err)
			}

			opts.Page = page + 1
		}

		var err error
		image, err = findImageByName(imageList, name.(string))

		if err != nil {
			return err
		}
	} else {
		var (
			err  error
			resp *godo.Response
		)

		image, resp, err = client.Images.GetBySlug(context.Background(), slug.(string))
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return fmt.Errorf("image not found: %s", err)
			}
			return fmt.Errorf("Error retrieving image: %s", err)
		}
	}

	flattenedImage, err := flattenDigitalOceanImage(*image, meta)
	if err != nil {
		return err
	}

	if err := setResourceDataFromMap(d, flattenedImage); err != nil {
		return err
	}

	d.SetId(strconv.Itoa(image.ID))

	return nil
}

func findImageByName(images []godo.Image, name string) (*godo.Image, error) {
	results := make([]godo.Image, 0)
	for _, v := range images {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no user image found with name %s", name)
	}
	return nil, fmt.Errorf("too many user images found with name %s (found %d, expected 1)", name, len(results))
}
