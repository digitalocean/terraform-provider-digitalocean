package vpcpeering

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanVPCPeering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanVPCPeeringCreate,
		ReadContext:   resourceDigitalOceanVPCPeeringRead,
		UpdateContext: resourceDigitalOceanVPCPeeringUpdate,
		DeleteContext: resourceDigitalOceanVPCPeeringDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the VPC Peering",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the VPC Peering",
				ValidateFunc: validation.NoZeroValues,
			},
			"vpc_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "The list of VPCs to be peered",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the VPC Peering",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the VPC was created",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
	}
}

func resourceDigitalOceanVPCPeeringCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	name := d.Get("name").(string)
	vpcIDs := d.Get("vpc_ids").(*schema.Set).List()

	vpcIDsString := make([]string, len(vpcIDs))
	for i, v := range vpcIDs {
		vpcIDsString[i] = v.(string)
	}

	vpcPeeringRequest := &godo.VPCPeeringCreateRequest{
		Name:   name,
		VPCIDs: vpcIDsString,
	}

	log.Printf("[DEBUG] VPC Peering create request: %#v", vpcPeeringRequest)
	vpcPeering, _, err := client.VPCs.CreateVPCPeering(context.Background(), vpcPeeringRequest)
	if err != nil {
		return diag.Errorf("Error creating VPC Peering: %s", err)
	}

	d.SetId(vpcPeering.ID)
	log.Printf("[INFO] VPC Peering created, ID: %s", d.Id())

	return resourceDigitalOceanVPCPeeringRead(ctx, d, meta)
}

func resourceDigitalOceanVPCPeeringUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChange("name") {
		vpcPeeringUpdateRequest := &godo.VPCPeeringUpdateRequest{
			Name: d.Get("name").(string),
		}

		_, _, err := client.VPCs.UpdateVPCPeering(context.Background(), d.Id(), vpcPeeringUpdateRequest)
		if err != nil {
			return diag.Errorf("Error updating VPC Peering: %s", err)
		}
		log.Printf("[INFO] Updated VPC Peering")
	}

	return resourceDigitalOceanVPCPeeringRead(ctx, d, meta)
}

func resourceDigitalOceanVPCPeeringDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	vpcPeeringID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		resp, err := client.VPCs.DeleteVPCPeering(context.Background(), vpcPeeringID)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusForbidden {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(fmt.Errorf("error deleting VPC Peering: %s", err))
		}

		if err := waitForVPCPeeringDeletion(ctx, client, vpcPeeringID, d.Timeout(schema.TimeoutDelete)); err != nil {
			return retry.NonRetryableError(err)
		}

		d.SetId("")
		log.Printf("[INFO] VPC Peering deleted, ID: %s", vpcPeeringID)

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDigitalOceanVPCPeeringRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	vpcPeering, resp, err := client.VPCs.GetVPCPeering(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] VPC Peering (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading VPC Peering: %s", err)
	}

	d.SetId(vpcPeering.ID)
	d.Set("name", vpcPeering.Name)
	d.Set("status", vpcPeering.Status)
	d.Set("vpc_ids", vpcPeering.VPCIDs)
	d.Set("created_at", vpcPeering.CreatedAt.UTC().String())

	return nil
}

func waitForVPCPeeringDeletion(ctx context.Context, client *godo.Client, vpcPeeringID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		_, resp, err := client.VPCs.GetVPCPeering(ctx, vpcPeeringID)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				// VPC peering is fully deleted
				return nil
			}
			return fmt.Errorf("error checking VPC Peering status: %s", err)
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for VPC Peering to be deleted")
		}

		time.Sleep(2 * time.Second)
	}
}
