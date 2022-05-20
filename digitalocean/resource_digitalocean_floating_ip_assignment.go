package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanFloatingIPAssignment() *schema.Resource {
	return &schema.Resource{
		// TODO: Uncomment when dates for deprecation timeline are set.
		// DeprecationMessage: "This resource is deprecated and will be removed in a future release. Please use digitalocean_reserved_ip_assignment instead.",
		CreateContext: resourceDigitalOceanReservedIPAssignmentCreate,
		ReadContext:   resourceDigitalOceanReservedIPAssignmentRead,
		DeleteContext: resourceDigitalOceanReservedIPAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanReservedIPAssignmentImport,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"droplet_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
