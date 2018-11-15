package digitalocean

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanKubernetes() *schema.Resource {
	return &schema.Resource{
		Create:        resourceDigitalOceanKubernetesCreate,
		Read:          resourceDigitalOceanKubernetesRead,
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},

			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"cluster_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),

			"node_pools": nodePoolSchema(),

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kube_config": kubernetesConfig(),
		},
	}
}

func nodePoolSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: false,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.NoZeroValues,
				},

				"size": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.NoZeroValues,
				},

				"count": {
					Type:     schema.TypeInt,
					Required: false,
					Default:  1,
				},

				"tags": tagsSchema(),

				"nodes": nodeSchema(),
			},
		},
	}
}

func nodeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: false,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.NoZeroValues,
				},

				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"updated_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func kubernetesConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: false,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"raw_config": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"client_key": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"client_certificate": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"cluster_ca_certificate": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"host": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"username": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"password": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func resourceDigitalOceanKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts := &godo.KubernetesClusterCreateRequest{
		Name:       d.Get("name").(string),
		RegionSlug: d.Get("region").(string),
		Tags:       d.Get("tags").([]string),
	}

	cluster, _, err := client.Kubernetes.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating Kubernetes cluster: %s", err)
	}

	d.SetId(cluster.ID)

	return resourceDigitalOceanKubernetesRead(d, meta)
}

func resourceDigitalOceanKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	cluster, resp, err := client.Kubernetes.Get(context.Background(), d.Id())
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Kubernetes cluster: %s", err)
	}

	d.Set("name", cluster.Name)

	return nil
}
