package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDomainCreate,
		Read:   resourceDigitalOceanDomainRead,
		Delete: resourceDigitalOceanDomainDelete,
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
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	// Build up our creation options

	opts := &godo.DomainCreateRequest{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("ip_address"); ok {
		opts.IPAddress = v.(string)
	}

	log.Printf("[DEBUG] Domain create configuration: %#v", opts)
	domain, _, err := client.Domains.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating Domain: %s", err)
	}

	d.SetId(domain.Name)
	log.Printf("[INFO] Domain Name: %s", domain.Name)

	return resourceDigitalOceanDomainRead(d, meta)
}

func resourceDigitalOceanDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	domain, resp, err := client.Domains.Get(context.Background(), d.Id())
	if err != nil {
		// If the domain is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving domain: %s", err)
	}

	d.Set("name", domain.Name)
	d.Set("urn", domain.URN())

	return nil
}

func resourceDigitalOceanDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting Domain: %s", d.Id())
	_, err := client.Domains.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Domain: %s", err)
	}

	d.SetId("")
	return nil
}
