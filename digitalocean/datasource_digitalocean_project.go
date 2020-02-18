package digitalocean

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanProjectRead,
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
			"resources": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
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
	} else {
		defaultProject, _, err := client.Projects.GetDefault(context.Background())
		if err != nil {
			return fmt.Errorf("Unable to load default project: %s", err)
		}
		project = defaultProject
	}

	if err := d.Set("id", project.ID); err != nil {
		return fmt.Errorf("Unable to set `id` attribute: %s", err)
	}
	if err := d.Set("name", project.Name); err != nil {
		return fmt.Errorf("Unable to set `name` attribute: %s", err)
	}
	if err := d.Set("purpose", strings.TrimPrefix(project.Purpose, "Other: ")); err != nil {
		return fmt.Errorf("Unable to set `purpose` attribute: %s", err)
	}
	if err := d.Set("description", project.Description); err != nil {
		return fmt.Errorf("Unable to set `description` attribute: %s", err)
	}
	if err := d.Set("environment", project.Environment); err != nil {
		return fmt.Errorf("Unable to set `environment` attribute: %s", err)
	}
	if err := d.Set("owner_uuid", project.OwnerUUID); err != nil {
		return fmt.Errorf("Unable to set `owner_uuid` attribute: %s", err)
	}
	if err := d.Set("owner_id", project.OwnerID); err != nil {
		return fmt.Errorf("Unable to set `owner_id` attribute: %s", err)
	}
	if err := d.Set("is_default", project.IsDefault); err != nil {
		return fmt.Errorf("Unable to set `is_default` attribute: %s", err)
	}
	if err := d.Set("created_at", project.CreatedAt); err != nil {
		return fmt.Errorf("Unable to set `created_at` attribute: %s", err)
	}
	if err := d.Set("updated_at", project.UpdatedAt); err != nil {
		return fmt.Errorf("Unable to set `updated_at` attribute: %s", err)
	}

	urns, err := loadResourceURNs(client, project.ID)
	if err != nil {
		return fmt.Errorf("Error loading project resource URNs: %s", err)
	}

	if err := d.Set("resources", urns); err != nil {
		return fmt.Errorf("Unable to set `resources` attribute: %s", err)
	}

	d.SetId(project.ID)

	return nil
}
