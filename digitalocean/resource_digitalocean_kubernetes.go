package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/kr/pretty"
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
				ForceNew:     true,
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
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		Set:      hashNodePool,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},

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
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},

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

	pools := expandNodePools(d.Get("node_pool").(*schema.Set).List())
	poolCreateRequests := make([]*godo.KubernetesNodePoolCreateRequest, len(pools))
	for i, pool := range pools {
		poolCreateRequests[i] = &godo.KubernetesNodePoolCreateRequest{
			Name:  pool.Name,
			Size:  pool.Size,
			Tags:  pool.Tags,
			Count: pool.Count,
		}
	}

	opts := &godo.KubernetesClusterCreateRequest{
		Name:        d.Get("name").(string),
		RegionSlug:  d.Get("region").(string),
		VersionSlug: d.Get("version").(string),
		Tags:        expandTags(d.Get("tags").(*schema.Set).List()),
		NodePools:   poolCreateRequests,
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
	d.Set("tags", flattenTags(filterTags(cluster.Tags)))
	d.Set("status", cluster.Status.State)
	d.Set("created_at", cluster.CreatedAt.UTC().String())
	d.Set("updated_at", cluster.UpdatedAt.UTC().String())

	if err := d.Set("node_pool", flattenNodePools(cluster.NodePools, cluster.Tags...)); err != nil {
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

func resourceDigitalOceanKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	// Figure out the changes and then call the appropriate API methods
	if d.HasChange("name") || d.HasChange("tags") {

		opts := &godo.KubernetesClusterUpdateRequest{
			Name: d.Get("name").(string),
			Tags: expandTags(d.Get("tags").(*schema.Set).List()),
		}

		_, resp, err := client.Kubernetes.Update(context.Background(), d.Id(), opts)
		if err != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("Unable to update cluster: %s", err)
		}
	}

	// Update node pools
	if d.HasChange("node_pool") {
		old, new := d.GetChange("node_pool")
		fmt.Println("old:")
		pretty.Println(old)

		fmt.Println("new:")
		pretty.Println(new)

		/*
			// process deleted pools
			poolsToDelete := make([]godo.KubernetesNodePool, 0)
			for i, p := range old. {

			}
		*/
	}

	return resourceDigitalOceanKubernetesRead(d, meta)
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

func expandNodePools(nodePools []interface{}) []*godo.KubernetesNodePool {
	expandedNodePools := make([]*godo.KubernetesNodePool, 0, len(nodePools))
	for _, rawPool := range nodePools {
		pool := rawPool.(map[string]interface{})
		cr := &godo.KubernetesNodePool{
			ID:    pool["id"].(string),
			Name:  pool["name"].(string),
			Size:  pool["size"].(string),
			Count: pool["count"].(int),
			Tags:  expandTags(pool["tags"].(*schema.Set).List()),
			Nodes: expandNodes(pool["nodes"].([]interface{})),
		}

		expandedNodePools = append(expandedNodePools, cr)
	}

	return expandedNodePools
}

func expandNodes(nodes []interface{}) []*godo.KubernetesNode {
	expandedNodes := make([]*godo.KubernetesNode, 0, len(nodes))
	for _, rawNode := range nodes {
		node := rawNode.(map[string]interface{})
		n := &godo.KubernetesNode{
			ID:   node["id"].(string),
			Name: node["name"].(string),
		}

		expandedNodes = append(expandedNodes, n)
	}

	return expandedNodes
}

func flattenNodePools(pools []*godo.KubernetesNodePool, parentTags ...string) *schema.Set {
	if pools == nil {
		return nil
	}

	flattenedPools := schema.NewSet(hashNodePool, []interface{}{})
	for _, pool := range pools {
		rawPool := map[string]interface{}{
			"id":    pool.ID,
			"name":  pool.Name,
			"size":  pool.Size,
			"count": pool.Count,
		}

		if pool.Tags != nil {
			rawPool["tags"] = flattenTags(filterTags(pool.Tags, parentTags...))
		}

		if pool.Nodes != nil {
			rawPool["nodes"] = flattenNodes(pool.Nodes)
		}

		flattenedPools.Add(rawPool)
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
			"id":         node.ID,
			"name":       node.Name,
			"status":     node.Status.State,
			"created_at": node.CreatedAt.UTC().String(),
			"updated_at": node.UpdatedAt.UTC().String(),
		}

		flattenedNodes = append(flattenedNodes, rawNode)
	}

	return flattenedNodes
}

// custom hashing function for the set index
func hashNodePool(v interface{}) int {
	pool := v.(map[string]interface{})
	hash := hashcode.String(pool["name"].(string))

	//fmt.Printf("id: %s, hash: %d", pool["name"], hash)

	return hash
}

// we need to filter tags to remove any automatically added to avoid state problems,
// these are tags starting with "k8s:", named "k8s" or duplicates of the cluster tags
func filterTags(tags []string, parentTags ...string) []string {
	filteredTags := make([]string, 0)
	for _, t := range tags {
		if !strings.HasPrefix(t, "k8s:") && t != "k8s" && !tagsContain(parentTags, t) {
			filteredTags = append(filteredTags, t)
		}
	}

	return filteredTags
}

func tagsContain(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}

	return false
}
