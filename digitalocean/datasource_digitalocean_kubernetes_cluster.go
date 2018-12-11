package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceDigitalOceanKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
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

			"node_pool": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"tags": tagsSchema(),

						"nodes": nodeSchema(),
					},
				},
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

			"kube_config": kubernetesConfigSchema(),
		},
	}
}

func dataSourceDigitalOceanKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	clusters, resp, err := client.Kubernetes.List(context.Background(), &godo.ListOptions{})
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("No clusters found")
		}

		return fmt.Errorf("Error listing Kuberentes clusters: %s", err)
	}

	// select the correct cluster
	for _, c := range clusters {
		if c.Name == d.Get("name").(string) {
			d.SetId(c.ID)

			return digitaloceanKubernetesClusterRead(client, c, d)
		}
	}

	return fmt.Errorf("Unable to find cluster with name: %s", d.Get("name").(string))
}
