package digitalocean

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// to distinguish between a node pool resource and the default pool from the cluster
// we automatically add this tag to the default pool
const digitaloceanKubernetesDefaultNodePoolTag = "terraform:default-node-pool"

func resourceDigitalOceanKubernetesNodePool() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceDigitalOceanKubernetesNodePoolCreate,
		ReadContext:   resourceDigitalOceanKubernetesNodePoolRead,
		UpdateContext: resourceDigitalOceanKubernetesNodePoolUpdate,
		DeleteContext: resourceDigitalOceanKubernetesNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanKubernetesNodePoolImportState,
		},
		SchemaVersion: 1,

		Schema: nodePoolSchema(true),
	}
}

func resourceDigitalOceanKubernetesNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	rawPool := map[string]interface{}{
		"name":       d.Get("name"),
		"size":       d.Get("size"),
		"tags":       d.Get("tags"),
		"labels":     d.Get("labels"),
		"node_count": d.Get("node_count"),
		"auto_scale": d.Get("auto_scale"),
		"min_nodes":  d.Get("min_nodes"),
		"max_nodes":  d.Get("max_nodes"),
		"taint":      d.Get("taint"),
	}

	pool, err := digitaloceanKubernetesNodePoolCreate(client, rawPool, d.Get("cluster_id").(string))
	if err != nil {
		return diag.Errorf("Error creating Kubernetes node pool: %s", err)
	}

	d.SetId(pool.ID)

	return resourceDigitalOceanKubernetesNodePoolRead(ctx, d, meta)
}

func resourceDigitalOceanKubernetesNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	pool, resp, err := client.Kubernetes.GetNodePool(context.Background(), d.Get("cluster_id").(string), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Kubernetes node pool: %s", err)
	}

	d.Set("name", pool.Name)
	d.Set("size", pool.Size)
	d.Set("node_count", pool.Count)
	d.Set("actual_node_count", pool.Count)
	d.Set("tags", flattenTags(filterTags(pool.Tags)))
	d.Set("labels", flattenLabels(pool.Labels))
	d.Set("auto_scale", pool.AutoScale)
	d.Set("min_nodes", pool.MinNodes)
	d.Set("max_nodes", pool.MaxNodes)
	d.Set("nodes", flattenNodes(pool.Nodes))
	d.Set("taint", flattenNodePoolTaints(pool.Taints))

	// Assign a node_count only if it's been set explicitly, since it's
	// optional and we don't want to update with a 0 if it's not set.
	if _, ok := d.GetOk("node_count"); ok {
		d.Set("node_count", pool.Count)
	}

	return nil
}

func resourceDigitalOceanKubernetesNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	rawPool := map[string]interface{}{
		"name": d.Get("name"),
		"tags": d.Get("tags"),
	}

	if _, ok := d.GetOk("node_count"); ok {
		rawPool["node_count"] = d.Get("node_count")
	}

	rawPool["labels"] = d.Get("labels")
	rawPool["auto_scale"] = d.Get("auto_scale")
	rawPool["min_nodes"] = d.Get("min_nodes")
	rawPool["max_nodes"] = d.Get("max_nodes")

	_, newTaint := d.GetChange("taint")
	rawPool["taint"] = newTaint

	_, err := digitaloceanKubernetesNodePoolUpdate(client, rawPool, d.Get("cluster_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("Error updating node pool: %s", err)
	}

	return resourceDigitalOceanKubernetesNodePoolRead(ctx, d, meta)
}

func resourceDigitalOceanKubernetesNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	return digitaloceanKubernetesNodePoolDelete(client, d.Get("cluster_id").(string), d.Id())
}

