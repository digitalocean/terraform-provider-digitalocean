package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanLoadbalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanLoadbalancerCreate,
		ReadContext:   resourceDigitalOceanLoadbalancerRead,
		UpdateContext: resourceDigitalOceanLoadbalancerUpdate,
		DeleteContext: resourceDigitalOceanLoadbalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceDigitalOceanLoadBalancerV0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateLoadBalancerStateV0toV1,
				Version: 0,
			},
		},

		Schema: resourceDigitalOceanLoadBalancerV1(),

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {

			if _, hasHealthCheck := diff.GetOk("healthcheck"); hasHealthCheck {

				healthCheckProtocol := diff.Get("healthcheck.0.protocol").(string)
				_, hasPath := diff.GetOk("healthcheck.0.path")
				if healthCheckProtocol == "http" {
					if !hasPath {
						return fmt.Errorf("health check `path` is required for when protocol is `http`")
					}
				} else if healthCheckProtocol == "https" {
					if !hasPath {
						return fmt.Errorf("health check `path` is required for when protocol is `https`")
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

			if err := loadbalancerDiffCheck(ctx, diff, v); err != nil {
				return err
			}

			return nil
		},
	}
}

func resourceDigitalOceanLoadBalancerV1() map[string]*schema.Schema {
	loadBalancerV0Schema := resourceDigitalOceanLoadBalancerV0().Schema
	loadBalancerV1Schema := map[string]*schema.Schema{}

	forwardingRuleSchema := map[string]*schema.Schema{
		"certificate_name": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.NoZeroValues,
		},
	}

	for k, v := range loadBalancerV0Schema["forwarding_rule"].Elem.(*schema.Resource).Schema {
		forwardingRuleSchema[k] = v
	}
	forwardingRuleSchema["certificate_id"].Computed = true
	forwardingRuleSchema["certificate_id"].Deprecated = "Certificate IDs may change, for example when a Let's Encrypt certificate is auto-renewed. Please specify 'certificate_name' instead."

	for k, v := range loadBalancerV0Schema {
		loadBalancerV1Schema[k] = v
	}
	loadBalancerV1Schema["forwarding_rule"].Elem.(*schema.Resource).Schema = forwardingRuleSchema

	return loadBalancerV1Schema
}

func loadbalancerDiffCheck(ctx context.Context, d *schema.ResourceDiff, v interface{}) error {
	typ, typSet := d.GetOk("type")
	region, regionSet := d.GetOk("region")

	if !typSet && !regionSet {
		return fmt.Errorf("missing 'region' value")
	}

	typStr := typ.(string)
	switch strings.ToUpper(typStr) {
	case "GLOBAL":
		if regionSet && region.(string) != "" {
			return fmt.Errorf("'region' must be empty or not set when 'type' is '%s'", typStr)
		}
	case "REGIONAL":
		if !regionSet || region.(string) == "" {
			return fmt.Errorf("'region' must be set and not be empty when 'type' is '%s'", typStr)
		}
	}

	return nil
}

