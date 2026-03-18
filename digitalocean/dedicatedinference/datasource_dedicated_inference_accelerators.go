package dedicatedinference

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanDedicatedInferenceAccelerators() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dedicatedInferenceAcceleratorSchema(),
		ResultAttributeName: "accelerators",
		ExtraQuerySchema: map[string]*schema.Schema{
			"dedicated_inference_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of the dedicated inference endpoint to list accelerators for.",
				ValidateFunc: validation.NoZeroValues,
			},
		},
		GetRecords:    getDigitalOceanDedicatedInferenceAccelerators,
		FlattenRecord: flattenDigitalOceanDedicatedInferenceAccelerator,
	}

	return datalist.NewResource(dataListConfig)
}

func dedicatedInferenceAcceleratorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The unique ID of the accelerator.",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the accelerator.",
		},
		"slug": {
			Type:        schema.TypeString,
			Description: "The slug identifier for the accelerator type.",
		},
		"status": {
			Type:        schema.TypeString,
			Description: "The current status of the accelerator.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "The date and time when the accelerator was created.",
		},
	}
}

func getDigitalOceanDedicatedInferenceAccelerators(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	diID := extra["dedicated_inference_id"].(string)

	var allItems []godo.DedicatedInferenceAcceleratorInfo
	opts := &godo.DedicatedInferenceListAcceleratorsOptions{
		ListOptions: godo.ListOptions{
			Page:    1,
			PerPage: 200,
		},
	}

	for {
		items, resp, err := client.DedicatedInference.ListAccelerators(context.Background(), diID, opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving accelerators for dedicated inference endpoint (%s): %s", diID, err)
		}
		allItems = append(allItems, items...)
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error retrieving accelerators for dedicated inference endpoint (%s): %s", diID, err)
		}
		opts.Page = page + 1
	}

	records := make([]interface{}, len(allItems))
	for i, item := range allItems {
		records[i] = item
	}
	return records, nil
}

func flattenDigitalOceanDedicatedInferenceAccelerator(record, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	item := record.(godo.DedicatedInferenceAcceleratorInfo)

	flat := map[string]interface{}{
		"id":     item.ID,
		"name":   item.Name,
		"slug":   item.Slug,
		"status": item.Status,
	}

	if !item.CreatedAt.IsZero() {
		flat["created_at"] = item.CreatedAt.UTC().String()
	} else {
		flat["created_at"] = ""
	}

	return flat, nil
}
