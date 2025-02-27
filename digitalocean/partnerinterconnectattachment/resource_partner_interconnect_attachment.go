package partnerinterconnectattachment

import (
	"context"
	"fmt"
	"log"
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

func ResourceDigitalOceanPartnerInterconnectAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanPartnerInterconnectAttachmentCreate,
		ReadContext:   resourceDigitalOceanPartnerInterconnectAttachmentRead,
		UpdateContext: resourceDigitalOceanPartnerInterconnectAttachmentUpdate,
		DeleteContext: resourceDigitalOceanPartnerInterconnectAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Partner Interconnect Attachment",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the Partner Interconnect Attachment",
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
				Description:  "The region where the Partner Interconnect Attachment will be created",
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
			"vpc_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The list of VPC IDs to attach the Partner Interconnect Attachment to",
				Set:         schema.HashString,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the Partner Interconnect Attachment",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the Partner Interconnect Attachment was created",
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
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
	}
}

func resourceDigitalOceanPartnerInterconnectAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	name := d.Get("name").(string)
	connectionBandwidthInMbps := d.Get("connection_bandwidth_in_mbps").(int)
	region := d.Get("region").(string)
	naasProvider := d.Get("naas_provider").(string)
	vpcIDs := d.Get("vpc_ids").(*schema.Set).List()

	vpcIDsString := make([]string, len(vpcIDs))
	for i, v := range vpcIDs {
		vpcIDsString[i] = v.(string)
	}

	partnerInterconnectAttachmentRequest := &godo.PartnerInterconnectAttachmentCreateRequest{
		Name:                      name,
		ConnectionBandwidthInMbps: connectionBandwidthInMbps,
		Region:                    region,
		NaaSProvider:              naasProvider,
		VPCIDs:                    vpcIDsString,
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
			partnerInterconnectAttachmentRequest.BGP = bgp
		}
	}

	log.Printf("[DEBUG] Partner Interconnect Attachment create request: %#v", partnerInterconnectAttachmentRequest)

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		partnerInterconnectAttachment, resp, err := client.PartnerInterconnectAttachments.Create(context.Background(), partnerInterconnectAttachmentRequest)
		if err != nil {
			if resp != nil {
				switch resp.StatusCode {
				case http.StatusBadRequest, http.StatusUnprocessableEntity, http.StatusConflict:
					return retry.NonRetryableError(fmt.Errorf("failed to create Partner Interconnect Attachment: %s", err))
				}
			}
			return retry.RetryableError(fmt.Errorf("error creating Partner Interconnect Attachment: %s", err))
		}

		d.SetId(partnerInterconnectAttachment.ID)

		log.Printf("[DEBUG] Waiting for Partner Interconnect Attachment (%s) to become active", d.Get("name"))
		stateConf := &retry.StateChangeConf{
			Delay:      5 * time.Second,
			Pending:    []string{"CREATING"},
			Target:     []string{"CREATED"},
			Refresh:    partnerInterconnectAttachmentStateRefreshFunc(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			MinTimeout: 5 * time.Second,
		}
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return retry.NonRetryableError(fmt.Errorf("error waiting for Partner Interconnect Attachment (%s) to become active: %s", d.Get("name"), err))
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Partner Interconnect Attachment created, ID: %s", d.Id())

	return resourceDigitalOceanPartnerInterconnectAttachmentRead(ctx, d, meta)
}

func resourceDigitalOceanPartnerInterconnectAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	updateRequest := &godo.PartnerInterconnectAttachmentUpdateRequest{}

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

	_, _, err := client.PartnerInterconnectAttachments.Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating Partner Interconnect Attachment: %s", err)
	}

	if updateRequest.Name != "" {
		log.Printf("[INFO] Updated Partner Interconnect Attachment Name")
	}

	if len(updateRequest.VPCIDs) != 0 {
		log.Printf("[INFO] Updated Partner Interconnect Attachment VPC IDs")
	}

	return nil
}

func resourceDigitalOceanPartnerInterconnectAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	partnerInterconnectAttachmentID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		resp, err := client.PartnerInterconnectAttachments.Delete(context.Background(), partnerInterconnectAttachmentID)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusForbidden {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(fmt.Errorf("error deleting Partner Interconnect Attachment: %s", err))
		}

		log.Printf("[DEBUG] Waiting for Partner Interconnect Attachment (%s) to be deleted", d.Get("name"))
		stateConf := &retry.StateChangeConf{
			Delay:      5 * time.Second,
			Pending:    []string{"DELETING"},
			Target:     []string{http.StatusText(http.StatusNotFound)},
			Refresh:    partnerInterconnectAttachmentStateRefreshFunc(client, partnerInterconnectAttachmentID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			MinTimeout: 5 * time.Second,
		}
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return retry.NonRetryableError(fmt.Errorf("error waiting for Partner Interconnect Attachment (%s) to be deleted: %s", d.Get("name"), err))
		}

		d.SetId("")
		log.Printf("[INFO] Partner Interconnect Attachment deleted, ID: %s", partnerInterconnectAttachmentID)

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDigitalOceanPartnerInterconnectAttachmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	partnerInterconectAttachment, resp, err := client.PartnerInterconnectAttachments.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Partner Interconnect Attachment (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading Partner Interconnect Attachment: %s", err)
	}

	d.SetId(partnerInterconectAttachment.ID)
	d.Set("name", partnerInterconectAttachment.Name)
	d.Set("state", partnerInterconectAttachment.State)
	d.Set("created_at", partnerInterconectAttachment.CreatedAt.UTC().String())
	d.Set("region", strings.ToLower(partnerInterconectAttachment.Region))
	d.Set("connection_bandwidth_in_mbps", partnerInterconectAttachment.ConnectionBandwidthInMbps)
	d.Set("naas_provider", partnerInterconectAttachment.NaaSProvider)
	d.Set("vpc_ids", partnerInterconectAttachment.VPCIDs)

	bgp := partnerInterconectAttachment.BGP
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

func partnerInterconnectAttachmentStateRefreshFunc(client *godo.Client, id string) retry.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		partnerInterconnectAttachment, resp, err := client.PartnerInterconnectAttachments.Get(context.Background(), id)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return partnerInterconnectAttachment, http.StatusText(resp.StatusCode), nil
			}
			return nil, "", fmt.Errorf("error issuingn read request in partnerInterconnectAttachmentStateRefreshFunc to DigitalOcean for Partner Interconnect Attachment '%s': %s", id, err)
		}

		return partnerInterconnectAttachment, partnerInterconnectAttachment.State, nil
	}
}
