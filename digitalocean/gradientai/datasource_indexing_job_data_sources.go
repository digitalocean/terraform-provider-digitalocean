package gradientai

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanIndexingJobDataSources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanIndexingJobDataSourcesRead,
		Schema: map[string]*schema.Schema{
			"indexing_job_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the indexing job",
			},
			"indexed_data_sources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of indexed data sources for the indexing job",
				Elem:        IndexedDataSourceSchema(),
			},
		},
	}
}

func dataSourceDigitalOceanIndexingJobDataSourcesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	indexingJobUUIDRaw, ok := d.GetOk("indexing_job_uuid")
	if !ok || indexingJobUUIDRaw == nil {
		return diag.Errorf("indexing_job_uuid must be provided")
	}
	indexingJobUUID := indexingJobUUIDRaw.(string)

	// Call the API to list indexed data sources for the indexing job
	indexedDataSourcesResponse, _, err := client.GradientAI.ListIndexingJobDataSources(ctx, indexingJobUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Flatten and set indexed data sources
	var indexedDataSources []godo.IndexedDataSource
	if indexedDataSourcesResponse != nil {
		indexedDataSources = indexedDataSourcesResponse.IndexedDataSources
	}

	flattened := flattenIndexedDataSources(indexedDataSources)
	if err := d.Set("indexed_data_sources", flattened); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(indexingJobUUID)
	return nil
}
