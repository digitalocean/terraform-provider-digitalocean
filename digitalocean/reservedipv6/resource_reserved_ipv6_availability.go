package reservedipv6

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

func IsReservedIPV6NotFound(resp *godo.Response, err error) bool {
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

func newReservedIPV6AvailableStateRefreshFunc(client *godo.Client, ipAddress string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		reservedIP, resp, err := client.ReservedIPV6s.Get(context.Background(), ipAddress)

		if IsReservedIPV6NotFound(resp, err) {
			log.Printf("[DEBUG] Reserved IPv6 (%s) not yet available", ipAddress)
			return nil, "not-found", nil
		}

		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving reserved IPv6 (%s): %s", ipAddress, err)
		}

		return reservedIP, "available", nil
	}
}

// waitForReservedIPV6Availability polls the reserved IPv6 until it is visible via
// Get, and returns the resolved object. Similar to Reserved IPv4, the API may
// experience eventual-consistency propagation delay after creation -- a single
// Get immediately after Create can race with API indexing, return 404, and clear
// the ID. Callers use the resolved *godo.ReservedIPV6 directly.
func waitForReservedIPV6Availability(ctx context.Context, d *schema.ResourceData, meta interface{}) (*godo.ReservedIPV6, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	stateConf := &retry.StateChangeConf{
		Pending:        []string{"not-found"},
		Target:         []string{"available"},
		Refresh:        newReservedIPV6AvailableStateRefreshFunc(client, d.Id()),
		Delay:          0,
		MinTimeout:     3 * time.Second,
		Timeout:        90 * time.Second,
		NotFoundChecks: 30,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	reservedIP, ok := result.(*godo.ReservedIPV6)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T for reserved IPv6 (%s) availability result", result, d.Id())
	}

	return reservedIP, nil
}
