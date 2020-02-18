package digitalocean

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var (
	digitalOceanProjectsFilterKeys = []string{
		"name",
		"purpose",
		"description",
		"environment",
		"is_default",
	}

	digitalOceanProjectsSortKeys = []string{
		"name",
		"purpose",
		"description",
		"environment",
	}
)

func dataSourceDigitalOceanProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanProjectsRead,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema(digitalOceanProjectsFilterKeys),
			"sort":   sortSchema(digitalOceanProjectsSortKeys),
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"purpose": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"environment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_id": {
							Type:     schema.TypeInt,
							Computed: true,
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
						//"resources": {
						//	Type:     schema.TypeSet,
						//	Computed: true,
						//	Elem:     &schema.Schema{Type: schema.TypeString},
						//},
					},
				},
			},
		},
	}
}

func getDigitalOceanProjects(client *godo.Client) ([]*godo.Project, error) {
	allProjects := []*godo.Project{}

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
			allProjects = append(allProjects, &project)
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

func dataSourceDigitalOceanProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projects, err := getDigitalOceanProjects(client)
	if err != nil {
		return err
	}
	log.Printf("projects = %+v", projects)

	if v, ok := d.GetOk("filter"); ok {
		filters := expandFilters(v.(*schema.Set).List())
		log.Printf("filters = %+v", filters)
		log.Printf("projects (before filter) = %+v", projects)
		projects = filterDigitalOceanProjects(projects, filters)
		log.Printf("projects (after filter) = %+v", projects)
	}

	if v, ok := d.GetOk("sort"); ok {
		sorts := expandSorts(v.([]interface{}))
		log.Printf("sorts = %+v", sorts)
		log.Printf("projects (before sort) = %+v", projects)
		projects = sortDigitalOceanProjects(projects, sorts)
		log.Printf("projects (after sort) = %+v", projects)
	}

	d.SetId(resource.UniqueId())

	flattenedProjects := make([]map[string]interface{}, len(projects))
	for i, project := range projects {
		flattenedProjects[i] = flattenProject(project)
	}

	if err := d.Set("projects", flattenedProjects); err != nil {
		return fmt.Errorf("Unable to set `project_ids` attribute: %s", err)
	}

	return nil
}

func flattenProject(project *godo.Project) map[string]interface{} {
	flattenedProject := map[string]interface{}{}

	flattenedProject["id"] = project.ID
	flattenedProject["name"] = project.Name
	flattenedProject["description"] = project.Description
	flattenedProject["purpose"] = project.Purpose
	flattenedProject["environment"] = project.Environment
	flattenedProject["owner_uuid"] = project.OwnerUUID
	flattenedProject["owner_id"] = project.OwnerID
	flattenedProject["created_at"] = project.CreatedAt
	flattenedProject["updated_at"] = project.UpdatedAt
	flattenedProject["is_default"] = project.IsDefault

	return flattenedProject
}

func filterDigitalOceanProjects(projects []*godo.Project, filters []commonFilter) []*godo.Project {
	for _, f := range filters {
		var filteredProjects []*godo.Project

		filterFunc := func(project *godo.Project) bool {
			result := false

			for _, filterValue := range f.values {
				switch f.key {
				case "name":
					result = result || strings.EqualFold(filterValue, project.Name)

				case "purpose":
					result = result || strings.EqualFold(filterValue, project.Purpose)

				case "description":
					result = result || strings.EqualFold(filterValue, project.Description)

				case "environment":
					result = result || strings.EqualFold(filterValue, project.Environment)

				case "is_default":
					if isDefault, err := strconv.ParseBool(filterValue); err == nil {
						result = result || isDefault == project.IsDefault
					}
				default:
				}
			}

			return result
		}

		for _, project := range projects {
			if filterFunc(project) {
				filteredProjects = append(filteredProjects, project)
			}
		}

		projects = filteredProjects
	}

	return projects
}

func sortDigitalOceanProjects(projects []*godo.Project, sorts []commonSort) []*godo.Project {
	sort.Slice(projects, func(_i, _j int) bool {
		for _, s := range sorts {
			// Handle multiple sorts by applying them in order

			i := _i
			j := _j
			if strings.EqualFold(s.direction, "desc") {
				// If the direction is desc, reverse index to compare
				i = _j
				j = _i
			}

			switch s.key {
			case "name":
				if !strings.EqualFold(projects[i].Name, projects[j].Name) {
					return strings.Compare(projects[i].Name, projects[j].Name) <= 0
				}

			case "purpose":
				if !strings.EqualFold(projects[i].Purpose, projects[j].Purpose) {
					return strings.Compare(projects[i].Purpose, projects[j].Purpose) <= 0
				}

			case "description":
				if !strings.EqualFold(projects[i].Description, projects[j].Description) {
					return strings.Compare(projects[i].Description, projects[j].Description) <= 0
				}

			case "environment":
				if !strings.EqualFold(projects[i].Environment, projects[j].Environment) {
					return strings.Compare(projects[i].Environment, projects[j].Environment) <= 0
				}
			}
		}

		return true
	})

	return projects
}
