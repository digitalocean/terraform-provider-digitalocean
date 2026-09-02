package microdroplet

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// stateValues is the closed set accepted by the settable `state` attribute on
// digitalocean_microdroplet. The API also exposes transient values (creating,
// pausing, resuming, ...) but users should never set those directly.
var stateValues = []string{
	string(godo.MicroDropletStateRunning),
	string(godo.MicroDropletStatePaused),
}

// networkingValues enumerates the accepted values for the `networking`
// attribute on digitalocean_microdroplet.
var networkingValues = []string{
	string(godo.MicroDropletNetworkingPublic),
	string(godo.MicroDropletNetworkingVPC),
}

// httpProtocolValues enumerates the accepted values for the `http_protocol`
// attribute on digitalocean_microdroplet. The control plane accepts only
// `http` (HTTP/1.1) and `http2`; `https` is not a valid platform value.
var httpProtocolValues = []string{
	string(godo.MicroDropletHTTPProtocolHTTP),
	string(godo.MicroDropletHTTPProtocolHTTP2),
}

// tagsSchemaForceNew returns tag.TagsSchema() with ForceNew set. Tags are
// accepted by MicroDropletCreateRequest but the MicroDroplets API exposes no
// endpoint to mutate them afterwards, so any change has to recreate the
// resource. Marking ForceNew keeps Terraform's plan honest — without it,
// changes would silently no-op on apply.
func tagsSchemaForceNew() *schema.Schema {
	s := tag.TagsSchema()
	s.ForceNew = true
	return s
}

// microDropletResourceSchema returns the resource-side schema used by
// ResourceDigitalOceanMicroDroplet. The data source schemas reuse this via
// microDropletDataSourceSchema which recasts everything to Computed.
func microDropletResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Name of the MicroDroplet",
			ValidateFunc: validation.NoZeroValues,
		},
		"region": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			Description:  "DigitalOcean region slug. Required when creating from oci_ref; optional when restoring from a checkpoint (inherited).",
			ValidateFunc: validation.NoZeroValues,
		},
		"size": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "Compute size. Required when creating from oci_ref; optional when restoring from a checkpoint (inherited).",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cpu": {
						Type:         schema.TypeInt,
						Required:     true,
						ForceNew:     true,
						Description:  "Number of vCPUs",
						ValidateFunc: validation.IntAtLeast(1),
					},
					"memory": {
						Type:         schema.TypeInt,
						Required:     true,
						ForceNew:     true,
						Description:  "Memory in MiB",
						ValidateFunc: validation.IntAtLeast(1),
					},
					"disk": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Attached disk in GB (provisioned with the size)",
					},
				},
			},
		},
		"source": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "Workload source. Exactly one of oci_ref or checkpoint_id must be set.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"oci_ref": {
						Type:          schema.TypeString,
						Optional:      true,
						ForceNew:      true,
						ConflictsWith: []string{"source.0.checkpoint_id"},
						Description:   "OCI reference for the workload container",
						ValidateFunc:  validation.NoZeroValues,
					},
					"checkpoint_id": {
						Type:          schema.TypeString,
						Optional:      true,
						ForceNew:      true,
						ConflictsWith: []string{"source.0.oci_ref"},
						Description:   "Checkpoint UUID to restore",
						ValidateFunc:  validation.NoZeroValues,
					},
				},
			},
		},
		"networking": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			Description:  "Networking mode: 'public' or 'vpc'",
			ValidateFunc: validation.StringInSlice(networkingValues, false),
		},
		"vpc_uuid": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			Description:  "UUID of the VPC to attach the MicroDroplet to. Only valid when networking is 'vpc'.",
			ValidateFunc: validation.NoZeroValues,
		},
		"http_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Description:  "Port the MicroDroplet exposes over HTTP",
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"http_protocol": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "HTTP protocol: 'http' or 'http2'",
			ValidateFunc: validation.StringInSlice(httpProtocolValues, false),
		},
		"ports": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Guest ports open for ingress. Defaults to just http_port when omitted.",
			Elem: &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
		},
		"environment": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "Environment variables passed to the MicroDroplet",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"auto_pause": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "Auto-pause configuration. Forces recreation on change: the MicroDroplets API has no in-place update path for auto_pause.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Required:    true,
						ForceNew:    true,
						Description: "Whether auto-pause is enabled. Forces recreation on change (no in-place API path).",
					},
					"idle_timeout": {
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						Description:  "Idle timeout as a Go duration string (e.g. '5m', '30s'). Forces recreation on change (no in-place API path).",
						ValidateFunc: validation.NoZeroValues,
					},
				},
			},
		},
		"auto_resume": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Whether the MicroDroplet should auto-resume on request. Forces recreation on change: the MicroDroplets API has no in-place update path for auto_resume.",
		},
		"tags": tagsSchemaForceNew(),
		"state": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          string(godo.MicroDropletStateRunning),
			Description:      "Desired lifecycle state: 'running' or 'paused'. Changes are applied by calling the microdroplet pause / resume action endpoints.",
			ValidateFunc:     validation.StringInSlice(stateValues, false),
			DiffSuppressFunc: suppressStateDiffWhenAutoPause,
		},
		"current_state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Observed lifecycle state of the MicroDroplet",
		},
		"failure_reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Human-readable explanation when current_state is failed",
		},
		"urls": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Ingress URLs for the MicroDroplet",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hostname": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Hostname (no scheme)",
					},
					"port": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Guest port this URL forwards to",
					},
					"default": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "Whether this is the system default URL",
					},
					"status": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "URL lifecycle status (PENDING or ACTIVE)",
					},
				},
			},
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation timestamp for the MicroDroplet",
		},
		"urn": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The uniform resource name (URN) for the MicroDroplet",
		},
	}
}