func resourceDigitalOceanLoadBalancerV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},
			"size": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"lb-small",
					"lb-medium",
					"lb-large",
				}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if sizeUnit, ok := d.GetOk("size_unit"); ok {
						switch {
						case new == "lb-small" && sizeUnit.(int) == 1:
							return true
						case new == "lb-medium" && sizeUnit.(int) == 3:
							return true
						case new == "lb-large" && sizeUnit.(int) == 6:
							return true
						}
						return false
					}
					return old == new
				},
			},
			"size_unit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 100),
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
				Deprecated: "This field has been deprecated. You can no longer specify an algorithm for load balancers.",
			},

			"forwarding_rule": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"glb_settings"},
				MinItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry_protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
								"http2",
								"http3",
								"tcp",
								"udp",
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
								"udp",
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
				Set: hashForwardingRules,
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
								"https",
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
				DiffSuppressFunc: util.CaseSensitive,
				ValidateFunc:     tag.ValidateTag,
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

			"enable_backend_keepalive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"disable_lets_encrypt_dns_records": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"vpc_uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"http_idle_timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"firewall": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "the rules for ALLOWING traffic to the LB (strings in the form: 'ip:1.2.3.4' or 'cidr:1.2.0.0/16')",
						},
						"deny": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "the rules for DENYING traffic to the LB (strings in the form: 'ip:1.2.3.4' or 'cidr:1.2.0.0/16')",
						},
					},
				},
			},

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"REGIONAL", "GLOBAL"}, true),
				Description:  "the type of the load balancer (GLOBAL or REGIONAL)",
			},

			"domains": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				MinItems:    1,
				Description: "the list of domains required to ingress traffic to global load balancer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
							Description:  "domain name",
						},
						"is_managed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "flag indicating if domain is managed by DigitalOcean",
						},
						"certificate_name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.NoZeroValues,
							Description:  "name of certificate required for TLS handshaking",
						},
						"verification_error_reasons": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "list of domain verification errors",
						},
						"ssl_validation_error_reasons": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "list of domain SSL validation errors",
						},
					},
				},
			},

			"glb_settings": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"forwarding_rule"},
				Description:   "configuration options for global load balancer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
							}, false),
							Description: "target protocol rules",
						},
						"target_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{80, 443}),
							Description:  "target port rules",
						},
						"region_priorities": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "region priority map",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"failover_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 99),
							Description:  "fail-over threshold",
						},
						"cdn": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "CDN specific configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "cache enable flag",
									},
								},
							},
						},
					},
				},
			},

			"target_load_balancer_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "list of load balancer IDs to put behind a global load balancer",
			},

			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"EXTERNAL", "INTERNAL"}, true),
				Description:  "the network type of the load balancer (INTERNAL or EXTERNAL)",
			},
		},
	}
}

func migrateLoadBalancerStateV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if len(rawState) == 0 {
		log.Println("[DEBUG] Empty state; nothing to migrate.")
		return rawState, nil
	}
	log.Println("[DEBUG] Migrating load balancer schema from v0 to v1.")
	client := meta.(*config.CombinedConfig).GodoClient()

	// When the certificate type is lets_encrypt, the certificate
	// ID will change when it's renewed, so we have to rely on the
	// certificate name as the primary identifier instead.
	for _, forwardingRule := range rawState["forwarding_rule"].([]interface{}) {
		fw := forwardingRule.(map[string]interface{})
		if fw["certificate_id"].(string) == "" {
			continue
		}

		cert, _, err := client.Certificates.Get(context.Background(), fw["certificate_id"].(string))
		if err != nil {
			return rawState, err
		}

		fw["certificate_id"] = cert.Name
		fw["certificate_name"] = cert.Name
	}

	return rawState, nil
}

