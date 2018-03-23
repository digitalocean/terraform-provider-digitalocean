package digitalocean

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/godo/context"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDigitalOceanImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanImageRead,
		Schema: map[string]*schema.Schema{
			"source": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "images source (user, application or distribution)",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of the image",
			},
			"name_regex": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
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

func dataSourceDigitalOceanImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts := &godo.ListOptions{PerPage: 200}

	name, nameOk := d.GetOk("name")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if !nameOk && !nameRegexOk {
		return fmt.Errorf("One of name or name_regex must be assigned")
	}

	var match func(string) bool

	if nameRegexOk {
		match = regexp.MustCompile(nameRegex.(string)).MatchString
	} else {
		match = func(s string) bool {
			return s == name.(string)
		}
	}

	source := d.Get("source").(string)

	var listFn func(ctx context.Context, opt *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	switch source {
	case "user":
		listFn = client.Images.ListUser
	case "application":
		listFn = client.Images.ListApplication
	case "distribution":
		listFn = client.Images.ListDistribution
	default:
		return fmt.Errorf("source must be one of user, application, or distribution")
	}

	var images []*godo.Image

	for {
		images, res, err := listFn(context.Background(), opts)
		if err != nil {
			d.SetId("")
			return err
		}

		for _, image := range images {
			if match(image.Name) {
				images = append(images, image)
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

	if len(images) == 0 {
		return fmt.Errorf("no %s image found with name %s", source, name)
	}
	if len(images) > 1 {
		return fmt.Errorf("too many %s images found with name %s (found %d, expected 1)", source, name, len(images))
	}

	d.SetId(images[0].Name)
	d.Set("name", images[0].Name)
	d.Set("image", strconv.Itoa(images[0].ID))
	d.Set("min_disk_size", images[0].MinDiskSize)
	d.Set("private", !images[0].Public)
	d.Set("regions", images[0].Regions)
	d.Set("type", images[0].Type)

	return nil
}
