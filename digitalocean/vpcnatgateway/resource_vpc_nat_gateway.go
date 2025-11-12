package vpcnatgateway

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanVPCNATGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanVPCNATGatewayCreate,
		ReadContext:   resourceDigitalOceanVPCNATGatewayRead,
		UpdateContext: resourceDigitalOceanVPCNATGatewayUpdate,
		DeleteContext: resourceDigitalOceanVPCNATGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the VPC NAT Gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the VPC NAT Gateway",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the VPC NAT Gateway",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the VPC NAT Gateway",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of the VPC NAT Gateway",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Size of the VPC NAT Gateway",
			},
			"vpcs": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of ingress VPCs",
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_uuid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the ingress VPC",
							ForceNew:    true,
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP of the VPC NAT Gateway",
						},
						"default_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "UDP connection timeout (in seconds)",
				ValidateFunc: validation.All(validation.NoZeroValues),
			},
			"icmp_timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "ICMP connection timeout (in seconds)",
				ValidateFunc: validation.All(validation.NoZeroValues),
			},
			"tcp_timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "TCP connection timeout (in seconds)",
				ValidateFunc: validation.All(validation.NoZeroValues),
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
				Optional:    true,
				Computed:    true,
				Description: "ID of the project to which the VPC NAT Gateway will be assigned.",
			},
		},
	}
}

func egressPublicGatewaysSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv4 address",
			},
		},
	}
}

func resourceDigitalOceanVPCNATGatewayCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	createReq := &godo.VPCNATGatewayRequest{
		Name:   d.Get("name").(string),
		Type:   d.Get("type").(string),
		Region: d.Get("region").(string),
		Size:   uint32(d.Get("size").(int)),
		VPCs:   expandVPCs(d.Get("vpcs").(*schema.Set).List()),
	}
	if v, ok := d.GetOk("udp_timeout_seconds"); ok {
		createReq.UDPTimeoutSeconds = uint32(v.(int))
	}
	if v, ok := d.GetOk("icmp_timeout_seconds"); ok {
		createReq.ICMPTimeoutSeconds = uint32(v.(int))
	}
	if v, ok := d.GetOk("tcp_timeout_seconds"); ok {
		createReq.TCPTimeoutSeconds = uint32(v.(int))
	}
	gateway, _, err := client.VPCNATGateways.Create(context.Background(), createReq)
	if err != nil {
		return diag.Errorf("Error creating VPC NAT Gateway: %v", err)
	}
	d.SetId(gateway.ID)

	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"NEW"},
		Refresh:    vpcNatGatewayRefreshFunc(client, d.Id()),
		Target:     []string{"ACTIVE"},
		Timeout:    15 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for VPC NAT Gateway (%s) to become active: %v", gateway.Name, err)
	}

	return resourceDigitalOceanVPCNATGatewayRead(ctx, d, meta)
}

func resourceDigitalOceanVPCNATGatewayRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	gateway, _, err := client.VPCNATGateways.Get(context.Background(), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("nat gateway with id %s not found", d.Id())) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving VPC NAT Gateway: %v", err)
	}

	d.Set("name", gateway.Name)
	d.Set("type", gateway.Type)
	d.Set("state", gateway.State)
	d.Set("region", gateway.Region)
	d.Set("size", gateway.Size)
	d.Set("vpcs", flattenVPCs(gateway.VPCs))
	d.Set("egresses", flattenEgresses(gateway.Egresses))
	d.Set("udp_timeout_seconds", int(gateway.UDPTimeoutSeconds))
	d.Set("icmp_timeout_seconds", int(gateway.ICMPTimeoutSeconds))
	d.Set("tcp_timeout_seconds", int(gateway.TCPTimeoutSeconds))
	d.Set("created_at", gateway.CreatedAt.UTC().String())
	d.Set("updated_at", gateway.UpdatedAt.UTC().String())
	d.Set("project_id", gateway.ProjectID)

	return nil
}

func resourceDigitalOceanVPCNATGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	updateReq := &godo.VPCNATGatewayRequest{Name: d.Get("name").(string)}
	if v, ok := d.GetOk("type"); ok {
		updateReq.Type = v.(string)
	}
	if v, ok := d.GetOk("region"); ok {
		updateReq.Region = v.(string)
	}
	if v, ok := d.GetOk("size"); ok {
		updateReq.Size = uint32(v.(int))
	}
	if v, ok := d.GetOk("udp_timeout_seconds"); ok {
		updateReq.UDPTimeoutSeconds = uint32(v.(int))
	}
	if v, ok := d.GetOk("icmp_timeout_seconds"); ok {
		updateReq.ICMPTimeoutSeconds = uint32(v.(int))
	}
	if v, ok := d.GetOk("tcp_timeout_seconds"); ok {
		updateReq.TCPTimeoutSeconds = uint32(v.(int))
	}
	_, _, err := client.VPCNATGateways.Update(context.Background(), d.Id(), updateReq)
	if err != nil {
		return diag.Errorf("Error updating VPC NAT Gateway: %v", err)
	}

	return resourceDigitalOceanVPCNATGatewayRead(ctx, d, meta)
}

func resourceDigitalOceanVPCNATGatewayDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, err := client.VPCNATGateways.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting VPC NAT Gateway: %v", err)
	}

	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"DELETING"},
		Refresh:    vpcNatGatewayRefreshFunc(client, d.Id()),
		Target:     []string{http.StatusText(http.StatusNotFound)},
		Timeout:    15 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for VPC NAT Gateway (%s) to be deleted: %v", d.Get("name"), err)
	}

	d.SetId("")
	return nil
}

func vpcNatGatewayRefreshFunc(client *godo.Client, gatewayID string) retry.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		gateway, _, err := client.VPCNATGateways.Get(context.Background(), gatewayID)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("nat gateway with id %s not found", gatewayID)) {
				return gateway, http.StatusText(http.StatusNotFound), nil
			}
			return nil, "", fmt.Errorf("Error retrieving VPC NAT Gateway: %v", err)
		}
		if gateway.State != "ACTIVE" {
			return gateway, gateway.State, nil
		}
		return gateway, "ACTIVE", nil
	}
}

func expandVPCs(vpcs []interface{}) []*godo.IngressVPC {
	ingressVPCs := make([]*godo.IngressVPC, 0, len(vpcs))
	for i := range vpcs {
		vpc := vpcs[i].(map[string]interface{})
		ingressVPCs = append(ingressVPCs, &godo.IngressVPC{
			VpcUUID:        vpc["vpc_uuid"].(string),
			GatewayIP:      vpc["gateway_ip"].(string),
			DefaultGateway: vpc["default_gateway"].(bool),
		})
	}
	return ingressVPCs
}
