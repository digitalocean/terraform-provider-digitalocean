package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func firewallSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "the firewall identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "name of the firewall",
		},
		"status": {
			Type:        schema.TypeString,
			Description: "indicates the current state of the firewall (e.g `waiting`, `succeded`, `failed`)",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "ISO8601 date and time representing when the firewall was created",
		},
		"inbound_rules": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     firewallInboundRuleSchema(),
			Description: "List of inbound access rules specifying the protocol, ports, and sources for inbound traffic allowed through the firewall",
		},
		"outbound_rules": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     firewallOutboundRuleSchema(),
			Description: "List of inbound access rules specifying the protocol, ports, and sources for inbound traffic allowed through the firewall",
		},
		"droplet_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     schema.TypeInt,
		},
		"tags": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     schema.TypeString,
		},
	}
}

func firewallInboundRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:        schema.TypeString,
				Description: "The type of traffic to be allowed. This may be one of `tcp`, `udp`, or `icmp`",
			},
			"ports": {
				Type:        schema.TypeString,
				Description: "The ports on which traffic will be allowed specified as a string containing a single port, a range (e.g. `8000-9000`), or `0` when all ports are open for a protocol. For ICMP rules this parameter will always return `0`",
			},
			"sources": {
				Type:        schema.TypeString,
				Description: "An object specifying locations from which inbound traffic will be accepted",
			},
		},
	}
}

func firewallOutboundRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:        schema.TypeString,
				Description: "The type of traffic to be allowed. This may be one of `tcp`, `udp`, or `icmp`",
			},
			"ports": {
				Type:        schema.TypeString,
				Description: "The ports on which traffic will be allowed specified as a string containing a single port, a range (e.g. `8000-9000`), or `0` when all ports are open for a protocol. For ICMP rules this parameter will always return `0`",
			},
			"sources": {
				Type:        schema.TypeString,
				Description: "An object specifying locations from which inbound traffic will be accepted",
			},
		},
	}
}

func firewallRuleTargetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"addresses": {
				Type:        schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeString,
				Description: "List of the IPv4 addresses, IPv6 addresses, IPv4 CIDRs, and/or IPv6 CIDRs to which the firewall will allow traffic",
			},
			"droplet_ids": {
				Type:        schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeInt,
				Description: "List of the IDs of the Droplets to which the firewall will allow traffic",
			},
			"load_balancer_ids": {
				Type:        schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeString,
				Description: "List of the IDs (UUID) of the load balancers to which the firewall will allow traffic",
			},
			"tags": {
				Type: schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeString,
				Description: "List of Tags corresponding to groups of Droplets to which the firewall will allow traffic",
			},
		},
	}
}

func getDigitalOceanFirewalls(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*CombinedConfig).godoClient()

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allFirewalls []interface{}

	for {

		firewalls, resp, err := client.Firewalls.List(context.Background(), opts)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving firewaqlls: %s", err)
		}

		for _, firewall := range firewalls {
			allFirewalls = append(allFirewalls, firewall)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving firewalls: %s", err)
		}

		opts.Page = page + 1
	}

	return allFirewalls, nil
}