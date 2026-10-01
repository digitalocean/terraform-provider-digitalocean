package reservedip

import (
	"errors"
	"net/http"
	"testing"

	"github.com/digitalocean/godo"
)

func TestIsReservedIPNotFound(t *testing.T) {
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

			if got := IsReservedIPNotFound(tt.resp, tt.err); got != tt.want {
				t.Fatalf(
					"isReservedIPNotFound() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}