// microDropletDataSourceSchema returns the resource schema recast so every
// attribute is Computed and safe to expose on the datasource. Filter/select
// attributes (id, name) are re-set to Optional+Computed by the caller.
func microDropletDataSourceSchema() map[string]*schema.Schema {
	base := microDropletResourceSchema()
	for k, v := range base {
		clone := *v
		clone.Required = false
		clone.Optional = false
		clone.ForceNew = false
		clone.Default = nil
		clone.ValidateFunc = nil
		clone.DiffSuppressFunc = nil
		clone.ConflictsWith = nil
		clone.Computed = true
		clone.MaxItems = 0
		clone.MinItems = 0
		base[k] = &clone
	}
	return base
}

// suppressStateDiffWhenAutoPause suppresses spurious `state` diffs when
// auto_pause is enabled: the API can flip the observed state to `paused` at
// any time, and Terraform should not fight the platform in that case. The
// diff is only suppressed when moving from `paused` (observed) to `running`
// (config) — the reverse (user explicitly requesting `paused`) always applies.
func suppressStateDiffWhenAutoPause(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if !autoPauseEnabled(d) {
		return false
	}
	return oldValue == string(godo.MicroDropletStatePaused) &&
		newValue == string(godo.MicroDropletStateRunning)
}

// autoPauseEnabled returns true when the resource config declares an
// `auto_pause` block with `enabled = true`.
func autoPauseEnabled(d *schema.ResourceData) bool {
	raw, ok := d.GetOk("auto_pause")
	if !ok {
		return false
	}
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return false
	}
	entry, ok := list[0].(map[string]interface{})
	if !ok {
		return false
	}
	enabled, _ := entry["enabled"].(bool)
	return enabled
}

// expandAutoPause turns the `auto_pause` HCL block into a godo AutoPauseConfig.
// Returns nil when the block is empty.
func expandAutoPause(raw interface{}) *godo.AutoPauseConfig {
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	entry, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := &godo.AutoPauseConfig{}
	if v, ok := entry["enabled"].(bool); ok {
		cfg.Enabled = godo.PtrTo(v)
	}
	if v, ok := entry["idle_timeout"].(string); ok && v != "" {
		cfg.IdleTimeout = v
	}
	return cfg
}

// flattenAutoPause converts a godo AutoPauseConfig back into the list-of-map
// shape that Terraform expects for a single-item block.
func flattenAutoPause(cfg *godo.AutoPauseConfig) []interface{} {
	if cfg == nil {
		return nil
	}
	entry := map[string]interface{}{}
	if cfg.Enabled != nil {
		entry["enabled"] = *cfg.Enabled
	} else {
		entry["enabled"] = false
	}
	entry["idle_timeout"] = cfg.IdleTimeout
	return []interface{}{entry}
}

// expandEnvironment converts the `environment` map[string]interface{} that
// Terraform produces into the map[string]string that godo expects.
func expandEnvironment(raw interface{}) map[string]string {
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}

func expandSource(raw interface{}) (*godo.MicroDropletSource, error) {
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("source is required")
	}
	entry, ok := list[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid source block")
	}
	ociRef, _ := entry["oci_ref"].(string)
	checkpointID, _ := entry["checkpoint_id"].(string)
	if (ociRef == "") == (checkpointID == "") {
		return nil, fmt.Errorf("source must set exactly one of oci_ref or checkpoint_id")
	}
	src := &godo.MicroDropletSource{}
	if ociRef != "" {
		src.OCIRef = ociRef
	} else {
		src.CheckpointID = checkpointID
	}
	return src, nil
}

func flattenSource(src *godo.MicroDropletSource) []interface{} {
	if src == nil {
		return nil
	}
	entry := map[string]interface{}{}
	if src.OCIRef != "" {
		entry["oci_ref"] = src.OCIRef
	}
	if src.CheckpointID != "" {
		entry["checkpoint_id"] = src.CheckpointID
	}
	return []interface{}{entry}
}

func expandSizeRequest(raw interface{}) *godo.MicroDropletSizeRequest {
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	entry, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return &godo.MicroDropletSizeRequest{
		CPU:    uint32(entry["cpu"].(int)),
		Memory: uint32(entry["memory"].(int)),
	}
}

