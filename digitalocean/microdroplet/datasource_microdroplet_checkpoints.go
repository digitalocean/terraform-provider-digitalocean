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
// over the sibling checkpoints collection. Optionally filter by the
// MicroDroplet the checkpoints were captured from.
func DataSourceDigitalOceanMicroDropletCheckpoints() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        microDropletCheckpointSchema(),
		ResultAttributeName: "checkpoints",
		GetRecords:          getDigitalOceanMicroDropletCheckpoints,
		FlattenRecord:       flattenMicroDropletCheckpoint,
		ExtraQuerySchema: map[string]*schema.Schema{
			"microdroplet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Filter checkpoints captured from this MicroDroplet UUID. Omit to list all checkpoints for the team.",
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
			Description: "ID of the MicroDroplet the checkpoint was captured from",
		},
		"microdroplet_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the MicroDroplet the checkpoint was captured from",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Checkpoint name",
		},
		"region": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Region slug where the checkpoint is stored",
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
		"id":                c.ID,
		"microdroplet_id":   c.MicroDropletID,
		"microdroplet_name": c.MicroDropletName,
		"name":              c.Name,
		"region":            c.Region,
		"status":            string(c.Status),
		"memory_bytes":      int(c.MemoryBytes),
		"disk_bytes":        int(c.DiskBytes),
		"created_at":        c.Created,
	}, nil
}

func getDigitalOceanMicroDropletCheckpoints(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	microdropletID, _ := extra["microdroplet_id"].(string)

	opts := &godo.ListMicroDropletCheckpointsOptions{
		ListOptions:    godo.ListOptions{Page: 1, PerPage: 200},
		MicroDropletID: microdropletID,
	}

	var records []interface{}
	for {
		batch, resp, err := client.MicroDroplets.ListCheckpoints(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving MicroDroplet checkpoints: %w", err)
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
