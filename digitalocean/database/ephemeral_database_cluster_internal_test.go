package database

import (
	"context"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

func TestDatabaseClusterEphemeralResourceOpenWithoutClient(t *testing.T) {
	resource := &databaseClusterEphemeralResource{}

	var resp ephemeral.OpenResponse
	resource.Open(context.Background(), ephemeral.OpenRequest{}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when the provider client is nil")
	}
}

func TestPopulateDatabaseClusterEphemeralModel(t *testing.T) {
	database := &godo.Database{
		Name:        "example",
		EngineSlug:  "pg",
		VersionSlug: "15",
		RegionSlug:  "nyc1",
		Connection: &godo.DatabaseConnection{
			Host:     "public.example.com",
			Port:     25060,
			Database: "defaultdb",
			User:     "doadmin",
			Password: "secret",
			URI:      "postgres://doadmin:secret@public.example.com:25060/defaultdb?sslmode=require",
		},
		PrivateConnection: &godo.DatabaseConnection{
			Host:     "private.example.com",
			Port:     25060,
			Database: "defaultdb",
			User:     "doadmin",
			Password: "secret",
			URI:      "postgres://doadmin:secret@private.example.com:25060/defaultdb?sslmode=require",
		},
		UIConnection: &godo.DatabaseConnection{
			Host:     "ui.example.com",
			Port:     8080,
			Database: "defaultdb",
			User:     "doadmin",
			Password: "ui-secret",
			URI:      "postgres://doadmin:ui-secret@ui.example.com:8080/defaultdb",
		},
	}

	var data databaseClusterEphemeralModel
	if err := populateDatabaseClusterEphemeralModel(&data, database); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := data.Password.ValueString(); got != "secret" {
		t.Fatalf("expected password %q, got %q", "secret", got)
	}
	if got := data.UIPassword.ValueString(); got != "ui-secret" {
		t.Fatalf("expected UI password %q, got %q", "ui-secret", got)
	}
	if got := data.URI.ValueString(); !strings.Contains(got, "doadmin:secret@") {
		t.Fatalf("expected public URI to contain credentials, got %q", got)
	}
	if got := data.PrivateURI.ValueString(); !strings.Contains(got, "doadmin:secret@") {
		t.Fatalf("expected private URI to contain credentials, got %q", got)
	}
	if got := data.UIURI.ValueString(); !strings.Contains(got, "doadmin:ui-secret@") {
		t.Fatalf("expected UI URI to contain credentials, got %q", got)
	}
}
