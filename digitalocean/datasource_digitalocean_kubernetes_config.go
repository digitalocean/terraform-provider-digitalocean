package digitalocean

import (
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kr/pretty"
	"gopkg.in/yaml.v2"
)

type kubernetesConfig struct {
	Clusters []kubernetesConfigCluster `yaml:"clusters"`
	Users    []kubernetesConfigUser    `yaml:"users"`
}

type kubernetesConfigCluster struct {
	Cluster kubernetesConfigClusterData `yaml:"cluster"`
	Name    string                      `yaml:"name"`
}
type kubernetesConfigClusterData struct {
	ClusterCACertificate string `yaml:"certificate-authority-data"`
	Server               string `yaml:"server"`
}

type kubernetesConfigUser struct {
	Name string                   `yaml:"name"`
	User kubernetesConfigUserData `yaml:"user"`
}

type kubernetesConfigUserData struct {
	ClientKeyData     string `yaml:"client-key-data"`
	ClientCertificate string `yaml:"client-certificate-data"`
}

func kubernetesConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"raw_config": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"host": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"cluster_ca_certificate": {
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
			},
		},
	}
}

func flattenKubeConfig(config *godo.KubernetesClusterConfig) []interface{} {
	rawConfig := map[string]interface{}{
		"raw_config": string(config.KubeconfigYAML),
	}

	// parse the yaml into an object
	var c kubernetesConfig
	err := yaml.Unmarshal(config.KubeconfigYAML, &c)
	if err != nil {
		log.Printf("[DEBUG] error unmarshalling config: %s", err)
		return nil
	}

	pretty.Println(c)

	if len(c.Clusters) < 1 {
		return []interface{}{rawConfig}
	}

	rawConfig["cluster_ca_certificate"] = c.Clusters[0].Cluster.ClusterCACertificate
	rawConfig["host"] = c.Clusters[0].Cluster.Server

	if len(c.Users) < 1 {
		return []interface{}{rawConfig}
	}

	rawConfig["client_key"] = c.Users[0].User.ClientKeyData
	rawConfig["client_certificate"] = c.Users[0].User.ClientCertificate

	return []interface{}{rawConfig}
}
