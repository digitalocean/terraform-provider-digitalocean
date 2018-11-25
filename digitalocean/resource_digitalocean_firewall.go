package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanFirewallCreate,
		Read:   resourceDigitalOceanFirewallRead,
		Update: resourceDigitalOceanFirewallUpdate,
		Delete: resourceDigitalOceanFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
				Elem: &schema.Resource{
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
						"source_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
							Optional: true,
						},
						"source_droplet_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Optional: true,
						},
						"source_load_balancer_uids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
							Optional: true,
						},
						"source_tags": tagsSchema(),
					},
				},
			},

			"outbound_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
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
						"destination_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
							Optional: true,
						},
						"destination_droplet_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Optional: true,
						},
						"destination_load_balancer_uids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
							Optional: true,
						},
						"destination_tags": tagsSchema(),
					},
				},
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
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			inboundRules, hasInbound := diff.GetOk("inbound_rule")
			outboundRules, hasOutbound := diff.GetOk("outbound_rule")

			if !hasInbound && !hasOutbound {
				return fmt.Errorf("At least one rule must be specified")
			}

			for _, v := range inboundRules.(*schema.Set).List() {
				inbound := v.(map[string]interface{})
				protocol := inbound["protocol"]

				port := inbound["port_range"]
				if protocol != "icmp" && port == "" {
					return fmt.Errorf("`port_range` of inbound rules is required if protocol is `tcp` or `udp`")
				}
			}

			for _, v := range outboundRules.(*schema.Set).List() {
				inbound := v.(map[string]interface{})
				protocol := inbound["protocol"]

				port := inbound["port_range"]
				if protocol != "icmp" && port == "" {
					return fmt.Errorf("`port_range` of outbound rules is required if protocol is `tcp` or `udp`")
				}
			}

			return nil
		},
	}
}

func resourceDigitalOceanFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts, err := firewallRequest(d, client)
	if err != nil {
		return fmt.Errorf("Error in firewall request: %s", err)
	}

	log.Printf("[DEBUG] Firewall create configuration: %#v", opts)

	firewall, _, err := client.Firewalls.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating firewall: %s", err)
	}

	// Assign the firewall id
	d.SetId(firewall.ID)

	log.Printf("[INFO] Firewall ID: %s", d.Id())

	return resourceDigitalOceanFirewallRead(d, meta)
}

func resourceDigitalOceanFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	// Retrieve the firewall properties for updating the state
	firewall, resp, err := client.Firewalls.Get(context.Background(), d.Id())
	if err != nil {
		// check if the firewall no longer exists.
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Firewall (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving firewall: %s", err)
	}

	d.Set("status", firewall.Status)
	d.Set("create_at", firewall.Created)
	d.Set("pending_changes", firewallPendingChanges(d, firewall))
	d.Set("name", firewall.Name)

	if err := d.Set("droplet_ids", flattenFirewallDropletIds(firewall.DropletIDs)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `droplet_ids`: %+v", err)
	}

	if err := d.Set("inbound_rule", flattenFirewallInboundRules(firewall.InboundRules)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Firewall inbound_rule error: %#v", err)
	}

	if err := d.Set("outbound_rule", flattenFirewallOutboundRules(firewall.OutboundRules)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Firewall outbound_rule error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(firewall.Tags)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `tags`: %+v", err)
	}

	return nil
}

func resourceDigitalOceanFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts, err := firewallRequest(d, client)
	if err != nil {
		return fmt.Errorf("Error in firewall request: %s", err)
	}

	log.Printf("[DEBUG] Firewall update configuration: %#v", opts)

	_, _, err = client.Firewalls.Update(context.Background(), d.Id(), opts)
	if err != nil {
		return fmt.Errorf("Error updating firewall: %s", err)
	}

	return resourceDigitalOceanFirewallRead(d, meta)
}

func resourceDigitalOceanFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	log.Printf("[INFO] Deleting firewall: %s", d.Id())

	// Destroy the droplet
	_, err := client.Firewalls.Delete(context.Background(), d.Id())

	// Handle remotely destroyed droplets
	if err != nil && strings.Contains(err.Error(), "404 Not Found") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error deleting firewall: %s", err)
	}

	return nil
}

func firewallRequest(d *schema.ResourceData, client *godo.Client) (*godo.FirewallRequest, error) {
	// Build up our firewall request
	opts := &godo.FirewallRequest{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("droplet_ids"); ok {
		opts.DropletIDs = expandFirewallDropletIds(v.(*schema.Set).List())
	}

	// Get inbound_rules
	if v, ok := d.GetOk("inbound_rule"); ok {
		opts.InboundRules = expandFirewallInboundRules(v.(*schema.Set).List())
	}

	// Get outbound_rules
	if v, ok := d.GetOk("outbound_rule"); ok {
		opts.OutboundRules = expandFirewallOutboundRules(v.(*schema.Set).List())
	}

	// Get tags
	opts.Tags = expandTags(d.Get("tags").(*schema.Set).List())

	return opts, nil
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

		flattenedRules[i] = rawRule
	}

	return flattenedRules
}
