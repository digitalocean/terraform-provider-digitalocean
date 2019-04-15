package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanLoadbalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanLoadbalancerCreate,
		Read:   resourceDigitalOceanLoadbalancerRead,
		Update: resourceDigitalOceanLoadbalancerUpdate,
		Delete: resourceDigitalOceanLoadbalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the load balancer",
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "round_robin",
				ValidateFunc: validation.StringInSlice([]string{
					"round_robin",
					"least_connections",
				}, false),
			},

			"forwarding_rule": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry_protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
								"http2",
								"tcp",
							}, false),
						},
						"entry_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"target_protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
								"http2",
								"tcp",
							}, false),
						},
						"target_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"certificate_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"tls_passthrough": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"healthcheck": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"tcp",
							}, false),
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"path": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"check_interval_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: validation.IntBetween(3, 300),
						},
						"response_timeout_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(3, 300),
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3,
							ValidateFunc: validation.IntBetween(2, 10),
						},
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(2, 10),
						},
					},
				},
			},

			"sticky_sessions": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true, //this needs to be computed as the API returns a struct with none as the type
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "none",
							ValidateFunc: validation.StringInSlice([]string{
								"cookies",
								"none",
							}, false),
						},
						"cookie_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(2, 40),
						},
						"cookie_ttl_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},

			"droplet_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeInt},
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"droplet_tag"},
			},

			"droplet_tag": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: CaseSensitive,
				ValidateFunc:     validateTag,
			},

			"redirect_http_to_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enable_proxy_protocol": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			if _, hasHealthCheck := diff.GetOk("healthcheck"); hasHealthCheck {

				healthCheckProtocol := diff.Get("healthcheck.0.protocol").(string)
				_, hasPath := diff.GetOk("healthcheck.0.path")
				if healthCheckProtocol == "http" {
					if !hasPath {
						return fmt.Errorf("health check `path` is required for when protocol is `http`")
					}
				} else {
					if hasPath {
						return fmt.Errorf("health check `path` is not allowed for when protocol is `tcp`")
					}
				}
			}

			if _, hasStickySession := diff.GetOk("sticky_sessions.#"); hasStickySession {

				sessionType := diff.Get("sticky_sessions.0.type").(string)
				_, hasCookieName := diff.GetOk("sticky_sessions.0.cookie_name")
				_, hasTtlSeconds := diff.GetOk("sticky_sessions.0.cookie_ttl_seconds")
				if sessionType == "cookies" {
					if !hasCookieName {
						return fmt.Errorf("sticky sessions `cookie_name` is required for when type is `cookie`")
					}
					if !hasTtlSeconds {
						return fmt.Errorf("sticky sessions `cookie_ttl_seconds` is required for when type is `cookie`")
					}
				} else {
					if hasCookieName {
						return fmt.Errorf("sticky sessions `cookie_name` is not allowed for when type is `none`")
					}
					if hasTtlSeconds {
						return fmt.Errorf("sticky sessions `cookie_ttl_seconds` is not allowed for when type is `none`")
					}
				}
			}

			return nil
		},
	}
}

func buildLoadBalancerRequest(d *schema.ResourceData) (*godo.LoadBalancerRequest, error) {
	opts := &godo.LoadBalancerRequest{
		Name:                d.Get("name").(string),
		Region:              d.Get("region").(string),
		Algorithm:           d.Get("algorithm").(string),
		RedirectHttpToHttps: d.Get("redirect_http_to_https").(bool),
		EnableProxyProtocol: d.Get("enable_proxy_protocol").(bool),
		ForwardingRules:     expandForwardingRules(d.Get("forwarding_rule").([]interface{})),
	}

	if v, ok := d.GetOk("droplet_tag"); ok {
		opts.Tag = v.(string)
	} else if v, ok := d.GetOk("droplet_ids"); ok {
		var droplets []int
		for _, id := range v.(*schema.Set).List() {
			droplets = append(droplets, id.(int))
		}

		opts.DropletIDs = droplets
	}

	if v, ok := d.GetOk("healthcheck"); ok {
		opts.HealthCheck = expandHealthCheck(v.([]interface{}))
	}

	if v, ok := d.GetOk("sticky_sessions"); ok {
		opts.StickySessions = expandStickySessions(v.([]interface{}))
	}

	return opts, nil
}

func resourceDigitalOceanLoadbalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Create a Loadbalancer Request")

	lbOpts, err := buildLoadBalancerRequest(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Loadbalancer Create: %#v", lbOpts)
	loadbalancer, _, err := client.LoadBalancers.Create(context.Background(), lbOpts)
	if err != nil {
		return fmt.Errorf("Error creating Load Balancer: %s", err)
	}

	d.SetId(loadbalancer.ID)

	log.Printf("[DEBUG] Waiting for Load Balancer (%s) to become active", d.Get("name"))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"new"},
		Target:     []string{"active"},
		Refresh:    loadbalancerStateRefreshFunc(client, d.Id()),
		Timeout:    10 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Load Balancer (%s) to become active: %s", d.Get("name"), err)
	}

	return resourceDigitalOceanLoadbalancerRead(d, meta)
}

func resourceDigitalOceanLoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Reading the details of the Loadbalancer %s", d.Id())
	loadbalancer, resp, err := client.LoadBalancers.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Load Balancer (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Loadbalancer: %s", err)
	}

	d.Set("name", loadbalancer.Name)
	d.Set("urn", loadbalancer.URN())
	d.Set("ip", loadbalancer.IP)
	d.Set("status", loadbalancer.Status)
	d.Set("algorithm", loadbalancer.Algorithm)
	d.Set("region", loadbalancer.Region.Slug)
	d.Set("redirect_http_to_https", loadbalancer.RedirectHttpToHttps)
	d.Set("enable_proxy_protocol", loadbalancer.EnableProxyProtocol)
	d.Set("droplet_tag", loadbalancer.Tag)

	if err := d.Set("droplet_ids", flattenDropletIds(loadbalancer.DropletIDs)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Load Balancer droplet_ids - error: %#v", err)
	}

	if err := d.Set("sticky_sessions", flattenStickySessions(loadbalancer.StickySessions)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Load Balancer sticky_sessions - error: %#v", err)
	}

	if err := d.Set("healthcheck", flattenHealthChecks(loadbalancer.HealthCheck)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Load Balancer healthcheck - error: %#v", err)
	}

	if err := d.Set("forwarding_rule", flattenForwardingRules(loadbalancer.ForwardingRules)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Load Balancer forwarding_rule - error: %#v", err)
	}

	return nil

}

func resourceDigitalOceanLoadbalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	lbOpts, err := buildLoadBalancerRequest(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Load Balancer Update: %#v", lbOpts)
	_, _, err = client.LoadBalancers.Update(context.Background(), d.Id(), lbOpts)
	if err != nil {
		return fmt.Errorf("Error updating Load Balancer: %s", err)
	}

	return resourceDigitalOceanLoadbalancerRead(d, meta)
}

func resourceDigitalOceanLoadbalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting Load Balancer: %s", d.Id())
	_, err := client.LoadBalancers.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Load Balancer: %s", err)
	}

	d.SetId("")
	return nil

}
