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

func ResourceDigitalOceanVPC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanVPCPeeringCreate,
		ReadContext:   resourceDigitalOceanVPCPeeringRead,
		UpdateContext: resourceDigitalOceanVPCPeeringUpdate,
		DeleteContext: resourceDigitalOceanVPCPeeringDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the VPC Peering",
				ValidateFunc: validation.NoZeroValues,
			},
			"vpc_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The list of VPCs to be peered",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: validateVPCIDs,
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
	vpcIDs := d.Get("vpc_ids").([]any)

	vpcIDStrings := make([]string, len(vpcIDs))
	for i, v := range vpcIDs {
		vpcIDStrings[i] = v.(string)
	}

	vpcPeeringRequest := &godo.VPCPeeringCreateRequest{
		Name:   name,
		VPCIDs: vpcIDStrings,
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
			if resp.StatusCode == http.StatusForbidden {
				return retry.RetryableError(err)
			} else {
				return retry.NonRetryableError(fmt.Errorf("error deleting VPC Peering: %s", err))
			}
		}

		d.SetId("")
		log.Printf("[INFO] VPC Peering deleted, ID: %s", vpcPeeringID)

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	} else {
		return nil
	}
}

func resourceDigitalOceanVPCPeeringRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func validateVPCIDs(val any, key string) ([]string, []error) {
	v := val.([]any)
	if len(v) != 2 {
		return nil, []error{fmt.Errorf("%q must contain exactly 2 VPC IDs, got %d", key, len(v))}
	}

	vpcIDMap := make(map[string]struct{})
	for _, id := range v {
		vpcID := id.(string)
		if _, ok := vpcIDMap[vpcID]; ok {
			return nil, []error{fmt.Errorf("%q must contain unique VPC IDs, duplicate found: %s", key, vpcID)}
		}
		vpcIDMap[vpcID] = struct{}{}
	}

	return nil, nil
}
