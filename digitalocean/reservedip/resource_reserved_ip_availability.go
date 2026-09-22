package reservedip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func IsReservedIPNotFound(resp *godo.Response, err error) bool {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}

	var errResp *godo.ErrorResponse
	if errors.As(err, &errResp) &&
		errResp.Response != nil &&
		errResp.Response.StatusCode == http.StatusNotFound {
		return true
	}

	return false
}

func newReservedIPAvailableStateRefreshFunc(client *godo.Client, ipAddress string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		reservedIP, resp, err := client.ReservedIPs.Get(context.Background(), ipAddress)

		if IsReservedIPNotFound(resp, err) {
			log.Printf("[DEBUG] Reserved IP (%s) not yet available", ipAddress)
			return nil, "not-found", nil
		}

		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving reserved IP (%s): %s", ipAddress, err)
		}

		return reservedIP, "available", nil
	}
}

// waitForReservedIPAvailability polls the reserved IP until it is visible via
// Get, and returns the resolved object. The Reserved IP API has been observed
// to flap between 404 and 200 multiple times within a few seconds of a
// create or assign/unassign action -- a single successful check is not a
// reliable signal that a *subsequent, separate* Get will also succeed. To
// avoid reopening that race, callers must use the *godo.ReservedIP returned
// here directly instead of issuing another Get once this returns
// successfully.
func waitForReservedIPAvailability(ctx context.Context, d *schema.ResourceData, meta interface{}) (*godo.ReservedIP, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"not-found"},
		Target:     []string{"available"},
		Refresh:    newReservedIPAvailableStateRefreshFunc(client, d.Id()),
		Delay:      0,
		MinTimeout: 3 * time.Second,

		// Observed against a real account: propagation usually clears in a
		// few seconds, but occasionally exceeds the previous 30s budget.
		// Widened to 90s / 30 checks (at the 3s MinTimeout poll interval)
		// to absorb that tail without masking a genuinely deleted resource
		// for an unreasonable amount of time.
		Timeout:        90 * time.Second,
		NotFoundChecks: 30,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	reservedIP, ok := result.(*godo.ReservedIP)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T for reserved IP (%s) availability result", result, d.Id())
	}

	return reservedIP, nil
}
