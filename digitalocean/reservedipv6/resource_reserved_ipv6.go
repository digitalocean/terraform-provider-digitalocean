package reservedipv6

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanReservedIPV6() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanReservedIPV6Create,
		UpdateContext: resourceDigitalOceanReservedIPV6Update,
		ReadContext:   resourceDigitalOceanReservedIPV6Read,
		DeleteContext: resourceDigitalOceanReservedIPV6Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanReservedIPV6Import,
		},

		Schema: map[string]*schema.Schema{
			"region_slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the reserved ipv6",
			},
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPv6Address,
			},
			"droplet_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDigitalOceanReservedIPV6Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Creating a reserved IPv6 in a region")
	regionOpts := &godo.ReservedIPV6CreateRequest{
		Region: d.Get("region_slug").(string),
	}

	log.Printf("[DEBUG] Reserved IPv6 create: %#v", regionOpts)
	reservedIP, _, err := client.ReservedIPV6s.Create(context.Background(), regionOpts)
	if err != nil {
		return diag.Errorf("Error creating reserved IPv6: %s", err)
	}

	d.SetId(reservedIP.IP)

	return resourceDigitalOceanReservedIPV6Read(ctx, d, meta)
}

func resourceDigitalOceanReservedIPV6Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChange("droplet_id") {
		if v, ok := d.GetOk("droplet_id"); ok {
			log.Printf("[INFO] Assigning the reserved IPv6 %s to the Droplet %d", d.Id(), v.(int))
			action, _, err := client.ReservedIPV6Actions.Assign(context.Background(), d.Id(), v.(int))
			if err != nil {
				return diag.Errorf(
					"Error assigning reserved IPv6 (%s) to the Droplet: %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPV6Ready(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IP (%s) to be Assigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[INFO] Unassigning the reserved IP %s", d.Id())
			action, _, err := client.ReservedIPV6Actions.Unassign(context.Background(), d.Id())
			if err != nil {
				return diag.Errorf(
					"Error unassigning reserved IP (%s): %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPV6Ready(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IP (%s) to be Unassigned: %s", d.Id(), unassignedErr)
			}
		}
	}

	return resourceDigitalOceanReservedIPV6Read(ctx, d, meta)
}

func resourceDigitalOceanReservedIPV6Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Reading the details of the reserved IPv6 %s", d.Id())
	reservedIP, resp, err := client.ReservedIPV6s.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] Reserved IPv6 (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving reserved IPv6: %s", err)
	}

	if _, ok := d.GetOk("droplet_id"); ok && reservedIP.Droplet != nil {
		d.Set("region_slug", reservedIP.Droplet.Region.Slug)
		d.Set("droplet_id", reservedIP.Droplet.ID)
	} else {
		d.Set("region_slug", reservedIP.RegionSlug)
	}

	d.Set("ip", reservedIP.IP)
	d.Set("urn", reservedIP.URN())

	return nil
}

func resourceDigitalOceanReservedIPV6Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if _, ok := d.GetOk("droplet_id"); ok {
		log.Printf("[INFO] Unassigning the reserved IPv6 from the Droplet")
		action, resp, err := client.ReservedIPV6Actions.Unassign(context.Background(), d.Id())
		if resp.StatusCode != 422 {
			if err != nil {
				return diag.Errorf(
					"Error unassigning reserved IPv6 (%s) from the droplet: %s", d.Id(), err)
			}

			_, unassignedErr := waitForReservedIPV6Ready(ctx, d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return diag.Errorf(
					"Error waiting for reserved IPv6 (%s) to be unassigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[DEBUG] Couldn't unassign reserved IPv6 (%s) from droplet, possibly out of sync: %s", d.Id(), err)
		}
	}

	log.Printf("[INFO] Deleting reserved IPv6: %s", d.Id())
	_, err := client.ReservedIPV6s.Delete(context.Background(), d.Id())
	if err != nil && strings.Contains(err.Error(), "404") {
		return diag.Errorf("Error deleting reserved IPv6: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanReservedIPV6Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*config.CombinedConfig).GodoClient()
	reservedIP, resp, err := client.ReservedIPV6s.Get(context.Background(), d.Id())
	if resp.StatusCode != 404 {
		if err != nil {
			return nil, err
		}

		d.Set("ip", reservedIP.IP)
		d.Set("urn", reservedIP.URN())
		d.Set("region_slug", reservedIP.RegionSlug)

		if reservedIP.Droplet != nil {
			d.Set("droplet_id", reservedIP.Droplet.ID)
		}
	}

	return []*schema.ResourceData{d}, nil
}

func waitForReservedIPV6Ready(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, actionID int) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for reserved IPv6 (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newReservedIPV6StateRefreshFunc(d, meta, actionID),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newReservedIPV6StateRefreshFunc(
	d *schema.ResourceData, meta interface{}, actionID int) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Assigning the reserved IPv6 to the Droplet")
		action, _, err := client.Actions.Get(context.Background(), actionID)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving reserved IPv6 (%s) ActionId (%d): %s", d.Id(), actionID, err)
		}

		log.Printf("[INFO] The reserved IPv6 Action Status is %s", action.Status)
		return &action, action.Status, nil
	}
}
