package certificate

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/digitalocean/godo"
)

// godoErr builds a *godo.ErrorResponse suitable for unit testing the error
// classification helpers without making any network calls.
func godoErr(status int, message string) error {
	return &godo.ErrorResponse{
		Response: &http.Response{StatusCode: status},
		Message:  message,
	}
}

func TestCertificateDeleteBlockedByLoadBalancer(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		blocked bool
	}{
		{
			name:    "nil error is not blocked",
			err:     nil,
			blocked: false,
		},
		{
			name:    "plain Go error is not blocked",
			err:     errors.New("network unreachable"),
			blocked: false,
		},
		{
			name:    "403 with primary in-use message is blocked",
			err:     godoErr(http.StatusForbidden, "Make sure the certificate is not in use before deleting it"),
			blocked: true,
		},
		{
			name:    "403 with alternate in-use message is blocked",
			err:     godoErr(http.StatusForbidden, "Certificate is being used by one or more load balancers"),
			blocked: true,
		},
		{
			name:    "match is case-insensitive",
			err:     godoErr(http.StatusForbidden, "MAKE SURE THE CERTIFICATE IS NOT IN USE BEFORE DELETING IT"),
			blocked: true,
		},
		{
			name:    "403 with unrelated message is not blocked",
			err:     godoErr(http.StatusForbidden, "rate limit exceeded"),
			blocked: false,
		},
		{
			name:    "404 with matching message is not blocked",
			err:     godoErr(http.StatusNotFound, "Make sure the certificate is not in use before deleting it"),
			blocked: false,
		},
		{
			name:    "500 is not blocked even with matching message",
			err:     godoErr(http.StatusInternalServerError, "being used by one or more load balancers"),
			blocked: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := certificateDeleteBlockedByLoadBalancer(tc.err); got != tc.blocked {
				t.Fatalf("certificateDeleteBlockedByLoadBalancer(%v) = %v, want %v", tc.err, got, tc.blocked)
			}
		})
	}
}

// TestResourceDigitalOceanCertificateTimeouts asserts that the resource ships
// with the explicit create/delete timeouts the GLB rotation fix relies on.
// If a future change drops or shortens these defaults the rotation runbook
// stops working, so we guard them here.
func TestResourceDigitalOceanCertificateTimeouts(t *testing.T) {
	r := ResourceDigitalOceanCertificate()
	if r.Timeouts == nil {
		t.Fatal("expected ResourceTimeout to be configured on digitalocean_certificate")
	}

	if r.Timeouts.Create == nil || *r.Timeouts.Create != 10*time.Minute {
		t.Fatalf("expected default Create timeout of 10m, got %v", r.Timeouts.Create)
	}
	if r.Timeouts.Delete == nil || *r.Timeouts.Delete != certificateDeleteDefaultTimeout {
		t.Fatalf("expected default Delete timeout of %v, got %v", certificateDeleteDefaultTimeout, r.Timeouts.Delete)
	}

	if certificateDeleteInitialBackoff >= certificateDeleteMaxBackoff {
		t.Fatalf("initial backoff %v must be smaller than max backoff %v", certificateDeleteInitialBackoff, certificateDeleteMaxBackoff)
	}

	if _, ok := r.Schema["name"]; !ok {
		t.Fatal("expected 'name' attribute on certificate schema")
	}
}
