package project

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanProjectCreate,
		ReadContext:   resourceDigitalOceanProjectRead,
		UpdateContext: resourceDigitalOceanProjectUpdate,
		DeleteContext: resourceDigitalOceanProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description:  "the description of the project",
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
				Description:  "the environment of the project's resources",
				ValidateFunc: validation.StringInSlice([]string{"development", "staging", "production"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) == strings.ToLower(new)
				},
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
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "determine if the project is the default or not.",
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
			"resources": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "the resources associated with the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
	}
}

func resourceDigitalOceanProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

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
		return diag.Errorf("Error creating Project: %s", err)
	}

	if v, ok := d.GetOk("resources"); ok {

		resources, err := assignResourcesToProject(client, project.ID, v.(*schema.Set))
		if err != nil {
			return diag.Errorf("Error creating project: %s", err)
		}

		d.Set("resources", resources)
	}

	if v, ok := d.GetOk("is_default"); ok {
		updateReq := &godo.UpdateProjectRequest{
			Name:        project.Name,
			Description: project.Description,
			Purpose:     project.Purpose,
			Environment: project.Environment,
			IsDefault:   v.(bool),
		}

		_, _, err := client.Projects.Update(context.Background(), project.ID, updateReq)
		if err != nil {
			return diag.Errorf("Error setting project as default: %s", err)
		}
	}

	d.SetId(project.ID)
	log.Printf("[INFO] Project created, ID: %s", d.Id())

	return resourceDigitalOceanProjectRead(ctx, d, meta)
}

func resourceDigitalOceanProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	project, resp, err := client.Projects.Get(context.Background(), d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] Project  (%s) was not found - removing from state", d.Id())
			d.SetId("")
		}

		return diag.Errorf("Error reading Project: %s", err)
	}

	d.SetId(project.ID)
	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("purpose", strings.TrimPrefix(project.Purpose, "Other: ")); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", project.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("environment", project.Environment); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("is_default", project.IsDefault); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("owner_uuid", project.OwnerUUID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("owner_id", project.OwnerID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("created_at", project.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("updated_at", project.UpdatedAt); err != nil {
		return diag.FromErr(err)
	}

	urns, err := LoadResourceURNs(client, project.ID)
	if err != nil {
		return diag.Errorf("Error reading Project: %s", err)
	}

	if err = d.Set("resources", urns); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDigitalOceanProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	projectId := d.Id()

	d.Partial(true)

	projectRequest := &godo.UpdateProjectRequest{
		Name:        d.Get("name"),
		Description: d.Get("description"),
		Purpose:     d.Get("purpose"),
		Environment: d.Get("environment"),
		IsDefault:   d.Get("is_default"),
	}

	_, _, err := client.Projects.Update(context.Background(), projectId, projectRequest)

	if err != nil {
		return diag.Errorf("Error updating Project: %s", err)
	}

	// The API requires project resources to be reassigned to another project if the association needs to be deleted.
	if d.HasChange("resources") {
		oldURNs, newURNs := d.GetChange("resources")
		remove, add := util.GetSetChanges(oldURNs.(*schema.Set), newURNs.(*schema.Set))

		if remove.Len() > 0 {
			_, err = assignResourcesToDefaultProject(client, remove)
			if err != nil {
				return diag.Errorf("Error assigning resources to default project: %s", err)
			}
		}

		if add.Len() > 0 {
			_, err = assignResourcesToProject(client, projectId, add)
			if err != nil {
				return diag.Errorf("Error Updating project: %s", err)
			}
		}

		d.Set("resources", newURNs)
	}

	log.Printf("[INFO] Updated Project, ID: %s", projectId)
	d.Partial(false)

	return resourceDigitalOceanProjectRead(ctx, d, meta)
}

func resourceDigitalOceanProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	projectID := d.Id()

	if v, ok := d.GetOk("resources"); ok {
		_, err := assignResourcesToDefaultProject(client, v.(*schema.Set))
		if err != nil {
			return diag.Errorf("Error assigning resource to default project: %s", err)
		}

		d.Set("resources", nil)
		log.Printf("[DEBUG] Resources assigned to default project.")
	}

	// Moving resources is async and projects can not be deleted till empty. Retries may be required.
	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.Projects.Delete(context.Background(), projectID)
		if err != nil {
			if util.IsDigitalOceanError(err, http.StatusPreconditionFailed, "cannot delete a project with resources") {
				log.Printf("[DEBUG] Received %s, retrying project deletion", err.Error())
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error deleting project (%s): %s", projectID, err)
	}

	d.SetId("")
	log.Printf("[INFO] Project deleted, ID: %s", projectID)

	return nil
}

func assignResourcesToDefaultProject(client *godo.Client, resources *schema.Set) (*[]interface{}, error) {
	defaultProject, _, defaultProjErr := client.Projects.GetDefault(context.Background())
	if defaultProjErr != nil {
		return nil, fmt.Errorf("Error locating default project %s", defaultProjErr)
	}

	return assignResourcesToProject(client, defaultProject.ID, resources)
}

func assignResourcesToProject(client *godo.Client, projectID string, resources *schema.Set) (*[]interface{}, error) {
	var urns []interface{}

	for _, resource := range resources.List() {
		if resource == nil {
			continue
		}

		if resource == "" {
			continue
		}

		urns = append(urns, resource.(string))
	}

	_, _, err := client.Projects.AssignResources(context.Background(), projectID, urns...)
	if err != nil {
		return nil, fmt.Errorf("Error assigning resources: %s", err)
	}

	return &urns, nil
}
