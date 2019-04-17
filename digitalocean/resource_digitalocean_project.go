package digitalocean

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"strings"
)

func resourceDigitalOceanProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanProjectCreate,
		Read:   resourceDigitalOceanProjectRead,
		Update: resourceDigitalOceanProjectUpdate,
		Delete: resourceDigitalOceanProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the human-readable name for the project",
				ValidateFunc: validation.All(
					validation.NoZeroValues,
					validation.StringLenBetween(1, 175),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "the descirption of the project",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"purpose": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Web Application",
				Description:  "the purpose of the project",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"environment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Development",
				Description:  "the environment of the project's resources",
				ValidateFunc: validation.StringInSlice([]string{"development", "staging", "production"}, true),
			},
			"owner_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the unique universal identifier of the project owner.",
			},
			"owner_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the id of the project owner.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the date and time when the project was created, (ISO8601)",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the date and time when the project was last updated, (ISO8601)",
			},
		},
	}
}

func resourceDigitalOceanProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectRequest := &godo.CreateProjectRequest{
		Name:    d.Get("name").(string),
		Purpose: d.Get("purpose").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		projectRequest.Description = v.(string)
	}

	if v, ok := d.GetOk("environment"); ok {
		projectRequest.Environment = v.(string)
	}

	log.Printf("[DEBUG] Project create request: %#v", projectRequest)
	project, _, err := client.Projects.Create(context.Background(), projectRequest)

	if err != nil {
		return fmt.Errorf("Error creating Project: %s", err)
	}

	d.SetId(project.ID)
	log.Printf("[INFO] Project created, ID: %s", d.Id())

	return resourceDigitalOceanProjectRead(d, meta)
}

func resourceDigitalOceanProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	project, resp, err := client.Projects.Get(context.Background(), d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] Project  (%s) was not found - removing from state", d.Id())
			d.SetId("")
		}

		return fmt.Errorf("Error reading Project: %s", err)
	}

	d.Set("id", project.ID)
	d.Set("name", project.Name)
	d.Set("purpose", strings.TrimPrefix(project.Purpose, "Other: "))
	d.Set("description", project.Description)
	d.Set("environment", project.Environment)
	d.Set("is_default", project.IsDefault)
	d.Set("owner_uuid", project.OwnerUUID)
	d.Set("owner_id", project.OwnerID)
	d.Set("created_at", project.CreatedAt)
	d.Set("updated_at", project.UpdatedAt)

	return err
}

func resourceDigitalOceanProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectRequest := &godo.UpdateProjectRequest{
		Name:        d.Get("name"),
		Description: d.Get("description"),
		Purpose:     d.Get("purpose"),
		Environment: d.Get("environment"),
		IsDefault:   d.Get("is_default"),
	}

	_, _, err := client.Projects.Update(context.Background(), d.Id(), projectRequest)
	if err != nil {
		return fmt.Errorf("Error updating Project: %s", err)
	}

	log.Printf("[INFO] Updated Project")

	return resourceDigitalOceanProjectRead(d, meta)
}

func resourceDigitalOceanProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	resourceId := d.Id()

	_, err := client.Projects.Delete(context.Background(), resourceId)
	if err != nil {
		return fmt.Errorf("Error deleteing Project %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] Project deleted, ID: %s", resourceId)

	return nil
}
