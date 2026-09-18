package microvm

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDigitalOceanMicroVMCheckpoints returns a plural data source
// over the sibling checkpoints collection. Optionally filter by the
// MicroVM the checkpoints were captured from.
func DataSourceDigitalOceanMicroVMCheckpoints() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        microVMCheckpointSchema(),
		ResultAttributeName: "checkpoints",
		GetRecords:          getDigitalOceanMicroVMCheckpoints,
		FlattenRecord:       flattenMicroVMCheckpoint,
		ExtraQuerySchema: map[string]*schema.Schema{
			"microvm_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Filter checkpoints captured from this MicroVM UUID. Omit to list all checkpoints for the team.",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
	return datalist.NewResource(dataListConfig)
}

func microVMCheckpointSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Checkpoint ID",
		},
		"microvm_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the MicroVM the checkpoint was captured from",
		},
		"microvm_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the MicroVM the checkpoint was captured from",
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

func flattenMicroVMCheckpoint(rawRecord, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	c, ok := rawRecord.(godo.MicroVMCheckpoint)
	if !ok {
		return nil, fmt.Errorf("unexpected record type %T", rawRecord)
	}
	return map[string]interface{}{
		"id":           c.ID,
		"microvm_id":   c.MicroVMID,
		"microvm_name": c.MicroVMName,
		"name":         c.Name,
		"region":       c.Region,
		"status":       string(c.Status),
		"memory_bytes": int(c.MemoryBytes),
		"disk_bytes":   int(c.DiskBytes),
		"created_at":   c.Created,
	}, nil
}

func getDigitalOceanMicroVMCheckpoints(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	microvmID, _ := extra["microvm_id"].(string)

	opts := &godo.ListMicroVMCheckpointsOptions{
		ListOptions: godo.ListOptions{Page: 1, PerPage: 200},
		MicroVMID:   microvmID,
	}

	var records []interface{}
	for {
		batch, resp, err := client.MicroVMs.ListCheckpoints(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving MicroVM checkpoints: %w", err)
		}
		for _, c := range batch {
			records = append(records, c)
		}
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error paging MicroVM checkpoints: %w", err)
		}
		opts.Page = page + 1
	}
	return records, nil
}
