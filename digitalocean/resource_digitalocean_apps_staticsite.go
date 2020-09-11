package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanAppsStaticSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanAppsStaticSiteCreate,
		Read:   resourceDigitalOceanAppsStaticSiteRead,
		Update: resourceDigitalOceanAppsStaticSiteUpdate,
		Delete: resourceDigitalOceanAppsStaticSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"git": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo": {
							Type:     schema.TypeString,
							Required: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Required: true,
						},
						"deploy_on_push": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"github": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo": {
							Type:     schema.TypeString,
							Required: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Required: true,
						},
						"deploy_on_push": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"build_command": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"env": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"scope": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"RUN_TIME",
								"BUILD_TIME",
								"RUN_AND_BUILD_TIME",
							}, false),
						},
						"value": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GENERAL",
								"SECRET",
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceDigitalOceanAppsStaticSiteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	opts, err := firewallRequest(d, client)
	if err != nil {
		return fmt.Errorf("Error in firewall request: %s", err)
	}

	log.Printf("[DEBUG] Firewall create configuration: %#v", opts)

	firewall, _, err := client.Firewalls.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating firewall: %s", err)
	}

	// Assign the firewall id
	d.SetId(firewall.ID)

	log.Printf("[INFO] Firewall ID: %s", d.Id())

	return resourceDigitalOceanAppsStaticSiteRead(d, meta)
}

func resourceDigitalOceanAppsStaticSiteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	// Retrieve the firewall properties for updating the state
	firewall, resp, err := client.Firewalls.Get(context.Background(), d.Id())
	if err != nil {
		// check if the firewall no longer exists.
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Firewall (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving firewall: %s", err)
	}

	d.Set("status", firewall.Status)
	d.Set("create_at", firewall.Created)
	d.Set("pending_changes", firewallPendingChanges(d, firewall))
	d.Set("name", firewall.Name)

	if err := d.Set("droplet_ids", flattenFirewallDropletIds(firewall.DropletIDs)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `droplet_ids`: %+v", err)
	}

	if err := d.Set("inbound_rule", flattenFirewallInboundRules(firewall.InboundRules)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Firewall inbound_rule error: %#v", err)
	}

	if err := d.Set("outbound_rule", flattenFirewallOutboundRules(firewall.OutboundRules)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Firewall outbound_rule error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(firewall.Tags)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `tags`: %+v", err)
	}

	return nil
}

func resourceDigitalOceanAppsStaticSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	opts, err := firewallRequest(d, client)
	if err != nil {
		return fmt.Errorf("Error in firewall request: %s", err)
	}

	log.Printf("[DEBUG] Firewall update configuration: %#v", opts)

	_, _, err = client.Firewalls.Update(context.Background(), d.Id(), opts)
	if err != nil {
		return fmt.Errorf("Error updating firewall: %s", err)
	}

	return resourceDigitalOceanAppsStaticSiteRead(d, meta)
}

func resourceDigitalOceanAppsStaticSiteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting firewall: %s", d.Id())

	// Destroy the droplet
	_, err := client.Firewalls.Delete(context.Background(), d.Id())

	// Handle remotely destroyed droplets
	if err != nil && strings.Contains(err.Error(), "404 Not Found") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error deleting firewall: %s", err)
	}

	return nil
}
