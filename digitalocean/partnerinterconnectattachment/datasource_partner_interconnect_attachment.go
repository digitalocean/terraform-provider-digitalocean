package partnerinterconnectattachment

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

func DataSourceDigitalOceanPartnerInterconnectAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanPartnerInterconnectAttachmentRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The ID of the Partner Interconnect Attachment",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the Partner Interconnect Attachment",
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

func dataSourceDigitalOceanPartnerInterconnectAttachmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	var foundPartnerInterconnectAttachment *godo.PartnerInterconnectAttachment

	if id, ok := d.GetOk("id"); ok {
		partnerInterconnectAttachment, _, err := client.PartnerInterconnectAttachments.Get(ctx, id.(string))
		if err != nil {
			return diag.Errorf("error retrieving Partner Interconnect Attachment: %s", err)
		}

		foundPartnerInterconnectAttachment = partnerInterconnectAttachment
	} else if name, ok := d.GetOk("name"); ok {
		partnerInterconnectAttachments, err := listPartnerInterconnectAttachments(client)
		if err != nil {
			return diag.Errorf("error retrieving Partner Interconnect Attachment: %s", err)
		}

		partnerInterconnectAttachment, err := findPartnerInterconnectAttachmentByName(partnerInterconnectAttachments, name.(string))
		if err != nil {
			return diag.Errorf("error retrieving Partner Interconnect Attachment: %s", err)
		}

		foundPartnerInterconnectAttachment = partnerInterconnectAttachment
	}

	if foundPartnerInterconnectAttachment == nil {
		return diag.Errorf("Bad Request: %s", fmt.Errorf("'name' or 'id' must be provided"))
	}

	d.SetId(foundPartnerInterconnectAttachment.ID)
	d.Set("name", foundPartnerInterconnectAttachment.Name)
	d.Set("connection_bandwidth_in_mbps", foundPartnerInterconnectAttachment.ConnectionBandwidthInMbps)
	d.Set("region", strings.ToLower(foundPartnerInterconnectAttachment.Region))
	d.Set("naas_provider", foundPartnerInterconnectAttachment.NaaSProvider)
	d.Set("vpc_ids", foundPartnerInterconnectAttachment.VPCIDs)
	if bgp := foundPartnerInterconnectAttachment.BGP; bgp.PeerRouterIP != "" || bgp.LocalRouterIP != "" || bgp.PeerASN != 0 {
		bgpMap := map[string]interface{}{
			"local_router_ip": bgp.LocalRouterIP,
			"peer_router_asn": bgp.PeerASN,
			"peer_router_ip":  bgp.PeerRouterIP,
		}
		if err := d.Set("bgp", []interface{}{bgpMap}); err != nil {
			return diag.FromErr(err)
		}
	}
	d.Set("state", foundPartnerInterconnectAttachment.State)
	d.Set("created_at", foundPartnerInterconnectAttachment.CreatedAt.UTC().String())

	return nil
}

func listPartnerInterconnectAttachments(client *godo.Client) ([]*godo.PartnerInterconnectAttachment, error) {
	piaList := []*godo.PartnerInterconnectAttachment{}
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		partnerInterconnectAttachments, resp, err := client.PartnerInterconnectAttachments.List(context.Background(), opts)
		if err != nil {
			return nil, err
		}

		piaList = append(piaList, partnerInterconnectAttachments...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return piaList, fmt.Errorf("error retrieving Partner Interconnect Attachments: %s", err)
		}

		opts.Page = page + 1
	}

	return piaList, nil
}

func findPartnerInterconnectAttachmentByName(partnerInterconnectAttachments []*godo.PartnerInterconnectAttachment, name string) (*godo.PartnerInterconnectAttachment, error) {
	for _, partnerInterconnectAttachment := range partnerInterconnectAttachments {
		if partnerInterconnectAttachment.Name == name {
			return partnerInterconnectAttachment, nil
		}
	}

	return nil, fmt.Errorf("no Partner Interconnect Attachment found with name: %s", name)
}
