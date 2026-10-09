package microvm

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

// ResourceDigitalOceanMicroVM returns the digitalocean_microvm
// resource schema.
func ResourceDigitalOceanMicroVM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanMicroVMCreate,
		ReadContext:   resourceDigitalOceanMicroVMRead,
		UpdateContext: resourceDigitalOceanMicroVMUpdate,
		DeleteContext: resourceDigitalOceanMicroVMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: microVMResourceSchema(),
	}
}

func resourceDigitalOceanMicroVMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	src, err := expandSource(d.Get("source"))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &godo.MicroVMCreateRequest{
		Name:   d.Get("name").(string),
		Source: src,
	}

	if v, ok := d.GetOk("region"); ok {
		req.Region = v.(string)
	}
	if v, ok := d.GetOk("size"); ok {
		req.Size = expandSizeRequest(v)
	}
	if src.OCIRef != "" {
		if req.Region == "" {
			return diag.Errorf("region is required when source.oci_ref is set")
		}
		if req.Size == nil {
			return diag.Errorf("size is required when source.oci_ref is set")
		}
	}

	if v, ok := d.GetOk("networking"); ok {
		req.Networking = godo.MicroVMNetworking(v.(string))
	}
	if v, ok := d.GetOk("vpc_uuid"); ok {
		req.VPCUUID = v.(string)
	}
	if v, ok := d.GetOk("http_port"); ok {
		req.HTTPPort = uint32(v.(int))
	}
	if v, ok := d.GetOk("http_protocol"); ok {
		req.HTTPProtocol = godo.MicroVMHTTPProtocol(v.(string))
	}
	if v, ok := d.GetOk("ports"); ok {
		req.Ports = expandPorts(v)
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

	log.Printf("[DEBUG] MicroVM create request: %+v", req)

	m, _, err := client.MicroVMs.Create(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating MicroVM: %s", err)
	}

	d.SetId(m.ID)
	log.Printf("[INFO] MicroVM created, ID: %s", d.Id())

	if _, err := waitForMicroVMState(
		ctx, client, d.Id(),
		godo.MicroVMStateRunning,
		// `unknown` is a legal transient right after POST — treat it as
		// pending so the waiter doesn't bail with `unexpected state 'unknown'`.
		[]godo.MicroVMState{godo.MicroVMStateUnknown, godo.MicroVMStateCreating},
		d.Timeout(schema.TimeoutCreate),
	); err != nil {
		return diag.Errorf("Error waiting for MicroVM (%s) to become running: %s", d.Id(), err)
	}

	// Config asked for paused — transition after create so we deterministically
	// observe the "running -> paused" flow rather than racing the API.
	if desired := d.Get("state").(string); desired == string(godo.MicroVMStatePaused) {
		if err := transitionMicroVMState(ctx, client, d.Id(), godo.MicroVMStatePaused, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDigitalOceanMicroVMRead(ctx, d, meta)
}

func resourceDigitalOceanMicroVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	m, resp, err := client.MicroVMs.Get(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[WARN] MicroVM (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving MicroVM: %s", err)
	}

	if err := setMicroVMAttributes(d, m); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// resourceDigitalOceanMicroVMUpdate handles the only in-place mutation
// the MicroVMs API supports: the desired lifecycle `state` (running /
// paused) via the Pause and Resume action endpoints.
//
// Every other user-configurable attribute on the resource — auto_pause,
// auto_resume, tags, environment, http_port, http_protocol, networking, etc.
// — is marked ForceNew because godo's MicroVMsService exposes no update
// endpoint for them. Without ForceNew, plan changes would silently no-op on
// apply (clean plan, "1 changed", no API call, then Read overwrites state)
// causing invisible drift like a user lowering auto_pause.idle_timeout to
// cut cost but the platform keeping the old timeout.
func resourceDigitalOceanMicroVMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChange("state") {
		target := godo.MicroVMState(d.Get("state").(string))
		if err := transitionMicroVMState(ctx, client, d.Id(), target, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDigitalOceanMicroVMRead(ctx, d, meta)
}

func resourceDigitalOceanMicroVMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	resp, err := client.MicroVMs.Delete(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.Errorf("Error deleting MicroVM: %s", err)
	}

	log.Printf("[INFO] MicroVM deleted, ID: %s", d.Id())
	d.SetId("")
	return nil
}

// transitionMicroVMState invokes the dedicated Pause / Resume action
// endpoint on the MicroVMs API and then polls until the observed state
// matches. Used both by CreateContext (create-then-pause) and UpdateContext
// (user-driven pause/resume).
//
// godo's Pause / Resume are documented as synchronous, but we still poll
// afterwards as a belt-and-suspenders check: fake API servers used in tests
// and edge cases in the platform can return before the state is fully
// reflected on Get.
func transitionMicroVMState(ctx context.Context, client *godo.Client, id string, target godo.MicroVMState, timeout time.Duration) error {
	var pending []godo.MicroVMState
	switch target {
	case godo.MicroVMStatePaused:
		// Include `unknown` — the API can briefly report it between the
		// action returning and the state settling on Get.
		pending = []godo.MicroVMState{godo.MicroVMStateUnknown, godo.MicroVMStateRunning, godo.MicroVMStatePausing}
		log.Printf("[INFO] Pausing MicroVM (%s)", id)
		if _, _, err := client.MicroVMs.Pause(ctx, id); err != nil {
			return fmt.Errorf("error pausing MicroVM (%s): %w", id, err)
		}
	case godo.MicroVMStateRunning:
		pending = []godo.MicroVMState{godo.MicroVMStateUnknown, godo.MicroVMStatePaused, godo.MicroVMStateResuming}
		log.Printf("[INFO] Resuming MicroVM (%s)", id)
		if _, _, err := client.MicroVMs.Resume(ctx, id); err != nil {
			return fmt.Errorf("error resuming MicroVM (%s): %w", id, err)
		}
	default:
		return fmt.Errorf("unsupported target MicroVM state %q", target)
	}

	if _, err := waitForMicroVMState(ctx, client, id, target, pending, timeout); err != nil {
		return fmt.Errorf("error waiting for MicroVM (%s) to reach state %s: %w", id, target, err)
	}
	return nil
}

// waitForMicroVMState polls MicroVMs.Get until the observed state
// matches target, or the timeout elapses. `pending` should list every state
// that is legally a step on the way to `target`.
func waitForMicroVMState(ctx context.Context, client *godo.Client, id string, target godo.MicroVMState, pending []godo.MicroVMState, timeout time.Duration) (interface{}, error) {
	log.Printf("[INFO] Waiting for MicroVM (%s) to reach state %s", id, target)

	pendingStrs := make([]string, len(pending))
	for i, s := range pending {
		pendingStrs[i] = string(s)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    pendingStrs,
		Target:     []string{string(target)},
		Refresh:    microVMStateRefreshFunc(ctx, client, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func microVMStateRefreshFunc(ctx context.Context, client *godo.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		m, _, err := client.MicroVMs.Get(ctx, id)
		if err != nil {
			return nil, "", err
		}
		if m.State == godo.MicroVMStateFailed {
			reason := m.FailureReason
			if reason == "" {
				reason = "unknown"
			}
			return m, string(m.State), fmt.Errorf("MicroVM %s entered failed state: %s", id, reason)
		}
		return m, string(m.State), nil
	}
}
