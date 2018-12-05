package digitalocean

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a schema.Provider for DigitalOcean.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DIGITALOCEAN_TOKEN", nil),
				Description: "The token key for API operations.",
			},
			"access_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DO_ACCESS_KEY_ID", nil),
				Description: "The access key ID for Spaces API operations.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DO_SECRET_ACCESS_KEY", nil),
				Description: "The secret access key for Spaces API operations.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"digitalocean_certificate":      dataSourceDigitalOceanCertificate(),
			"digitalocean_domain":           dataSourceDigitalOceanDomain(),
			"digitalocean_droplet":          dataSourceDigitalOceanDroplet(),
			"digitalocean_droplet_snapshot": dataSourceDigitalOceanDropletSnapshot(),
			"digitalocean_floating_ip":      dataSourceDigitalOceanFloatingIp(),
			"digitalocean_image":            dataSourceDigitalOceanImage(),
			"digitalocean_loadbalancer":     dataSourceDigitalOceanLoadbalancer(),
			"digitalocean_record":           dataSourceDigitalOceanRecord(),
			"digitalocean_ssh_key":          dataSourceDigitalOceanSSHKey(),
			"digitalocean_tag":              dataSourceDigitalOceanTag(),
			"digitalocean_volume_snapshot":  dataSourceDigitalOceanVolumeSnapshot(),
			"digitalocean_volume":           dataSourceDigitalOceanVolume(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"digitalocean_bucket":                 resourceDigitalOceanBucket(),
			"digitalocean_certificate":            resourceDigitalOceanCertificate(),
			"digitalocean_domain":                 resourceDigitalOceanDomain(),
			"digitalocean_droplet":                resourceDigitalOceanDroplet(),
			"digitalocean_droplet_snapshot":       resourceDigitalOceanDropletSnapshot(),
			"digitalocean_firewall":               resourceDigitalOceanFirewall(),
			"digitalocean_floating_ip":            resourceDigitalOceanFloatingIp(),
			"digitalocean_floating_ip_assignment": resourceDigitalOceanFloatingIpAssignment(),
			"digitalocean_loadbalancer":           resourceDigitalOceanLoadbalancer(),
			"digitalocean_record":                 resourceDigitalOceanRecord(),
			"digitalocean_ssh_key":                resourceDigitalOceanSSHKey(),
			"digitalocean_tag":                    resourceDigitalOceanTag(),
			"digitalocean_volume":                 resourceDigitalOceanVolume(),
			"digitalocean_volume_attachment":      resourceDigitalOceanVolumeAttachment(),
			"digitalocean_volume_snapshot":        resourceDigitalOceanVolumeSnapshot(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token:     d.Get("token").(string),
		AccessID:  d.Get("access_id").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	return config.Client()
}
