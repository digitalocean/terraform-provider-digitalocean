package microdroplet

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDigitalOceanMicroDroplet returns the digitalocean_microdroplet
// resource schema.
func ResourceDigitalOceanMicroDroplet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanMicroDropletCreate,
		ReadContext:   resourceDigitalOceanMicroDropletRead,
		UpdateContext: resourceDigitalOceanMicroDropletUpdate,
		DeleteContext: resourceDigitalOceanMicroDropletDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: microDropletResourceSchema(),
	}
}

func resourceDigitalOceanMicroDropletCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := &godo.MicroDropletCreateRequest{
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(string),
		Image:  d.Get("image").(string),
	}

	if v, ok := d.GetOk("networking"); ok {
		req.Networking = godo.MicroDropletNetworking(v.(string))
	}
	if v, ok := d.GetOk("vpc_uuid"); ok {
		req.VPCUUID = v.(string)
	}
	if v, ok := d.GetOk("http_port"); ok {
		req.HTTPPort = uint32(v.(int))
	}
	if v, ok := d.GetOk("http_protocol"); ok {
		req.HTTPProtocol = godo.MicroDropletHTTPProtocol(v.(string))
	}
	if v, ok := d.GetOk("environment"); ok {
		req.Environment = expandEnvironment(v)
	}
	if v, ok := d.GetOk("auto_pause"); ok {
		req.AutoPause = expandAutoPause(v)
	}
	if v, ok := d.GetOkExists("auto_resume"); ok {
		req.AutoResume = godo.PtrTo(v.(bool))
	}
	if v, ok := d.GetOk("tags"); ok {
		req.Tags = tag.ExpandTags(v.(*schema.Set).List())
	}

	log.Printf("[DEBUG] MicroDroplet create request: %+v", req)

	m, _, err := client.MicroDroplets.Create(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating MicroDroplet: %s", err)
	}

	d.SetId(m.ID)
	log.Printf("[INFO] MicroDroplet created, ID: %s", d.Id())

	if _, err := waitForMicroDropletState(
		ctx, client, d.Id(),
		godo.MicroDropletStateRunning,
		[]godo.MicroDropletState{godo.MicroDropletStateCreating},
		d.Timeout(schema.TimeoutCreate),
	); err != nil {
		return diag.Errorf("Error waiting for MicroDroplet (%s) to become running: %s", d.Id(), err)
	}

	// Config asked for paused — transition after create so we deterministically
	// observe the "running -> paused" flow rather than racing the API.
	if desired := d.Get("state").(string); desired == string(godo.MicroDropletStatePaused) {
		if err := transitionMicroDropletState(ctx, client, d.Id(), godo.MicroDropletStatePaused, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDigitalOceanMicroDropletRead(ctx, d, meta)
}

func resourceDigitalOceanMicroDropletRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	m, resp, err := client.MicroDroplets.Get(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[WARN] MicroDroplet (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving MicroDroplet: %s", err)
	}

	if err := setMicroDropletAttributes(d, m); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDigitalOceanMicroDropletUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChange("state") {
		target := godo.MicroDropletState(d.Get("state").(string))
		if err := transitionMicroDropletState(ctx, client, d.Id(), target, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDigitalOceanMicroDropletRead(ctx, d, meta)
}

func resourceDigitalOceanMicroDropletDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	resp, err := client.MicroDroplets.Delete(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.Errorf("Error deleting MicroDroplet: %s", err)
	}

	log.Printf("[INFO] MicroDroplet deleted, ID: %s", d.Id())
	d.SetId("")
	return nil
}

// transitionMicroDropletState issues a PATCH with the target state and waits
// for the API to reflect it. Used both by CreateContext (create-then-pause)
// and UpdateContext (user-driven pause/resume).
func transitionMicroDropletState(ctx context.Context, client *godo.Client, id string, target godo.MicroDropletState, timeout time.Duration) error {
	var pending []godo.MicroDropletState
	switch target {
	case godo.MicroDropletStatePaused:
		pending = []godo.MicroDropletState{godo.MicroDropletStateRunning, godo.MicroDropletStatePausing}
	case godo.MicroDropletStateRunning:
		pending = []godo.MicroDropletState{godo.MicroDropletStatePaused, godo.MicroDropletStateResuming}
	default:
		return fmt.Errorf("unsupported target MicroDroplet state %q", target)
	}

	log.Printf("[INFO] Transitioning MicroDroplet (%s) to state %s", id, target)
	_, _, err := client.MicroDroplets.Update(ctx, id, &godo.MicroDropletUpdateRequest{State: target})
	if err != nil {
		return fmt.Errorf("error transitioning MicroDroplet (%s) to %s: %w", id, target, err)
	}

	if _, err := waitForMicroDropletState(ctx, client, id, target, pending, timeout); err != nil {
		return fmt.Errorf("error waiting for MicroDroplet (%s) to reach state %s: %w", id, target, err)
	}
	return nil
}

// waitForMicroDropletState polls MicroDroplets.Get until the observed state
// matches target, or the timeout elapses. `pending` should list every state
// that is legally a step on the way to `target`.
func waitForMicroDropletState(ctx context.Context, client *godo.Client, id string, target godo.MicroDropletState, pending []godo.MicroDropletState, timeout time.Duration) (interface{}, error) {
	log.Printf("[INFO] Waiting for MicroDroplet (%s) to reach state %s", id, target)

	pendingStrs := make([]string, len(pending))
	for i, s := range pending {
		pendingStrs[i] = string(s)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    pendingStrs,
		Target:     []string{string(target)},
		Refresh:    microDropletStateRefreshFunc(ctx, client, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func microDropletStateRefreshFunc(ctx context.Context, client *godo.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		m, _, err := client.MicroDroplets.Get(ctx, id)
		if err != nil {
			return nil, "", err
		}
		if m.State == godo.MicroDropletStateFailed {
			return m, string(m.State), fmt.Errorf("MicroDroplet %s entered failed state", id)
		}
		return m, string(m.State), nil
	}
}
