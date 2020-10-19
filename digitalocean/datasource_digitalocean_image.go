package digitalocean

import (
	"context"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

	recordSchema["source"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Default:       "user",
		ValidateFunc:  validation.StringInSlice([]string{"all", "applications", "distributions", "user"}, true),
		ConflictsWith: []string{"id", "slug"},
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanImageRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	var foundImage *godo.Image

	if id, ok := d.GetOk("id"); ok {
		image, resp, err := client.Images.GetByID(context.Background(), id.(int))
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return diag.Errorf("image ID %d not found: %s", id.(int), err)
			}
			return diag.Errorf("Error retrieving image: %s", err)
		}
		foundImage = image
	} else if slug, ok := d.GetOk("slug"); ok {
		image, resp, err := client.Images.GetBySlug(context.Background(), slug.(string))
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return diag.Errorf("image not found: %s", err)
			}
			return diag.Errorf("Error retrieving image: %s", err)
		}
		foundImage = image
	} else if name, ok := d.GetOk("name"); ok {
		source := strings.ToLower(d.Get("source").(string))

		var listImages imageListFunc
		if source == "all" {
			listImages = client.Images.List
		} else if source == "distributions" {
			listImages = client.Images.ListDistribution
		} else if source == "applications" {
			listImages = client.Images.ListApplication
		} else if source == "user" {
			listImages = client.Images.ListUser
		} else {
			return diag.Errorf("Illegal state: source=%s", source)
		}

		images, err := listDigitalOceanImages(listImages)
		if err != nil {
			return diag.FromErr(err)
		}

		var results []interface{}

		for _, image := range images {
			if image.(godo.Image).Name == name {
				results = append(results, image)
			}
		}

		if len(results) == 0 {
			return diag.Errorf("no image found with name %s", name)
		} else if len(results) > 1 {
			return diag.Errorf("too many images found with name %s (found %d, expected 1)", name, len(results))
		}

		result := results[0].(godo.Image)
		foundImage = &result
	} else {
		return diag.Errorf("Illegal state: one of id, name, or slug must be set")
	}

	flattenedImage, err := flattenDigitalOceanImage(*foundImage, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setResourceDataFromMap(d, flattenedImage); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(foundImage.ID))

	return nil
}
