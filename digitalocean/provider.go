package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/digitalocean/terraform-provider-digitalocean/internal/mutexkv"
)

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for DigitalOcean.
func Provider() *schema.Provider {
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
			"spaces_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPACES_ENDPOINT_URL", "https://{{.Region}}.digitaloceanspaces.com"),
				Description: "The URL to use for the DigitalOcean Spaces API.",
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
			"digitalocean_account":               dataSourceDigitalOceanAccount(),
			"digitalocean_app":                   dataSourceDigitalOceanApp(),
			"digitalocean_certificate":           dataSourceDigitalOceanCertificate(),
			"digitalocean_container_registry":    dataSourceDigitalOceanContainerRegistry(),
			"digitalocean_database_cluster":      dataSourceDigitalOceanDatabaseCluster(),
			"digitalocean_database_ca":           dataSourceDigitalOceanDatabaseCA(),
			"digitalocean_database_replica":      dataSourceDigitalOceanDatabaseReplica(),
			"digitalocean_domain":                dataSourceDigitalOceanDomain(),
			"digitalocean_domains":               dataSourceDigitalOceanDomains(),
			"digitalocean_droplet":               dataSourceDigitalOceanDroplet(),
			"digitalocean_droplets":              dataSourceDigitalOceanDroplets(),
			"digitalocean_droplet_snapshot":      dataSourceDigitalOceanDropletSnapshot(),
			"digitalocean_firewall":              dataSourceDigitalOceanFirewall(),
			"digitalocean_floating_ip":           dataSourceDigitalOceanFloatingIP(),
			"digitalocean_image":                 dataSourceDigitalOceanImage(),
			"digitalocean_images":                dataSourceDigitalOceanImages(),
			"digitalocean_kubernetes_cluster":    dataSourceDigitalOceanKubernetesCluster(),
			"digitalocean_kubernetes_versions":   dataSourceDigitalOceanKubernetesVersions(),
			"digitalocean_loadbalancer":          dataSourceDigitalOceanLoadbalancer(),
			"digitalocean_project":               dataSourceDigitalOceanProject(),
			"digitalocean_projects":              dataSourceDigitalOceanProjects(),
			"digitalocean_record":                dataSourceDigitalOceanRecord(),
			"digitalocean_records":               dataSourceDigitalOceanRecords(),
			"digitalocean_region":                dataSourceDigitalOceanRegion(),
			"digitalocean_regions":               dataSourceDigitalOceanRegions(),
			"digitalocean_reserved_ip":           dataSourceDigitalOceanReservedIP(),
			"digitalocean_sizes":                 dataSourceDigitalOceanSizes(),
			"digitalocean_spaces_bucket":         dataSourceDigitalOceanSpacesBucket(),
			"digitalocean_spaces_buckets":        dataSourceDigitalOceanSpacesBuckets(),
			"digitalocean_spaces_bucket_object":  dataSourceDigitalOceanSpacesBucketObject(),
			"digitalocean_spaces_bucket_objects": dataSourceDigitalOceanSpacesBucketObjects(),
			"digitalocean_ssh_key":               dataSourceDigitalOceanSSHKey(),
			"digitalocean_ssh_keys":              dataSourceDigitalOceanSSHKeys(),
			"digitalocean_tag":                   dataSourceDigitalOceanTag(),
			"digitalocean_tags":                  dataSourceDigitalOceanTags(),
			"digitalocean_volume_snapshot":       dataSourceDigitalOceanVolumeSnapshot(),
			"digitalocean_volume":                dataSourceDigitalOceanVolume(),
			"digitalocean_vpc":                   dataSourceDigitalOceanVPC(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"digitalocean_app":                                   resourceDigitalOceanApp(),
			"digitalocean_certificate":                           resourceDigitalOceanCertificate(),
			"digitalocean_container_registry":                    resourceDigitalOceanContainerRegistry(),
			"digitalocean_container_registry_docker_credentials": resourceDigitalOceanContainerRegistryDockerCredentials(),
			"digitalocean_cdn":                                   resourceDigitalOceanCDN(),
			"digitalocean_database_cluster":                      resourceDigitalOceanDatabaseCluster(),
			"digitalocean_database_connection_pool":              resourceDigitalOceanDatabaseConnectionPool(),
			"digitalocean_database_db":                           resourceDigitalOceanDatabaseDB(),
			"digitalocean_database_firewall":                     resourceDigitalOceanDatabaseFirewall(),
			"digitalocean_database_replica":                      resourceDigitalOceanDatabaseReplica(),
			"digitalocean_database_user":                         resourceDigitalOceanDatabaseUser(),
			"digitalocean_domain":                                resourceDigitalOceanDomain(),
			"digitalocean_droplet":                               resourceDigitalOceanDroplet(),
			"digitalocean_droplet_snapshot":                      resourceDigitalOceanDropletSnapshot(),
			"digitalocean_firewall":                              resourceDigitalOceanFirewall(),
			"digitalocean_floating_ip":                           resourceDigitalOceanFloatingIP(),
			"digitalocean_floating_ip_assignment":                resourceDigitalOceanFloatingIPAssignment(),
			"digitalocean_kubernetes_cluster":                    resourceDigitalOceanKubernetesCluster(),
			"digitalocean_kubernetes_node_pool":                  resourceDigitalOceanKubernetesNodePool(),
			"digitalocean_loadbalancer":                          resourceDigitalOceanLoadbalancer(),
			"digitalocean_monitor_alert":                         resourceDigitalOceanMonitorAlert(),
			"digitalocean_project":                               resourceDigitalOceanProject(),
			"digitalocean_project_resources":                     resourceDigitalOceanProjectResources(),
			"digitalocean_record":                                resourceDigitalOceanRecord(),
			"digitalocean_reserved_ip":                           resourceDigitalOceanReservedIP(),
			"digitalocean_reserved_ip_assignment":                resourceDigitalOceanReservedIPAssignment(),
			"digitalocean_spaces_bucket":                         resourceDigitalOceanBucket(),
			"digitalocean_spaces_bucket_object":                  resourceDigitalOceanSpacesBucketObject(),
			"digitalocean_spaces_bucket_policy":                  resourceDigitalOceanSpacesBucketPolicy(),
			"digitalocean_ssh_key":                               resourceDigitalOceanSSHKey(),
			"digitalocean_tag":                                   resourceDigitalOceanTag(),
			"digitalocean_volume":                                resourceDigitalOceanVolume(),
			"digitalocean_volume_attachment":                     resourceDigitalOceanVolumeAttachment(),
			"digitalocean_volume_snapshot":                       resourceDigitalOceanVolumeSnapshot(),
			"digitalocean_vpc":                                   resourceDigitalOceanVPC(),
			"digitalocean_custom_image":                          resourceDigitalOceanCustomImage(),
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

	if endpoint, ok := d.GetOk("spaces_endpoint"); ok {
		config.SpacesAPIEndpoint = endpoint.(string)
	}

	return config.Client()
}