func buildLoadBalancerRequest(client *godo.Client, d *schema.ResourceData) (*godo.LoadBalancerRequest, error) {
	forwardingRules, err := expandForwardingRules(client, d.Get("forwarding_rule").(*schema.Set).List())
	if err != nil {
		return nil, err
	}

	opts := &godo.LoadBalancerRequest{
		Name:                         d.Get("name").(string),
		Region:                       d.Get("region").(string),
		Algorithm:                    d.Get("algorithm").(string),
		RedirectHttpToHttps:          d.Get("redirect_http_to_https").(bool),
		EnableProxyProtocol:          d.Get("enable_proxy_protocol").(bool),
		EnableBackendKeepalive:       d.Get("enable_backend_keepalive").(bool),
		ForwardingRules:              forwardingRules,
		DisableLetsEncryptDNSRecords: godo.Bool(d.Get("disable_lets_encrypt_dns_records").(bool)),
		ProjectID:                    d.Get("project_id").(string),
	}
	sizeUnit, ok := d.GetOk("size_unit")
	if ok {
		opts.SizeUnit = uint32(sizeUnit.(int))
	} else {
		opts.SizeSlug = d.Get("size").(string)
	}

	idleTimeout, ok := d.GetOk("http_idle_timeout_seconds")
	if ok {
		t := uint64(idleTimeout.(int))
		opts.HTTPIdleTimeoutSeconds = &t
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

	if v, ok := d.GetOk("firewall"); ok {
		opts.Firewall = expandLBFirewall(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("vpc_uuid"); ok {
		opts.VPCUUID = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		opts.Type = v.(string)
	}

	if v, ok := d.GetOk("domains"); ok {
		domains, err := expandDomains(client, v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}

		opts.Domains = domains
	}

	if v, ok := d.GetOk("glb_settings"); ok {
		opts.GLBSettings = expandGLBSettings(v.([]interface{}))
	}

	if v, ok := d.GetOk("target_load_balancer_ids"); ok {
		var lbIDs []string
		for _, id := range v.(*schema.Set).List() {
			lbIDs = append(lbIDs, id.(string))
		}

		opts.TargetLoadBalancerIDs = lbIDs
	}

	if v, ok := d.GetOk("network_type"); ok {
		opts.Network = v.(string)
	}

	return opts, nil
}

func resourceDigitalOceanLoadbalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Create a Loadbalancer Request")

	lbOpts, err := buildLoadBalancerRequest(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Loadbalancer Create: %#v", lbOpts)
	loadbalancer, _, err := client.LoadBalancers.Create(context.Background(), lbOpts)
	if err != nil {
		return diag.Errorf("Error creating Load Balancer: %s", err)
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
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for Load Balancer (%s) to become active: %s", d.Get("name"), err)
	}

	return resourceDigitalOceanLoadbalancerRead(ctx, d, meta)
}

func resourceDigitalOceanLoadbalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Reading the details of the Loadbalancer %s", d.Id())
	loadbalancer, resp, err := client.LoadBalancers.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Load Balancer (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Loadbalancer: %s", err)
	}

	d.Set("name", loadbalancer.Name)
	d.Set("urn", loadbalancer.URN())
	d.Set("ip", loadbalancer.IP)
	d.Set("status", loadbalancer.Status)
	d.Set("algorithm", loadbalancer.Algorithm)
	d.Set("redirect_http_to_https", loadbalancer.RedirectHttpToHttps)
	d.Set("enable_proxy_protocol", loadbalancer.EnableProxyProtocol)
	d.Set("enable_backend_keepalive", loadbalancer.EnableBackendKeepalive)
	d.Set("droplet_tag", loadbalancer.Tag)
	d.Set("vpc_uuid", loadbalancer.VPCUUID)
	d.Set("http_idle_timeout_seconds", loadbalancer.HTTPIdleTimeoutSeconds)
	d.Set("project_id", loadbalancer.ProjectID)

	if loadbalancer.SizeUnit > 0 {
		d.Set("size_unit", loadbalancer.SizeUnit)
	} else {
		d.Set("size", loadbalancer.SizeSlug)
	}

	if loadbalancer.Region != nil {
		d.Set("region", loadbalancer.Region.Slug)
	}

	d.Set("disable_lets_encrypt_dns_records", loadbalancer.DisableLetsEncryptDNSRecords)

	if err := d.Set("droplet_ids", flattenDropletIds(loadbalancer.DropletIDs)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer droplet_ids - error: %#v", err)
	}

	if err := d.Set("sticky_sessions", flattenStickySessions(loadbalancer.StickySessions)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer sticky_sessions - error: %#v", err)
	}

	if err := d.Set("healthcheck", flattenHealthChecks(loadbalancer.HealthCheck)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer healthcheck - error: %#v", err)
	}

	forwardingRules, err := flattenForwardingRules(client, loadbalancer.ForwardingRules)
	if err != nil {
		return diag.Errorf("[DEBUG] Error building Load Balancer forwarding rules - error: %#v", err)
	}

	if err := d.Set("forwarding_rule", forwardingRules); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer forwarding_rule - error: %#v", err)
	}

	if err := d.Set("firewall", flattenLBFirewall(loadbalancer.Firewall)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer firewall - error: %#v", err)
	}

	return nil
}

func resourceDigitalOceanLoadbalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	lbOpts, err := buildLoadBalancerRequest(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Load Balancer Update: %#v", lbOpts)
	_, _, err = client.LoadBalancers.Update(context.Background(), d.Id(), lbOpts)
	if err != nil {
		return diag.Errorf("Error updating Load Balancer: %s", err)
	}

	return resourceDigitalOceanLoadbalancerRead(ctx, d, meta)
}

func resourceDigitalOceanLoadbalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Deleting Load Balancer: %s", d.Id())
	resp, err := client.LoadBalancers.Delete(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error deleting Load Balancer: %s", err)
	}

	d.SetId("")
	return nil

}
