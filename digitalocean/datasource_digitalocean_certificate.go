package digitalocean

import (
	"context"
	"fmt"
	"net/http"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanCertificateRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the certificate",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes

			// When the certificate type is lets_encrypt, the certificate
			// ID will change when it's renewed, so we have to rely on the
			// certificate name as the primary identifier instead.
			// We include the UUID as another computed field for use in the
			// short-term refresh function that waits for it to be ready.
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "uuid of the certificate",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of the certificate",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "current state of the certificate",
			},
			"domains": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "domains for which the certificate was issued",
			},
			"not_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "expiration date and time of the certificate",
			},
			"sha1_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA1 fingerprint of the certificate",
			},
		},
	}
}

func dataSourceDigitalOceanCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	// When the certificate type is lets_encrypt, the certificate
	// ID will change when it's renewed, so we have to rely on the
	// certificate name as the primary identifier instead.
	name := d.Get("name").(string)
	cert, err := findCertificateByName(client, name)
	if err != nil {
		return err
	}

	d.SetId(cert.Name)
	d.Set("name", cert.Name)
	d.Set("uuid", cert.ID)
	d.Set("type", cert.Type)
	d.Set("state", cert.State)
	d.Set("not_after", cert.NotAfter)
	d.Set("sha1_fingerprint", cert.SHA1Fingerprint)

	if err := d.Set("domains", flattenDigitalOceanCertificateDomains(cert.DNSNames)); err != nil {
		return fmt.Errorf("Error setting `domain`: %+v", err)
	}

	return nil
}

func findCertificateByName(client *godo.Client, name string) (*godo.Certificate, error) {
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		certs, resp, err := client.Certificates.List(context.Background(), opts)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("Error retrieving certificates: %s", err)
		}

		for _, cert := range certs {
			if cert.Name == name {
				return &cert, nil
			}
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving certificates: %s", err)
		}

		opts.Page = page + 1
	}

	return nil, fmt.Errorf("Certificate %s not found", name)
}
