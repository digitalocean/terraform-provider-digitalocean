package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDigitalOceanFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanFirewallCreate,
		ReadContext:   resourceDigitalOceanFirewallRead,
		UpdateContext: resourceDigitalOceanFirewallUpdate,
		DeleteContext: resourceDigitalOceanFirewallDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: firewallSchema(),

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {

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

func resourceDigitalOceanFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	opts, err := firewallRequest(d, client)
	if err != nil {
		return diag.Errorf("Error in firewall request: %s", err)
	}

	log.Printf("[DEBUG] Firewall create configuration: %#v", opts)

	firewall, _, err := client.Firewalls.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating firewall: %s", err)
	}

	// Assign the firewall id
	d.SetId(firewall.ID)

	log.Printf("[INFO] Firewall ID: %s", d.Id())

	return resourceDigitalOceanFirewallRead(ctx, d, meta)
}

func resourceDigitalOceanFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	// Retrieve the firewall properties for updating the state
	firewall, resp, err := client.Firewalls.Get(context.Background(), d.Id())
	if err != nil {
		// check if the firewall no longer exists.
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Firewall (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving firewall: %s", err)
	}

	d.Set("status", firewall.Status)
	d.Set("created_at", firewall.Created)
	d.Set("pending_changes", firewallPendingChanges(d, firewall))
	d.Set("name", firewall.Name)

	if err := d.Set("droplet_ids", flattenFirewallDropletIds(firewall.DropletIDs)); err != nil {
		return diag.Errorf("[DEBUG] Error setting `droplet_ids`: %+v", err)
	}

	if err := d.Set("inbound_rule", flattenFirewallInboundRules(firewall.InboundRules)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Firewall inbound_rule error: %#v", err)
	}

	if err := d.Set("outbound_rule", flattenFirewallOutboundRules(firewall.OutboundRules)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Firewall outbound_rule error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(firewall.Tags)); err != nil {
		return diag.Errorf("[DEBUG] Error setting `tags`: %+v", err)
	}

	return nil
}

func resourceDigitalOceanFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	opts, err := firewallRequest(d, client)
	if err != nil {
		return diag.Errorf("Error in firewall request: %s", err)
	}

	log.Printf("[DEBUG] Firewall update configuration: %#v", opts)

	_, _, err = client.Firewalls.Update(context.Background(), d.Id(), opts)
	if err != nil {
		return diag.Errorf("Error updating firewall: %s", err)
	}

	return resourceDigitalOceanFirewallRead(ctx, d, meta)
}

func resourceDigitalOceanFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting firewall: %s", d.Id())

	// Destroy the droplet
	_, err := client.Firewalls.Delete(context.Background(), d.Id())

	// Handle remotely destroyed droplets
	if err != nil && strings.Contains(err.Error(), "404 Not Found") {
		return nil
	}

	if err != nil {
		return diag.Errorf("Error deleting firewall: %s", err)
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
