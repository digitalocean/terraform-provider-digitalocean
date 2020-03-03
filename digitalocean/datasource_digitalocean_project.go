package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanProject() *schema.Resource {
	recordSchema := projectSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["id"].ConflictsWith = []string{"name"}
	recordSchema["id"].Optional = true
	recordSchema["id"].Computed = false
	recordSchema["name"].ConflictsWith = []string{"id"}
	recordSchema["name"].Optional = true
	recordSchema["name"].Computed = false

	return &schema.Resource{
		Read:   dataSourceDigitalOceanProjectRead,
		Schema: recordSchema,
	}
}

func dataSourceDigitalOceanProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	// Load the specified project, otherwise load the default project.
	var project *godo.Project
	if projectId, ok := d.GetOk("id"); ok {
		thisProject, _, err := client.Projects.Get(context.Background(), projectId.(string))
		if err != nil {
			return fmt.Errorf("Unable to load project ID %s: %s", projectId, err)
		}
		project = thisProject
	} else if name, ok := d.GetOk("name"); ok {
		projects, err := getDigitalOceanProjects(client)
		if err != nil {
			return fmt.Errorf("Unable to load projects: %s", err)
		}

		var projectsWithName []*godo.Project
		for _, project := range projects {
			if project.Name == name {
				projectsWithName = append(projectsWithName, project)
			}
		}
		if len(projectsWithName) == 0 {
			return fmt.Errorf("No projects found with name '%s'", name)
		} else if len(projectsWithName) > 1 {
			return fmt.Errorf("Multiple projects found with name '%s'", name)
		}

		// Single result so choose that project.
		project = projects[0]
	} else {
		defaultProject, _, err := client.Projects.GetDefault(context.Background())
		if err != nil {
			return fmt.Errorf("Unable to load default project: %s", err)
		}
		project = defaultProject
	}

	flattenedProject, err := flattenDigitalOceanProject(project, meta)
	if err != nil {
		return err
	}

	if err := setResourceDataFromMap(d, flattenedProject); err != nil {
		return err
	}

	d.SetId(project.ID)
	return nil
}
