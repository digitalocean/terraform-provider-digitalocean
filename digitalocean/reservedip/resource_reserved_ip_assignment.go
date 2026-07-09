package reservedip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

	if err := waitOnDroplet(ctx, client, dropletID); err != nil {
		return diag.Errorf("Error waiting for droplet (%d) to be ready for reserved IP assign: %s", dropletID, err)
	}

	log.Printf("[INFO] Assigning the reserved IP (%s) to the Droplet %d", ipAddress, dropletID)
	action, _, err := client.ReservedIPActions.Assign(context.Background(), ipAddress, dropletID)
	if err != nil {
		return diag.Errorf(
			"Error Assigning reserved IP (%s) to the droplet: %s", ipAddress, err)
	}

	_, unassignedErr := waitForReservedIPAssignmentReady(ctx, d, meta, action.ID, reservedIPActionAssign)
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

		_, unassignedErr := waitForReservedIPAssignmentReady(ctx, d, meta, action.ID, reservedIPActionUnassign)
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
	ctx context.Context, d *schema.ResourceData, meta interface{}, actionID int, op reservedIPActionOperation) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IP (%s) action (%d) to complete",
		d.Get("ip_address").(string), actionID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"new", "in-progress"},
		Target:     []string{"completed"},
		Refresh:    newReservedIPAssignmentActionStateRefreshFunc(d, meta, actionID, op),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
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
