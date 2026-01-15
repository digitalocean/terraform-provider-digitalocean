package gradientai

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDigitalOceanIndexingJobCancel defines the DigitalOcean GradientAI Indexing Job Cancel resource
func ResourceDigitalOceanIndexingJobCancel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanIndexingJobCancelCreate,
		ReadContext:   resourceDigitalOceanIndexingJobCancelRead,
		DeleteContext: resourceDigitalOceanIndexingJobCancelDelete,
		// No update context since this is a one-time action

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the indexing job to cancel.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the indexing job after cancellation.",
			},
			"knowledge_base_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the knowledge base associated with this indexing job.",
			},
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current phase of the indexing job.",
			},
			"completed_datasources": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of data sources that were completed before cancellation.",
			},
			"total_datasources": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of data sources in the indexing job.",
			},
			"tokens": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tokens processed before cancellation.",
			},
			"total_items_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of items that failed during indexing.",
			},
			"total_items_indexed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of items that were successfully indexed.",
			},
			"total_items_skipped": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of items that were skipped during indexing.",
			},
			"data_source_uuids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of data source UUIDs associated with this indexing job.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the indexing job was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the indexing job was last updated.",
			},
			"started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the indexing job was started.",
			},
			"finished_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the indexing job was finished.",
			},
		},
	}
}

func resourceDigitalOceanIndexingJobCancelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	uuid := d.Get("uuid").(string)

	// Cancel the indexing job
	jobResponse, _, err := client.GradientAI.CancelIndexingJob(ctx, uuid)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error canceling indexing job (%s): %s", uuid, err))
	}

	// Set the ID to the UUID
	d.SetId(uuid)

	// Set all the job attributes
	return setIndexingJobCancelAttributes(d, &jobResponse.Job)
}

func resourceDigitalOceanIndexingJobCancelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	uuid := d.Id()

	// Get the current status of the indexing job
	jobResponse, _, err := client.GradientAI.GetIndexingJob(ctx, uuid)
	if err != nil {
		// If the job is not found, remove it from state
		if godoErr, ok := err.(*godo.ErrorResponse); ok && godoErr.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading indexing job (%s): %s", uuid, err))
	}

	return setIndexingJobCancelAttributes(d, &jobResponse.Job)
}

func resourceDigitalOceanIndexingJobCancelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Canceling an indexing job is idempotent - there's nothing to "delete"
	// The job cancellation itself is permanent, so we just remove it from state
	d.SetId("")
	return nil
}

func setIndexingJobCancelAttributes(d *schema.ResourceData, job *godo.LastIndexingJob) diag.Diagnostics {
	d.Set("uuid", job.Uuid)
	d.Set("status", job.Status)
	d.Set("knowledge_base_uuid", job.KnowledgeBaseUuid)
	d.Set("phase", job.Phase)
	d.Set("completed_datasources", job.CompletedDatasources)
	d.Set("total_datasources", job.TotalDatasources)
	d.Set("tokens", job.Tokens)
	d.Set("total_items_failed", job.TotalItemsFailed)
	d.Set("total_items_indexed", job.TotalItemsIndexed)
	d.Set("total_items_skipped", job.TotalItemsSkipped)

	// Handle data source UUIDs
	if job.DataSourceUuids != nil {
		dataSourceUuids := make([]interface{}, len(job.DataSourceUuids))
		for i, uuid := range job.DataSourceUuids {
			dataSourceUuids[i] = uuid
		}
		d.Set("data_source_uuids", dataSourceUuids)
	} else {
		d.Set("data_source_uuids", []interface{}{})
	}

	// Handle timestamps
	if job.CreatedAt != nil {
		d.Set("created_at", job.CreatedAt.UTC().String())
	} else {
		d.Set("created_at", "")
	}

	if job.UpdatedAt != nil {
		d.Set("updated_at", job.UpdatedAt.UTC().String())
	} else {
		d.Set("updated_at", "")
	}

	if job.StartedAt != nil {
		d.Set("started_at", job.StartedAt.UTC().String())
	} else {
		d.Set("started_at", "")
	}

	if job.FinishedAt != nil {
		d.Set("finished_at", job.FinishedAt.UTC().String())
	} else {
		d.Set("finished_at", "")
	}

	return nil
}
