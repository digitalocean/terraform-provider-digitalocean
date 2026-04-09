package nfs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func isTransientNfsActionError(err error) bool {
	if strings.Contains(strings.ToLower(err.Error()), "timeout") {
		return true
	}

	if doErr, ok := err.(*godo.ErrorResponse); ok && doErr.Response != nil {
		status := doErr.Response.StatusCode
		return status == 409 || status == 429 || status >= 500
	}

	return false
}

func ResourceDigitalOceanNfsAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanNfsAttachmentCreate,
		ReadContext:   resourceDigitalOceanNfsAttachmentRead,
		DeleteContext: resourceDigitalOceanNfsAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceDigitalOceanNfsAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	vpcId := d.Get("vpc_id").(string)
	shareId := d.Get("share_id").(string)
	region := ""
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	share, _, err := client.Nfs.Get(context.Background(), shareId, region)
	if err != nil {
		return diag.Errorf("Error retrieving share: %s", err)
	}

	// If share is attached to a different VPC, use reassign
	if len(share.VpcIDs) > 0 && share.VpcIDs[0] != vpcId {
		err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
			log.Printf("[DEBUG] Reassigning share (%s) from VPC (%s) to VPC (%s)", shareId, share.VpcIDs[0], vpcId)

			err := reassignNfs(ctx, client, shareId, region, share.VpcIDs[0], vpcId)
			if err != nil {
				retryErr := fmt.Errorf("[WARN] Error reassigning share (%s) from VPC (%s) to VPC (%s): %s", shareId, share.VpcIDs[0], vpcId, err)
				if isTransientNfsActionError(err) {
					return retry.RetryableError(retryErr)
				}

				return retry.NonRetryableError(retryErr)
			}

			return nil
		})

		if err != nil {
			return diag.Errorf("Error reassigning share after retry timeout: %s", err)
		}
	} else if len(share.VpcIDs) == 0 {
		// Share is not attached to any VPC, attach it to the target VPC
		err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {

			log.Printf("[DEBUG] Attaching Share (%s) to VPC (%s)", shareId, vpcId)
			action, _, err := client.NfsActions.Attach(context.Background(), shareId, vpcId, region)
			if err != nil {

				return retry.NonRetryableError(
					fmt.Errorf("[WARN] Error attaching share (%s) to VPC (%s): %s", shareId, vpcId, err))
			}

			log.Printf("[DEBUG] Share attach action id: %d", action.ID)

			// Poll the share to check VPC Id
			if err = waitForNfsAttach(ctx, client, shareId, region, vpcId); err != nil {

				return retry.NonRetryableError(
					fmt.Errorf("[DEBUG] Error waiting for attach share (%s) to VPC (%s) to finish: %s", shareId, vpcId, err))
			}

			return nil
		})

		if err != nil {
			return diag.Errorf("Error attaching share to vpc after retry timeout: %s", err)
		}
	}

	d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%s-%s-", vpcId, shareId)))

	return nil
}

func waitForNfsAttach(ctx context.Context, client *godo.Client, id, region string, expectedVpcID string) error {
	for i := 0; i < 60; i++ {
		share, _, err := client.Nfs.Get(ctx, id, region)
		if err != nil {
			return err
		}

		if share.Status == "ACTIVE" && len(share.VpcIDs) != 0 && share.VpcIDs[0] == expectedVpcID {
			return nil
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for NFS attach to complete")
}

func waitForNfsDetach(ctx context.Context, client *godo.Client, id, region string, expectedVpcID string) error {
	for i := 0; i < 60; i++ {
		share, _, err := client.Nfs.Get(ctx, id, region)
		if err != nil {
			return err
		}

		if share.Status == "INACTIVE" && len(share.VpcIDs) == 0 {
			return nil
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for NFS detach to complete")
}

func resourceDigitalOceanNfsAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	vpcId := d.Get("vpc_id").(string)
	shareId := d.Get("share_id").(string)
	region := ""
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	share, resp, err := client.Nfs.Get(context.Background(), shareId, region)
	if err != nil {
		// If the share is already destroyed, mark as
		// successfully removed
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving share: %s", err)
	}

	if share.Status == "ACTIVE" && len(share.VpcIDs) == 0 || share.VpcIDs[0] != vpcId {
		log.Printf("[DEBUG] Share Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

func resourceDigitalOceanNfsAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	vpcId := d.Get("vpc_id").(string)
	shareId := d.Get("share_id").(string)
	region := ""
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	// Only one nfs can be detached at one time to a single vpc.
	err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {

		log.Printf("[DEBUG] Detaching Share (%s) from VPC (%s)", shareId, vpcId)
		action, _, err := client.NfsActions.Detach(context.Background(), shareId, vpcId, region)
		if err != nil {

			return retry.NonRetryableError(
				fmt.Errorf("[WARN] Error detaching share (%s) from VPC (%s): %s", shareId, vpcId, err))
		}

		log.Printf("[DEBUG] Share detach action id: %d", action.ID)
		// Poll the share to check
		if err = waitForNfsDetach(ctx, client, shareId, region, vpcId); err != nil {

			return retry.NonRetryableError(
				fmt.Errorf("[DEBUG] Error waiting for detach share (%s) to VPC (%s) to finish: %s", shareId, vpcId, err))
		}

		return nil
	})

	if err != nil {
		return diag.Errorf("Error detaching share from vpc after retry timeout: %s", err)
	}

	return nil
}
