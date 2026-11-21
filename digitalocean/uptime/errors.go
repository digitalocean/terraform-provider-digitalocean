package uptime

import (
	"net/http"
	"strings"

	"github.com/digitalocean/godo"
)

// isUptimeResourceGone reports whether an Uptime API error means the check
// or alert no longer exists. Unlike the rest of the DigitalOcean API, the
// Uptime endpoints answer a request for a deleted check with a 403
// "forbidden" / "you are not authorized to access this resource" error
// instead of a 404, so both are treated as "gone". Other 403 errors are
// still reported.
func isUptimeResourceGone(resp *godo.Response, err error) bool {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}

	errResp, ok := err.(*godo.ErrorResponse)
	if !ok || errResp.Response == nil || errResp.Response.StatusCode != http.StatusForbidden {
		return false
	}

	return strings.Contains(strings.ToLower(errResp.Message), "not authorized to access this resource")
}
