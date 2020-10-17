package digitalocean

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanProjectResources() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanProjectResourcesUpdate,
		UpdateContext: resourceDigitalOceanProjectResourcesUpdate,
		ReadContext:   resourceDigitalOceanProjectResourcesRead,
		DeleteContext: resourceDigitalOceanProjectResourcesDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "project ID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resources": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "the resources associated with the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDigitalOceanProjectResourcesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Get("project").(string)

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error while retrieving project %s: %v", projectId, err)
	}

	if d.HasChange("resources") {
		oldURNs, newURNs := d.GetChange("resources")

		if oldURNs.(*schema.Set).Len() > 0 {
			_, err = assignResourcesToDefaultProject(client, oldURNs.(*schema.Set))
			if err != nil {
				return diag.Errorf("Error assigning resources to default project: %s", err)
			}
		}

		var urns *[]interface{}

		if newURNs.(*schema.Set).Len() > 0 {
			urns, err = assignResourcesToProject(client, projectId, newURNs.(*schema.Set))
			if err != nil {
				return diag.Errorf("Error assigning resources to project %s: %s", projectId, err)
			}
		}

		if err = d.Set("resources", urns); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(projectId)

	return resourceDigitalOceanProjectResourcesRead(ctx, d, meta)
}

func resourceDigitalOceanProjectResourcesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Id()

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error while retrieving project: %v", err)
	}

	if err = d.Set("project", projectId); err != nil {
		return diag.FromErr(err)
	}

	apiURNs, err := loadResourceURNs(client, projectId)
	if err != nil {
		return diag.Errorf("Error while retrieving project resources: %s", err)
	}

	var newURNs []string

	configuredURNs := d.Get("resources").(*schema.Set).List()
	for _, rawConfiguredURN := range configuredURNs {
		configuredURN := rawConfiguredURN.(string)

		for _, apiURN := range *apiURNs {
			if configuredURN == apiURN {
				newURNs = append(newURNs, configuredURN)
			}
		}
	}

	if err = d.Set("resources", newURNs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDigitalOceanProjectResourcesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Get("project").(string)
	urns := d.Get("resources").(*schema.Set)

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error while retrieving project: %s", err)
	}

	if urns.Len() > 0 {
		if _, err = assignResourcesToDefaultProject(client, urns); err != nil {
			return diag.Errorf("Error assigning resources to default project: %s", err)
		}
	}

	d.SetId("")
	return nil
}
