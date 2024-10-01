package certificate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanCertificateCreate,
		ReadContext:   resourceDigitalOceanCertificateRead,
		DeleteContext: resourceDigitalOceanCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceDigitalOceanCertificateV1(),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceDigitalOceanCertificateV0().CoreConfigSchema().ImpliedType(),
				Upgrade: MigrateCertificateStateV0toV1,
				Version: 0,
			},
		},
	}
}

func resourceDigitalOceanCertificateV1() map[string]*schema.Schema {
	certificateV1Schema := map[string]*schema.Schema{
		// Note that this UUID will change on auto-renewal of a
		// lets_encrypt certificate.
		"uuid": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range resourceDigitalOceanCertificateV0().Schema {
		certificateV1Schema[k] = v
	}

	return certificateV1Schema
}

func resourceDigitalOceanCertificateV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"private_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				StateFunc:    util.HashStringStateFunc(),
				// In order to support older statefiles with fully saved private_key
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new != "" && old == d.Get("private_key")
				},
			},

			"leaf_certificate": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				StateFunc:    util.HashStringStateFunc(),
				// In order to support older statefiles with fully saved leaf_certificate
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new != "" && old == d.Get("leaf_certificate")
				},
			},

			"certificate_chain": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				StateFunc:    util.HashStringStateFunc(),
				// In order to support older statefiles with fully saved certificate_chain
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new != "" && old == d.Get("certificate_chain")
				},
			},

			"domains": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"private_key", "leaf_certificate", "certificate_chain"},
				// The domains attribute is computed for custom certs and should be ignored in diffs.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type") == "custom"
				},
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "custom",
				ValidateFunc: validation.StringInSlice([]string{
					"custom",
					"lets_encrypt",
				}, false),
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"not_after": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sha1_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func MigrateCertificateStateV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if len(rawState) == 0 {
		log.Println("[DEBUG] Empty state; nothing to migrate.")
		return rawState, nil
	}
	log.Println("[DEBUG] Migrating certificate schema from v0 to v1.")

	// When the certificate type is lets_encrypt, the certificate
	// ID will change when it's renewed, so we have to rely on the
	// certificate name as the primary identifier instead.
	rawState["uuid"] = rawState["id"]
	rawState["id"] = rawState["name"]

	return rawState, nil
}

func buildCertificateRequest(d *schema.ResourceData) (*godo.CertificateRequest, error) {
	req := &godo.CertificateRequest{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	if v, ok := d.GetOk("private_key"); ok {
		req.PrivateKey = v.(string)
	}
	if v, ok := d.GetOk("leaf_certificate"); ok {
		req.LeafCertificate = v.(string)
	}
	if v, ok := d.GetOk("certificate_chain"); ok {
		req.CertificateChain = v.(string)
	}

	if v, ok := d.GetOk("domains"); ok {
		req.DNSNames = expandDigitalOceanCertificateDomains(v.(*schema.Set).List())
	}

	return req, nil
}

func resourceDigitalOceanCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	certificateType := d.Get("type").(string)
	if certificateType == "custom" {
		if _, ok := d.GetOk("private_key"); !ok {
			return diag.Errorf("`private_key` is required for when type is `custom` or empty")
		}

		if _, ok := d.GetOk("leaf_certificate"); !ok {
			return diag.Errorf("`leaf_certificate` is required for when type is `custom` or empty")
		}
	} else if certificateType == "lets_encrypt" {

		if _, ok := d.GetOk("domains"); !ok {
			return diag.Errorf("`domains` is required for when type is `lets_encrypt`")
		}
	}

	log.Printf("[INFO] Create a Certificate Request")

	certReq, err := buildCertificateRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Certificate Create: %#v", certReq)
	cert, _, err := client.Certificates.Create(context.Background(), certReq)
	if err != nil {
		return diag.Errorf("Error creating Certificate: %s", err)
	}

	// When the certificate type is lets_encrypt, the certificate
	// ID will change when it's renewed, so we have to rely on the
	// certificate name as the primary identifier instead.
	d.SetId(cert.Name)

	// We include the UUID as another computed field for use in the
	// short-term refresh function that waits for it to be ready.
	err = d.Set("uuid", cert.ID)
	if err != nil {
		return diag.Errorf("Error setting key UUID with value cert ID: %s", cert.ID)
	}

	log.Printf("[INFO] Waiting for certificate (%s) to have state 'verified'", cert.Name)
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"verified"},
		Refresh:    newCertificateStateRefreshFunc(d, meta),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for certificate (%s) to become active: %s", d.Get("name"), err)
	}

	return resourceDigitalOceanCertificateRead(ctx, d, meta)
}

func resourceDigitalOceanCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	// When the certificate type is lets_encrypt, the certificate
	// ID will change when it's renewed, so we have to rely on the
	// certificate name as the primary identifier instead.
	log.Printf("[INFO] Reading the details of the Certificate %s", d.Id())
	cert, err := FindCertificateByName(client, d.Id())
	// check if the certificate no longer exists.
	if cert == nil && strings.Contains(err.Error(), "not found") {
		log.Printf("[WARN] DigitalOcean Certificate (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("Error retrieving Certificate: %s", err)
	}

	d.Set("name", cert.Name)
	d.Set("uuid", cert.ID)
	d.Set("type", cert.Type)
	d.Set("state", cert.State)
	d.Set("not_after", cert.NotAfter)
	d.Set("sha1_fingerprint", cert.SHA1Fingerprint)

	if err := d.Set("domains", flattenDigitalOceanCertificateDomains(cert.DNSNames)); err != nil {
		return diag.Errorf("Error setting `domains`: %+v", err)
	}

	return nil

}

func resourceDigitalOceanCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Deleting Certificate: %s", d.Id())
	cert, err := FindCertificateByName(client, d.Id())
	if err != nil {
		return diag.Errorf("Error retrieving Certificate: %s", err)
	}
	if cert == nil {
		return nil
	}

	timeout := 30 * time.Second
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err = client.Certificates.Delete(context.Background(), cert.ID)
		if err != nil {
			if util.IsDigitalOceanError(err, http.StatusForbidden, "Make sure the certificate is not in use before deleting it") {
				log.Printf("[DEBUG] Received %s, retrying certificate deletion", err.Error())
				time.Sleep(1 * time.Second)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error deleting Certificate: %s", err)
	}

	return nil

}

func expandDigitalOceanCertificateDomains(domains []interface{}) []string {
	expandedDomains := make([]string, len(domains))
	for i, v := range domains {
		expandedDomains[i] = v.(string)
	}

	return expandedDomains
}

func flattenDigitalOceanCertificateDomains(domains []string) *schema.Set {
	if domains == nil {
		return nil
	}

	flattenedDomains := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range domains {
		if v != "" {
			flattenedDomains.Add(v)
		}
	}

	return flattenedDomains
}

func newCertificateStateRefreshFunc(d *schema.ResourceData, meta interface{}) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()
	return func() (interface{}, string, error) {

		// Retrieve the certificate properties
		uuid := d.Get("uuid").(string)
		cert, _, err := client.Certificates.Get(context.Background(), uuid)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving certificate: %s", err)
		}

		return cert, cert.State, nil
	}
}
