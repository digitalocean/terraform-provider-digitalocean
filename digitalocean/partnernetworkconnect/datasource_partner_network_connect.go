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

func DataSourceDigitalOceanPartnerAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanPartnerAttachmentRead,
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

func dataSourceDigitalOceanPartnerAttachmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	var foundPartnerAttachment *godo.PartnerAttachment

	if id, ok := d.GetOk("id"); ok {
		partnerAttachment, _, err := client.PartnerAttachment.Get(ctx, id.(string))
		if err != nil {
			return diag.Errorf("error retrieving Partner Attachment: %s", err)
		}

		foundPartnerAttachment = partnerAttachment
	} else if name, ok := d.GetOk("name"); ok {
		partnerAttachments, err := listPartnerAttachments(client)
		if err != nil {
			return diag.Errorf("error retrieving Partner Attachments: %s", err)
		}

		partnerAttachment, err := findPartnerAttachmentByName(partnerAttachments, name.(string))
		if err != nil {
			return diag.Errorf("error retrieving Partner Attachment: %s", err)
		}

		foundPartnerAttachment = partnerAttachment
	}

	if foundPartnerAttachment == nil {
		return diag.Errorf("Bad Request: %s", fmt.Errorf("'name' or 'id' must be provided"))
	}

	d.SetId(foundPartnerAttachment.ID)
	d.Set("name", foundPartnerAttachment.Name)
	d.Set("connection_bandwidth_in_mbps", foundPartnerAttachment.ConnectionBandwidthInMbps)
	d.Set("region", strings.ToLower(foundPartnerAttachment.Region))
	d.Set("naas_provider", foundPartnerAttachment.NaaSProvider)
	d.Set("vpc_ids", foundPartnerAttachment.VPCIDs)
	if bgp := foundPartnerAttachment.BGP; bgp.PeerRouterIP != "" || bgp.LocalRouterIP != "" || bgp.PeerASN != 0 {
		bgpMap := map[string]interface{}{
			"local_router_ip": bgp.LocalRouterIP,
			"peer_router_asn": bgp.PeerASN,
			"peer_router_ip":  bgp.PeerRouterIP,
		}
		if err := d.Set("bgp", []interface{}{bgpMap}); err != nil {
			return diag.FromErr(err)
		}
	}
	d.Set("state", foundPartnerAttachment.State)
	d.Set("created_at", foundPartnerAttachment.CreatedAt.UTC().String())

	return nil
}

func listPartnerAttachments(client *godo.Client) ([]*godo.PartnerAttachment, error) {
	paList := []*godo.PartnerAttachment{}
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		partnerAttachments, resp, err := client.PartnerAttachment.List(context.Background(), opts)
		if err != nil {
			return nil, err
		}

		paList = append(paList, partnerAttachments...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return paList, fmt.Errorf("error retrieving Partner Attachments: %s", err)
		}

		opts.Page = page + 1
	}

	return paList, nil
}

func findPartnerAttachmentByName(partnerAttachments []*godo.PartnerAttachment, name string) (*godo.PartnerAttachment, error) {
	for _, partnerAttachment := range partnerAttachments {
		if partnerAttachment.Name == name {
			return partnerAttachment, nil
		}
	}

	return nil, fmt.Errorf("no Partner Attachment found with name: %s", name)
}
