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

func newReservedIPActionStateRefreshFunc(client *godo.Client, ipAddress string, actionID int) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action, resp, err := client.ReservedIPActions.Get(context.Background(), ipAddress, actionID)
		if isReservedIPActionNotFound(resp, err) {
			// The assign/unassign call returns an action ID immediately, but the
			// action may not be queryable yet (notably for BYOIP-prefix reserved
			// IPs). Treat 404 as a transient condition and keep polling.
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

func newReservedIPAssignmentActionStateRefreshFunc(d *schema.ResourceData, meta interface{}, actionID int) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()
	ipAddress := d.Get("ip_address").(string)

	return func() (interface{}, string, error) {
		log.Printf("[INFO] Refreshing the reserved IP state")
		return newReservedIPActionStateRefreshFunc(client, ipAddress, actionID)()
	}
}

func newReservedIPResourceActionStateRefreshFunc(d *schema.ResourceData, meta interface{}, actionID int) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()

	return func() (interface{}, string, error) {
		log.Printf("[INFO] Assigning the reserved IP to the Droplet")
		return newReservedIPActionStateRefreshFunc(client, d.Id(), actionID)()
	}
}
