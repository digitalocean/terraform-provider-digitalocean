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
			Required:     true,
			ForceNew:     true,
			Description:  "DigitalOcean region slug where the MicroDroplet is deployed",
			ValidateFunc: validation.NoZeroValues,
		},
		"size": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "MicroDroplet size slug",
			ValidateFunc: validation.NoZeroValues,
		},
		"image": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "MicroDroplet image UUID or URN",
			ValidateFunc: validation.NoZeroValues,
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
						Description: "Whether auto-pause is enabled",
					},
					"idle_timeout": {
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
						Description:  "Idle timeout as a Go duration string (e.g. '5m', '30s')",
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
		"endpoint": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Public endpoint URL for the MicroDroplet",
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
		clone.Computed = true
		// MaxItems/MinItems are only meaningful for configurable attributes;
		// the plugin SDK rejects them on Computed-only fields.
		clone.MaxItems = 0
		clone.MinItems = 0
		base[k] = &clone
	}
	// godo.MicroDroplet has no Tags field, so tags cannot round-trip through
	// the datasource. Dropping the attribute is more honest than exposing an
	// always-empty set that would silently break filter { key = "tags" }.
	delete(base, "tags")
	return base
}

// microDropletImageResourceSchema returns the resource-side schema used by
// ResourceDigitalOceanMicroDropletImage.
func microDropletImageResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Name of the MicroDroplet image",
			ValidateFunc: validation.NoZeroValues,
		},
		"source": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Source OCI reference for the MicroDroplet image",
			ValidateFunc: validation.NoZeroValues,
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Lifecycle status of the MicroDroplet image",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation timestamp for the MicroDroplet image",
		},
		"urn": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The uniform resource name (URN) for the MicroDroplet image",
		},
	}
}

func microDropletImageDataSourceSchema() map[string]*schema.Schema {
	base := microDropletImageResourceSchema()
	for k, v := range base {
		clone := *v
		clone.Required = false
		clone.Optional = false
		clone.ForceNew = false
		clone.Default = nil
		clone.ValidateFunc = nil
		clone.DiffSuppressFunc = nil
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

// setMicroDropletAttributes writes the state observed on a godo.MicroDroplet
// into the ResourceData without touching the settable `state` attribute
// (which reflects user intent, not observed state).
//
// `tags` are deliberately not written back: godo.MicroDroplet has no Tags
// field, so we cannot verify what the platform actually stored. The resource
// keeps whatever tags the user configured at create time in state; changing
// them recreates the resource (see tagsSchemaForceNew).
func setMicroDropletAttributes(d *schema.ResourceData, m *godo.MicroDroplet) error {
	if m == nil {
		return fmt.Errorf("cannot set attributes from nil MicroDroplet")
	}
	d.SetId(m.ID)
	d.Set("name", m.Name)
	d.Set("region", m.Region)
	d.Set("size", m.Size)
	d.Set("image", m.Image)
	d.Set("networking", string(m.Networking))
	d.Set("endpoint", m.Endpoint)
	d.Set("current_state", string(m.State))
	d.Set("created_at", m.Created)
	d.Set("urn", m.URN())

	if err := d.Set("auto_pause", flattenAutoPause(m.AutoPause)); err != nil {
		return fmt.Errorf("error setting auto_pause: %w", err)
	}
	if m.AutoResume != nil {
		d.Set("auto_resume", *m.AutoResume)
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
		"id":            m.ID,
		"name":          m.Name,
		"region":        m.Region,
		"size":          m.Size,
		"image":         m.Image,
		"networking":    string(m.Networking),
		"endpoint":      m.Endpoint,
		"current_state": string(m.State),
		"state":         string(m.State),
		"created_at":    m.Created,
		"urn":           m.URN(),
		"auto_pause":    flattenAutoPause(m.AutoPause),
	}
	if m.AutoResume != nil {
		out["auto_resume"] = *m.AutoResume
	} else {
		out["auto_resume"] = false
	}
	return out, nil
}

// flattenMicroDropletImage flattens a godo.MicroDropletImage for the datalist
// datasource.
func flattenMicroDropletImage(rawRecord, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	i, ok := rawRecord.(godo.MicroDropletImage)
	if !ok {
		return nil, fmt.Errorf("unexpected record type %T", rawRecord)
	}
	return map[string]interface{}{
		"id":         i.ID,
		"name":       i.Name,
		"source":     i.Source,
		"status":     string(i.Status),
		"created_at": i.Created,
		"urn":        i.URN(),
	}, nil
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

// getDigitalOceanMicroDropletImages is the GetRecords callback for the plural
// MicroDroplet image datasource.
func getDigitalOceanMicroDropletImages(meta interface{}, _ map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.ListOptions{Page: 1, PerPage: 200}
	var records []interface{}
	for {
		batch, resp, err := client.MicroDropletImages.List(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving MicroDroplet images: %w", err)
		}
		for _, i := range batch {
			records = append(records, i)
		}
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error paging MicroDroplet images: %w", err)
		}
		opts.Page = page + 1
	}
	return records, nil
}
