package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanReservedIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanReservedIPCreate,
		UpdateContext: resourceDigitalOceanReservedIPUpdate,
		ReadContext:   resourceDigitalOceanReservedIPRead,
		DeleteContext: resourceDigitalOceanReservedIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanReservedIPImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the reserved ip",
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"droplet_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDigitalOceanReservedIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Create a reserved IP In a Region")
	regionOpts := &godo.ReservedIPCreateRequest{
		Region: d.Get("region").(string),
	}

	log.Printf("[DEBUG] ReservedIP Create: %#v", regionOpts)
	reservedIP, _, err := client.ReservedIPs.Create(context.Background(), regionOpts)
	if err != nil {
		return diag.Errorf("Error creating reserved IP: %s", err)
	}

	d.SetId(reservedIP.IP)

	if v, ok := d.GetOk("droplet_id"); ok {

		log.Printf("[INFO] Assigning the reserved IP to the Droplet %d", v.(int))
		action, _, err := client.ReservedIPActions.Assign(context.Background(), d.Id(), v.(int))
		if err != nil {
			return diag.Errorf(
				"Error Assigning reserved IP (%s) to the droplet: %s", d.Id(), err)
		}

		_, unassignedErr := waitForReservedIPReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
		if unassignedErr != nil {
			return diag.Errorf(
				"Error waiting for reserved IP (%s) to be Assigned: %s", d.Id(), unassignedErr)
		}
	}

	return resourceDigitalOceanReservedIPRead(ctx, d, meta)
}

func resourceDigitalOceanReservedIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	if d.HasChange("droplet_id") {
		if v, ok := d.GetOk("droplet_id"); ok {
			log.Printf("[INFO] Assigning the reserved IP %s to the Droplet %d", d.Id(), v.(int))
			action, _, err := client.ReservedIPActions.Assign(context.Background(), d.Id(), v.(int))
			if err != nil {
				return diag.Errorf(
					"Error Assigning reserved IP (%s) to the droplet: %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IP (%s) to be Assigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[INFO] Unassigning the reserved IP %s", d.Id())
			action, _, err := client.ReservedIPActions.Unassign(context.Background(), d.Id())
			if err != nil {
				return diag.Errorf(
					"Error unassigning reserved IP (%s): %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IP (%s) to be Unassigned: %s", d.Id(), unassignedErr)
			}
		}
	}

	return resourceDigitalOceanReservedIPRead(ctx, d, meta)
}

func resourceDigitalOceanReservedIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Reading the details of the reserved IP %s", d.Id())
	reservedIP, resp, err := client.ReservedIPs.Get(context.Background(), d.Id())
	if resp.StatusCode != 404 {
		if err != nil {
			return diag.Errorf("Error retrieving reserved IP: %s", err)
		}

		if _, ok := d.GetOk("droplet_id"); ok && reservedIP.Droplet != nil {
			d.Set("region", reservedIP.Droplet.Region.Slug)
			d.Set("droplet_id", reservedIP.Droplet.ID)
		} else {
			d.Set("region", reservedIP.Region.Slug)
		}

		d.Set("ip_address", reservedIP.IP)
		d.Set("urn", reservedIP.URN())
	} else {
		d.SetId("")
	}

	return nil
}

func resourceDigitalOceanReservedIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	if _, ok := d.GetOk("droplet_id"); ok {
		log.Printf("[INFO] Unassigning the reserved IP from the Droplet")
		action, resp, err := client.ReservedIPActions.Unassign(context.Background(), d.Id())
		if resp.StatusCode != 422 {
			if err != nil {
				return diag.Errorf(
					"Error unassigning reserved IP (%s) from the droplet: %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IP (%s) to be unassigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[DEBUG] Couldn't unassign reserved IP (%s) from droplet, possibly out of sync: %s", d.Id(), err)
		}
	}

	log.Printf("[INFO] Deleting reserved IP: %s", d.Id())
	_, err := client.ReservedIPs.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting reserved IP: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanReservedIPImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*CombinedConfig).godoClient()
	reservedIP, resp, err := client.ReservedIPs.Get(context.Background(), d.Id())
	if resp.StatusCode != 404 {
		if err != nil {
			return nil, err
		}

		d.Set("ip_address", reservedIP.IP)
		d.Set("urn", reservedIP.URN())
		d.Set("region", reservedIP.Region.Slug)

		if reservedIP.Droplet != nil {
			d.Set("droplet_id", reservedIP.Droplet.ID)
		}
	}

	return []*schema.ResourceData{d}, nil
}

func waitForReservedIPReady(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IP (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPStateRefreshFunc(d, attribute, meta, actionID),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPStateRefreshFunc(
	d *schema.ResourceData, attribute string, meta interface{}, actionID int) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).godoClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Assigning the reserved IP to the Droplet")
		action, _, err := client.ReservedIPActions.Get(context.Background(), d.Id(), actionID)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving reserved IP (%s) ActionId (%d): %s", d.Id(), actionID, err)
		}

		log.Printf("[INFO] The reserved IP Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}
