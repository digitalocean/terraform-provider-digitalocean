package util

import (
	"strings"

	"github.com/digitalocean/godo"
)

// IsDigitalOceanError detects if a given error is a *godo.ErrorResponse for
// the specified code and message.
func IsDigitalOceanError(err error, code int, message string) bool {
	if err, ok := err.(*godo.ErrorResponse); ok {
		return err.Response.StatusCode == code &&
			strings.Contains(strings.ToLower(err.Message), strings.ToLower(message))
	}
	return false
}
