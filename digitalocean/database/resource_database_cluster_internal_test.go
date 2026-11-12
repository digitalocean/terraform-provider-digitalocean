package database

import (
	"net/url"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestBuildDBConnectionURIWithOptions(t *testing.T) {
	conn := &godo.DatabaseConnection{
		URI:      "postgresql://doadmin:s3cret@db.example.com:25060/defaultdb?sslmode=require",
		Host:     "db.example.com",
		Port:     25060,
		User:     "doadmin",
		Password: "s3cret",
		Database: "defaultdb",
	}

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanDatabaseCluster().Schema, map[string]interface{}{
		"engine":   "pg",
		"password": "s3cret",
	})

	t.Run("with credentials", func(t *testing.T) {
		got, err := buildDBConnectionURIWithOptions(conn, d, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		parsed, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse uri: %v", err)
		}

		if parsed.User == nil {
			t.Fatalf("expected userinfo in uri, got: %s", got)
		}

		password, ok := parsed.User.Password()
		if !ok || password != "s3cret" {
			t.Fatalf("expected password in uri, got: %s", got)
		}

		if parsed.User.Username() != "doadmin" {
			t.Fatalf("expected username doadmin, got: %s", parsed.User.Username())
		}
	})

	t.Run("without credentials", func(t *testing.T) {
		got, err := buildDBConnectionURIWithOptions(conn, d, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		parsed, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse uri: %v", err)
		}

		if parsed.User != nil {
			t.Fatalf("expected no userinfo in uri, got: %s", got)
		}

		if parsed.Host != "db.example.com:25060" {
			t.Fatalf("expected host to be retained, got: %s", parsed.Host)
		}

		if parsed.Path != "/defaultdb" {
			t.Fatalf("expected path to be retained, got: %s", parsed.Path)
		}

		if parsed.Query().Get("sslmode") != "require" {
			t.Fatalf("expected query to be retained, got: %s", parsed.RawQuery)
		}
	})
}
