package digitalocean

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const digitaloceanKubernetesDefaultNodePoolTag = "k8s:default-node-pool"

func resourceDigitalOceanKubernetesCreate() *schema.Resource {

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

	// for a node pool resource name is not computed
	s["name"].Computed = false

	// add the cluster id
	s["cluster_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.NoZeroValues,
		ForceNew:     true,
	}

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
			ValidateFunc: validation.NoZeroValues,
		},

		"count": {
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
	client := meta.(*godo.Client)

	rawPool := map[string]interface{}{
		"size":  d.Get("size"),
		"count": d.Get("count"),
		"tags":  d.Get("tags"),
	}

	pool, err := digitaloceanKubernetesNodePoolCreate(client, rawPool, d.Get("cluster_id").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error creating Kubernetes node pool: %s", err)
	}

	d.SetId(pool.ID)

	return resourceDigitalOceanKubernetesNodePoolRead(d, meta)
}

func resourceDigitalOceanKubernetesNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	pool, resp, err := client.Kubernetes.GetNodePool(context.Background(), d.Get("cluster_id").(string), d.Id())
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Kubernetes node pool: %s", err)
	}

	cluster, resp, err := client.Kubernetes.Get(context.Background(), d.Get("cluster_id").(string))
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Kubernetes cluster: %s", err)
	}

	d.Set("name", pool.Name)
	d.Set("size", pool.Size)
	d.Set("count", pool.Count)
	d.Set("tags", flattenTags(filterTags(pool.Tags, cluster.Tags...)))

	if pool.Nodes != nil {
		d.Set("nodes", flattenNodes(pool.Nodes))
	}

	return nil
}

func resourceDigitalOceanKubernetesNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	rawPool := map[string]interface{}{
		"name":  d.Get("name"),
		"count": d.Get("count"),
		"tags":  d.Get("tags"),
	}

	_, err := digitaloceanKubernetesNodePoolUpdate(client, rawPool, d.Get("cluster_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error updating node pool: %s", err)
	}

	return resourceDigitalOceanKubernetesNodePoolRead(d, meta)
}

func resourceDigitalOceanKubernetesNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	return digitaloceanKubernetesNodePoolDelete(client, d.Get("cluster_id").(string), d.Id())
}

func digitaloceanKubernetesNodePoolCreate(client *godo.Client, pool map[string]interface{}, clusterID string, customTags ...string) (*godo.KubernetesNodePool, error) {
	// append any custom tags
	tags := expandTags(pool["tags"].(*schema.Set).List())
	tags = append(tags, customTags...)

	p, _, err := client.Kubernetes.CreateNodePool(context.Background(), clusterID, &godo.KubernetesNodePoolCreateRequest{
		Name:  pool["name"].(string),
		Size:  pool["size"].(string),
		Count: pool["count"].(int),
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

	p, resp, err := client.Kubernetes.UpdateNodePool(context.Background(), clusterID, poolID, &godo.KubernetesNodePoolUpdateRequest{
		Name:  pool["name"].(string),
		Count: pool["count"].(int),
		Tags:  tags,
	})

	if err != nil {
		if resp.StatusCode == 404 {
			return nil, nil
		}

		return nil, fmt.Errorf("Unable to update nodepool: %s", err)
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

		allRunning := true
		for _, n := range pool.Nodes {
			if n.Status.State != "running" {
				allRunning = false
			} else {
				fmt.Println(n.Status.State)
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
	ticker := time.NewTicker(10 * time.Second)
	timeout := 120
	n := 0

	for range ticker.C {
		_, resp, err := client.Kubernetes.GetNodePool(context.Background(), id, poolID)
		if err != nil {
			ticker.Stop()

			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return fmt.Errorf("Error trying to read nodepool state: %s", err)
		}

		if n > timeout {
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

func flattenNodePool(pool *godo.KubernetesNodePool, parentTags ...string) []interface{} {
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
