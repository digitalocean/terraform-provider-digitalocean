package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanCDN() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanCDNCreate,
		Read:   resourceDigitalOceanCDNRead,
		Update: resourceDigitalOceanCDNUpdate,
		Delete: resourceDigitalOceanCDNDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"origin": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "fully qualified domain name (FQDN) for the origin server",
				ValidateFunc: validation.NoZeroValues,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The amount of time the content is cached in the CDN",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of a DigitalOcean managed TLS certificate for use with custom domains",
			},
			"custom_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "fully qualified domain name (FQDN) for custom subdomain, (requires certificate_id)",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "fully qualified domain name (FQDN) to serve the CDN content",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time (ISO8601) of when the CDN endpoint was created.",
			},
		},
	}
}

func resourceDigitalOceanCDNCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	cdnRequest := &godo.CDNCreateRequest{
		Origin: d.Get("origin").(string),
	}

	if v, ok := d.GetOk("ttl"); ok {
		cdnRequest.TTL = uint32(v.(int))
	}

	if v, ok := d.GetOk("custom_domain"); ok {
		cdnRequest.CustomDomain = v.(string)
	}

	if v, ok := d.GetOk("certificate_id"); ok {
		cdnRequest.CertificateID = v.(string)
	}

	log.Printf("[DEBUG] CDN create request: %#v", cdnRequest)
	cdn, _, err := client.CDNs.Create(context.Background(), cdnRequest)
	if err != nil {
		return fmt.Errorf("Error creating CDN: %s", err)
	}

	d.SetId(cdn.ID)
	log.Printf("[INFO] CDN created, ID: %s", d.Id())

	return resourceDigitalOceanCDNRead(d, meta)
}

func resourceDigitalOceanCDNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	cdn, resp, err := client.CDNs.Get(context.Background(), d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] CDN  (%s) was not found - removing from state", d.Id())
			d.SetId("")
		}
		return fmt.Errorf("Error reading CDN: %s", err)
	}

	d.SetId(cdn.ID)
	d.Set("origin", cdn.Origin)
	d.Set("ttl", cdn.TTL)
	d.Set("endpoint", cdn.Endpoint)
	d.Set("created_at", cdn.CreatedAt)
	d.Set("custom_domain", cdn.CustomDomain)
	d.Set("certificate_id", cdn.CertificateID)
	return err
}

func resourceDigitalOceanCDNUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	d.Partial(true)

	if d.HasChange("ttl") {
		ttlUpdateRequest := &godo.CDNUpdateTTLRequest{
			TTL: uint32(d.Get("ttl").(int)),
		}
		_, _, err := client.CDNs.UpdateTTL(context.Background(), d.Id(), ttlUpdateRequest)

		if err != nil {
			return fmt.Errorf("Error updating CDN TTL: %s", err)
		}
		log.Printf("[INFO] Updated TTL on CDN")
		d.SetPartial("ttl")
	}

	if d.HasChange("certificate_id") || d.HasChange("custom_domain") {
		cdnUpdateRequest := &godo.CDNUpdateCustomDomainRequest{
			CustomDomain:  d.Get("custom_domain").(string),
			CertificateID: d.Get("certificate_id").(string),
		}

		_, _, err := client.CDNs.UpdateCustomDomain(context.Background(), d.Id(), cdnUpdateRequest)

		if err != nil {
			return fmt.Errorf("Error updating CDN custom domain: %s", err)
		}
		log.Printf("[INFO] Updated custom domain/certificate on CDN")
		d.SetPartial("custom_domain_and_certificate_id")
	}

	d.Partial(false)
	return resourceDigitalOceanCDNRead(d, meta)
}

func resourceDigitalOceanCDNDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	resourceId := d.Id()

	_, err := client.CDNs.Delete(context.Background(), resourceId)
	if err != nil {
		return fmt.Errorf("Error deleting CDN: %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] CDN deleted, ID: %s", resourceId)

	return nil
}
