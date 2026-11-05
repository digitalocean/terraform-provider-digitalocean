package spaces

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceDigitalOceanBucketCorsConfigurationRead_gone(t *testing.T) {
	cases := map[string]string{
		"cors configuration removed": "NoSuchCORSConfiguration",
		"bucket deleted":             "NoSuchBucket",
	}

	for name, code := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Error><Code>` + code + `</Code><Message>not found</Message></Error>`))
			}))
			defer server.Close()

			// The S3 client addresses the bucket as a subdomain of the endpoint,
			// so send every connection to the test server.
			orig := http.DefaultClient.Transport
			http.DefaultClient.Transport = &http.Transport{
				DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
					return (&net.Dialer{}).DialContext(ctx, network, server.Listener.Addr().String())
				},
			}
			t.Cleanup(func() { http.DefaultClient.Transport = orig })

			c := &config.Config{
				APIEndpoint:       server.URL,
				SpacesAPIEndpoint: server.URL,
				AccessID:          "access",
				SecretKey:         "secret",
			}
			meta, err := c.Client()
			if err != nil {
				t.Fatal(err)
			}

			d := schema.TestResourceDataRaw(t, ResourceDigitalOceanBucketCorsConfiguration().Schema, map[string]interface{}{
				"bucket": "my-bucket",
				"region": "nyc3",
			})
			d.SetId("my-bucket")

			if diags := resourceDigitalOceanBucketCorsConfigurationRead(context.Background(), d, meta); diags.HasError() {
				t.Fatalf("expected no error, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected the resource to be removed from state, id is %q", d.Id())
			}
		})
	}
}
