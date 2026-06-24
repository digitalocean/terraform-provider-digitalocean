package reservedip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func ResourceDigitalOceanReservedIPAssignment() *schema.Resource {
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
	client := meta.(*config.CombinedConfig).GodoClient()

	ipAddress := d.Get("ip_address").(string)
	dropletID := d.Get("droplet_id").(int)

	log.Printf("[INFO] Assigning the reserved IP (%s) to the Droplet %d", ipAddress, dropletID)
	action, _, err := client.ReservedIPActions.Assign(context.Background(), ipAddress, dropletID)
	if err != nil {
		return diag.Errorf(
			"Error Assigning reserved IP (%s) to the droplet: %s", ipAddress, err)
	}

	_, unassignedErr := waitForReservedIPAssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID, reservedIPAssignmentCheck{dropletID: dropletID, assign: true})
	if unassignedErr != nil {
		return diag.Errorf(
			"Error waiting for reserved IP (%s) to be Assigned: %s", ipAddress, unassignedErr)
	}

	d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%d-%s-", dropletID, ipAddress)))
	return resourceDigitalOceanReservedIPAssignmentRead(ctx, d, meta)
}

func resourceDigitalOceanReservedIPAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	ipAddress := d.Get("ip_address").(string)
	dropletID := d.Get("droplet_id").(int)

	log.Printf("[INFO] Reading the details of the reserved IP %s", ipAddress)
	reservedIP, _, err := client.ReservedIPs.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IP: %s", err)
	}

	if reservedIP.Droplet == nil || reservedIP.Droplet.ID != dropletID {
		log.Printf("[INFO] A Droplet was detected on the reserved IP.")
		d.SetId("")
	}

	return nil
}

func resourceDigitalOceanReservedIPAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

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

		_, unassignedErr := waitForReservedIPAssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID, reservedIPAssignmentCheck{dropletID: dropletID, assign: false})
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

// reservedIPAssignmentCheck describes the end-state we expect once an
// assign/unassign action settles. It lets us confirm completion by reading the
// reserved IP directly when the action-status endpoint returns 404 - which is
// persistent for BYOIP-prefix reserved IPs, whose returned action ID never
// becomes queryable. Without this fallback, polling the missing action loops
// until the timeout even though DigitalOcean already completed the assignment.
type reservedIPAssignmentCheck struct {
	dropletID int
	assign    bool // true: expect the IP attached to dropletID; false: expect it detached
}

// reservedIPActionNotFound reports whether the action-status lookup came back as
// a 404 (the action ID is not, or not yet, queryable).
func reservedIPActionNotFound(resp *godo.Response, err error) bool {
	if err == nil {
		return false
	}
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	if godoErr, ok := err.(*godo.ErrorResponse); ok && godoErr.Response != nil && godoErr.Response.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}

func waitForReservedIPAssignmentReady(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int, check reservedIPAssignmentCheck) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IP (%s) to have %s of %s",
		d.Get("ip_address").(string), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPAssignmentStateRefreshFunc(d, attribute, meta, actionID, check),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPAssignmentStateRefreshFunc(
	d *schema.ResourceData, attribute string, meta interface{}, actionID int, check reservedIPAssignmentCheck) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()
	ipAddress := d.Get("ip_address").(string)
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Refreshing the reserved IP state")
		action, resp, err := client.ReservedIPActions.Get(context.Background(), ipAddress, actionID)
		if reservedIPActionNotFound(resp, err) {
			// The assign/unassign call returns an action ID immediately, but the
			// action may not be queryable yet - and for BYOIP-prefix reserved IPs it
			// never becomes queryable. Rather than poll the missing action until the
			// timeout, fall back to the reserved IP itself as the source of truth.
			reservedIP, _, ripErr := client.ReservedIPs.Get(context.Background(), ipAddress)
			if ripErr != nil {
				log.Printf("[DEBUG] Reserved IP (%s) action (%d) returned 404; reserved IP read failed (%s); will retry", ipAddress, actionID, ripErr)
				return nil, "", nil
			}

			attachedToTarget := reservedIP.Droplet != nil && reservedIP.Droplet.ID == check.dropletID
			if attachedToTarget == check.assign {
				log.Printf("[INFO] Reserved IP (%s) action (%d) not queryable; confirmed desired state on the reserved IP -> completed", ipAddress, actionID)
				return reservedIP, "completed", nil
			}

			// Not yet in the desired state. Return a non-nil result so the poll stays
			// "in-progress" without consuming a NotFoundChecks tick.
			log.Printf("[DEBUG] Reserved IP (%s) action (%d) returned 404; reserved IP not yet in desired state; will retry", ipAddress, actionID)
			return reservedIP, "in-progress", nil
		}
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving reserved IP (%s) ActionId (%d): %s", ipAddress, actionID, err)
		}

		log.Printf("[INFO] The reserved IP Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}

func resourceDigitalOceanReservedIPAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%s-%s-", s[1], s[0])))
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
