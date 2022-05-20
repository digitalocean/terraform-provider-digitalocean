package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanReservedIPAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanReservedIPAssignmentCreate,
		ReadContext:   resourceDigitalOceanReservedIPAssignmentRead,
		DeleteContext: resourceDigitalOceanReservedIPAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanReservedIPAssignmentImport,
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

func resourceDigitalOceanReservedIPAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	ipAddress := d.Get("ip_address").(string)
	dropletID := d.Get("droplet_id").(int)

	log.Printf("[INFO] Assigning the reserved IP (%s) to the Droplet %d", ipAddress, dropletID)
	action, _, err := client.ReservedIPActions.Assign(context.Background(), ipAddress, dropletID)
	if err != nil {
		return diag.Errorf(
			"Error Assigning reserved IP (%s) to the droplet: %s", ipAddress, err)
	}

	_, unassignedErr := waitForReservedIPAssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
	if unassignedErr != nil {
		return diag.Errorf(
			"Error waiting for reserved IP (%s) to be Assigned: %s", ipAddress, unassignedErr)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%d-%s-", dropletID, ipAddress)))
	return resourceDigitalOceanReservedIPAssignmentRead(ctx, d, meta)
}

func resourceDigitalOceanReservedIPAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	ipAddress := d.Get("ip_address").(string)
	dropletID := d.Get("droplet_id").(int)

	log.Printf("[INFO] Reading the details of the reserved IP %s", ipAddress)
	reservedIP, _, err := client.ReservedIPs.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IP: %s", err)
	}

	if reservedIP.Droplet == nil || reservedIP.Droplet.ID != dropletID {
		log.Printf("[INFO] A droplet was detected on the reserved IP.")
		d.SetId("")
	}

	return nil
}

func resourceDigitalOceanReservedIPAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	ipAddress := d.Get("ip_address").(string)
	dropletID := d.Get("droplet_id").(int)

	log.Printf("[INFO] Reading the details of the reserved IP %s", ipAddress)
	reservedIP, _, err := client.ReservedIPs.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IP: %s", err)
	}

	if reservedIP.Droplet.ID == dropletID {
		log.Printf("[INFO] Unassigning the reserved IP from the Droplet")
		action, _, err := client.ReservedIPActions.Unassign(context.Background(), ipAddress)
		if err != nil {
			return diag.Errorf("Error unassigning reserved IP (%s) from the droplet: %s", ipAddress, err)
		}

		_, unassignedErr := waitForReservedIPAssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
		if unassignedErr != nil {
			return diag.Errorf(
				"Error waiting for reserved IP (%s) to be unassigned: %s", ipAddress, unassignedErr)
		}
	} else {
		log.Printf("[INFO] reserved IP already unassigned, removing from state.")
	}

	d.SetId("")
	return nil
}

func waitForReservedIPAssignmentReady(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IP (%s) to have %s of %s",
		d.Get("ip_address").(string), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPAssignmentStateRefreshFunc(d, attribute, meta, actionID),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPAssignmentStateRefreshFunc(
	d *schema.ResourceData, attribute string, meta interface{}, actionID int) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).godoClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Refreshing the reserved IP state")
		action, _, err := client.ReservedIPActions.Get(context.Background(), d.Get("ip_address").(string), actionID)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving reserved IP (%s) ActionId (%d): %s", d.Get("ip_address").(string), actionID, err)
		}

		log.Printf("[INFO] The reserved IP Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}

func resourceDigitalOceanReservedIPAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-%s-", s[1], s[0])))
		d.Set("ip_address", s[0])
		dropletID, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		d.Set("droplet_id", dropletID)
	} else {
		return nil, errors.New("must use the reserved IP and the ID of the Droplet joined with a comma (e.g. `ip_address,droplet_id`)")
	}

	return []*schema.ResourceData{d}, nil
}
