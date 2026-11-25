package database

import (
	"context"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

func TestDatabaseUserEphemeralResourceOpenWithoutClient(t *testing.T) {
	resource := &databaseUserEphemeralResource{}

	var resp ephemeral.OpenResponse
	resource.Open(context.Background(), ephemeral.OpenRequest{}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when the provider client is nil")
	}
}

func TestPopulateDatabaseUserEphemeralModel(t *testing.T) {
	user := &godo.DatabaseUser{
		Role:       "normal",
		Password:   "password",
		AccessCert: "certificate",
		AccessKey:  "key",
		MySQLSettings: &godo.DatabaseMySQLUserSettings{
			AuthPlugin: "mysql_native_password",
		},
	}

	var data databaseUserEphemeralModel
	populateDatabaseUserEphemeralModel(&data, user)

	if got := data.Role.ValueString(); got != "normal" {
		t.Fatalf("expected role %q, got %q", "normal", got)
	}
	if got := data.Password.ValueString(); got != "password" {
		t.Fatalf("expected password %q, got %q", "password", got)
	}
	if got := data.AccessCert.ValueString(); got != "certificate" {
		t.Fatalf("expected access certificate %q, got %q", "certificate", got)
	}
	if got := data.AccessKey.ValueString(); got != "key" {
		t.Fatalf("expected access key %q, got %q", "key", got)
	}
	if got := data.MySQLAuthPlugin.ValueString(); got != "mysql_native_password" {
		t.Fatalf("expected MySQL auth plugin %q, got %q", "mysql_native_password", got)
	}
}

func TestPopulateDatabaseUserEphemeralModelMissingCredentials(t *testing.T) {
	user := &godo.DatabaseUser{Role: "normal"}

	var data databaseUserEphemeralModel
	populateDatabaseUserEphemeralModel(&data, user)

	if !data.Password.IsNull() {
		t.Fatalf("expected null password, got %q", data.Password.ValueString())
	}
	if !data.AccessCert.IsNull() {
		t.Fatalf("expected null access certificate, got %q", data.AccessCert.ValueString())
	}
	if !data.AccessKey.IsNull() {
		t.Fatalf("expected null access key, got %q", data.AccessKey.ValueString())
	}
	if !data.MySQLAuthPlugin.IsNull() {
		t.Fatalf("expected null MySQL auth plugin, got %q", data.MySQLAuthPlugin.ValueString())
	}
}
