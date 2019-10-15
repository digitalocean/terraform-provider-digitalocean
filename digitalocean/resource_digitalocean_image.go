package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanImageCreate,
		Read:   resourceDigitalOceanImageRead,
		Update: resourceDigitalOceanImageUpdate,
		Delete: resourceDigitalOceanImageDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanImageImport,
		},
		SchemaVersion: 1,

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

func resourceDigitalOceanImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	// Build up our creation options
	opts := &godo.CustomImageCreateRequest{
		Name:         d.Get("name").(string),
		Url:          d.Get("url").(string),
		Region:       d.Get("region").(string),
		Distribution: d.Get("distribution").(string),
		Description:  d.Get("description").(string),
		Tags:         expandTags(d.Get("tags").(*schema.Set).List()),
	}

	log.Printf("[DEBUG] Image create configuration: %#v", opts)

	image, _, err := client.Images.Create(context.Background(), opts)

	if err != nil {
		return fmt.Errorf("Error creating image: %s", err)
	}

	// Assign the image id
	d.SetId(strconv.Itoa(image.ID))

	log.Printf("[INFO] Droplet ID: %s", d.Id())

	_, err = waitForDropletAttribute(d, "active", []string{"new"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for image (%s) to become ready: %s", d.Id(), err)
	}
	return resourceDigitalOceanImageRead(d, meta)
}

func resourceDigitalOceanImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid image id: %v", err)
	}

	// Retrieve the image properties for updating the state
	image, resp, err := client.Images.GetByID(context.Background(), id)
	if err != nil {
		// check if the image no longer exists.
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Image (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving image: %s", err)
	}

	// Image can drift once the image is build if a remote drift is detected
	// as can cause issues with slug changes due image patch that shoudn't be sync.
	// See: https://github.com/terraform-providers/terraform-provider-digitalocean/issues/152

	d.Set("name", image.Name)
	d.Set("slug", image.Slug)
	d.Set("image", strconv.Itoa(image.ID))
	d.Set("distribution", image.Distribution)
	d.Set("min_disk_size", image.MinDiskSize)
	d.Set("private", !image.Public)
	d.Set("regions", image.Regions)
	d.Set("type", image.Type)

	if err := d.Set("tags", flattenTags(image.Tags)); err != nil {
		return fmt.Errorf("Error setting `tags`: %+v", err)
	}

	return nil
}

func resourceDigitalOceanImageImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Retrieve the image from API during import
	client := meta.(*CombinedConfig).godoClient()
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("Invalid image id: %v", err)
	}

	image, resp, err := client.Images.GetByID(context.Background(), id)
	if resp.StatusCode != 404 {
		if err != nil {
			return nil, fmt.Errorf("Error importing droplet: %s", err)
		}

		d.Set("slug", image.Slug)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDigitalOceanImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid image id: %v", err)
	}

	d.Set("id", id)

	if d.HasChange("tags") {
		err = setTags(client, d)
		if err != nil {
			return fmt.Errorf("Error updating tags: %s", err)
		}
	}

	return resourceDigitalOceanImageRead(d, meta)
}

func resourceDigitalOceanImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid image id: %v", err)
	}

	_, err = waitForImageAttribute(
		d, "false", []string{"", "true"}, "locked", meta)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for image to be unlocked for destroy (%s): %s", d.Id(), err)
	}

	log.Printf("[INFO] Deleting image: %s", d.Id())

	// Destroy the image
	resp, err := client.Images.Delete(context.Background(), id)

	// Handle already destroyed images
	if err != nil && resp.StatusCode == 404 {
		return nil
	}

	_, err = waitForImageDestroy(d, meta)
	if err != nil && strings.Contains(err.Error(), "404") {
		return nil
	} else if err != nil {
		return fmt.Errorf("Error deleting image: %s", err)
	}

	return nil
}

func waitForImageDestroy(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for image (%s) to be destroyed", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "off"},
		Target:     []string{"archived"},
		Refresh:    newImageStateRefreshFunc(d, "status", meta),
		Timeout:    60 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForImageAttribute(
	d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	// Wait for the image so we can get the networking attributes
	// that show up after a while
	log.Printf(
		"[INFO] Waiting for image (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newImageStateRefreshFunc(d, attribute, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		// This is a hack around DO API strangeness.
		// https://github.com/hashicorp/terraform/issues/481
		//
		NotFoundChecks: 60,
	}

	return stateConf.WaitForState()
}

// TODO This function still needs a little more refactoring to make it
// cleaner and more efficient
func newImageStateRefreshFunc(
	d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).godoClient()
	return func() (interface{}, string, error) {
		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, "", err
		}

		err = resourceDigitalOceanImageRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		// If the image is locked, continue waiting. We can
		// only perform actions on unlocked images, so it's
		// pointless to look at that status
		if d.Get("locked").(bool) {
			log.Println("[DEBUG] Image is locked, skipping status check and retrying")
			return nil, "", nil
		}

		// See if we can access our attribute
		if attr, ok := d.GetOkExists(attribute); ok {
			// Retrieve the image properties
			image, _, err := client.Images.GetByID(context.Background(), id)
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving image: %s", err)
			}

			switch attr.(type) {
			case bool:
				return &image, strconv.FormatBool(attr.(bool)), nil
			default:
				return &image, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}