func flattenSize(size *godo.MicroDropletSize) []interface{} {
	if size == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"cpu":    int(size.CPU),
		"memory": int(size.Memory),
		"disk":   int(size.Disk),
	}}
}

func expandPorts(raw interface{}) []uint32 {
	set, ok := raw.(*schema.Set)
	if !ok || set == nil || set.Len() == 0 {
		return nil
	}
	ports := make([]uint32, 0, set.Len())
	for _, v := range set.List() {
		ports = append(ports, uint32(v.(int)))
	}
	return ports
}

func flattenPorts(ports []uint32) *schema.Set {
	vals := make([]interface{}, len(ports))
	for i, p := range ports {
		vals[i] = int(p)
	}
	return schema.NewSet(schema.HashInt, vals)
}

func flattenURLs(urls []godo.MicroDropletURL) []interface{} {
	out := make([]interface{}, len(urls))
	for i, u := range urls {
		out[i] = map[string]interface{}{
			"hostname": u.Hostname,
			"port":     u.Port,
			"default":  u.Default,
			"status":   string(u.Status),
		}
	}
	return out
}

// setMicroDropletAttributes writes the state observed on a godo.MicroDroplet
// into the ResourceData without touching the settable `state` attribute
// (which reflects user intent, not observed state).
func setMicroDropletAttributes(d *schema.ResourceData, m *godo.MicroDroplet) error {
	if m == nil {
		return fmt.Errorf("cannot set attributes from nil MicroDroplet")
	}
	d.SetId(m.ID)
	d.Set("name", m.Name)
	d.Set("region", m.Region)
	d.Set("networking", string(m.Networking))
	d.Set("current_state", string(m.State))
	d.Set("failure_reason", m.FailureReason)
	d.Set("created_at", m.Created)
	d.Set("urn", m.URN())

	if err := d.Set("size", flattenSize(m.Size)); err != nil {
		return fmt.Errorf("error setting size: %w", err)
	}
	if err := d.Set("source", flattenSource(m.Source)); err != nil {
		return fmt.Errorf("error setting source: %w", err)
	}
	if err := d.Set("urls", flattenURLs(m.URLs)); err != nil {
		return fmt.Errorf("error setting urls: %w", err)
	}
	if err := d.Set("ports", flattenPorts(m.Ports)); err != nil {
		return fmt.Errorf("error setting ports: %w", err)
	}
	if err := d.Set("auto_pause", flattenAutoPause(m.AutoPause)); err != nil {
		return fmt.Errorf("error setting auto_pause: %w", err)
	}
	if m.AutoResume != nil {
		d.Set("auto_resume", *m.AutoResume)
	}
	if err := d.Set("tags", tag.FlattenTags(m.Tags)); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}
	return nil
}

// flattenMicroDroplet flattens a godo.MicroDroplet into the map shape the
// datalist datasource expects.
func flattenMicroDroplet(rawRecord, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	m, ok := rawRecord.(godo.MicroDroplet)
	if !ok {
		return nil, fmt.Errorf("unexpected record type %T", rawRecord)
	}
	out := map[string]interface{}{
		"id":             m.ID,
		"name":           m.Name,
		"region":         m.Region,
		"size":           flattenSize(m.Size),
		"source":         flattenSource(m.Source),
		"networking":     string(m.Networking),
		"urls":           flattenURLs(m.URLs),
		"ports":          flattenPorts(m.Ports),
		"failure_reason": m.FailureReason,
		"current_state":  string(m.State),
		"state":          string(m.State),
		"created_at":     m.Created,
		"urn":            m.URN(),
		"auto_pause":     flattenAutoPause(m.AutoPause),
		"tags":           tag.FlattenTags(m.Tags),
	}
	if m.AutoResume != nil {
		out["auto_resume"] = *m.AutoResume
	} else {
		out["auto_resume"] = false
	}
	return out, nil
}

// getDigitalOceanMicroDroplets is the GetRecords callback for the plural
// MicroDroplet datasource. It supports optional region and name filters via
// the ExtraQuerySchema.
func getDigitalOceanMicroDroplets(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	region, _ := extra["region"].(string)
	name, _ := extra["name"].(string)

	opts := &godo.ListOptions{Page: 1, PerPage: 200}

	var records []interface{}
	for {
		var (
			batch []godo.MicroDroplet
			resp  *godo.Response
			err   error
		)
		switch {
		case region != "":
			batch, resp, err = client.MicroDroplets.ListByRegion(context.Background(), region, opts)
		case name != "":
			batch, resp, err = client.MicroDroplets.ListByName(context.Background(), name, opts)
		default:
			batch, resp, err = client.MicroDroplets.List(context.Background(), opts)
		}
		if err != nil {
			return nil, fmt.Errorf("error retrieving MicroDroplets: %w", err)
		}
		for _, m := range batch {
			records = append(records, m)
		}
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error paging MicroDroplets: %w", err)
		}
		opts.Page = page + 1
	}
	return records, nil
}
