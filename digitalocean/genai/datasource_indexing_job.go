package genai

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanIndexingJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanIndexingJobRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the indexing job",
			},
			"knowledge_base_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Knowledge base UUID",
			},
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current phase of the batch job",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the indexing job",
			},
			"completed_datasources": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of datasources indexed completed",
			},
			"total_datasources": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of datasources being indexed",
			},
			"tokens": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tokens",
			},
			"total_items_failed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Total items failed",
			},
			"total_items_indexed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Total items indexed",
			},
			"total_items_skipped": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Total items skipped",
			},
			"data_source_uuids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of data source UUIDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp",
			},
			"started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Start timestamp",
			},
			"finished_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Finish timestamp",
			},
		},
	}
}

func dataSourceDigitalOceanIndexingJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	indexingJobUUIDRaw, ok := d.GetOk("uuid")
	if !ok || indexingJobUUIDRaw == nil {
		return diag.Errorf("uuid must be provided")
	}
	indexingJobUUID := indexingJobUUIDRaw.(string)

	// Call the API to get the specific indexing job
	indexingJobResponse, _, err := client.GenAI.GetIndexingJob(ctx, indexingJobUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the job details directly on the resource
	if indexingJobResponse != nil {
		job := &indexingJobResponse.Job

		d.Set("knowledge_base_uuid", job.KnowledgeBaseUuid)
		d.Set("phase", job.Phase)
		d.Set("status", job.Status)
		d.Set("completed_datasources", job.CompletedDatasources)
		d.Set("total_datasources", job.TotalDatasources)
		d.Set("tokens", job.Tokens)
		d.Set("total_items_failed", job.TotalItemsFailed)
		d.Set("total_items_indexed", job.TotalItemsIndexed)
		d.Set("total_items_skipped", job.TotalItemsSkipped)

		// Handle data source UUIDs
		if job.DataSourceUuids != nil {
			dataSourceUuids := make([]interface{}, len(job.DataSourceUuids))
			for j, uuid := range job.DataSourceUuids {
				dataSourceUuids[j] = uuid
			}
			d.Set("data_source_uuids", dataSourceUuids)
		} else {
			d.Set("data_source_uuids", []interface{}{})
		}

		// Handle timestamps
		if job.CreatedAt != nil {
			d.Set("created_at", job.CreatedAt.UTC().String())
		}
		if job.UpdatedAt != nil {
			d.Set("updated_at", job.UpdatedAt.UTC().String())
		}
		if job.StartedAt != nil {
			d.Set("started_at", job.StartedAt.UTC().String())
		}
		if job.FinishedAt != nil {
			d.Set("finished_at", job.FinishedAt.UTC().String())
		}
	}

	d.SetId(indexingJobUUID)
	return nil
}
