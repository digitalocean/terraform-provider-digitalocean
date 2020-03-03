package digitalocean

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:     schema.TypeBool,
				Computed: true,
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
				Description: "the resources associated with the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
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

	if v, ok := d.GetOk("resources"); ok {

		resources, err := assignResourcesToProject(client, project.ID, v.(*schema.Set))
		if err != nil {
			return fmt.Errorf("Error creating project: %s", err)
		}

		d.Set("resources", resources)
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

	urns, err := loadResourceURNs(client, project.ID)
	if err != nil {
		return fmt.Errorf("Error reading Project: %s", err)
	}

	d.Set("resources", urns)

	return err
}

func resourceDigitalOceanProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
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
		return fmt.Errorf("Error updating Project: %s", err)
	}

	d.SetPartial("project_updated")

	// The API requires project resources to be reassigned to another project if the association needs to be deleted.
	// a diff of the resource could be implemented instead of removing all, (bulk) and adding the back again.
	if d.HasChange("resources") {
		oldURNs, newURNs := d.GetChange("resources")

		assignResourcesToDefaultProject(client, oldURNs.(*schema.Set))

		var urns *[]interface{}

		if newURNs.(*schema.Set).Len() != 0 {
			urns, err = assignResourcesToProject(client, projectId, newURNs.(*schema.Set))
			if err != nil {
				return fmt.Errorf("Error Updating project: %s", err)
			}
		}

		d.Set("resources", urns)
		d.SetPartial("project_resources_updated")
	}

	log.Printf("[INFO] Updated Project, ID: %s", projectId)
	d.Partial(false)

	return resourceDigitalOceanProjectRead(d, meta)
}

func resourceDigitalOceanProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Id()

	if v, ok := d.GetOk("resources"); ok {

		_, err := assignResourcesToDefaultProject(client, v.(*schema.Set))
		if err != nil {
			return fmt.Errorf("Error assigning resource to default project: %s", err)
		}

		d.Set("resources", nil)
		log.Printf("[DEBUG] Resources assigned to default project.")
	}

	_, err := client.Projects.Delete(context.Background(), projectId)
	if err != nil {
		return fmt.Errorf("Error deleteing project %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] Project deleted, ID: %s", projectId)

	return nil
}

func assignResourcesToDefaultProject(client *godo.Client, resources *schema.Set) (*[]interface{}, error) {

	defaultProject, _, defaultProjErr := client.Projects.GetDefault(context.Background())
	if defaultProjErr != nil {
		return nil, fmt.Errorf("Error locating default project %s", defaultProjErr)
	}

	return assignResourcesToProject(client, defaultProject.ID, resources)
}

func assignResourcesToProject(client *godo.Client, projectId string, resources *schema.Set) (*[]interface{}, error) {

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

	_, _, err := client.Projects.AssignResources(context.Background(), projectId, urns...)

	if err != nil {
		return nil, fmt.Errorf("Error assigning resources: %s", err)
	}

	return &urns, nil
}

func loadResourceURNs(client *godo.Client, projectId string) (*[]string, error) {
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	resourceList := []godo.ProjectResource{}
	for {
		resources, resp, err := client.Projects.ListResources(context.Background(), projectId, opts)
		if err != nil {
			return nil, fmt.Errorf("Error loading project resources: %s", err)
		}

		for _, r := range resources {
			resourceList = append(resourceList, r)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error loading project resources: %s", err)
		}

		opts.Page = page + 1
	}

	var urns []string
	for _, rsrc := range resourceList {
		urns = append(urns, rsrc.URN)
	}

	return &urns, nil
}
