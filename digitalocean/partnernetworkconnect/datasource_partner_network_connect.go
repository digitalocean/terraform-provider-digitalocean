package partnernetworkconnect

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanPartnerNetworkConnect() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanPartnerNetworkConnectRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The ID of the Partner Attachment",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the Partner Attachment",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"connection_bandwidth_in_mbps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"naas_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bgp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_router_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"peer_router_asn": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"peer_router_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanPartnerNetworkConnectRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	var foundPartnerNetworkConnect *godo.PartnerNetworkConnect

	if id, ok := d.GetOk("id"); ok {
		partnerNetworkConnect, _, err := client.PartnerNetworkConnect.Get(ctx, id.(string))
		if err != nil {
			return diag.Errorf("error retrieving Partner Network Connect: %s", err)
		}

		foundPartnerNetworkConnect = partnerNetworkConnect
	} else if name, ok := d.GetOk("name"); ok {
		partnerNetworkConnects, err := listPartnerNetworkConnect(client)
		if err != nil {
			return diag.Errorf("error retrieving Partner Network Connect: %s", err)
		}

		partnerNetworkConnect, err := findPartnerNetworkConnectByName(partnerNetworkConnects, name.(string))
		if err != nil {
			return diag.Errorf("error retrieving Partner Network Connect: %s", err)
		}

		foundPartnerNetworkConnect = partnerNetworkConnect
	}

	if foundPartnerNetworkConnect == nil {
		return diag.Errorf("Bad Request: %s", fmt.Errorf("'name' or 'id' must be provided"))
	}

	d.SetId(foundPartnerNetworkConnect.ID)
	d.Set("name", foundPartnerNetworkConnect.Name)
	d.Set("connection_bandwidth_in_mbps", foundPartnerNetworkConnect.ConnectionBandwidthInMbps)
	d.Set("region", strings.ToLower(foundPartnerNetworkConnect.Region))
	d.Set("naas_provider", foundPartnerNetworkConnect.NaaSProvider)
	d.Set("vpc_ids", foundPartnerNetworkConnect.VPCIDs)
	if bgp := foundPartnerNetworkConnect.BGP; bgp.PeerRouterIP != "" || bgp.LocalRouterIP != "" || bgp.PeerASN != 0 {
		bgpMap := map[string]interface{}{
			"local_router_ip": bgp.LocalRouterIP,
			"peer_router_asn": bgp.PeerASN,
			"peer_router_ip":  bgp.PeerRouterIP,
		}
		if err := d.Set("bgp", []interface{}{bgpMap}); err != nil {
			return diag.FromErr(err)
		}
	}
	d.Set("state", foundPartnerNetworkConnect.State)
	d.Set("created_at", foundPartnerNetworkConnect.CreatedAt.UTC().String())

	return nil
}

func listPartnerNetworkConnect(client *godo.Client) ([]*godo.PartnerNetworkConnect, error) {
	pncList := []*godo.PartnerNetworkConnect{}
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		partnerNetworkConnects, resp, err := client.PartnerNetworkConnect.List(context.Background(), opts)
		if err != nil {
			return nil, err
		}

		pncList = append(pncList, partnerNetworkConnects...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return pncList, fmt.Errorf("error retrieving Partner Network Connects: %s", err)
		}

		opts.Page = page + 1
	}

	return pncList, nil
}

func findPartnerNetworkConnectByName(partnerNetworkConnects []*godo.PartnerNetworkConnect, name string) (*godo.PartnerNetworkConnect, error) {
	for _, partnerNetworkConnect := range partnerNetworkConnects {
		if partnerNetworkConnect.Name == name {
			return partnerNetworkConnect, nil
		}
	}

	return nil, fmt.Errorf("no Partner Network Connect found with name: %s", name)
}
