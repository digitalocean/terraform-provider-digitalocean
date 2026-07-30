package microdroplet

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDigitalOceanMicroDropletCheckpoints returns a plural data source
// over the checkpoints belonging to a given MicroDroplet. Checkpoints are
// created automatically by DigitalOcean when a MicroDroplet pauses; they are
// read-only from the customer API, so no matching resource is provided.
func DataSourceDigitalOceanMicroDropletCheckpoints() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        microDropletCheckpointSchema(),
		ResultAttributeName: "checkpoints",
		GetRecords:          getDigitalOceanMicroDropletCheckpoints,
		FlattenRecord:       flattenMicroDropletCheckpoint,
		ExtraQuerySchema: map[string]*schema.Schema{
			"microdroplet_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of the MicroDroplet whose checkpoints should be listed",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
	return datalist.NewResource(dataListConfig)
}

func microDropletCheckpointSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Checkpoint ID",
		},
		"microdroplet_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the parent MicroDroplet",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Checkpoint name",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Lifecycle status of the checkpoint",
		},
		"memory_bytes": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Size of the persisted memory image, in bytes",
		},
		"disk_bytes": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Size of the persisted disk image, in bytes",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation timestamp for the checkpoint",
		},
	}
}

func flattenMicroDropletCheckpoint(rawRecord, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	c, ok := rawRecord.(godo.MicroDropletCheckpoint)
	if !ok {
		return nil, fmt.Errorf("unexpected record type %T", rawRecord)
	}
	return map[string]interface{}{
		"id":              c.ID,
		"microdroplet_id": c.MicroDropletID,
		"name":            c.Name,
		"status":          string(c.Status),
		"memory_bytes":    int(c.MemoryBytes),
		"disk_bytes":      int(c.DiskBytes),
		"created_at":      c.Created,
	}, nil
}

// getDigitalOceanMicroDropletCheckpoints paginates the ListCheckpoints godo
// endpoint for the MicroDroplet identified by the required `microdroplet_id`
// query attribute.
func getDigitalOceanMicroDropletCheckpoints(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	microdropletID, _ := extra["microdroplet_id"].(string)
	if microdropletID == "" {
		return nil, fmt.Errorf("microdroplet_id is required")
	}

	opts := &godo.ListOptions{Page: 1, PerPage: 200}

	var records []interface{}
	for {
		batch, resp, err := client.MicroDroplets.ListCheckpoints(context.Background(), microdropletID, opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving checkpoints for MicroDroplet %s: %w", microdropletID, err)
		}
		for _, c := range batch {
			records = append(records, c)
		}
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error paging MicroDroplet checkpoints: %w", err)
		}
		opts.Page = page + 1
	}
	return records, nil
}
