package digitalocean

import (
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/godo/context"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDigitalOceanImage(source string) *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return dataSourceDigitalOceanImageRead(source, d, meta)
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the image",
			},
			// computed attributes
			"image": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "slug or id of the image",
			},
			"min_disk_size": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "minimum disk size required by the image",
			},
			"private": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the image private or non-private",
			},
			"regions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of the regions that the image is available in",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of the image",
			},
		},
	}
}

func dataSourceDigitalOceanImageRead(source string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts := &godo.ListOptions{}

	name := d.Get("name").(string)

	var listFn func(ctx context.Context, opt *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	switch source {
	case "user":
		listFn = client.Images.ListUser
	case "application":
		listFn = client.Images.ListApplication
	case "distribution":
		listFn = client.Images.ListDistribution
	}

	var image *godo.Image

outer:
	for {
		// This method should be configurable
		images, res, err := listFn(context.Background(), opts)
		if err != nil {
			d.SetId("")
			return err
		}

		for _, v := range images {
			if v.Name == name {
				image = &v
				break outer
			}
		}

		// if we are at the last page, break out the for loop
		if res.Links == nil || res.Links.IsLastPage() {
			break
		}

		page, err := res.Links.CurrentPage()
		if err != nil {
			return fmt.Errorf("cannot retrieve current page: %v", err)
		}

		// set the page we want for the next request
		opts.Page = page + 1
	}

	if image == nil {
		return fmt.Errorf("no %s image found with name %s", source, name)
	}

	d.SetId(image.Name)
	d.Set("name", image.Name)
	d.Set("image", strconv.Itoa(image.ID))
	d.Set("min_disk_size", image.MinDiskSize)
	d.Set("private", !image.Public)
	d.Set("regions", image.Regions)
	d.Set("type", image.Type)

	return nil
}