func resourceDigitalOceanKubernetesNodePoolImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if _, ok := d.GetOk("cluster_id"); ok {
		// Short-circuit: The resource already has a cluster ID, no need to search for it.
		return []*schema.ResourceData{d}, nil
	}

	client := meta.(*CombinedConfig).godoClient()

	nodePoolId := d.Id()

	// Scan all of the Kubernetes clusters to recover the node pool's cluster ID.
	var clusterId string
	var nodePool *godo.KubernetesNodePool
	listOptions := godo.ListOptions{}
	for {
		clusters, response, err := client.Kubernetes.List(context.Background(), &listOptions)
		if err != nil {
			return nil, fmt.Errorf("Unable to list Kubernetes clusters: %v", err)
		}

		for _, cluster := range clusters {
			for _, np := range cluster.NodePools {
				if np.ID == nodePoolId {
					if clusterId != "" {
						// This should never happen but good practice to assert that it does not occur.
						return nil, fmt.Errorf("Illegal state: node pool ID %s is associated with multiple clusters", nodePoolId)
					}
					clusterId = cluster.ID
					nodePool = np
				}
			}
		}

		if response.Links == nil || response.Links.IsLastPage() {
			break
		}

		page, err := response.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		listOptions.Page = page + 1
	}

	if clusterId == "" {
		return nil, fmt.Errorf("Did not find the cluster owning the node pool %s", nodePoolId)
	}

	// Ensure that the node pool does not have the default tag set.
	for _, tag := range nodePool.Tags {
		if tag == digitaloceanKubernetesDefaultNodePoolTag {
			return nil, fmt.Errorf("Node pool %s has the default node pool tag set; import the owning digitalocean_kubernetes_cluster resource instead (cluster ID=%s)",
				nodePoolId, clusterId)
		}
	}

	// Set the cluster_id attribute with the cluster's ID.
	d.Set("cluster_id", clusterId)
	return []*schema.ResourceData{d}, nil
}

func digitaloceanKubernetesNodePoolCreate(client *godo.Client, pool map[string]interface{}, clusterID string, customTags ...string) (*godo.KubernetesNodePool, error) {
	// append any custom tags
	tags := expandTags(pool["tags"].(*schema.Set).List())
	tags = append(tags, customTags...)

	req := &godo.KubernetesNodePoolCreateRequest{
		Name:      pool["name"].(string),
		Size:      pool["size"].(string),
		Count:     pool["node_count"].(int),
		Tags:      tags,
		Labels:    expandLabels(pool["labels"].(map[string]interface{})),
		AutoScale: pool["auto_scale"].(bool),
		MinNodes:  pool["min_nodes"].(int),
		MaxNodes:  pool["max_nodes"].(int),
		Taints:    expandNodePoolTaints(pool["taint"].(*schema.Set).List()),
	}

	p, _, err := client.Kubernetes.CreateNodePool(context.Background(), clusterID, req)

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

	req := &godo.KubernetesNodePoolUpdateRequest{
		Name: pool["name"].(string),
		Tags: tags,
	}

	if pool["node_count"] != nil {
		req.Count = intPtr(pool["node_count"].(int))
	}

	if pool["auto_scale"] == nil {
		pool["auto_scale"] = false
	}
	req.AutoScale = boolPtr(pool["auto_scale"].(bool))

	if pool["min_nodes"] != nil {
		req.MinNodes = intPtr(pool["min_nodes"].(int))
	}

	if pool["max_nodes"] != nil {
		req.MaxNodes = intPtr(pool["max_nodes"].(int))
	}

	if pool["labels"] != nil {
		req.Labels = expandLabels(pool["labels"].(map[string]interface{}))
	}

	if pool["taint"] != nil {
		t := expandNodePoolTaints(pool["taint"].(*schema.Set).List())
		req.Taints = &t
	}

	p, resp, err := client.Kubernetes.UpdateNodePool(context.Background(), clusterID, poolID, req)

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

func digitaloceanKubernetesNodePoolDelete(client *godo.Client, clusterID, poolID string) diag.Diagnostics {
	// delete the old pool
	_, err := client.Kubernetes.DeleteNodePool(context.Background(), clusterID, poolID)
	if err != nil {
		return diag.Errorf("Unable to delete node pool %s", err)
	}

	err = waitForKubernetesNodePoolDelete(client, clusterID, poolID)
	if err != nil {
		return diag.FromErr(err)
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
