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

func spacesBucketPolicyTestMeta(t *testing.T, status int, body string) interface{} {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

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
	return meta
}

func spacesBucketPolicyTestData(t *testing.T) *schema.ResourceData {
	t.Helper()

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanSpacesBucketPolicy().Schema, map[string]interface{}{
		"bucket": "my-bucket",
		"region": "nyc3",
		"policy": `{"Version":"2012-10-17","Statement":[]}`,
	})
	d.SetId("my-bucket")
	return d
}

func s3ErrorBody(code string) string {
	return `<?xml version="1.0" encoding="UTF-8"?><Error><Code>` + code + `</Code><Message>error</Message></Error>`
}

func TestResourceDigitalOceanBucketPolicyRead_gone(t *testing.T) {
	for _, code := range []string{"NoSuchBucketPolicy", "NoSuchBucket"} {
		t.Run(code, func(t *testing.T) {
			meta := spacesBucketPolicyTestMeta(t, http.StatusNotFound, s3ErrorBody(code))
			d := spacesBucketPolicyTestData(t)

			if diags := resourceDigitalOceanBucketPolicyRead(context.Background(), d, meta); diags.HasError() {
				t.Fatalf("expected no error, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected the resource to be removed from state, id is %q", d.Id())
			}
		})
	}
}

func TestResourceDigitalOceanBucketPolicyRead_error(t *testing.T) {
	meta := spacesBucketPolicyTestMeta(t, http.StatusForbidden, s3ErrorBody("AccessDenied"))
	d := spacesBucketPolicyTestData(t)

	if diags := resourceDigitalOceanBucketPolicyRead(context.Background(), d, meta); !diags.HasError() {
		t.Fatal("expected the read error to be returned")
	}
	if d.Id() != "my-bucket" {
		t.Fatalf("expected the resource to stay in state, id is %q", d.Id())
	}
	if got := d.Get("policy").(string); got == "" {
		t.Fatal("expected the policy in state to be kept on a read error")
	}
}

func TestResourceDigitalOceanBucketPolicyDelete_gone(t *testing.T) {
	for _, code := range []string{"NoSuchBucketPolicy", "NoSuchBucket"} {
		t.Run(code, func(t *testing.T) {
			meta := spacesBucketPolicyTestMeta(t, http.StatusNotFound, s3ErrorBody(code))
			d := spacesBucketPolicyTestData(t)

			if diags := resourceDigitalOceanBucketPolicyDelete(context.Background(), d, meta); diags.HasError() {
				t.Fatalf("expected no error, got: %v", diags)
			}
		})
	}
}
