package digitalocean

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanCertificateCreate,
		Read:   resourceDigitalOceanCertificateRead,
		Delete: resourceDigitalOceanCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			},

			"leaf_certificate": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"certificate_chain": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
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

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			certificateType := diff.Get("type").(string)
			if certificateType == "custom" {
				if _, ok := diff.GetOk("private_key"); !ok {
					return fmt.Errorf("`private_key` is required for when type is `custom` or empty")
				}

				if _, ok := diff.GetOk("leaf_certificate"); !ok {
					return fmt.Errorf("`leaf_certificate` is required for when type is `custom` or empty")
				}
			} else if certificateType == "lets_encrypt" {

				if _, ok := diff.GetOk("domains"); !ok {
					return fmt.Errorf("`domains` is required for when type is `lets_encrypt`")
				}
			}

			return nil
		},
	}
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

func resourceDigitalOceanCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Create a Certificate Request")

	certReq, err := buildCertificateRequest(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Certificate Create: %#v", certReq)
	cert, _, err := client.Certificates.Create(context.Background(), certReq)
	if err != nil {
		return fmt.Errorf("Error creating Certificate: %s", err)
	}

	d.SetId(cert.ID)

	log.Printf("[INFO] Waiting for certificate (%s) to have state 'verified'", cert.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"verified"},
		Refresh:    newCertificateStateRefreshFunc(d, meta),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for certificate (%s) to become active: %s", d.Get("name"), err)
	}

	return resourceDigitalOceanCertificateRead(d, meta)
}

func resourceDigitalOceanCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Reading the details of the Certificate %s", d.Id())
	cert, resp, err := client.Certificates.Get(context.Background(), d.Id())
	if err != nil {
		// check if the certificate no longer exists.
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Certificate (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Certificate: %s", err)
	}

	d.Set("name", cert.Name)
	d.Set("type", cert.Type)
	d.Set("state", cert.State)
	d.Set("not_after", cert.NotAfter)
	d.Set("sha1_fingerprint", cert.SHA1Fingerprint)

	if err := d.Set("domains", flattenDigitalOceanCertificateDomains(cert.DNSNames)); err != nil {
		return fmt.Errorf("Error setting `domains`: %+v", err)
	}

	return nil

}

func resourceDigitalOceanCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting Certificate: %s", d.Id())
	_, err := client.Certificates.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Certificate: %s", err)
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

func newCertificateStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).godoClient()
	return func() (interface{}, string, error) {

		// Retrieve the certificate properties
		cert, _, err := client.Certificates.Get(context.Background(), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving certifica: %s", err)
		}

		return cert, cert.State, nil
	}
}
