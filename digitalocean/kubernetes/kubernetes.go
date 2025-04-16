package kubernetes

import (
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func nodePoolSchema(isResource bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"size": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"actual_node_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"node_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			DiffSuppressFunc: func(key, old, new string, d *schema.ResourceData) bool {
				nodeCountKey := "node_count"
				actualNodeCountKey := "actual_node_count"

				// Since this schema is shared between the node pool resource
				// and as the node pool sub-element of the cluster resource,
				// we need to check for both variants of the incoming key.
				keyParts := strings.Split(key, ".")
				if keyParts[0] == "node_pool" {
					npKeyParts := keyParts[:len(keyParts)-1]
					nodeCountKeyParts := append(npKeyParts, "node_count")
					nodeCountKey = strings.Join(nodeCountKeyParts, ".")
					actualNodeCountKeyParts := append(npKeyParts, "actual_node_count")
					actualNodeCountKey = strings.Join(actualNodeCountKeyParts, ".")
				}

				// If node_count equals actual_node_count already, then
				// suppress the diff.
				if d.Get(nodeCountKey).(int) == d.Get(actualNodeCountKey).(int) {
					return true
				}

				// Otherwise suppress the diff only if old equals new.
				return old == new
			},
		},

		"auto_scale": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"min_nodes": {
			Type:     schema.TypeInt,
			Optional: true,
		},

		"max_nodes": {
			Type:     schema.TypeInt,
			Optional: true,
		},

		"tags": tag.TagsSchema(),

		"labels": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nodes": nodeSchema(),

		"taint": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     nodePoolTaintSchema(),
		},
	}

	if isResource {
		// add the cluster id
		s["cluster_id"] = &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		}

		// remove the id when this is used in a specific resource
		// not as a child
		delete(s, "id")
	}

	return s
}

func nodePoolTaintSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"effect": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NoSchedule",
					"PreferNoSchedule",
					"NoExecute",
				}, false),
			},
		},
	}
}

func nodeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"droplet_id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"updated_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandLabels(labels map[string]interface{}) map[string]string {
	expandedLabels := make(map[string]string)
	for key, value := range labels {
		expandedLabels[key] = value.(string)
	}
	return expandedLabels
}

func flattenLabels(labels map[string]string) map[string]interface{} {
	flattenedLabels := make(map[string]interface{})
	for key, value := range labels {
		flattenedLabels[key] = value
	}
	return flattenedLabels
}

func expandNodePools(nodePools []interface{}) []*godo.KubernetesNodePool {
	expandedNodePools := make([]*godo.KubernetesNodePool, 0, len(nodePools))
	for _, rawPool := range nodePools {
		pool := rawPool.(map[string]interface{})
		cr := &godo.KubernetesNodePool{
			ID:        pool["id"].(string),
			Name:      pool["name"].(string),
			Size:      pool["size"].(string),
			Count:     pool["node_count"].(int),
			AutoScale: pool["auto_scale"].(bool),
			MinNodes:  pool["min_nodes"].(int),
			MaxNodes:  pool["max_nodes"].(int),
			Tags:      tag.ExpandTags(pool["tags"].(*schema.Set).List()),
			Labels:    expandLabels(pool["labels"].(map[string]interface{})),
			Nodes:     expandNodes(pool["nodes"].([]interface{})),
			Taints:    expandNodePoolTaints(pool["taint"].(*schema.Set).List()),
		}

		expandedNodePools = append(expandedNodePools, cr)
	}

	return expandedNodePools
}

func expandMaintPolicyOpts(config []interface{}) (*godo.KubernetesMaintenancePolicy, error) {
	maintPolicy := &godo.KubernetesMaintenancePolicy{}
	configMap := config[0].(map[string]interface{})

	if v, ok := configMap["day"]; ok {
		day, err := godo.KubernetesMaintenanceToDay(v.(string))
		if err != nil {
			return nil, err
		}
		maintPolicy.Day = day
	}

	if v, ok := configMap["start_time"]; ok {
		maintPolicy.StartTime = v.(string)
	}

	return maintPolicy, nil
}

func expandControlPlaneFirewallOpts(raw []interface{}) *godo.KubernetesControlPlaneFirewall {
	if len(raw) == 0 || raw[0] == nil {
		return &godo.KubernetesControlPlaneFirewall{}
	}
	controlPlaneFirewallConfig := raw[0].(map[string]interface{})

	controlPlaneFirewall := &godo.KubernetesControlPlaneFirewall{
		Enabled:          godo.PtrTo(controlPlaneFirewallConfig["enabled"].(bool)),
		AllowedAddresses: expandAllowedAddresses(controlPlaneFirewallConfig["allowed_addresses"].([]interface{})),
	}
	return controlPlaneFirewall
}

func expandAllowedAddresses(addrs []interface{}) []string {
	var expandedAddrs []string
	for _, item := range addrs {
		if str, ok := item.(string); ok {
			expandedAddrs = append(expandedAddrs, str)
		}
	}
	return expandedAddrs
}

func expandRoutingAgentOpts(raw []interface{}) *godo.KubernetesRoutingAgent {
	if len(raw) == 0 || raw[0] == nil {
		return &godo.KubernetesRoutingAgent{}
	}

	rawRoutingAgentObj := raw[0].(map[string]interface{})

	routingAgent := &godo.KubernetesRoutingAgent{
		Enabled: godo.PtrTo(rawRoutingAgentObj["enabled"].(bool)),
	}

	return routingAgent
}

