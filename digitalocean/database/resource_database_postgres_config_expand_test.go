package database

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandPgBouncerIgnoreStartupParametersFromSet(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"server_reset_query_always": false,
			"ignore_startup_parameters": schema.NewSet(schema.HashString, []interface{}{"extra_float_digits"}),
			"min_pool_size":             0,
			"server_lifetime":           0,
			"server_idle_timeout":       0,
			"autodb_pool_size":          0,
			"autodb_pool_mode":          "",
			"autodb_max_db_connections": 0,
			"autodb_idle_timeout":       0,
		},
	}

	got := expandPgBouncer(raw)
	if got.IgnoreStartupParameters == nil {
		t.Fatal("expected ignore startup parameters")
	}
	if len(*got.IgnoreStartupParameters) != 1 || (*got.IgnoreStartupParameters)[0] != "extra_float_digits" {
		t.Fatalf("got %#v", got.IgnoreStartupParameters)
	}
}
