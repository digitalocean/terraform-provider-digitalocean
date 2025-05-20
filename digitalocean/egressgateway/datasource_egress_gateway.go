package egressgateway

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanEgressGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanEgressGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ID of the Egress Gateway",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the Egress Gateway",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the Egress Gateway",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the Egress Gateway",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the Egress Gateway",
			},
			"vpcs": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of ingress VPCs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the ingress VPC",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP of the Egress Gateway",
						},
						"default_egress_gateway": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if this is the default Egress Gateway in the VPC",
						},
					},
				},
			},
			"egresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of egresses",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_gateways": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of public gateway IPs",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPv4 address",
									},
								},
							},
						},
					},
				},
			},
			"udp_timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "UDP connection timeout (in seconds)",
			},
			"icmp_timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ICMP connection timeout (in seconds)",
			},
			"tcp_timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "TCP connection timeout (in seconds)",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Egress Gateway create timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Egress Gateway update timestamp",
			},
		},
	}
}

func dataSourceDigitalOceanEgressGatewayRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var foundEgressGateway *godo.EgressGateway
	if id, ok := d.GetOk("id"); ok {
		gateway, _, err := client.EgressGateways.Get(context.Background(), id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving Egress Gateway: %v", err)
		}
		foundEgressGateway = gateway
	} else if name, ok := d.GetOk("name"); ok {
		gateways, _, err := client.EgressGateways.List(context.Background(), &godo.EgressGatewaysListOptions{
			ListOptions: godo.ListOptions{Page: 1},
			Name:        []string{name.(string)},
		})
		if err != nil || len(gateways) < 1 {
			return diag.Errorf("Error listing Egress Gateway by name: %v", err)
		}
		if len(gateways) > 1 {
			return diag.Errorf("Too many Egress Gateways found by name: %v (%s)", len(gateways), name.(string))
		}
		foundEgressGateway = gateways[0]
	}

	d.SetId(foundEgressGateway.ID)
	d.Set("name", foundEgressGateway.Name)
	d.Set("type", foundEgressGateway.Type)
	d.Set("state", foundEgressGateway.State)
	d.Set("region", foundEgressGateway.Region)
	d.Set("vpcs", flattenVPCs(foundEgressGateway.VPCs))
	d.Set("egresses", flattenEgresses(foundEgressGateway.Egresses))
	d.Set("udp_timeout_seconds", int(foundEgressGateway.UDPTimeoutSeconds))
	d.Set("icmp_timeout_seconds", int(foundEgressGateway.ICMPTimeoutSeconds))
	d.Set("tcp_timeout_seconds", int(foundEgressGateway.TCPTimeoutSeconds))
	d.Set("created_at", foundEgressGateway.CreatedAt.UTC().String())
	d.Set("updated_at", foundEgressGateway.UpdatedAt.UTC().String())

	return nil
}

func flattenVPCs(vpcs []*godo.IngressVPC) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(vpcs))
	for _, vpc := range vpcs {
		r := make(map[string]interface{})
		r["vpc_uuid"] = vpc.VpcUUID
		r["gateway_ip"] = vpc.GatewayIP
		r["default_egress_gateway"] = vpc.DefaultEgressGateway
		result = append(result, r)
	}
	return result
}

func flattenEgresses(egresses *godo.Egresses) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if egresses != nil {
		for _, egress := range egresses.PublicGateways {
			r := make(map[string]interface{})
			r["public_gateways"] = map[string]interface{}{
				"ipv4": egress.IPv4,
			}
			result = append(result, r)
		}
	}
	return result
}
