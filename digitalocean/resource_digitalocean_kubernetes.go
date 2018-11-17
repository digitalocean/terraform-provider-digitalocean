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
		Update:        resourceDigitalOceanKubernetesUpdate,
		Delete:        resourceDigitalOceanKubernetesDelete,
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

			"node_pool": nodePoolSchema(),

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
		Optional: true,
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
					Required: true,
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
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
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
		Computed: true,
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
		Name:        d.Get("name").(string),
		RegionSlug:  d.Get("region").(string),
		VersionSlug: d.Get("version").(string),
		Tags:        expandTags(d.Get("tags").(*schema.Set).List()),
		NodePools:   expandNodePools(d.Get("node_pool").(*schema.Set).List()),
	}

	cluster, _, err := client.Kubernetes.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating Kubernetes cluster: %s", err)
	}

	d.SetId(cluster.ID)

	return resourceDigitalOceanKubernetesRead(d, meta)
}

func expandNodePools(nodePools []interface{}) []*godo.KubernetesNodePoolCreateRequest {
	expandedNodePools := make([]*godo.KubernetesNodePoolCreateRequest, 0, len(nodePools))
	for _, rawPool := range nodePools {

		pool := rawPool.(map[string]interface{})
		cr := &godo.KubernetesNodePoolCreateRequest{
			Name:  pool["name"].(string),
			Size:  pool["size"].(string),
			Count: pool["count"].(int),
			Tags:  expandTags(pool["tags"].(*schema.Set).List()),
		}

		expandedNodePools = append(expandedNodePools, cr)
	}

	return expandedNodePools
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

func resourceDigitalOceanKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDigitalOceanKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	resp, err := client.Kubernetes.Delete(context.Background(), d.Id())
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Unable to delete cluster: %s", err)
	}

	d.SetId("")

	return nil
}
