package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanVPCCreate,
		Read:   resourceDigitalOceanVPCRead,
		Update: resourceDigitalOceanVPCUpdate,
		Delete: resourceDigitalOceanVPCDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
	}
}

func resourceDigitalOceanVPCCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	vpcRequest := &godo.VPCCreateRequest{
		Name:       d.Get("name").(string),
		RegionSlug: d.Get("region").(string),
	}

	// if v, ok := d.GetOk("description"); ok {
	// 	vpcRequest.Description = v.(string)
	// }

	// if v, ok := d.GetOk("ip_range"); ok {
	// 	vpcRequest.IPRange = v.(string)
	// }

	log.Printf("[DEBUG] VPC create request: %#v", vpcRequest)
	vpc, _, err := client.VPCs.Create(context.Background(), vpcRequest)
	if err != nil {
		return fmt.Errorf("Error creating VPC: %s", err)
	}

	d.SetId(vpc.ID)
	log.Printf("[INFO] VPC created, ID: %s", d.Id())

	return resourceDigitalOceanVPCRead(d, meta)
}

func resourceDigitalOceanVPCRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	vpc, resp, err := client.VPCs.Get(context.Background(), d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] VPC  (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading VPC: %s", err)
	}

	d.SetId(vpc.ID)
	d.Set("name", vpc.Name)
	d.Set("region", vpc.RegionSlug)
	//d.Set("description", vpc.Description)
	//d.Set("ip_range", vpc.IPRange)
	//d.Set("urn", vpc.URN)
	d.Set("default", vpc.Default)
	d.Set("created_at", vpc.CreatedAt.UTC().String())
	return err
}

func resourceDigitalOceanVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	if d.HasChange("name") { // || d.HasChange("description")
		vpcUpdateRequest := &godo.VPCUpdateRequest{
			Name: d.Get("name").(string),
			//Description: d.Get("description").(string),
		}
		_, _, err := client.VPCs.Update(context.Background(), d.Id(), vpcUpdateRequest)

		if err != nil {
			return fmt.Errorf("Error updating VPC : %s", err)
		}
		log.Printf("[INFO] Updated VPC")
	}

	return resourceDigitalOceanVPCRead(d, meta)
}

func resourceDigitalOceanVPCDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	vpcID := d.Id()

	_, err := client.VPCs.Delete(context.Background(), vpcID)
	if err != nil {
		return fmt.Errorf("Error deleting VPC: %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] VPC deleted, ID: %s", vpcID)

	return nil
}
