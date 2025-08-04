package loadbalancer

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanLoadbalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanLoadbalancerRead,
		Schema: map[string]*schema.Schema{

			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "id of the load balancer",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "name of the load balancer",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			// computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the load balancer",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the region that the load balancer is deployed in",
			},
			"size_unit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the size of the load balancer.",
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the size of the load balancer",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "public-facing IP address of the load balancer",
			},
			"algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "algorithm used to determine which backend Droplet will be selected by a client",
				Deprecated:  "This field has been deprecated. You can no longer specify an algorithm for load balancers.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "current state of the Load Balancer",
			},
			"ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forwarding_rule": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the protocol used for traffic to the load balancer",
						},
						"entry_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the port on which the load balancer instance will listen",
						},
						"target_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the protocol used for traffic to the backend droplets",
						},
						"target_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the port on the backend Droplets to which the load balancer will send traffic",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the id of the tls certificate used for ssl termination if enabled",
						},
						"certificate_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the name of the tls certificate used for ssl termination if enabled",
						},
						"tls_passthrough": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether ssl encrypted traffic will be passed through to the backend droplets",
						},
					},
				},
				Description: "list of forwarding rules of the load balancer",
				Set:         hashForwardingRules,
			},
			"healthcheck": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the protocol used for health checks sent to the backend droplets",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the port on the backend droplets on which the health check will attempt a connection",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the path on the backend Droplets to which the Load Balancer will send a request",
						},
						"check_interval_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the number of seconds between between two consecutive health checks",
						},
						"response_timeout_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the number of seconds to wait for a response until marking a health check as failed",
						},
						"unhealthy_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of times a health check must fail for a backend droplet to be marked 'unhealthy' and be removed from the pool",
						},
						"healthy_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the number of times a health check must pass for a backend droplet to be marked 'healthy' and be re-added to the pool",
						},
					},
				},
				Description: "health check settings for the load balancer",
			},

			"sticky_sessions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "how and if requests from a client will be persistently served by the same backend droplet",
						},
						"cookie_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the name of the cookie sent to the client",
						},
						"cookie_ttl_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the number of seconds until the cookie set by the Load Balancer expires",
						},
					},
				},
				Description: "sticky sessions settings for the load balancer",
			},
			"droplet_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
				Description: "ids of the droplets assigned to the load balancer",
			},
			"droplet_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the name of a tag corresponding to droplets assigned to the load balancer",
			},
			"redirect_http_to_https": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether http requests will be redirected to https",
			},
			"enable_proxy_protocol": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether PROXY Protocol should be used to pass information from connecting client requests to the backend service",
			},
			"enable_backend_keepalive": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether HTTP keepalive connections are maintained to target Droplets",
			},
			"disable_lets_encrypt_dns_records": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether to disable automatic DNS record creation for Let's Encrypt certificates that are added to the load balancer",
			},
			"vpc_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the VPC in which the load balancer is located",
			},
			"http_idle_timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: " Specifies the idle timeout for HTTPS connections on the load balancer.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the project that the load balancer is associated with.",
			},
			"firewall": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "the firewall rules for allowing/denying traffic to the load balancer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "the rules for ALLOWING traffic to the LB (strings in the form: 'ip:1.2.3.4' or 'cidr:1.2.0.0/16')",
						},
						"deny": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "the rules for DENYING traffic to the LB (strings in the form: 'ip:1.2.3.4' or 'cidr:1.2.0.0/16')",
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the type of the load balancer (GLOBAL, REGIONAL, or REGIONAL_NETWORK)",
			},
			"domains": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "the list of domains required to ingress traffic to global load balancer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "domain name",
						},
						"is_managed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "flag indicating if domain is managed by DigitalOcean",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "certificate ID for TLS handshaking",
						},
						"certificate_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of certificate required for TLS handshaking",
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
				Type:        schema.TypeList,
				Computed:    true,
				Description: "configuration options for global load balancer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "target protocol rules",
						},
						"target_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "target port rules",
						},
						"region_priorities": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "region priority map",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"failover_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "fail-over threshold",
						},
						"cdn": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CDN specific configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
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
				Computed:    true,
				Description: "list of load balancer IDs to put behind a global load balancer",
			},
			"network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the type of network the load balancer is accessible from (EXTERNAL or INTERNAL)",
			},
		},
	}
}

func dataSourceDigitalOceanLoadbalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var foundLoadbalancer *godo.LoadBalancer

	if id, ok := d.GetOk("id"); ok {
		loadbalancer, _, err := client.LoadBalancers.Get(context.Background(), id.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		foundLoadbalancer = loadbalancer
	} else if name, ok := d.GetOk("name"); ok {
		opts := &godo.ListOptions{
			Page:    1,
			PerPage: 200,
		}

		lbList := []godo.LoadBalancer{}

		for {
			lbs, resp, err := client.LoadBalancers.List(context.Background(), opts)

			if err != nil {
				return diag.Errorf("Error retrieving load balancers: %s", err)
			}

			lbList = append(lbList, lbs...)

			if resp.Links == nil || resp.Links.IsLastPage() {
				break
			}

			page, err := resp.Links.CurrentPage()
			if err != nil {
				return diag.Errorf("Error retrieving load balancers: %s", err)
			}

			opts.Page = page + 1
		}

		loadbalancer, err := findLoadBalancerByName(lbList, name.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		foundLoadbalancer = loadbalancer
	} else {
		return diag.Errorf("Error: specify either a name, or id to use to look up the load balancer")
	}

	d.SetId(foundLoadbalancer.ID)
	d.Set("name", foundLoadbalancer.Name)
	d.Set("urn", foundLoadbalancer.URN())
	if foundLoadbalancer.Region != nil {
		d.Set("region", foundLoadbalancer.Region.Slug)
	}
	if foundLoadbalancer.SizeUnit > 0 {
		d.Set("size_unit", foundLoadbalancer.SizeUnit)
	} else if foundLoadbalancer.SizeSlug != "" {
		d.Set("size", foundLoadbalancer.SizeSlug)
	}

	d.Set("ip", foundLoadbalancer.IP)
	d.Set("algorithm", foundLoadbalancer.Algorithm)
	d.Set("status", foundLoadbalancer.Status)
	d.Set("droplet_tag", foundLoadbalancer.Tag)
	d.Set("redirect_http_to_https", foundLoadbalancer.RedirectHttpToHttps)
	d.Set("enable_proxy_protocol", foundLoadbalancer.EnableProxyProtocol)
	d.Set("enable_backend_keepalive", foundLoadbalancer.EnableBackendKeepalive)
	d.Set("disable_lets_encrypt_dns_records", foundLoadbalancer.DisableLetsEncryptDNSRecords)
	d.Set("vpc_uuid", foundLoadbalancer.VPCUUID)
	d.Set("http_idle_timeout_seconds", foundLoadbalancer.HTTPIdleTimeoutSeconds)
	d.Set("project_id", foundLoadbalancer.ProjectID)
	d.Set("type", foundLoadbalancer.Type)
	d.Set("network", foundLoadbalancer.Network)

	if foundLoadbalancer.IPv6 != "" {
		d.Set("ipv6", foundLoadbalancer.IPv6)
	}

	if err := d.Set("droplet_ids", flattenDropletIds(foundLoadbalancer.DropletIDs)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer droplet_ids - error: %#v", err)
	}

	if err := d.Set("sticky_sessions", flattenStickySessions(foundLoadbalancer.StickySessions)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer sticky_sessions - error: %#v", err)
	}

	if err := d.Set("healthcheck", flattenHealthChecks(foundLoadbalancer.HealthCheck)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer healthcheck - error: %#v", err)
	}

	forwardingRules, err := flattenForwardingRules(client, foundLoadbalancer.ForwardingRules)
	if err != nil {
		return diag.Errorf("[DEBUG] Error building Load Balancer forwarding rules - error: %#v", err)
	}

	if err := d.Set("forwarding_rule", forwardingRules); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer forwarding_rule - error: %#v", err)
	}

	if err := d.Set("firewall", flattenLBFirewall(foundLoadbalancer.Firewall)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer firewall - error: %#v", err)
	}

	domains, err := flattenDomains(client, foundLoadbalancer.Domains)
	if err != nil {
		return diag.Errorf("[DEBUG] Error building Load Balancer domains - error: %#v", err)
	}

	if err := d.Set("domains", domains); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer domains - error: %#v", err)
	}

	if err := d.Set("glb_settings", flattenGLBSettings(foundLoadbalancer.GLBSettings)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Load Balancer glb settings - error: %#v", err)
	}

	if err := d.Set("target_load_balancer_ids", flattenLoadBalancerIds(foundLoadbalancer.TargetLoadBalancerIDs)); err != nil {
		return diag.Errorf("[DEBUG] Error setting target Load Balancer ids - error: %#v", err)
	}

	return nil
}

func findLoadBalancerByName(lbs []godo.LoadBalancer, name string) (*godo.LoadBalancer, error) {
	results := make([]godo.LoadBalancer, 0)
	for _, v := range lbs {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no load balancer found with name %s", name)
	}
	return nil, fmt.Errorf("too many load balancers found with name %s (found %d, expected 1)", name, len(results))
}
