package digitalocean

import (
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func firewallSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"droplet_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeInt},
			Optional: true,
		},

		"inbound_rule": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     firewallRuleSchema("source"),
		},

		"outbound_rule": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     firewallRuleSchema("destination"),
		},

		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"pending_changes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"droplet_id": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"removing": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"status": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},

		"tags": tagsSchema(),
	}
}

func firewallRuleSchema(prefix string) *schema.Resource {
	if prefix != "" && prefix[len(prefix)-1:] != "_" {
		prefix += "_"
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tcp",
					"udp",
					"icmp",
				}, false),
			},
			"port_range": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			prefix + "addresses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
			prefix + "droplet_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			prefix + "load_balancer_uids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
			prefix + "kubernetes_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
			prefix + "tags": tagsSchema(),
		},
	}
}

func expandFirewallDropletIds(droplets []interface{}) []int {
	expandedDroplets := make([]int, len(droplets))
	for i, v := range droplets {
		expandedDroplets[i] = v.(int)
	}

	return expandedDroplets
}

func expandFirewallRuleStringSet(strings []interface{}) []string {
	expandedStrings := make([]string, len(strings))
	for i, v := range strings {
		expandedStrings[i] = v.(string)
	}

	return expandedStrings
}

func expandFirewallInboundRules(rules []interface{}) []godo.InboundRule {
	expandedRules := make([]godo.InboundRule, 0, len(rules))
	for _, rawRule := range rules {
		var src godo.Sources

		rule := rawRule.(map[string]interface{})

		src.DropletIDs = expandFirewallDropletIds(rule["source_droplet_ids"].(*schema.Set).List())

		src.Addresses = expandFirewallRuleStringSet(rule["source_addresses"].(*schema.Set).List())

		src.LoadBalancerUIDs = expandFirewallRuleStringSet(rule["source_load_balancer_uids"].(*schema.Set).List())

		src.KubernetesIDs = expandFirewallRuleStringSet(rule["source_kubernetes_ids"].(*schema.Set).List())

		src.Tags = expandTags(rule["source_tags"].(*schema.Set).List())

		r := godo.InboundRule{
			Protocol:  rule["protocol"].(string),
			PortRange: rule["port_range"].(string),
			Sources:   &src,
		}

		expandedRules = append(expandedRules, r)
	}
	return expandedRules
}

func expandFirewallOutboundRules(rules []interface{}) []godo.OutboundRule {
	expandedRules := make([]godo.OutboundRule, 0, len(rules))
	for _, rawRule := range rules {
		var dest godo.Destinations

		rule := rawRule.(map[string]interface{})

		dest.DropletIDs = expandFirewallDropletIds(rule["destination_droplet_ids"].(*schema.Set).List())

		dest.Addresses = expandFirewallRuleStringSet(rule["destination_addresses"].(*schema.Set).List())

		dest.LoadBalancerUIDs = expandFirewallRuleStringSet(rule["destination_load_balancer_uids"].(*schema.Set).List())

		dest.KubernetesIDs = expandFirewallRuleStringSet(rule["destination_kubernetes_ids"].(*schema.Set).List())

		dest.Tags = expandTags(rule["destination_tags"].(*schema.Set).List())

		r := godo.OutboundRule{
			Protocol:     rule["protocol"].(string),
			PortRange:    rule["port_range"].(string),
			Destinations: &dest,
		}

		expandedRules = append(expandedRules, r)
	}
	return expandedRules
}

func firewallPendingChanges(d *schema.ResourceData, firewall *godo.Firewall) []interface{} {
	remote := make([]interface{}, 0, len(firewall.PendingChanges))
	for _, change := range firewall.PendingChanges {
		rawChange := map[string]interface{}{
			"droplet_id": change.DropletID,
			"removing":   change.Removing,
			"status":     change.Status,
		}
		remote = append(remote, rawChange)
	}
	return remote
}

func flattenFirewallDropletIds(droplets []int) *schema.Set {
	if droplets == nil {
		return nil
	}

	flattenedDroplets := schema.NewSet(schema.HashInt, []interface{}{})
	for _, v := range droplets {
		flattenedDroplets.Add(v)
	}

	return flattenedDroplets
}

func flattenFirewallRuleStringSet(strings []string) *schema.Set {
	flattenedStrings := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range strings {
		flattenedStrings.Add(v)
	}

	return flattenedStrings
}

func flattenFirewallInboundRules(rules []godo.InboundRule) []interface{} {
	if rules == nil {
		return nil
	}

	flattenedRules := make([]interface{}, len(rules))
	for i, rule := range rules {
		sources := rule.Sources
		protocol := rule.Protocol
		portRange := rule.PortRange

		rawRule := map[string]interface{}{
			"protocol": protocol,
		}

		// The API returns 0 when the port range was specified as all.
		// If protocol is `icmp` the API returns 0 for when port was
		// not specified.
		if portRange == "0" {
			if protocol != "icmp" {
				rawRule["port_range"] = "all"
			}
		} else {
			rawRule["port_range"] = portRange
		}

		if sources.Tags != nil {
			rawRule["source_tags"] = flattenTags(sources.Tags)
		}

		if sources.DropletIDs != nil {
			rawRule["source_droplet_ids"] = flattenFirewallDropletIds(sources.DropletIDs)
		}

		if sources.Addresses != nil {
			rawRule["source_addresses"] = flattenFirewallRuleStringSet(sources.Addresses)
		}

		if sources.LoadBalancerUIDs != nil {
			rawRule["source_load_balancer_uids"] = flattenFirewallRuleStringSet(sources.LoadBalancerUIDs)
		}

		if sources.KubernetesIDs != nil {
			rawRule["source_kubernetes_ids"] = flattenFirewallRuleStringSet(sources.KubernetesIDs)
		}

		flattenedRules[i] = rawRule
	}

	return flattenedRules
}

func flattenFirewallOutboundRules(rules []godo.OutboundRule) []interface{} {
	if rules == nil {
		return nil
	}

	flattenedRules := make([]interface{}, len(rules))
	for i, rule := range rules {
		destinations := rule.Destinations
		protocol := rule.Protocol
		portRange := rule.PortRange

		rawRule := map[string]interface{}{
			"protocol": protocol,
		}

		// The API returns 0 when the port range was specified as all.
		// If protocol is `icmp` the API returns 0 for when port was
		// not specified.
		if portRange == "0" {
			if protocol != "icmp" {
				rawRule["port_range"] = "all"
			}
		} else {
			rawRule["port_range"] = portRange
		}

		if destinations.Tags != nil {
			rawRule["destination_tags"] = flattenTags(destinations.Tags)
		}

		if destinations.DropletIDs != nil {
			rawRule["destination_droplet_ids"] = flattenFirewallDropletIds(destinations.DropletIDs)
		}

		if destinations.Addresses != nil {
			rawRule["destination_addresses"] = flattenFirewallRuleStringSet(destinations.Addresses)
		}

		if destinations.LoadBalancerUIDs != nil {
			rawRule["destination_load_balancer_uids"] = flattenFirewallRuleStringSet(destinations.LoadBalancerUIDs)
		}

		if destinations.KubernetesIDs != nil {
			rawRule["destination_kubernetes_ids"] = flattenFirewallRuleStringSet(destinations.KubernetesIDs)
		}

		flattenedRules[i] = rawRule
	}

	return flattenedRules
}
