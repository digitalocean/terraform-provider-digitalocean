package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceDigitalOceanCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanCertificateRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the certificate",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
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

func dataSourceDigitalOceanCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	name := d.Get("name").(string)

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	certList := []godo.Certificate{}

	for {
		certs, resp, err := client.Certificates.List(context.Background(), opts)

		if err != nil {
			return fmt.Errorf("Error retrieving ssh keys: %s", err)
		}

		for _, cert := range certs {
			certList = append(certList, cert)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return fmt.Errorf("Error retrieving ssh keys: %s", err)
		}

		opts.Page = page + 1
	}

	cert, err := findCertificateByName(certList, name)

	if err != nil {
		return err
	}

	d.SetId(cert.ID)
	d.Set("name", cert.Name)
	d.Set("type", cert.Type)
	d.Set("state", cert.State)
	d.Set("not_after", cert.NotAfter)
	d.Set("sha1_fingerprint", cert.SHA1Fingerprint)

	if err := d.Set("domains", flattenDigitalOceanCertificateDomains(cert.DNSNames)); err != nil {
		return fmt.Errorf("Error setting `domain`: %+v", err)
	}

	return nil
}

func findCertificateByName(certs []godo.Certificate, name string) (*godo.Certificate, error) {
	results := make([]godo.Certificate, 0)
	for _, v := range certs {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no certificate found with name %s", name)
	}
	return nil, fmt.Errorf("too many certificate found with name %s (found %d, expected 1)", name, len(results))
}
