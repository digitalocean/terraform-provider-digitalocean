package digitalocean

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

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
			Type:     schema.TypeString,
			Computed: true,
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
		"count": d.Get("count"),
		"tags":  d.Get("tags"),
	}

	_, err := digitaloceanKubernetesNodePoolUpdate(client, rawPool, d.Get("cluster_id").(string), d.Id(), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error updating node pool: %s", err)
	}

	return resourceDigitalOceanKubernetesNodePoolRead(d, meta)
}

func resourceDigitalOceanKubernetesNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	return digitaloceanKubernetesNodePoolDelete(client, d.Get("cluster_id").(string), d.Id())
}

func digitaloceanKubernetesNodePoolCreate(client *godo.Client, pool map[string]interface{}, clusterID, name string) (*godo.KubernetesNodePool, error) {
	p, _, err := client.Kubernetes.CreateNodePool(context.Background(), clusterID, &godo.KubernetesNodePoolCreateRequest{
		Name:  name,
		Size:  pool["size"].(string),
		Count: pool["count"].(int),
		Tags:  expandTags(pool["tags"].(*schema.Set).List()),
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

func digitaloceanKubernetesNodePoolUpdate(client *godo.Client, pool map[string]interface{}, clusterID, poolID, name string) (*godo.KubernetesNodePool, error) {
	p, _, err := client.Kubernetes.UpdateNodePool(context.Background(), clusterID, poolID, &godo.KubernetesNodePoolUpdateRequest{
		Name:  name,
		Count: pool["count"].(int),
		Tags:  expandTags(pool["tags"].(*schema.Set).List()),
	})

	err = waitForKubernetesNodePoolCreate(client, clusterID, poolID)
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
	ticker := time.NewTicker(10 * time.Second)
	timeout := 120
	n := 0

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

		if n > timeout {
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
			return fmt.Errorf("Error trying to read nodepool state: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			ticker.Stop()
			return nil
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

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func digitaloceanRandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
