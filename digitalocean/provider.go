package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a schema.Provider for DigitalOcean.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"DIGITALOCEAN_TOKEN",
					"DIGITALOCEAN_ACCESS_TOKEN",
				}, nil),
				Description: "The token key for API operations.",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DIGITALOCEAN_API_URL", "https://api.digitalocean.com"),
				Description: "The URL to use for the DigitalOcean API.",
			},
			"spaces_access_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPACES_ACCESS_KEY_ID", nil),
				Description: "The access key ID for Spaces API operations.",
			},
			"spaces_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPACES_SECRET_ACCESS_KEY", nil),
				Description: "The secret access key for Spaces API operations.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"digitalocean_account":             dataSourceDigitalOceanAccount(),
			"digitalocean_certificate":         dataSourceDigitalOceanCertificate(),
			"digitalocean_database_cluster":    dataSourceDigitalOceanDatabaseCluster(),
			"digitalocean_domain":              dataSourceDigitalOceanDomain(),
			"digitalocean_droplet":             dataSourceDigitalOceanDroplet(),
			"digitalocean_droplet_snapshot":    dataSourceDigitalOceanDropletSnapshot(),
			"digitalocean_floating_ip":         dataSourceDigitalOceanFloatingIp(),
			"digitalocean_image":               dataSourceDigitalOceanImage(),
			"digitalocean_kubernetes_cluster":  dataSourceDigitalOceanKubernetesCluster(),
			"digitalocean_kubernetes_versions": dataSourceDigitalOceanKubernetesVersions(),
			"digitalocean_loadbalancer":        dataSourceDigitalOceanLoadbalancer(),
			"digitalocean_project":             dataSourceDigitalOceanProject(),
			"digitalocean_projects":            dataSourceDigitalOceanProjects(),
			"digitalocean_record":              dataSourceDigitalOceanRecord(),
			"digitalocean_region":              dataSourceDigitalOceanRegion(),
			"digitalocean_regions":             dataSourceDigitalOceanRegions(),
			"digitalocean_sizes":               dataSourceDigitalOceanSizes(),
			"digitalocean_ssh_key":             dataSourceDigitalOceanSSHKey(),
			"digitalocean_tag":                 dataSourceDigitalOceanTag(),
			"digitalocean_volume_snapshot":     dataSourceDigitalOceanVolumeSnapshot(),
			"digitalocean_volume":              dataSourceDigitalOceanVolume(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"digitalocean_certificate":              resourceDigitalOceanCertificate(),
			"digitalocean_cdn":                      resourceDigitalOceanCDN(),
			"digitalocean_database_cluster":         resourceDigitalOceanDatabaseCluster(),
			"digitalocean_database_connection_pool": resourceDigitalOceanDatabaseConnectionPool(),
			"digitalocean_database_db":              resourceDigitalOceanDatabaseDB(),
			"digitalocean_database_firewall":        resourceDigitalOceanDatabaseFirewall(),
			"digitalocean_database_replica":         resourceDigitalOceanDatabaseReplica(),
			"digitalocean_database_user":            resourceDigitalOceanDatabaseUser(),
			"digitalocean_domain":                   resourceDigitalOceanDomain(),
			"digitalocean_droplet":                  resourceDigitalOceanDroplet(),
			"digitalocean_droplet_snapshot":         resourceDigitalOceanDropletSnapshot(),
			"digitalocean_firewall":                 resourceDigitalOceanFirewall(),
			"digitalocean_floating_ip":              resourceDigitalOceanFloatingIp(),
			"digitalocean_floating_ip_assignment":   resourceDigitalOceanFloatingIpAssignment(),
			"digitalocean_kubernetes_cluster":       resourceDigitalOceanKubernetesCluster(),
			"digitalocean_kubernetes_node_pool":     resourceDigitalOceanKubernetesNodePool(),
			"digitalocean_loadbalancer":             resourceDigitalOceanLoadbalancer(),
			"digitalocean_project":                  resourceDigitalOceanProject(),
			"digitalocean_record":                   resourceDigitalOceanRecord(),
			"digitalocean_spaces_bucket":            resourceDigitalOceanBucket(),
			"digitalocean_ssh_key":                  resourceDigitalOceanSSHKey(),
			"digitalocean_tag":                      resourceDigitalOceanTag(),
			"digitalocean_volume":                   resourceDigitalOceanVolume(),
			"digitalocean_volume_attachment":        resourceDigitalOceanVolumeAttachment(),
			"digitalocean_volume_snapshot":          resourceDigitalOceanVolumeSnapshot(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		Token:            d.Get("token").(string),
		APIEndpoint:      d.Get("api_endpoint").(string),
		AccessID:         d.Get("spaces_access_id").(string),
		SecretKey:        d.Get("spaces_secret_key").(string),
		TerraformVersion: terraformVersion,
	}

	return config.Client()
}
