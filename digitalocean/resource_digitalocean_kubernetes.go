package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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

			"kube_config": kubernetesConfigSchema(),
		},
	}
}

func nodePoolSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
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
		Type:     schema.TypeList,
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

func resourceDigitalOceanKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts := &godo.KubernetesClusterCreateRequest{
		Name:        d.Get("name").(string),
		RegionSlug:  d.Get("region").(string),
		VersionSlug: d.Get("version").(string),
		Tags:        expandTags(d.Get("tags").(*schema.Set).List()),
		NodePools:   expandNodePools(d.Get("node_pool").([]interface{})),
	}

	cluster, _, err := client.Kubernetes.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating Kubernetes cluster: %s", err)
	}

	// wait for completion
	cluster, err = waitForKubernetesClusterCreate(client, cluster.ID)
	if err != nil {
		return fmt.Errorf("Error creating Kubernetes cluster: %s", err)
	}

	// set the cluster id
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
	d.Set("region", cluster.RegionSlug)
	d.Set("version", cluster.VersionSlug)
	d.Set("cluster_subnet", cluster.ClusterSubnet)
	d.Set("service_subnet", cluster.ServiceSubnet)
	d.Set("ipv4_address", cluster.IPv4)
	d.Set("endpoint", cluster.Endpoint)
	d.Set("tags", flattenTags(cluster.Tags))
	d.Set("status", cluster.Status.State)
	d.Set("created_at", cluster.CreatedAt.UTC().String())
	d.Set("updated_at", cluster.UpdatedAt.UTC().String())

	if err := d.Set("node_pool", flattenNodePools(cluster.NodePools)); err != nil {
		log.Printf("[DEBUG] Error setting node pool attributes: %s %#v", err, cluster.NodePools)
	}

	// fetch the K8s config  and update the resource
	config, resp, err := client.Kubernetes.GetKubeConfig(context.Background(), cluster.ID)
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("Unable to fetch Kubernetes config: %s", err)
		}
	}
	d.Set("kube_config", flattenKubeConfig(config))

	return nil
}

func waitForKubernetesClusterCreate(client *godo.Client, id string) (*godo.KubernetesCluster, error) {
	ticker := time.NewTicker(10 * time.Second)
	timeout := 120
	n := 0

	for range ticker.C {
		cluster, _, err := client.Kubernetes.Get(context.Background(), id)
		if err != nil {
			ticker.Stop()
			return nil, fmt.Errorf("Error trying to read cluster state: %s", err)
		}

		if cluster.Status.State == "running" {
			ticker.Stop()
			return cluster, nil
		}

		if n > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return nil, fmt.Errorf("Timeout waiting to create cluster")
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

func flattenNodePools(pools []*godo.KubernetesNodePool) []interface{} {
	if pools == nil {
		return nil
	}

	flattenedPools := make([]interface{}, 0)
	for _, pool := range pools {
		rawPool := map[string]interface{}{
			"name":  pool.Name,
			"size":  pool.Size,
			"count": pool.Count,
		}

		if pool.Tags != nil {
			rawPool["tags"] = flattenTags(pool.Tags)
		}

		if pool.Nodes != nil {
			rawPool["nodes"] = flattenNodes(pool.Nodes)
		}

		flattenedPools = append(flattenedPools, rawPool)
	}

	return flattenedPools
}

func flattenNodes(nodes []*godo.KubernetesNode) []interface{} {
	if nodes == nil {
		return nil
	}

	flattenedNodes := make([]interface{}, 0)
	for _, node := range nodes {
		rawNode := map[string]interface{}{
			"name":       node.Name,
			"status":     node.Status.State,
			"created_at": node.CreatedAt.UTC().String(),
			"updated_at": node.UpdatedAt.UTC().String(),
		}

		flattenedNodes = append(flattenedNodes, rawNode)
	}

	return flattenedNodes
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
