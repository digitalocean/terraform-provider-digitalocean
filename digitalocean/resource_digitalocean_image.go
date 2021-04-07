package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Ref: https://developers.digitalocean.com/documentation/v2/#retrieve-an-existing-image-by-id
const (
	imageAvailableStatus = "available"
	imageDeletedStatus   = "deleted"
)

func resourceDigitalOceanCustomImage() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDigitalOceanCustomImageRead,
		CreateContext: resourceDigitalOceanCustomImageCreate,
		UpdateContext: resourceDigitalOceanCustomImageUpdate,
		DeleteContext: resourceDigitalOceanCustomImageDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"regions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"distribution": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Unknown",
				ValidateFunc: validation.StringInSlice(validImageDistributions(), false),
			},
			"tags": tagsSchema(),
			"image_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"min_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size_gigabytes": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanCustomImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	// TODO: Support multiple regions
	regions := d.Get("regions").([]interface{})
	region := regions[0].(string)

	imageCreateRequest := godo.CustomImageCreateRequest{
		Name:   d.Get("name").(string),
		Url:    d.Get("url").(string),
		Region: region,
	}

	if desc, ok := d.GetOk("description"); ok {
		imageCreateRequest.Description = desc.(string)
	}

	if dist, ok := d.GetOk("distribution"); ok {
		imageCreateRequest.Distribution = dist.(string)
	}

	if tags, ok := d.GetOk("tags"); ok {
		imageCreateRequest.Tags = expandTags(tags.(*schema.Set).List())
	}

	imageResponse, _, err := client.Images.Create(ctx, &imageCreateRequest)
	if err != nil {
		return diag.Errorf("Error creating custom image: %s", err)
	}
	id := strconv.Itoa(imageResponse.ID)
	d.SetId(id)
	_, err = waitForImage(ctx, d, imageAvailableStatus, imagePendingStatuses(), "status", meta)
	if err != nil {
		return diag.Errorf(
			"Error waiting for image (%s) to become ready: %s", d.Id(), err)
	}
	return resourceDigitalOceanCustomImageRead(ctx, d, meta)
}

func resourceDigitalOceanCustomImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	imageID := d.Id()

	id, err := strconv.Atoi(imageID)
	if err != nil {
		return diag.Errorf("Error converting id %s to string: %s", imageID, err)
	}

	imageResponse, _, err := client.Images.GetByID(ctx, id)
	if err != nil {
		return diag.Errorf("Error retrieving image with id %s: %s", imageID, err)
	}
	// Set status as deleted if image is deleted
	if imageResponse.Status == imageDeletedStatus {
		d.SetId("")
		return nil
	}
	d.Set("image_id", imageResponse.ID)
	d.Set("name", imageResponse.Name)
	d.Set("type", imageResponse.Type)
	d.Set("distribution", imageResponse.Distribution)
	d.Set("slug", imageResponse.Slug)
	d.Set("public", imageResponse.Public)
	d.Set("regions", imageResponse.Regions)
	d.Set("min_disk_size", imageResponse.MinDiskSize)
	d.Set("size_gigabytes", imageResponse.SizeGigaBytes)
	d.Set("created_at", imageResponse.Created)
	d.Set("description", imageResponse.Description)
	if err := d.Set("tags", flattenTags(imageResponse.Tags)); err != nil {
		return diag.Errorf("Error setting `tags`: %+v", err)
	}
	d.Set("status", imageResponse.Status)
	return nil
}

func resourceDigitalOceanCustomImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	imageID := d.Id()

	id, err := strconv.Atoi(imageID)
	if err != nil {
		return diag.Errorf("Error converting id %s to string: %s", imageID, err)
	}

	if d.HasChanges("name", "description", "distribution") {
		imageName := d.Get("name").(string)
		imageUpdateRequest := &godo.ImageUpdateRequest{
			Name:         imageName,
			Distribution: d.Get("distribution").(string),
			Description:  d.Get("description").(string),
		}

		_, _, err := client.Images.Update(ctx, id, imageUpdateRequest)
		if err != nil {
			return diag.Errorf("Error updating image %s, name %s: %s", imageID, imageName, err)
		}
	}

	return resourceDigitalOceanCustomImageRead(ctx, d, meta)
}

func resourceDigitalOceanCustomImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	imageID := d.Id()

	id, err := strconv.Atoi(imageID)
	if err != nil {
		return diag.Errorf("Error converting id %s to string: %s", imageID, err)
	}
	_, err = client.Images.Delete(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "Image Can not delete an already deleted image.") {
			log.Printf("[INFO] Image %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting image id %s: %s", imageID, err)
	}
	return nil
}

func waitForImage(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for image (%s) to have %s of %s", d.Id(), attribute, target)
	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    imageStateRefreshFunc(ctx, d, attribute, meta),
		Timeout:    120 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func imageStateRefreshFunc(ctx context.Context, d *schema.ResourceData, state string, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*CombinedConfig).godoClient()

		imageID := d.Id()

		id, err := strconv.Atoi(imageID)
		if err != nil {
			return nil, "", err
		}

		imageResponse, _, err := client.Images.GetByID(ctx, id)
		if err != nil {
			return nil, "", err
		}

		if imageResponse.Status == imageDeletedStatus {
			return nil, "", fmt.Errorf(imageResponse.ErrorMessage)
		}

		return imageResponse, imageResponse.Status, nil
	}
}

// Ref: https://developers.digitalocean.com/documentation/v2/#retrieve-an-existing-image-by-id
func imagePendingStatuses() []string {
	return []string{"new", "pending"}
}

// Ref:https://developers.digitalocean.com/documentation/v2/#create-a-custom-image
func validImageDistributions() []string {
	return []string{
		"Arch Linux",
		"CentOS",
		"CoreOS",
		"Debian",
		"Fedora",
		"Fedora Atomic",
		"FreeBSD",
		"Gentoo",
		"openSUSE",
		"RancherOS",
		"Ubuntu",
		"Unknown",
	}
}
