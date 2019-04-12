package digitalocean

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanProjectCreate,
		Read:   resourceDigitalOceanProjectRead,
		Update: resourceDigitalOceanProjectUpdate,
		Delete: resourceDigitalOceanProjectDelete,
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
			"purpose": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceDigitalOceanProjectCreate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*CombinedConfig).godoClient()

	return resourceDigitalOceanCertificateRead(d, meta)
}

func resourceDigitalOceanProjectRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*CombinedConfig).godoClient()

	return nil
}

func resourceDigitalOceanProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*CombinedConfig).godoClient()

	return resourceDigitalOceanCertificateRead(d, meta)
}

func resourceDigitalOceanProjectDelete(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*CombinedConfig).godoClient()

	return nil
}
