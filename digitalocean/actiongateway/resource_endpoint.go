package actiongateway

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	endpointPolicyAllow = "allow"
	endpointPolicyAsk   = "ask"
	endpointPolicyDeny  = "deny"

	// toolReferenceKindToolbelt is the REST kind of a version tool reference
	// that names a toolbelt rather than a single tool.
	toolReferenceKindToolbelt = "SESSION_TOOL_REFERENCE_KIND_TOOLBELT"
)

var endpointNameRegexp = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// validateEndpointName matches the registry's endpoint name contract, so a bad
// name fails at plan time instead of as a 400 on apply.
var validateEndpointName = validation.StringMatch(endpointNameRegexp, "must contain 1 to 64 letters, digits, dots, underscores or hyphens")

// ResourceDigitalOceanActionGatewayEndpoint returns the
// digitalocean_action_gateway_endpoint resource. The endpoint owns the version
// its configuration describes: changing the version block publishes a new head
// version in place, and the previous version stays live.
func ResourceDigitalOceanActionGatewayEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanActionGatewayEndpointCreate,
		ReadContext:   resourceDigitalOceanActionGatewayEndpointRead,
		UpdateContext: resourceDigitalOceanActionGatewayEndpointUpdate,
		DeleteContext: resourceDigitalOceanActionGatewayEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: endpointResourceSchema(),
	}
}

func endpointResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Name of the endpoint, unique among your team's endpoints",
			ValidateFunc: validateEndpointName,
		},
		"insights": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "Insights consent, which applies to every version of the endpoint. An omitted signal keeps its stored value.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"metrics": {
						Type:        schema.TypeBool,
						Optional:    true,
						Computed:    true,
						Description: "Whether the endpoint's metrics are sent to Insights",
					},
					"logs": {
						Type:        schema.TypeBool,
						Optional:    true,
						Computed:    true,
						Description: "Whether the endpoint's logs are sent to Insights. Your team can require this to be on.",
					},
					"traces": {
						Type:        schema.TypeBool,
						Optional:    true,
						Computed:    true,
						Description: "Whether the endpoint's traces are sent to Insights",
					},
				},
			},
		},
		"version": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "Version the endpoint's configuration describes. Changing it publishes a new head version.",
			Elem: &schema.Resource{
				Schema: endpointVersionSchema(),
			},
		},
		"mcp_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MCP URL that follows the endpoint's current head version",
		},
		"version_mcp_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MCP URL pinned to the current head version",
		},
		"is_default": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether this is the default endpoint that a generic /mcp request resolves to for your user",
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func endpointVersionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"actor_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "Actor the version's tool calls act as. No operation changes an endpoint's actor, so changing this replaces the endpoint.",
			ValidateFunc: validateUserID,
		},
		"tools": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Tools and toolbelts the version enables, each optionally pinned as <tool_slug>@<version> or toolbelt:<name>@<version>. Omit to enable every tool; set an empty list to enable none.",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
		"pinned_tools": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Tools and toolbelts the head version enables, each pinned by the API as <tool_slug>@<version> or toolbelt:<name>@<version>. Empty when every tool is enabled.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"policy": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Tool permission policy. Omitting it asks before every call.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_action": {
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
						Description:  "Action for calls no rule matches: allow, ask, or deny. Empty means ask.",
						ValidateFunc: validation.StringInSlice([]string{endpointPolicyAllow, endpointPolicyAsk, endpointPolicyDeny}, false),
					},
					"rules": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Rules for specific tools and toolbelts, in the order the configuration lists them",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"tool": {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "Tool slug, optionally pinned as <tool_slug>@<version>, or toolbelt:<name>@<version>",
									ValidateFunc: validation.NoZeroValues,
								},
								"action": {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "allow, ask, or deny",
									ValidateFunc: validation.StringInSlice([]string{endpointPolicyAllow, endpointPolicyAsk, endpointPolicyDeny}, false),
								},
								"match": {
									Type:        schema.TypeMap,
									Optional:    true,
									Description: "Argument values a call must match for the rule to apply, keyed by argument name",
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
		},
		"config": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "JSON object of options Gateway applies to the version, such as preloadTools or customInstructions",
			ValidateFunc:     validation.StringIsJSON,
			StateFunc:        normalizeJSONStateFunc,
			DiffSuppressFunc: structure.SuppressJsonDiff,
		},
		"network": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "VPC the version's tool calls are routed through",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"vpc_uuid": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "UUID of a VPC owned by your team",
						ValidateFunc: validation.IsUUID,
					},
					"region": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Region of the VPC",
					},
				},
			},
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the head version, <endpoint UUID>.v<version number>",
		},
		"version_number": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Version number of the head version",
		},
		"promoted_from_version_number": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Version number the head was promoted from, when it was created by a promotion rather than a publish",
		},
	}
}

