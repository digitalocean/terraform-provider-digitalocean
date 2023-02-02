package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func projectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type: schema.TypeString,
		},
		"name": {
			Type: schema.TypeString,
		},
		"description": {
			Type: schema.TypeString,
		},
		"purpose": {
			Type: schema.TypeString,
		},
		"environment": {
			Type: schema.TypeString,
		},
		"owner_uuid": {
			Type: schema.TypeString,
		},
		"owner_id": {
			Type: schema.TypeInt,
		},
		"is_default": {
			Type: schema.TypeBool,
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "the date and time when the project was created, (ISO8601)",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Description: "the date and time when the project was last updated, (ISO8601)",
		},
		"resources": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}
}

func getDigitalOceanProjects(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	var allProjects []interface{}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		projects, resp, err := client.Projects.List(context.Background(), opts)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving projects: %s", err)
		}

		for _, project := range projects {
			allProjects = append(allProjects, project)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving projects: %s", err)
		}

		opts.Page = page + 1
	}

	return allProjects, nil
}

func flattenDigitalOceanProject(rawProject interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	project, ok := rawProject.(godo.Project)
	if !ok {
		return nil, fmt.Errorf("Unable to convert to godo.Project")
	}

	flattenedProject := map[string]interface{}{}
	flattenedProject["id"] = project.ID
	flattenedProject["name"] = project.Name
	flattenedProject["purpose"] = strings.TrimPrefix(project.Purpose, "Other: ")
	flattenedProject["description"] = project.Description
	flattenedProject["environment"] = project.Environment
	flattenedProject["owner_uuid"] = project.OwnerUUID
	flattenedProject["owner_id"] = project.OwnerID
	flattenedProject["is_default"] = project.IsDefault
	flattenedProject["created_at"] = project.CreatedAt
	flattenedProject["updated_at"] = project.UpdatedAt

	urns, err := LoadResourceURNs(client, project.ID)
	if err != nil {
		return nil, fmt.Errorf("Error loading project resource URNs for project ID %s: %s", project.ID, err)
	}

	flattenedURNS := schema.NewSet(schema.HashString, []interface{}{})
	for _, urn := range *urns {
		flattenedURNS.Add(urn)
	}
	flattenedProject["resources"] = flattenedURNS

	return flattenedProject, nil
}

func LoadResourceURNs(client *godo.Client, projectId string) (*[]string, error) {
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

		resourceList = append(resourceList, resources...)

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
