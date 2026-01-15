package gradientai

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanKnowledgeBaseIndexingJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanKnowledgeBaseIndexingJobsRead,
		Schema: map[string]*schema.Schema{
			"knowledge_base_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the Knowledge Base",
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of indexing jobs for the Knowledge Base",
				Elem:        IndexingJobSchema(),
			},
			"meta": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Pagination metadata",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"page": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current page number",
						},
						"pages": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of pages",
						},
						"total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of items",
						},
					},
				},
			},
		},
	}
}

func dataSourceDigitalOceanKnowledgeBaseIndexingJobsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	kbUUIDRaw, ok := d.GetOk("knowledge_base_uuid")
	if !ok || kbUUIDRaw == nil {
		return diag.Errorf("knowledge_base_uuid must be provided")
	}
	kbUUID := kbUUIDRaw.(string)

	// Call the API to list indexing jobs for the KB
	// Note: The API may support filtering by knowledge_base_uuid as a query parameter
	jobsResponse, resp, err := client.GradientAI.ListIndexingJobs(ctx, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Filter jobs by knowledge base UUID
	var filteredJobs []godo.LastIndexingJob
	if jobsResponse != nil {
		for _, job := range jobsResponse.Jobs {
			if job.KnowledgeBaseUuid == kbUUID {
				filteredJobs = append(filteredJobs, job)
			}
		}
	}

	// Flatten and set jobs
	flattened := flattenKnowledgeBaseIndexingJobs(filteredJobs)
	if err := d.Set("jobs", flattened); err != nil {
		return diag.FromErr(err)
	}

	// Set pagination metadata if available
	if resp != nil && resp.Meta != nil {
		meta := []map[string]interface{}{
			{
				"page":  resp.Meta.Page,
				"pages": resp.Meta.Pages,
				"total": resp.Meta.Total,
			},
		}
		if err := d.Set("meta", meta); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(kbUUID)
	return nil
}