func normalizeJSONStateFunc(v interface{}) string {
	normalized, err := structure.NormalizeJsonString(v)
	if err != nil {
		return v.(string)
	}
	return normalized
}

// endpointVersionSpecChanged reports whether the configuration describes a
// different version than the one in state. Only the configured members are
// compared: the computed ones, such as the network region and the policy's
// resolved default action, become unknown on an unrelated change and must not
// turn that change into a publish.
func endpointVersionSpecChanged(d *schema.ResourceData) bool {
	return d.HasChanges(
		"version.0.actor_id",
		"version.0.tools",
		"version.0.policy.0.default_action",
		"version.0.policy.0.rules",
		"version.0.config",
		"version.0.network.0.vpc_uuid",
	)
}

func resourceDigitalOceanActionGatewayEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := &godo.ActionGatewayEndpointCreateRequest{
		Name:     d.Get("name").(string),
		Version:  expandEndpointVersion(d),
		Insights: expandEndpointInsights(d.GetRawConfig()),
	}

	created, _, err := client.ActionGateway.Endpoints.Create(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating Action Gateway endpoint: %s", err)
	}
	if created.Endpoint == nil || created.Endpoint.ID == "" {
		return diag.Errorf("Error creating Action Gateway endpoint: response contained no endpoint")
	}

	d.SetId(created.Endpoint.ID)
	log.Printf("[INFO] Action Gateway endpoint created: %s", d.Id())

	return resourceDigitalOceanActionGatewayEndpointRead(ctx, d, meta)
}

func resourceDigitalOceanActionGatewayEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	found, resp, err := client.ActionGateway.Endpoints.Get(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Action Gateway endpoint (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Action Gateway endpoint: %s", err)
	}
	endpoint := found.Endpoint
	if endpoint == nil {
		d.SetId("")
		return nil
	}

	// A version published outside Terraform has to appear as a diff, so the
	// resolved head's specification becomes the version block. The resource
	// never adopts an out-of-band configuration silently and never deletes a
	// version it did not publish.
	values := flattenEndpoint(endpoint)
	values["version"].([]interface{})[0].(map[string]interface{})["tools"] = reconcileEndpointTools(d, endpoint.HeadVersion)
	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	return nil
}

func resourceDigitalOceanActionGatewayEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	id := d.Id()

	if d.HasChanges("name", "insights") {
		req := &godo.ActionGatewayEndpointUpdateRequest{}
		if d.HasChange("name") {
			req.Name = d.Get("name").(string)
		}
		if d.HasChange("insights") {
			req.Insights = expandEndpointInsights(d.GetRawConfig())
		}
		// The API rejects an update that changes neither, which is reachable
		// only when the configuration repeated the stored consent exactly.
		if req.Name != "" || req.Insights != nil {
			if _, _, err := client.ActionGateway.Endpoints.Update(ctx, id, req); err != nil {
				return diag.Errorf("Error updating Action Gateway endpoint (%s): %s", id, err)
			}
		}
	}

	if endpointVersionSpecChanged(d) {
		if _, _, err := client.ActionGateway.Endpoints.CreateVersion(ctx, id, expandEndpointVersion(d)); err != nil {
			return diag.Errorf("Error publishing Action Gateway endpoint version (%s): %s", id, err)
		}
		log.Printf("[INFO] Action Gateway endpoint version published: %s", id)
	}

	return resourceDigitalOceanActionGatewayEndpointRead(ctx, d, meta)
}

func resourceDigitalOceanActionGatewayEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	resp, err := client.ActionGateway.Endpoints.Delete(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting Action Gateway endpoint: %s", err)
	}

	log.Printf("[INFO] Action Gateway endpoint deleted: %s", d.Id())
	d.SetId("")
	return nil
}

func flattenEndpoint(endpoint *godo.ActionGatewayEndpoint) map[string]interface{} {
	head := endpoint.HeadVersion
	version := map[string]interface{}{}
	versionMCPURL := ""

	if head != nil {
		pinned := flattenToolReferences(head.Tools)
		version = map[string]interface{}{
			// The endpoint's actor is authoritative and set once, at creation;
			// a version repeats it. Using it here keeps the value stable even
			// when every live version was deleted out of band.
			"actor_id":                     endpoint.ActorID,
			"pinned_tools":                 pinned,
			"tools":                        pinned,
			"policy":                       flattenEndpointPolicy(head.Policy),
			"config":                       formatEndpointConfig(head.Config),
			"network":                      flattenEndpointNetwork(head.Network),
			"id":                           head.ID,
			"version_number":               int(head.VersionNumber),
			"promoted_from_version_number": int(head.PromotedFromVersionNumber),
		}
		versionMCPURL = head.MCPURL
	} else {
		version["actor_id"] = endpoint.ActorID
	}

	return map[string]interface{}{
		"name":            endpoint.Name,
		"insights":        flattenEndpointInsights(endpoint.Insights),
		"version":         []interface{}{version},
		"mcp_url":         endpoint.MCPURL,
		"version_mcp_url": versionMCPURL,
		"is_default":      endpoint.IsDefault,
		"created_at":      formatTimestamp(endpoint.CreatedAt),
		"updated_at":      formatTimestamp(endpoint.UpdatedAt),
	}
}

// reconcileEndpointTools reports the head's pinned selection in the spelling
// the configuration used, so an unpinned tool does not show a diff against the
// version the API pinned it to. An omitted selection leaves the API's own
// spelling, which is also what an imported endpoint reports.
func reconcileEndpointTools(d *schema.ResourceData, head *godo.ActionGatewayEndpointVersion) []string {
	var pinned []string
	if head != nil {
		pinned = flattenToolReferences(head.Tools)
	}
	configured, known := versionToolsSelection(d.GetRawConfig())
	if !known {
		// Without a configuration the argument itself is the only record of
		// what was selected, and an import leaves it empty.
		configured = expandStringSet(d.Get("version.0.tools"))
	}
	return reconcilePinnedTools(configured, pinned)
}

func flattenEndpointInsights(insights *godo.ActionGatewaySessionInsights) []interface{} {
	if insights == nil || (insights.Metrics == nil && insights.Logs == nil && insights.Traces == nil) {
		return []interface{}{}
	}
	out := map[string]interface{}{}
	if insights.Metrics != nil {
		out["metrics"] = *insights.Metrics
	}
	if insights.Logs != nil {
		out["logs"] = *insights.Logs
	}
	if insights.Traces != nil {
		out["traces"] = *insights.Traces
	}
	return []interface{}{out}
}

func flattenEndpointNetwork(network *godo.ActionGatewayEndpointNetworkResponse) []interface{} {
	if network == nil {
		return []interface{}{}
	}
	return []interface{}{map[string]interface{}{
		"vpc_uuid": network.VPCUUID,
		"region":   network.Region,
	}}
}

