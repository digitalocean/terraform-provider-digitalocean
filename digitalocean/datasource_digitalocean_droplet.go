package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceDigitalOceanDroplet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanDropletRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the droplet",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the Droplet",
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the Droplet",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the region that the droplet instance is deployed in",
			},
			"image": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the image id or slug of the Droplet",
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the current size of the Droplet",
			},
			"disk": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the size of the droplets disk in gigabytes",
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the number of virtual cpus",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "memory of the droplet in megabytes",
			},
			"price_hourly": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "the droplets hourly price",
			},
			"price_monthly": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "the droplets monthly price",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the droplet instance",
			},
			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the droplet has been locked",
			},
			"ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the droplets public ipv4 address",
			},
			"ipv4_address_private": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the droplets private ipv4 address",
			},
			"ipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the droplets public ipv6 address",
			},
			"ipv6_address_private": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the droplets private ipv4 address",
			},
			"backups": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the droplet has backups enabled",
			},
			"ipv6": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the droplet has ipv6 enabled",
			},
			"private_networking": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the droplet has private networking enabled",
			},
			"monitoring": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the droplet has monitoring enabled",
			},
			"volume_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "list of volumes attached to the droplet",
			},

			"tags": tagsDataSourceSchema(),
		},
	}
}

func dataSourceDigitalOceanDropletRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	name := d.Get("name").(string)

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	dropletList := []godo.Droplet{}

	for {
		droplets, resp, err := client.Droplets.List(context.Background(), opts)

		if err != nil {
			return fmt.Errorf("Error retrieving droplets: %s", err)
		}

		for _, droplet := range droplets {
			dropletList = append(dropletList, droplet)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return fmt.Errorf("Error retrieving droplets: %s", err)
		}

		opts.Page = page + 1
	}

	droplet, err := findDropletByName(dropletList, name)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(droplet.ID))
	d.Set("name", droplet.Name)
	d.Set("urn", droplet.URN())
	d.Set("region", droplet.Region.Slug)
	d.Set("size", droplet.Size.Slug)
	d.Set("price_hourly", droplet.Size.PriceHourly)
	d.Set("price_monthly", droplet.Size.PriceMonthly)
	d.Set("disk", droplet.Disk)
	d.Set("vcpus", droplet.Vcpus)
	d.Set("memory", droplet.Memory)
	d.Set("status", droplet.Status)
	d.Set("locked", droplet.Locked)
	d.Set("created_at", droplet.Created)

	if droplet.Image.Slug == "" {
		d.Set("image", droplet.Image.ID)
	} else {
		d.Set("image", droplet.Image.Slug)
	}

	if publicIPv4 := findIPv4AddrByType(droplet, "public"); publicIPv4 != "" {
		d.Set("ipv4_address", publicIPv4)
	}

	if privateIPv4 := findIPv4AddrByType(droplet, "private"); privateIPv4 != "" {
		d.Set("ipv4_address_private", privateIPv4)
	}

	if publicIPv6 := findIPv6AddrByType(droplet, "public"); publicIPv6 != "" {
		d.Set("ipv6_address", strings.ToLower(publicIPv6))
	}

	if privateIPv6 := findIPv6AddrByType(droplet, "private"); privateIPv6 != "" {
		d.Set("ipv6_address_private", strings.ToLower(privateIPv6))
	}

	if features := droplet.Features; features != nil {
		d.Set("backups", containsDigitalOceanDropletFeature(features, "backups"))
		d.Set("ipv6", containsDigitalOceanDropletFeature(features, "ipv6"))
		d.Set("private_networking", containsDigitalOceanDropletFeature(features, "private_networking"))
		d.Set("monitoring", containsDigitalOceanDropletFeature(features, "monitoring"))
	}

	if err := d.Set("volume_ids", flattenDigitalOceanDropletVolumeIds(droplet.VolumeIDs)); err != nil {
		return fmt.Errorf("Error setting `volume_ids`: %+v", err)
	}

	if err := d.Set("tags", flattenTags(droplet.Tags)); err != nil {
		return fmt.Errorf("Error setting `tags`: %+v", err)
	}

	return nil
}

func findDropletByName(droplets []godo.Droplet, name string) (*godo.Droplet, error) {
	results := make([]godo.Droplet, 0)
	for _, v := range droplets {
		if v.Name == name {
			results = append(results, v)
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