func flattenRoutingAgentOpts(opts *godo.KubernetesRoutingAgent) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if opts == nil {
		return result
	}

	item := make(map[string]interface{})
	item["enabled"] = opts.Enabled

	result = append(result, item)

	return result
}

func flattenMaintPolicyOpts(opts *godo.KubernetesMaintenancePolicy) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	item := make(map[string]interface{})

	item["day"] = opts.Day.String()
	item["duration"] = opts.Duration
	item["start_time"] = opts.StartTime
	result = append(result, item)

	return result
}

func flattenControlPlaneFirewallOpts(opts *godo.KubernetesControlPlaneFirewall) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if opts == nil {
		return result
	}

	item := make(map[string]interface{})

	item["enabled"] = opts.Enabled
	item["allowed_addresses"] = opts.AllowedAddresses
	result = append(result, item)

	return result
}

func flattenNodePool(d *schema.ResourceData, keyPrefix string, pool *godo.KubernetesNodePool, parentTags ...string) []interface{} {
	rawPool := map[string]interface{}{
		"id":                pool.ID,
		"name":              pool.Name,
		"size":              pool.Size,
		"actual_node_count": pool.Count,
		"auto_scale":        pool.AutoScale,
		"min_nodes":         pool.MinNodes,
		"max_nodes":         pool.MaxNodes,
		"taint":             pool.Taints,
	}

	if pool.Tags != nil {
		rawPool["tags"] = tag.FlattenTags(FilterTags(pool.Tags))
	}

	if pool.Labels != nil {
		rawPool["labels"] = flattenLabels(pool.Labels)
	}

	if pool.Nodes != nil {
		rawPool["nodes"] = flattenNodes(pool.Nodes)
	}

	if pool.Taints != nil {
		rawPool["taint"] = flattenNodePoolTaints(pool.Taints)
	}

	// Assign a node_count only if it's been set explicitly, since it's
	// optional and we don't want to update with a 0 if it's not set.
	if _, ok := d.GetOk(keyPrefix + "node_count"); ok {
		rawPool["node_count"] = pool.Count
	}

	return []interface{}{rawPool}
}

func expandNodes(nodes []interface{}) []*godo.KubernetesNode {
	expandedNodes := make([]*godo.KubernetesNode, 0, len(nodes))
	for _, rawNode := range nodes {
		node := rawNode.(map[string]interface{})
		n := &godo.KubernetesNode{
			ID:   node["id"].(string),
			Name: node["name"].(string),
		}

		expandedNodes = append(expandedNodes, n)
	}

	return expandedNodes
}

func flattenNodes(nodes []*godo.KubernetesNode) []interface{} {
	flattenedNodes := make([]interface{}, 0)
	if nodes == nil {
		return flattenedNodes
	}

	for _, node := range nodes {
		rawNode := map[string]interface{}{
			"id":         node.ID,
			"name":       node.Name,
			"status":     node.Status.State,
			"droplet_id": node.DropletID,
			"created_at": node.CreatedAt.UTC().String(),
			"updated_at": node.UpdatedAt.UTC().String(),
		}

		flattenedNodes = append(flattenedNodes, rawNode)
	}

	return flattenedNodes
}

func expandNodePoolTaints(taints []interface{}) []godo.Taint {
	expandedTaints := make([]godo.Taint, 0, len(taints))
	for _, rawTaint := range taints {
		taint := rawTaint.(map[string]interface{})
		t := godo.Taint{
			Key:    taint["key"].(string),
			Value:  taint["value"].(string),
			Effect: taint["effect"].(string),
		}

		expandedTaints = append(expandedTaints, t)
	}

	return expandedTaints
}

func flattenNodePoolTaints(taints []godo.Taint) []interface{} {
	flattenedTaints := make([]interface{}, 0)
	if taints == nil {
		return flattenedTaints
	}

	for _, taint := range taints {
		rawTaint := map[string]interface{}{
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
		}

		flattenedTaints = append(flattenedTaints, rawTaint)
	}

	return flattenedTaints
}

// FilterTags filters tags to remove any automatically added to avoid state problems,
// these are tags starting with "k8s:" or named "k8s"
func FilterTags(tags []string) []string {
	filteredTags := make([]string, 0)
	for _, t := range tags {
		if !strings.HasPrefix(t, "k8s:") &&
			!strings.HasPrefix(t, "terraform:") &&
			t != "k8s" {
			filteredTags = append(filteredTags, t)
		}
	}

	return filteredTags
}

func expandCAConfigOpts(config []interface{}) *godo.KubernetesClusterAutoscalerConfiguration {
	caConfig := &godo.KubernetesClusterAutoscalerConfiguration{}
	configMap := config[0].(map[string]interface{})

	if v, ok := configMap["scale_down_utilization_threshold"]; ok {
		caConfig.ScaleDownUtilizationThreshold = godo.PtrTo(v.(float64))
	}

	if v, ok := configMap["scale_down_unneeded_time"]; ok {
		caConfig.ScaleDownUnneededTime = godo.PtrTo(v.(string))
	}

	return caConfig
}

func flattenCAConfigOpts(opts *godo.KubernetesClusterAutoscalerConfiguration) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if opts == nil {
		return result
	}

	item := make(map[string]interface{})
	item["scale_down_utilization_threshold"] = opts.ScaleDownUtilizationThreshold
	item["scale_down_unneeded_time"] = opts.ScaleDownUnneededTime
	result = append(result, item)

	return result
}
