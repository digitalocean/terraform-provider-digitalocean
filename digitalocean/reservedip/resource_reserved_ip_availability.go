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

func waitForReservedIPAvailability(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.CombinedConfig).GodoClient()

	stateConf := &retry.StateChangeConf{
		Pending:        []string{"not-found"},
		Target:         []string{"available"},
		Refresh:        newReservedIPAvailableStateRefreshFunc(client, d.Id()),
		Timeout:        30 * time.Second,
		Delay:          0,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 10,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
