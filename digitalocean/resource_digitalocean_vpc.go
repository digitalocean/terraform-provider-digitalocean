package digitalocean

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanVPC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanVPCCreate,
		ReadContext:   resourceDigitalOceanVPCRead,
		UpdateContext: resourceDigitalOceanVPCUpdate,
		DeleteContext: resourceDigitalOceanVPCDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the VPC",
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "DigitalOcean region slug for the VPC's location",
				ValidateFunc: validation.NoZeroValues,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A free-form description for the VPC",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"ip_range": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The range of IP addresses for the VPC in CIDR notation",
				ValidateFunc: validation.IsCIDR,
			},

			// Computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uniform resource name (URN) for the VPC",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not the VPC is the default one for the region",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the VPC was created",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
	}
}

func resourceDigitalOceanVPCCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	region := d.Get("region").(string)
	vpcRequest := &godo.VPCCreateRequest{
		Name:       d.Get("name").(string),
		RegionSlug: region,
	}

	if v, ok := d.GetOk("description"); ok {
		vpcRequest.Description = v.(string)
	}

	if v, ok := d.GetOk("ip_range"); ok {
		vpcRequest.IPRange = v.(string)
	}

	// Prevent parallel creation of VPCs in the same region to protect
	// against race conditions in IP range assignment.
	key := fmt.Sprintf("resource_digitalocean_vpc/%s", region)
	mutexKV.Lock(key)
	defer mutexKV.Unlock(key)

	log.Printf("[DEBUG] VPC create request: %#v", vpcRequest)
	vpc, _, err := client.VPCs.Create(context.Background(), vpcRequest)
	if err != nil {
		return diag.Errorf("Error creating VPC: %s", err)
	}

	d.SetId(vpc.ID)
	log.Printf("[INFO] VPC created, ID: %s", d.Id())

	return resourceDigitalOceanVPCRead(ctx, d, meta)
}

func resourceDigitalOceanVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	vpc, resp, err := client.VPCs.Get(context.Background(), d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] VPC  (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading VPC: %s", err)
	}

	d.SetId(vpc.ID)
	d.Set("name", vpc.Name)
	d.Set("region", vpc.RegionSlug)
	d.Set("description", vpc.Description)
	d.Set("ip_range", vpc.IPRange)
	d.Set("urn", vpc.URN)
	d.Set("default", vpc.Default)
	d.Set("created_at", vpc.CreatedAt.UTC().String())

	return nil
}

func resourceDigitalOceanVPCUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	if d.HasChanges("name", "description") {
		vpcUpdateRequest := &godo.VPCUpdateRequest{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Default:     godo.Bool(d.Get("default").(bool)),
		}
		_, _, err := client.VPCs.Update(context.Background(), d.Id(), vpcUpdateRequest)

		if err != nil {
			return diag.Errorf("Error updating VPC : %s", err)
		}
		log.Printf("[INFO] Updated VPC")
	}

	return resourceDigitalOceanVPCRead(ctx, d, meta)
}

func resourceDigitalOceanVPCDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()
	vpcID := d.Id()

	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		resp, err := client.VPCs.Delete(context.Background(), vpcID)
		if err != nil {
			// Retry if VPC still contains member resources to prevent race condition
			// with database cluster deletion.
			if resp.StatusCode == 403 {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(fmt.Errorf("Error deleting VPC: %s", err))
			}
		}

		d.SetId("")
		log.Printf("[INFO] VPC deleted, ID: %s", vpcID)

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	} else {
		return nil
	}
}
