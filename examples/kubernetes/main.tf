/*
    ID            string   `json:"id,omitempty"`
    Name          string   `json:"name,omitempty"`
    RegionSlug    string   `json:"region,omitempty"`
    VersionSlug   string   `json:"version,omitempty"`
    ClusterSubnet string   `json:"cluster_subnet,omitempty"`
    ServiceSubnet string   `json:"service_subnet,omitempty"`
    IPv4          string   `json:"ipv4,omitempty"`
    Endpoint      string   `json:"endpoint,omitempty"`
    Tags          []string `json:"tags,omitempty"`

    NodePools []*KubernetesNodePool `json:"node_pools,omitempty"`

    Status    *KubernetesClusterStatus `json:"status,omitempty"`
    CreatedAt time.Time                `json:"created_at,omitempty"`
    UpdatedAt time.Time                `json:"updated_at,omitempty"`
*/

/*
type KubernetesNodePool struct {
    ID    string   `json:"id,omitempty"`
    Name  string   `json:"name,omitempty"`
    Size  string   `json:"size,omitempty"`
    Count int      `json:"count,omitempty"`
    Tags  []string `json:"tags,omitempty"`

    Nodes []*KubernetesNode `json:"nodes,omitempty"`
}
*/

resource "digitalocean_kubernetes" "k8s" {
  name   = "example"
  region = "lon1"

  version = "v1.10.1"

  cluster_subnet = "10.1.0.0/24"
  service_subnet = "10.1.2.0/24"

  tags = ["foo", "bar"]

  node_pool {
    name  = "default"
    size  = "xlarge"
    count = 3

    tags = ["foo"]
  }
}