func flattenEndpointPolicy(policy *godo.ActionGatewaySessionPolicy) []interface{} {
	// An omitted policy asks before every call, and the API reports that as an
	// empty object rather than an absent one. Collapsing it keeps a
	// configuration that omits the block from diffing forever.
	if policy == nil || (policy.DefaultAction == "" && len(policy.Rules) == 0) {
		return []interface{}{}
	}
	return []interface{}{map[string]interface{}{
		"default_action": normalizePolicyAction(policy.DefaultAction),
		"rules":          flattenEndpointPolicyRules(policy.Rules),
	}}
}

func flattenEndpointPolicyRules(rules []godo.ActionGatewaySessionPolicyRule) []interface{} {
	out := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		match := make(map[string]interface{}, len(rule.Match))
		for k, v := range rule.Match {
			match[k] = v
		}
		out = append(out, map[string]interface{}{
			"tool":   rule.Tool,
			"action": normalizePolicyAction(rule.Action),
			"match":  match,
		})
	}
	return out
}

// normalizePolicyAction reports an action in the lowercase vocabulary the
// public API documents, so a deprecated prefixed spelling in a response cannot
// show a diff against a lowercase configuration. An unset action becomes empty,
// because the API reports omission as the unspecified enum value.
func normalizePolicyAction(action string) string {
	normalized := strings.ToLower(strings.TrimPrefix(action, "SESSION_POLICY_ACTION_"))
	if normalized == "unspecified" {
		return ""
	}
	return normalized
}

func formatEndpointConfig(config json.RawMessage) string {
	trimmed := strings.TrimSpace(string(config))
	if trimmed == "" || trimmed == "null" {
		return ""
	}
	return trimmed
}

// flattenToolReferences renders a head version's selected tools and toolbelts
// in the spelling the version spec accepts.
func flattenToolReferences(tools *godo.ActionGatewayEndpointVersionToolSet) []string {
	if tools == nil {
		return nil
	}
	out := make([]string, 0, len(tools.References))
	for _, ref := range tools.References {
		prefix := ""
		if ref.Kind == toolReferenceKindToolbelt {
			prefix = "toolbelt:"
		}
		pinned := prefix + ref.Name
		if ref.Version != "" {
			pinned += "@" + ref.Version
		}
		out = append(out, pinned)
	}
	sort.Strings(out)
	return out
}

// expandEndpointVersion builds the version spec the configuration describes.
func expandEndpointVersion(d *schema.ResourceData) *godo.ActionGatewayEndpointVersionSpec {
	spec := &godo.ActionGatewayEndpointVersionSpec{
		ActorID: d.Get("version.0.actor_id").(string),
		Tools:   expandEndpointVersionTools(d),
		Policy:  expandEndpointPolicy(d.Get("version.0.policy")),
		Network: expandEndpointNetwork(d.Get("version.0.network")),
	}
	if config := strings.TrimSpace(d.Get("version.0.config").(string)); config != "" {
		spec.Config = json.RawMessage(config)
	}
	return spec
}

// expandEndpointVersionTools returns a nil selection when the configuration
// omits tools, which enables every tool, and a possibly empty one when it names
// a selection, which enables exactly those. d.Get cannot tell the two apart
// because an omitted set and an empty set read as the same empty set.
func expandEndpointVersionTools(d *schema.ResourceData) []string {
	if tools, known := versionToolsSelection(d.GetRawConfig()); known {
		return tools
	}
	return expandStringSet(d.Get("version.0.tools"))
}

func expandEndpointPolicy(v interface{}) *godo.ActionGatewaySessionPolicy {
	raw, _ := v.([]interface{})
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	policy, ok := raw[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return &godo.ActionGatewaySessionPolicy{
		DefaultAction: policy["default_action"].(string),
		Rules:         expandEndpointPolicyRules(policy["rules"]),
	}
}

func expandEndpointPolicyRules(v interface{}) []godo.ActionGatewaySessionPolicyRule {
	raw, _ := v.([]interface{})
	out := make([]godo.ActionGatewaySessionPolicyRule, 0, len(raw))
	for _, item := range raw {
		rule, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, godo.ActionGatewaySessionPolicyRule{
			Tool:   rule["tool"].(string),
			Action: rule["action"].(string),
			Match:  expandStringMap(rule["match"]),
		})
	}
	return out
}

