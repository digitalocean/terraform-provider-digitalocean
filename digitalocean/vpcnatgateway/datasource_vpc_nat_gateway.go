package vpcnatgateway

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanVPCNATGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanVPCNATGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ID of the VPC NAT Gateway",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the VPC NAT Gateway",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the VPC NAT Gateway",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the VPC NAT Gateway",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the VPC NAT Gateway",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the VPC NAT Gateway",
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
							Description: "Gateway IP of the VPC NAT Gateway",
						},
						"default_gateway": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if this is the default VPC NAT Gateway in the VPC",
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
							Elem:        egressPublicGatewaysSchemaResource(),
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
				Description: "VPC NAT Gateway create timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC NAT Gateway update timestamp",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the project that the VPC NAT Gateway is associated with.",
			},
		},
	}
}

func dataSourceDigitalOceanVPCNATGatewayRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var foundGateway *godo.VPCNATGateway
	if id, ok := d.GetOk("id"); ok {
		gateway, _, err := client.VPCNATGateways.Get(context.Background(), id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving VPC NAT Gateway: %v", err)
		}
		foundGateway = gateway
	} else if name, ok := d.GetOk("name"); ok {
		gateways, _, err := client.VPCNATGateways.List(context.Background(), &godo.VPCNATGatewaysListOptions{
			ListOptions: godo.ListOptions{Page: 1},
			Name:        []string{name.(string)},
		})
		if err != nil || len(gateways) < 1 {
			return diag.Errorf("Error listing VPC NAT Gateway by name: %v", err)
		}
		if len(gateways) > 1 {
			return diag.Errorf("Too many VPC NAT Gateways found by name: %v (%s)", len(gateways), name.(string))
		}
		foundGateway = gateways[0]
	}

	d.SetId(foundGateway.ID)
	d.Set("name", foundGateway.Name)
	d.Set("type", foundGateway.Type)
	d.Set("state", foundGateway.State)
	d.Set("region", foundGateway.Region)
	d.Set("size", foundGateway.Size)
	d.Set("vpcs", flattenVPCs(foundGateway.VPCs))
	d.Set("egresses", flattenEgresses(foundGateway.Egresses))
	d.Set("udp_timeout_seconds", int(foundGateway.UDPTimeoutSeconds))
	d.Set("icmp_timeout_seconds", int(foundGateway.ICMPTimeoutSeconds))
	d.Set("tcp_timeout_seconds", int(foundGateway.TCPTimeoutSeconds))
	d.Set("created_at", foundGateway.CreatedAt.UTC().String())
	d.Set("updated_at", foundGateway.UpdatedAt.UTC().String())
	d.Set("project_id", foundGateway.ProjectID)

	return nil
}

func flattenVPCs(vpcs []*godo.IngressVPC) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(vpcs))
	for _, vpc := range vpcs {
		r := make(map[string]interface{})
		r["vpc_uuid"] = vpc.VpcUUID
		r["gateway_ip"] = vpc.GatewayIP
		r["default_gateway"] = vpc.DefaultGateway
		result = append(result, r)
	}
	return result
}

func flattenEgresses(egresses *godo.Egresses) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if egresses != nil {
		for _, egress := range egresses.PublicGateways {
			gatewaySet := schema.NewSet(schema.HashResource(egressPublicGatewaysSchemaResource()), []interface{}{})
			r := make(map[string]interface{})
			r["ipv4"] = egress.IPv4
			gatewaySet.Add(r)
			result = append(result, map[string]interface{}{
				"public_gateways": gatewaySet,
			})
		}
	}
	return result
}
