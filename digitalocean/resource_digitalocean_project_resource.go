package digitalocean

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanProjectResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanProjectResourceCreate,
		Read:   resourceDigitalOceanProjectResourceRead,
		Delete: resourceDigitalOceanProjectResourceDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceDigitalOceanProjectResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Get("project").(string)
	urn := d.Get("resource").(string)

	project, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return fmt.Errorf("Project ID %s does not exist", projectId)
		}

		return fmt.Errorf("Error while retrieving project: %v", err)
	}

	_, _, err = client.Projects.AssignResources(context.Background(), project.ID, urn)
	if err != nil {
		return fmt.Errorf("Error assigning resource %s to project %s: %s", urn, project.ID, err)
	}

	d.SetId(fmt.Sprintf("%s:%s", project.ID, urn))
	return resourceDigitalOceanProjectResourceRead(d, meta)
}

func decodeProjectResourceId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	} else {
		return "", "", fmt.Errorf("Expected ID for digitalocean_project_resource as PROJECT_ID:RESOURCE_URN, got: %s", id)
	}
}

func resourceDigitalOceanProjectResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectId, urn, err := decodeProjectResourceId(d.Id())
	if err != nil {
		return err
	}

	project, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error while retrieving project: %v", err)
	}

	resourceUrns, err := loadResourceURNs(client, project.ID)
	if err != nil {
		return err
	}

	foundUrn := false
	for _, resourceUrn := range *resourceUrns {
		if urn == resourceUrn {
			foundUrn = true
		}
	}

	if !foundUrn {
		// If the resource is no longer assigned to this project,
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project.ID); err != nil {
		return err
	}

	if err := d.Set("resource", urn); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanProjectResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Get("project").(string)
	urn := d.Get("resource").(string)

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error while retrieving project: %s", err)
	}

	defaultProject, _, err := client.Projects.GetDefault(context.Background())
	if err != nil {
		return fmt.Errorf("Error locating default project: %s", err)
	}

	_, _, err = client.Projects.AssignResources(context.Background(), defaultProject.ID, urn)
	if err != nil {
		return fmt.Errorf("Error assigning resource %s to default project: %s", urn, err)
	}

	d.SetId("")
	return nil
}
