package reservedipv6_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/reservedipv6"
)

func TestIsReservedIPV6NotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		resp *godo.Response
		err  error
		want bool
	}{
		{
			name: "no error",
			want: false,
		},
		{
			name: "404 response",
			resp: &godo.Response{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			err:  errors.New("not found"),
			want: true,
		},
		{
			name: "500 response",
			resp: &godo.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			err:  errors.New("server error"),
			want: false,
		},
		{
			name: "404 error response",
			err: &godo.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			want: true,
		},
		{
			name: "non-404 error response",
			err: &godo.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := reservedipv6.IsReservedIPV6NotFound(tt.resp, tt.err); got != tt.want {
				t.Fatalf(
					"IsReservedIPV6NotFound() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}
