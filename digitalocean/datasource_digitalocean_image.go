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
	return &schema.Resource{
		Read: dataSourceDigitalOceanImageRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "name of the image",
				ValidateFunc:  validation.NoZeroValues,
				ConflictsWith: []string{"slug"},
			},
			"slug": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "slug of the image",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"image": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "slug or id of the image",
			},
			"distribution": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "distribution of the OS of the image",
			},
			"private": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the image private or non-private",
			},
			"min_disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "minimum disk size required by the image",
			},
			"regions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "list of the regions that the image is available in",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of the image",
			},
		},
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

	d.SetId(strconv.Itoa(image.ID))
	d.Set("name", image.Name)
	d.Set("slug", image.Slug)
	d.Set("image", strconv.Itoa(image.ID))
	d.Set("distribution", image.Distribution)
	d.Set("min_disk_size", image.MinDiskSize)
	d.Set("private", !image.Public)
	d.Set("regions", image.Regions)
	d.Set("type", image.Type)

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
