package digitalocean

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The amount of time the content is cached in the CDN",
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

	return nil
}

func resourceDigitalOceanCDNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	return nil
}

func resourceDigitalOceanCDNUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	return nil
}

func resourceDigitalOceanCDNDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	return nil
}
