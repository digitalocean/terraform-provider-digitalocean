package reservedip

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type reservedIPActionOperation int

const (
	reservedIPActionAssign reservedIPActionOperation = iota
	reservedIPActionUnassign
)

func isReservedIPActionNotFound(resp *godo.Response, err error) bool {
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

func reservedIPActionComplete(reservedIP *godo.ReservedIP, op reservedIPActionOperation, dropletID int) bool {
	switch op {
	case reservedIPActionAssign:
		return reservedIP.Droplet != nil && reservedIP.Droplet.ID == dropletID
	case reservedIPActionUnassign:
		return reservedIP.Droplet == nil
	default:
		return false
	}
}

func newReservedIPActionStateRefreshFunc(client *godo.Client, ipAddress string, actionID int, op reservedIPActionOperation, dropletID int) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action, resp, err := client.ReservedIPActions.Get(context.Background(), ipAddress, actionID)
		if isReservedIPActionNotFound(resp, err) {
			// The assign/unassign call returns an action ID immediately, but the
			// action may not be queryable yet (notably for BYOIP-prefix reserved
			// IPs). Some BYOIP actions never appear on the actions endpoint even
			// after the assignment completes, so fall back to reading the IP.
			reservedIP, _, ipErr := client.ReservedIPs.Get(context.Background(), ipAddress)
			if ipErr == nil && reservedIPActionComplete(reservedIP, op, dropletID) {
				log.Printf("[INFO] Reserved IP (%s) action (%d) not found, but assignment state matches expected outcome", ipAddress, actionID)
				return reservedIP, "completed", nil
			}

			log.Printf("[DEBUG] Reserved IP (%s) action (%d) not found (404), will retry", ipAddress, actionID)
			return nil, "", nil
		}
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving reserved IP (%s) ActionId (%d): %s", ipAddress, actionID, err)
		}

		log.Printf("[INFO] The reserved IP Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}

func newReservedIPAssignmentActionStateRefreshFunc(d *schema.ResourceData, meta interface{}, actionID int, op reservedIPActionOperation) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()
	ipAddress := d.Get("ip_address").(string)
	dropletID := d.Get("droplet_id").(int)

	return func() (interface{}, string, error) {
		log.Printf("[INFO] Refreshing the reserved IP state")
		return newReservedIPActionStateRefreshFunc(client, ipAddress, actionID, op, dropletID)()
	}
}

func newReservedIPResourceActionStateRefreshFunc(d *schema.ResourceData, meta interface{}, actionID int, op reservedIPActionOperation) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()
	dropletID := 0
	if op == reservedIPActionAssign {
		if v, ok := d.GetOk("droplet_id"); ok {
			dropletID = v.(int)
		}
	}

	return func() (interface{}, string, error) {
		log.Printf("[INFO] Assigning the reserved IP to the Droplet")
		return newReservedIPActionStateRefreshFunc(client, d.Id(), actionID, op, dropletID)()
	}
}
