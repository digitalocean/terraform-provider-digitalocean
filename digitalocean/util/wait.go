package util

import (
	"context"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// WaitForAction waits for the action to finish using the resource.StateChangeConf.
func WaitForAction(client *godo.Client, action *godo.Action) error {
	var (
		pending   = "in-progress"
		target    = "completed"
		refreshfn = func() (result interface{}, state string, err error) {
			a, _, err := client.Actions.Get(context.Background(), action.ID)
			if err != nil {
				return nil, "", err
			}
			if a.Status == "errored" {
				return a, "errored", nil
			}
			if a.CompletedAt != nil {
				return a, target, nil
			}
			return a, pending, nil
		}
	)
	_, err := (&resource.StateChangeConf{
		Pending: []string{pending},
		Refresh: refreshfn,
		Target:  []string{target},

		Delay:      10 * time.Second,
		Timeout:    60 * time.Minute,
		MinTimeout: 3 * time.Second,

		// This is a hack around DO API strangeness.
		// https://github.com/hashicorp/terraform/issues/481
		//
		NotFoundChecks: 60,
	}).WaitForState()
	return err
}

func IsDigitalOceanError(err error, code int, message string) bool {
	if err, ok := err.(*godo.ErrorResponse); ok {
		return err.Response.StatusCode == code &&
			strings.Contains(strings.ToLower(err.Message), strings.ToLower(message))
	}
	return false
}
