package digitalocean

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// to distinguish between a node pool resource and the default pool from the cluster
// we automatically add this tag to the default pool
const digitaloceanKubernetesDefaultNodePoolTag = "terraform:default-node-pool"

func resourceDigitalOceanKubernetesNodePool() *schema.Resource {

	return &schema.Resource{
		Create:        resourceDigitalOceanKubernetesNodePoolCreate,
		Read:          resourceDigitalOceanKubernetesNodePoolRead,
		Update:        resourceDigitalOceanKubernetesNodePoolUpdate,
		Delete:        resourceDigitalOceanKubernetesNodePoolDelete,
		SchemaVersion: 1,

		Schema: nodePoolResourceSchema(),
	}
}

func nodePoolResourceSchema() map[string]*schema.Schema {
	s := nodePoolSchema()

	// add the cluster id
	s["cluster_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.NoZeroValues,
		ForceNew:     true,
	}

	// remove the id when this is used in a specific resource
	// not as a child
	delete(s, "id")

	return s
}

func nodePoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"node_count": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"tags": tagsSchema(),

		"nodes": nodeSchema(),
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

func resourceDigitalOceanKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	rawPool := map[string]interface{}{
		"name":       d.Get("name"),
		"size":       d.Get("size"),
		"node_count": d.Get("node_count"),
		"tags":       d.Get("tags"),
	}

	pool, err := digitaloceanKubernetesNodePoolCreate(client, rawPool, d.Get("cluster_id").(string))
	if err != nil {
		return fmt.Errorf("Error creating Kubernetes node pool: %s", err)
	}

	d.SetId(pool.ID)

	return resourceDigitalOceanKubernetesNodePoolRead(d, meta)
}

func resourceDigitalOceanKubernetesNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	pool, resp, err := client.Kubernetes.GetNodePool(context.Background(), d.Get("cluster_id").(string), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Kubernetes node pool: %s", err)
	}

	d.Set("name", pool.Name)
	d.Set("size", pool.Size)
	d.Set("node_count", pool.Count)
	d.Set("tags", flattenTags(filterTags(pool.Tags)))

	if pool.Nodes != nil {
		d.Set("nodes", flattenNodes(pool.Nodes))
	}

	return nil
}

func resourceDigitalOceanKubernetesNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	rawPool := map[string]interface{}{
		"name":       d.Get("name"),
		"node_count": d.Get("node_count"),
		"tags":       d.Get("tags"),
	}

	_, err := digitaloceanKubernetesNodePoolUpdate(client, rawPool, d.Get("cluster_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error updating node pool: %s", err)
	}

	return resourceDigitalOceanKubernetesNodePoolRead(d, meta)
}

func resourceDigitalOceanKubernetesNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	return digitaloceanKubernetesNodePoolDelete(client, d.Get("cluster_id").(string), d.Id())
}

func digitaloceanKubernetesNodePoolCreate(client *godo.Client, pool map[string]interface{}, clusterID string, customTags ...string) (*godo.KubernetesNodePool, error) {
	// append any custom tags
	tags := expandTags(pool["tags"].(*schema.Set).List())
	tags = append(tags, customTags...)

	p, _, err := client.Kubernetes.CreateNodePool(context.Background(), clusterID, &godo.KubernetesNodePoolCreateRequest{
		Name:  pool["name"].(string),
		Size:  pool["size"].(string),
		Count: pool["node_count"].(int),
		Tags:  tags,
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to create new default node pool %s", err)
	}

	err = waitForKubernetesNodePoolCreate(client, clusterID, p.ID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func digitaloceanKubernetesNodePoolUpdate(client *godo.Client, pool map[string]interface{}, clusterID, poolID string, customTags ...string) (*godo.KubernetesNodePool, error) {
	tags := expandTags(pool["tags"].(*schema.Set).List())
	tags = append(tags, customTags...)

	count := pool["node_count"].(int)
	p, resp, err := client.Kubernetes.UpdateNodePool(context.Background(), clusterID, poolID, &godo.KubernetesNodePoolUpdateRequest{
		Name:  pool["name"].(string),
		Count: &count,
		Tags:  tags,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil, nil
		}

		return nil, fmt.Errorf("Unable to update nodepool: %s", err)
	}

	err = waitForKubernetesNodePoolCreate(client, clusterID, p.ID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func digitaloceanKubernetesNodePoolDelete(client *godo.Client, clusterID, poolID string) error {
	// delete the old pool
	_, err := client.Kubernetes.DeleteNodePool(context.Background(), clusterID, poolID)
	if err != nil {
		return fmt.Errorf("Unable to delete node pool %s", err)
	}

	err = waitForKubernetesNodePoolDelete(client, clusterID, poolID)
	if err != nil {
		return err
	}

	return nil
}

func waitForKubernetesNodePoolCreate(client *godo.Client, id string, poolID string) error {
	tickerInterval := 10 //10s
	timeout := 1800      //1800s, 30min
	n := 0

	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
	for range ticker.C {
		pool, _, err := client.Kubernetes.GetNodePool(context.Background(), id, poolID)
		if err != nil {
			ticker.Stop()
			return fmt.Errorf("Error trying to read nodepool state: %s", err)
		}

		allRunning := len(pool.Nodes) == pool.Count
		for _, n := range pool.Nodes {
			if n.Status.State != "running" {
				allRunning = false
			}
		}

		if allRunning {
			ticker.Stop()
			return nil
		}

		if n*tickerInterval > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return fmt.Errorf("Timeout waiting to create nodepool")
}

func waitForKubernetesNodePoolDelete(client *godo.Client, id string, poolID string) error {
	tickerInterval := 10 //10s
	timeout := 1800      //1800s, 30min
	n := 0

	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
	for range ticker.C {
		_, resp, err := client.Kubernetes.GetNodePool(context.Background(), id, poolID)
		if err != nil {
			ticker.Stop()

			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return fmt.Errorf("Error trying to read nodepool state: %s", err)
		}

		if n*tickerInterval > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return fmt.Errorf("Timeout waiting to delete nodepool")
}

func expandNodePools(nodePools []interface{}) []*godo.KubernetesNodePool {
	expandedNodePools := make([]*godo.KubernetesNodePool, 0, len(nodePools))
	for _, rawPool := range nodePools {
		pool := rawPool.(map[string]interface{})
		cr := &godo.KubernetesNodePool{
			ID:    pool["id"].(string),
			Name:  pool["name"].(string),
			Size:  pool["size"].(string),
			Count: pool["node_count"].(int),
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

func flattenNodePool(pool *godo.KubernetesNodePool, parentTags ...string) []interface{} {
	rawPool := map[string]interface{}{
		"id":         pool.ID,
		"name":       pool.Name,
		"size":       pool.Size,
		"node_count": pool.Count,
	}

	if pool.Tags != nil {
		rawPool["tags"] = flattenTags(filterTags(pool.Tags))
	}

	if pool.Nodes != nil {
		rawPool["nodes"] = flattenNodes(pool.Nodes)
	}

	return []interface{}{rawPool}
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
