package reservedipv6

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

func ResourceDigitalOceanReservedIPV6Assignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanReservedIPV6AssignmentCreate,
		ReadContext:   resourceDigitalOceanReservedIPV6AssignmentRead,
		DeleteContext: resourceDigitalOceanReservedIPV6AssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanReservedIPV6AssignmentImport,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
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

func resourceDigitalOceanReservedIPV6AssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	ipAddress := d.Get("ip").(string)
	dropletID := d.Get("droplet_id").(int)

	log.Printf("[INFO] Assigning the reserved IPv6 (%s) to the Droplet %d", ipAddress, dropletID)
	action, _, err := client.ReservedIPV6Actions.Assign(context.Background(), ipAddress, dropletID)
	if err != nil {
		return diag.Errorf(
			"Error Assigning reserved IPv6 (%s) to the droplet: %s", ipAddress, err)
	}

	_, unassignedErr := waitForReservedIPV6AssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
	if unassignedErr != nil {
		return diag.Errorf(
			"Error waiting for reserved IPv6 (%s) to be Assigned: %s", ipAddress, unassignedErr)
	}

	d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%d-%s-", dropletID, ipAddress)))
	return resourceDigitalOceanReservedIPV6AssignmentRead(ctx, d, meta)
}

func resourceDigitalOceanReservedIPV6AssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	ipAddress := d.Get("ip").(string)
	dropletID := d.Get("droplet_id").(int)

	log.Printf("[INFO] Reading the details of the reserved IPv6 %s", ipAddress)
	reservedIPv6, _, err := client.ReservedIPV6s.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IPv6: %s", err)
	}

	if reservedIPv6.Droplet == nil || reservedIPv6.Droplet.ID != dropletID {
		log.Printf("[INFO] A Droplet was detected on the reserved IPv6.")
		d.SetId("")
	}

	return nil
}

func resourceDigitalOceanReservedIPV6AssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	ipAddress := d.Get("ip").(string)
	dropletID := d.Get("droplet_id").(int)

	log.Printf("[INFO] Reading the details of the reserved IPv6 %s", ipAddress)
	reservedIPv6, _, err := client.ReservedIPV6s.Get(context.Background(), ipAddress)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IPv6: %s", err)
	}

	if reservedIPv6.Droplet.ID == dropletID {
		log.Printf("[INFO] Unassigning the reserved IPv6 from the Droplet")
		action, _, err := client.ReservedIPActions.Unassign(context.Background(), ipAddress)
		if err != nil {
			return diag.Errorf("Error unassigning reserved IPv6 (%s) from the droplet: %s", ipAddress, err)
		}

		_, unassignedErr := waitForReservedIPV6AssignmentReady(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
		if unassignedErr != nil {
			return diag.Errorf(
				"Error waiting for reserved IPv6 (%s) to be unassigned: %s", ipAddress, unassignedErr)
		}
	} else {
		log.Printf("[INFO] reserved IPv6 already unassigned, removing from state.")
	}

	d.SetId("")
	return nil
}

func waitForReservedIPV6AssignmentReady(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IPv6 (%s) to have %s of %s",
		d.Get("ip").(string), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPV6AssignmentStateRefreshFunc(d, meta, actionID),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPV6AssignmentStateRefreshFunc(
	d *schema.ResourceData, meta interface{}, actionID int) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Refreshing the reserved IPv6 state")
		action, _, err := client.Actions.Get(context.Background(), actionID)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving reserved IPv6 (%s) ActionId (%d): %s", d.Get("ip_address").(string), actionID, err)
		}

		log.Printf("[INFO] The reserved IPv6 Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}

func resourceDigitalOceanReservedIPV6AssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(id.PrefixedUniqueId(fmt.Sprintf("%s-%s-", s[1], s[0])))
		d.Set("ip", s[0])
		dropletID, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		d.Set("droplet_id", dropletID)
	} else {
		return nil, errors.New("must use the reserved IPv6 and the ID of the Droplet joined with a comma (e.g. `ip,droplet_id`)")
	}

	return []*schema.ResourceData{d}, nil
}
