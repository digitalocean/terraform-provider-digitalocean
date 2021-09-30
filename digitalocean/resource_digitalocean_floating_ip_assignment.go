package digitalocean

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanFloatingIpAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanFloatingIpAssignmentCreate,
		ReadContext:   resourceDigitalOceanFloatingIpAssignmentRead,
		DeleteContext: resourceDigitalOceanFloatingIpAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"droplet_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceDigitalOceanFloatingIpAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	ip_address := d.Get("ip_address").(string)
	droplet_id := d.Get("droplet_id").(int)

	log.Printf("[INFO] Assigning the Floating IP (%s) to the Droplet %d", ip_address, droplet_id)
	action, _, err := client.FloatingIPActions.Assign(context.Background(), ip_address, droplet_id)
	if err != nil {
		return diag.Errorf(
			"Error Assigning FloatingIP (%s) to the droplet: %s", ip_address, err)
	}

	_, unassignedErr := waitForFloatingIPAssignmentReady(ctx, d, "done", []string{"new", "in-progress"}, "status", meta, action.ID)
	if unassignedErr != nil {
		return diag.Errorf(
			"Error waiting for FloatingIP (%s) to be Assigned: %s", ip_address, unassignedErr)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%d-%s-", droplet_id, ip_address)))
	return resourceDigitalOceanFloatingIpAssignmentRead(ctx, d, meta)
}

func resourceDigitalOceanFloatingIpAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	ip_address := d.Get("ip_address").(string)
	droplet_id := d.Get("droplet_id").(int)

	log.Printf("[INFO] Reading the details of the FloatingIP %s", ip_address)
	floatingIp, _, err := client.FloatingIPs.Get(context.Background(), ip_address)
	if err != nil {
		return diag.Errorf("Error retrieving FloatingIP: %s", err)
	}

	if floatingIp.Droplet == nil || floatingIp.Droplet.ID != droplet_id {
		log.Printf("[INFO] A droplet was detected on the FloatingIP.")
		d.SetId("")
	}

	return nil
}

func resourceDigitalOceanFloatingIpAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	ip_address := d.Get("ip_address").(string)
	droplet_id := d.Get("droplet_id").(int)

	log.Printf("[INFO] Reading the details of the FloatingIP %s", ip_address)
	floatingIp, _, err := client.FloatingIPs.Get(context.Background(), ip_address)
	if err != nil {
		return diag.Errorf("Error retrieving FloatingIP: %s", err)
	}

	if floatingIp.Droplet.ID == droplet_id {
		log.Printf("[INFO] Unassigning the Floating IP from the Droplet")
		action, _, err := client.FloatingIPActions.Unassign(context.Background(), ip_address)
		if err != nil {
			return diag.Errorf("Error unassigning FloatingIP (%s) from the droplet: %s", ip_address, err)
		}

		_, unassignedErr := waitForFloatingIPAssignmentReady(ctx, d, "done", []string{"new", "in-progress"}, "status", meta, action.ID)
		if unassignedErr != nil {
			return diag.Errorf(
				"Error waiting for FloatingIP (%s) to be unassigned: %s", ip_address, unassignedErr)
		}
	} else {
		log.Printf("[INFO] Floating IP already unassigned, removing from state.")
	}

	d.SetId("")
	return nil
}

func waitForFloatingIPAssignmentReady(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionId int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for FloatingIP (%s) to have %s of %s",
		d.Get("ip_address").(string), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newFloatingIPAssignmentStateRefreshFunc(d, attribute, meta, actionId),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newFloatingIPAssignmentStateRefreshFunc(
	d *schema.ResourceData, attribute string, meta interface{}, actionId int) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).godoClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Refreshing the Floating IP state")
		action, _, err := client.FloatingIPActions.Get(context.Background(), d.Get("ip_address").(string), actionId)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving FloatingIP (%s) ActionId (%d): %s", d.Get("ip_address").(string), actionId, err)
		}

		log.Printf("[INFO] The FloatingIP Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}
