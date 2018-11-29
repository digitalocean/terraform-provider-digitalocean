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
		Create:        resourceDigitalOceanKubernetesClusterCreate,
		Read:          resourceDigitalOceanKubernetesClusterRead,
		Update:        resourceDigitalOceanKubernetesClusterUpdate,
		Delete:        resourceDigitalOceanKubernetesClusterDelete,
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

			"node_pool": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: nodePoolSchema(),
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

func resourceDigitalOceanKubernetesClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	pools := expandNodePools(d.Get("node_pool").([]interface{}))
	poolCreateRequests := make([]*godo.KubernetesNodePoolCreateRequest, len(pools))
	for i, pool := range pools {
		tags := append(pool.Tags, digitaloceanKubernetesDefaultNodePoolTag)
		poolCreateRequests[i] = &godo.KubernetesNodePoolCreateRequest{
			Name:  pool.Name,
			Size:  pool.Size,
			Tags:  tags,
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

	return resourceDigitalOceanKubernetesClusterRead(d, meta)
}

func resourceDigitalOceanKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
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

	// find the default node pool from all the pools in the cluster
	// the default node pool has a custom tag k8s:default-node-pool
	for _, p := range cluster.NodePools {
		for _, t := range p.Tags {
			if t == digitaloceanKubernetesDefaultNodePoolTag {
				if err := d.Set("node_pool", flattenNodePool(p, cluster.Tags...)); err != nil {
					log.Printf("[DEBUG] Error setting node pool attributes: %s %#v", err, cluster.NodePools)
				}
			}
		}
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

func resourceDigitalOceanKubernetesClusterUpdate(d *schema.ResourceData, meta interface{}) error {
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

	// Update the node pool if necessary
	if !d.HasChange("node_pool") {
		return resourceDigitalOceanKubernetesClusterRead(d, meta)
	}

	old, new := d.GetChange("node_pool")
	fmt.Println(old)
	fmt.Println(new)

	// has the size changed if so we need to create a new pool and then destroy the old one
	oldPool := old.([]interface{})[0].(map[string]interface{})
	newPool := new.([]interface{})[0].(map[string]interface{})

	if oldPool["size"] != newPool["size"] {
		// create a new pool
		_, err := digitaloceanKubernetesNodePoolCreate(client, newPool, d.Id(), digitaloceanKubernetesDefaultNodePoolTag)
		if err != nil {
			return err
		}

		err = digitaloceanKubernetesNodePoolDelete(client, d.Id(), oldPool["id"].(string))
		if err != nil {
			return err
		}

		return resourceDigitalOceanKubernetesClusterRead(d, meta)
	}

	// only update the existing default pool
	_, err := digitaloceanKubernetesNodePoolUpdate(client, newPool, d.Id(), oldPool["id"].(string), digitaloceanKubernetesDefaultNodePoolTag)
	if err != nil {
		return err
	}

	return resourceDigitalOceanKubernetesClusterRead(d, meta)
}

func resourceDigitalOceanKubernetesClusterDelete(d *schema.ResourceData, meta interface{}) error {
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

// we need to filter tags to remove any automatically added to avoid state problems,
// these are tags starting with "k8s:", named "k8s" or duplicates of the cluster tags
func filterTags(tags []string, parentTags ...string) []string {
	filteredTags := make([]string, 0)
	for _, t := range tags {
		if !strings.HasPrefix(t, "k8s:") &&
			!strings.HasPrefix(t, "terraform:") &&
			t != "k8s" && !tagsContain(parentTags, t) {
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
