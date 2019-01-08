package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceDigitalOceanLoadbalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanLoadbalancerRead,
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the load balancer",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the region that the load blanacer is deployed in",
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
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "current state of the Load Balancer",
			},
			"forwarding_rule": {
				Type:     schema.TypeList,
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
						"tls_passthrough": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether ssl encrypted traffic will be passed through to the backend droplets",
						},
					},
				},
				Description: "list of forwarding rules of the load balancer",
			},
			"healthcheck": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
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
				MaxItems: 1,
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
		},
	}
}

func dataSourceDigitalOceanLoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	name := d.Get("name").(string)

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	lbList := []godo.LoadBalancer{}

	for {
		lbs, resp, err := client.LoadBalancers.List(context.Background(), opts)

		if err != nil {
			return fmt.Errorf("Error retrieving load balancers: %s", err)
		}

		for _, lb := range lbs {
			lbList = append(lbList, lb)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return fmt.Errorf("Error retrieving load balancers: %s", err)
		}

		opts.Page = page + 1
	}

	loadbalancer, err := findLoadBalancerByName(lbList, name)

	if err != nil {
		return err
	}

	d.SetId(loadbalancer.ID)
	d.Set("name", loadbalancer.Name)
	d.Set("region", loadbalancer.Region.Slug)
	d.Set("ip", loadbalancer.IP)
	d.Set("algorithm", loadbalancer.Algorithm)
	d.Set("status", loadbalancer.Status)
	d.Set("droplet_tag", loadbalancer.Tag)
	d.Set("redirect_http_to_https", loadbalancer.RedirectHttpToHttps)

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
