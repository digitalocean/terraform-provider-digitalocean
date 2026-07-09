package reservedip

import (
	"errors"
	"net/http"
	"testing"

	"github.com/digitalocean/godo"
)

func TestIsReservedIPActionNotFound(t *testing.T) {
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
			resp: &godo.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
			err:  errors.New("not found"),
			want: true,
		},
		{
			name: "500 response",
			resp: &godo.Response{Response: &http.Response{StatusCode: http.StatusInternalServerError}},
			err:  errors.New("server error"),
			want: false,
		},
		{
			name: "404 error response",
			err: &godo.ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusNotFound},
			},
			want: true,
		},
		{
			name: "non-404 error response",
			err: &godo.ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusForbidden},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isReservedIPActionNotFound(tt.resp, tt.err); got != tt.want {
				t.Fatalf("isReservedIPActionNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReservedIPActionComplete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		reservedIP *godo.ReservedIP
		op         reservedIPActionOperation
		dropletID  int
		want       bool
	}{
		{
			name:       "assign complete",
			reservedIP: &godo.ReservedIP{Droplet: &godo.Droplet{ID: 123}},
			op:         reservedIPActionAssign,
			dropletID:  123,
			want:       true,
		},
		{
			name:       "assign pending",
			reservedIP: &godo.ReservedIP{Droplet: nil},
			op:         reservedIPActionAssign,
			dropletID:  123,
			want:       false,
		},
		{
			name:       "assign wrong droplet",
			reservedIP: &godo.ReservedIP{Droplet: &godo.Droplet{ID: 456}},
			op:         reservedIPActionAssign,
			dropletID:  123,
			want:       false,
		},
		{
			name:       "unassign complete",
			reservedIP: &godo.ReservedIP{Droplet: nil},
			op:         reservedIPActionUnassign,
			want:       true,
		},
		{
			name:       "unassign pending",
			reservedIP: &godo.ReservedIP{Droplet: &godo.Droplet{ID: 123}},
			op:         reservedIPActionUnassign,
			want:       false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := reservedIPActionComplete(tt.reservedIP, tt.op, tt.dropletID); got != tt.want {
				t.Fatalf("reservedIPActionComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDropletActionPending(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status string
		want   bool
	}{
		{status: godo.ActionCompleted, want: false},
		{status: "errored", want: false},
		{status: godo.ActionInProgress, want: true},
		{status: "new", want: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.status, func(t *testing.T) {
			t.Parallel()

			if got := isDropletActionPending(tt.status); got != tt.want {
				t.Fatalf("isDropletActionPending(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestDropletHasPendingEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		actions []godo.Action
		want    bool
	}{
		{
			name: "no pending events",
			actions: []godo.Action{
				{Type: "create", Status: godo.ActionCompleted},
				{Type: "droplet_metadata_create", Status: godo.ActionCompleted},
			},
			want: false,
		},
		{
			name: "metadata create pending",
			actions: []godo.Action{
				{Type: "droplet_metadata_create", Status: godo.ActionInProgress},
			},
			want: true,
		},
		{
			name: "any other pending event",
			actions: []godo.Action{
				{Type: "resize", Status: godo.ActionInProgress},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := dropletHasPendingEvent(tt.actions); got != tt.want {
				t.Fatalf("dropletHasPendingEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