func expandEndpointNetwork(v interface{}) *godo.ActionGatewayEndpointNetwork {
	raw, _ := v.([]interface{})
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	network, ok := raw[0].(map[string]interface{})
	if !ok {
		return nil
	}
	if uuid, ok := network["vpc_uuid"].(string); ok && uuid != "" {
		return &godo.ActionGatewayEndpointNetwork{VPCUUID: uuid}
	}
	return nil
}

func expandStringMap(v interface{}) map[string]string {
	raw, _ := v.(map[string]interface{})
	if len(raw) == 0 {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, item := range raw {
		if s, ok := item.(string); ok {
			out[k] = s
		}
	}
	return out
}

// expandEndpointInsights returns the consent signals the configuration sets. A
// nil signal keeps the stored value: the API applies the team's stored consent
// to it, and rejects turning logs off for a team that requires them.
func expandEndpointInsights(raw cty.Value) *godo.ActionGatewaySessionInsights {
	block, ok := rawBlock(raw, "insights")
	if !ok {
		return nil
	}
	out := &godo.ActionGatewaySessionInsights{}
	found := false
	for name, target := range map[string]**bool{
		"metrics": &out.Metrics,
		"logs":    &out.Logs,
		"traces":  &out.Traces,
	} {
		value, ok := rawAttr(block, name)
		if !ok || value.IsNull() || !value.IsKnown() || value.Type() != cty.Bool {
			continue
		}
		enabled := value.True()
		*target = &enabled
		found = true
	}
	if !found {
		return nil
	}
	return out
}

// rawAttr returns the named attribute of a configuration object. known is
// false when the configuration has no such attribute, which is what a provider
// sees with no configuration, such as a refresh after an import.
func rawAttr(raw cty.Value, name string) (cty.Value, bool) {
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsObjectType() || !raw.Type().HasAttribute(name) {
		return cty.NilVal, false
	}
	return raw.GetAttr(name), true
}

// rawBlock returns the first element of the single nested block the
// configuration sets under name. known is false when the block is omitted.
func rawBlock(raw cty.Value, name string) (cty.Value, bool) {
	value, ok := rawAttr(raw, name)
	if !ok || value.IsNull() || !value.IsKnown() {
		return cty.NilVal, false
	}
	if !value.Type().IsListType() && !value.Type().IsSetType() {
		return cty.NilVal, false
	}
	if value.LengthInt() == 0 {
		return cty.NilVal, false
	}
	first := value.Index(cty.NumberIntVal(0))
	if first.IsNull() || !first.IsKnown() {
		return cty.NilVal, false
	}
	return first, true
}

// versionToolsSelection returns the configuration's tools selection. A nil
// selection is an omitted one, which enables every tool; a non-nil, possibly
// empty selection names the tools to enable. known is false when no
// configuration is available or when the value is not settled yet.
func versionToolsSelection(raw cty.Value) (tools []string, known bool) {
	version, ok := rawBlock(raw, "version")
	if !ok {
		return nil, false
	}
	value, ok := rawAttr(version, "tools")
	if !ok || !value.IsWhollyKnown() {
		return nil, false
	}
	if value.IsNull() {
		return nil, true
	}
	if !value.Type().IsSetType() && !value.Type().IsListType() {
		return nil, false
	}
	out := make([]string, 0, value.LengthInt())
	iter := value.ElementIterator()
	for iter.Next() {
		_, item := iter.Element()
		if item.Type() != cty.String || item.IsNull() || !item.IsKnown() {
			return nil, false
		}
		out = append(out, item.AsString())
	}
	sort.Strings(out)
	return out, true
}
