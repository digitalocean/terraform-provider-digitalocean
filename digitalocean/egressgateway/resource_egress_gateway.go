package egressgateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanEgressGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanEgressGatewayCreate,
		ReadContext:   resourceDigitalOceanEgressGatewayRead,
		UpdateContext: resourceDigitalOceanEgressGatewayUpdate,
		DeleteContext: resourceDigitalOceanEgressGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Egress Gateway",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Egress Gateway",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the Egress Gateway",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the Egress Gateway",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of the Egress Gateway",
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
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP of the Egress Gateway",
						},
						"default_egress_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "UDP connection timeout (in seconds)",
				ValidateFunc: validation.All(validation.NoZeroValues),
			},
			"icmp_timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ICMP connection timeout (in seconds)",
				ValidateFunc: validation.All(validation.NoZeroValues),
			},
			"tcp_timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "TCP connection timeout (in seconds)",
				ValidateFunc: validation.All(validation.NoZeroValues),
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

func resourceDigitalOceanEgressGatewayCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	createReq := &godo.EgressGatewayRequest{
		Name:   d.Get("name").(string),
		Type:   d.Get("type").(string),
		Region: d.Get("region").(string),
		VPCs:   expandVPCs(d.Get("vpcs").([]interface{})),
	}
	if v, ok := d.GetOk("udp_timeout_seconds"); ok {
		createReq.UDPTimeoutSeconds = uint32(v.(int))
	}
	if v, ok := d.GetOk("icmp_timeout_seconds"); ok {
		createReq.UDPTimeoutSeconds = uint32(v.(int))
	}
	if v, ok := d.GetOk("tcp_timeout_seconds"); ok {
		createReq.UDPTimeoutSeconds = uint32(v.(int))
	}
	gateway, _, err := client.EgressGateways.Create(context.Background(), createReq)
	if err != nil {
		return diag.Errorf("Error creating Egress Gateway: %v", err)
	}
	d.SetId(gateway.ID)

	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"STATE_PROVISIONING"},
		Refresh:    egressGatewayRefreshFunc(client, d.Id()),
		Target:     []string{"STATE_ACTIVE"},
		Timeout:    15 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for Egress Gateway (%s) to become active: %v", gateway.Name, err)
	}

	return resourceDigitalOceanEgressGatewayRead(ctx, d, meta)
}

func resourceDigitalOceanEgressGatewayRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	gateway, resp, err := client.EgressGateways.Get(context.Background(), d.Id())
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Egress Gateway: %v", err)
	}

	d.Set("name", gateway.Name)
	d.Set("type", gateway.Type)
	d.Set("state", gateway.State)
	d.Set("region", gateway.Region)
	d.Set("vpcs", flattenVPCs(gateway.VPCs))
	d.Set("egresses", flattenEgresses(gateway.Egresses))
	d.Set("udp_timeout_seconds", int(gateway.UDPTimeoutSeconds))
	d.Set("icmp_timeout_seconds", int(gateway.ICMPTimeoutSeconds))
	d.Set("tcp_timeout_seconds", int(gateway.TCPTimeoutSeconds))
	d.Set("created_at", gateway.CreatedAt.UTC().String())
	d.Set("updated_at", gateway.UpdatedAt.UTC().String())

	return nil
}

func resourceDigitalOceanEgressGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	updateReq := &godo.EgressGatewayRequest{Name: d.Get("name").(string)}
	if v, ok := d.GetOk("udp_timeout_seconds"); ok {
		updateReq.UDPTimeoutSeconds = uint32(v.(int))
	}
	if v, ok := d.GetOk("icmp_timeout_seconds"); ok {
		updateReq.UDPTimeoutSeconds = uint32(v.(int))
	}
	if v, ok := d.GetOk("tcp_timeout_seconds"); ok {
		updateReq.UDPTimeoutSeconds = uint32(v.(int))
	}
	_, _, err := client.EgressGateways.Update(context.Background(), d.Id(), updateReq)
	if err != nil {
		return diag.Errorf("Error updating Egress Gateway: %v", err)
	}

	return resourceDigitalOceanEgressGatewayRead(ctx, d, meta)
}

func resourceDigitalOceanEgressGatewayDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, err := client.EgressGateways.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting Egress Gateway: %v", err)
	}

	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"STATE_PROVISIONING", "STATE_NEW"},
		Refresh:    egressGatewayRefreshFunc(client, d.Id()),
		Target:     []string{"STATE_ACTIVE"},
		Timeout:    15 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for Egress Gateway (%s) to become active: %v", d.Get("name"), err)
	}

	d.SetId("")
	return nil
}

func egressGatewayRefreshFunc(client *godo.Client, gatewayID string) retry.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		gateway, resp, err := client.EgressGateways.Get(context.Background(), gatewayID)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return gateway, http.StatusText(resp.StatusCode), nil
			}
			return nil, "", fmt.Errorf("Error retrieving Egress Gateway: %v", err)
		}
		if gateway.State != "STATE_ACTIVE" {
			return gateway, gateway.State, nil
		}
		return gateway, "STATE_ACTIVE", nil
	}
}

func expandVPCs(vpcs []interface{}) []*godo.IngressVPC {
	ingressVPCs := make([]*godo.IngressVPC, 0, len(vpcs))
	for i := range vpcs {
		vpc := vpcs[i].(map[string]interface{})
		ingressVPCs = append(ingressVPCs, &godo.IngressVPC{
			VpcUUID:              vpc["vpc_uuid"].(string),
			GatewayIP:            vpc["gateway_ip"].(string),
			DefaultEgressGateway: vpc["default_egress_gateway"].(bool),
		})
	}
	return ingressVPCs
}
