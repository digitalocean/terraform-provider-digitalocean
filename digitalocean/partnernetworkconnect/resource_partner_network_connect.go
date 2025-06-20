package partnernetworkconnect

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
)

func ResourceDigitalOceanPartnerAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanPartnerAttachmentCreate,
		ReadContext:   resourceDigitalOceanPartnerAttachmentRead,
		UpdateContext: resourceDigitalOceanPartnerAttachmentUpdate,
		DeleteContext: resourceDigitalOceanPartnerAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Partner Attachment",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the Partner Attachment",
				ValidateFunc: validation.NoZeroValues,
			},
			"connection_bandwidth_in_mbps": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The connection bandwidth in Mbps",
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The region where the Partner Attachment will be created",
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			"naas_provider": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The NaaS provider",
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			"redundancy_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The redundancy zone for the NaaS",
				ForceNew:    true,
			},
			"vpc_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The list of VPC IDs to attach the Partner Attachment to",
				Set:         schema.HashString,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the Partner Attachment",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the Partner Attachment was created",
			},
			"bgp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
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
						"auth_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"parent_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Parent Partner Attachment ",
			},
			"children": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The children of Partner Attachment",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
	}
}

func resourceDigitalOceanPartnerAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	name := d.Get("name").(string)
	connectionBandwidthInMbps := d.Get("connection_bandwidth_in_mbps").(int)
	region := d.Get("region").(string)
	naasProvider := d.Get("naas_provider").(string)
	redundancyZone := d.Get("redundancy_zone").(string)
	vpcIDs := d.Get("vpc_ids").(*schema.Set).List()
	parentUUID := d.Get("parent_uuid").(string)

	vpcIDsString := make([]string, len(vpcIDs))
	for i, v := range vpcIDs {
		vpcIDsString[i] = v.(string)
	}

	partnerAttachmentRequest := &godo.PartnerAttachmentCreateRequest{
		Name:                      name,
		ConnectionBandwidthInMbps: connectionBandwidthInMbps,
		Region:                    region,
		NaaSProvider:              naasProvider,
		RedundancyZone:            redundancyZone,
		VPCIDs:                    vpcIDsString,
		ParentUuid:                parentUUID,
	}

	if bgpRaw, ok := d.GetOk("bgp"); ok {
		bgpList := bgpRaw.([]interface{})
		if len(bgpList) > 0 {
			bgpConfig := bgpList[0].(map[string]interface{})
			bgp := godo.BGP{
				LocalRouterIP: bgpConfig["local_router_ip"].(string),
				PeerASN:       bgpConfig["peer_router_asn"].(int),
				PeerRouterIP:  bgpConfig["peer_router_ip"].(string),
				AuthKey:       bgpConfig["auth_key"].(string),
			}
			partnerAttachmentRequest.BGP = bgp
		}
	}

	log.Printf("[DEBUG] Partner Attachment create request: %#v", partnerAttachmentRequest)

	partnerAttachment, resp, err := client.PartnerAttachment.Create(context.Background(), partnerAttachmentRequest)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusBadRequest, http.StatusUnprocessableEntity, http.StatusConflict:
				return diag.FromErr(fmt.Errorf("failed to create Partner Attachment: %s", err))
			}
		}
		return diag.FromErr(fmt.Errorf("error creating Partner Attachment: %s", err))
	}

	d.SetId(partnerAttachment.ID)

	log.Printf("[DEBUG] Waiting for Partner Attachment (%s) to become active", d.Get("name"))
	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"CREATING"},
		Target:     []string{"CREATED"},
		Refresh:    partnerAttachmentStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for Partner Attachment (%s) to become active: %s", d.Get("name"), err))
	}

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Partner Attachment created, ID: %s", d.Id())

	return resourceDigitalOceanPartnerAttachmentRead(ctx, d, meta)
}

func resourceDigitalOceanPartnerAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	updateRequest := &godo.PartnerAttachmentUpdateRequest{}

	if d.HasChange("name") {
		updateRequest.Name = d.Get("name").(string)
	}

	if d.HasChange("vpc_ids") {
		var vpcIDsString []string
		for _, v := range d.Get("vpc_ids").(*schema.Set).List() {
			vpcIDsString = append(vpcIDsString, v.(string))
		}
		updateRequest.VPCIDs = vpcIDsString
	}

	if updateRequest.Name == "" && len(updateRequest.VPCIDs) == 0 {
		return nil
	}

	_, _, err := client.PartnerAttachment.Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating Partner Attachment: %s", err)
	}

	if updateRequest.Name != "" {
		log.Printf("[INFO] Updated Partner Attachment Name")
	}

	if len(updateRequest.VPCIDs) != 0 {
		log.Printf("[INFO] Updated Partner Attachment VPC IDs")
	}

	return nil
}

func resourceDigitalOceanPartnerAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	partnerAttachmentID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		resp, err := client.PartnerAttachment.Delete(context.Background(), partnerAttachmentID)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusForbidden {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(fmt.Errorf("error deleting Partner Attachment: %s", err))
		}

		log.Printf("[DEBUG] Waiting for Partner Attachment (%s) to be deleted", d.Get("name"))
		stateConf := &retry.StateChangeConf{
			Delay:      5 * time.Second,
			Pending:    []string{"DELETING"},
			Target:     []string{http.StatusText(http.StatusNotFound)},
			Refresh:    partnerAttachmentStateRefreshFunc(client, partnerAttachmentID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			MinTimeout: 5 * time.Second,
		}
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return retry.NonRetryableError(fmt.Errorf("error waiting for Partner Attachment (%s) to be deleted: %s", d.Get("name"), err))
		}

		d.SetId("")
		log.Printf("[INFO] Partner Attachment deleted, ID: %s", partnerAttachmentID)

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDigitalOceanPartnerAttachmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	partnerAttachment, resp, err := client.PartnerAttachment.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Partner Attachment (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading Partner Attachment: %s", err)
	}

	d.SetId(partnerAttachment.ID)
	d.Set("name", partnerAttachment.Name)
	d.Set("state", partnerAttachment.State)
	d.Set("created_at", partnerAttachment.CreatedAt.UTC().String())
	d.Set("region", strings.ToLower(partnerAttachment.Region))
	d.Set("connection_bandwidth_in_mbps", partnerAttachment.ConnectionBandwidthInMbps)
	d.Set("naas_provider", partnerAttachment.NaaSProvider)
	d.Set("redundancy_zone", partnerAttachment.RedundancyZone)
	d.Set("vpc_ids", partnerAttachment.VPCIDs)
	d.Set("parent_uuid", partnerAttachment.ParentUuid)
	d.Set("children", partnerAttachment.Children)

	bgp := partnerAttachment.BGP
	if bgp.PeerRouterIP != "" || bgp.LocalRouterIP != "" || bgp.PeerASN != 0 {
		bgpMap := map[string]interface{}{
			"local_router_ip": bgp.LocalRouterIP,
			"peer_router_asn": bgp.PeerASN,
			"peer_router_ip":  bgp.PeerRouterIP,
			"auth_key":        d.Get("bgp.0.auth_key").(string),
		}
		if err := d.Set("bgp", []interface{}{bgpMap}); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func partnerAttachmentStateRefreshFunc(client *godo.Client, id string) retry.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		partnerAttachment, resp, err := client.PartnerAttachment.Get(context.Background(), id)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return partnerAttachment, http.StatusText(resp.StatusCode), nil
			}
			return nil, "", fmt.Errorf("error issuing read request in partnerAttachmentStateRefreshFunc to DigitalOcean for Partner Attachment '%s': %s", id, err)
		}

		return partnerAttachment, partnerAttachment.State, nil
	}
}
