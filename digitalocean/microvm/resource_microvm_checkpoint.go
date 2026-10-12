package microvm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceDigitalOceanMicroVMCheckpoint returns the
// digitalocean_microvm_checkpoint resource schema.
func ResourceDigitalOceanMicroVMCheckpoint() *schema.Resource {
	recordSchema := microVMCheckpointSchema()
	recordSchema["microvm_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "ID of the running MicroVM to checkpoint",
		ValidateFunc: validation.NoZeroValues,
	}
	recordSchema["name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ForceNew:     true,
		Description:  "Human-readable name for the checkpoint",
		ValidateFunc: validation.NoZeroValues,
	}
	delete(recordSchema, "id")

	return &schema.Resource{
		CreateContext: resourceDigitalOceanMicroVMCheckpointCreate,
		ReadContext:   resourceDigitalOceanMicroVMCheckpointRead,
		DeleteContext: resourceDigitalOceanMicroVMCheckpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: recordSchema,
	}
}

func resourceDigitalOceanMicroVMCheckpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	microVMID := d.Get("microvm_id").(string)
	req := &godo.MicroVMCheckpointCreateRequest{}
	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
	}

	cp, _, err := client.MicroVMs.CreateCheckpoint(ctx, microVMID, req)
	if err != nil {
		return diag.Errorf("Error creating checkpoint of MicroVM (%s): %s", microVMID, err)
	}

	d.SetId(cp.ID)
	log.Printf("[INFO] MicroVM checkpoint created, ID: %s", d.Id())

	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(godo.MicroVMCheckpointStatusUnknown),
			string(godo.MicroVMCheckpointStatusCreating),
		},
		Target:     []string{string(godo.MicroVMCheckpointStatusAvailable)},
		Refresh:    microVMCheckpointStatusRefreshFunc(ctx, client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for MicroVM checkpoint (%s) to become available: %s", d.Id(), err)
	}

	return resourceDigitalOceanMicroVMCheckpointRead(ctx, d, meta)
}

func resourceDigitalOceanMicroVMCheckpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	cp, resp, err := client.MicroVMs.GetCheckpoint(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[WARN] MicroVM checkpoint (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving MicroVM checkpoint: %s", err)
	}

	if cp.Status == godo.MicroVMCheckpointStatusDeleted || cp.Status == godo.MicroVMCheckpointStatusDeleting {
		log.Printf("[WARN] MicroVM checkpoint (%s) is %s - removing from state", d.Id(), cp.Status)
		d.SetId("")
		return nil
	}

	flat, err := flattenMicroVMCheckpoint(*cp, meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	delete(flat, "id")
	for k, v := range flat {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	return nil
}

func resourceDigitalOceanMicroVMCheckpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	resp, err := client.MicroVMs.DeleteCheckpoint(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.Errorf("Error deleting MicroVM checkpoint: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(godo.MicroVMCheckpointStatusAvailable),
			string(godo.MicroVMCheckpointStatusDeleting),
		},
		Target:     []string{string(godo.MicroVMCheckpointStatusDeleted)},
		Refresh:    microVMCheckpointDeleteRefreshFunc(ctx, client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for MicroVM checkpoint (%s) to be deleted: %s", d.Id(), err)
	}

	log.Printf("[INFO] MicroVM checkpoint deleted, ID: %s", d.Id())
	d.SetId("")
	return nil
}

func microVMCheckpointStatusRefreshFunc(ctx context.Context, client *godo.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cp, _, err := client.MicroVMs.GetCheckpoint(ctx, id)
		if err != nil {
			return nil, "", err
		}
		if cp.Status == godo.MicroVMCheckpointStatusFailed {
			return cp, string(cp.Status), fmt.Errorf("MicroVM checkpoint %s entered failed state", id)
		}
		return cp, string(cp.Status), nil
	}
}

// microVMCheckpointDeleteRefreshFunc reports a 404 as CHECKPOINT_DELETED so
// the waiter finishes whether the API hides deleted checkpoints or keeps
// returning them with that status.
func microVMCheckpointDeleteRefreshFunc(ctx context.Context, client *godo.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cp, resp, err := client.MicroVMs.GetCheckpoint(ctx, id)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return id, string(godo.MicroVMCheckpointStatusDeleted), nil
			}
			return nil, "", err
		}
		return cp, string(cp.Status), nil
	}
}
