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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"digitalocean_certificate":  dataSourceDigitalOceanCertificate(),
			"digitalocean_domain":       dataSourceDigitalOceanDomain(),
			"digitalocean_droplet":      dataSourceDigitalOceanDroplet(),
			"digitalocean_floating_ip":  dataSourceDigitalOceanFloatingIp(),
			"digitalocean_image":        dataSourceDigitalOceanImage(),
			"digitalocean_loadbalancer": dataSourceDigitalOceanLoadbalancer(),
			"digitalocean_record":       dataSourceDigitalOceanRecord(),
			"digitalocean_snapshot":     dataSourceDigitalOceanSnapshot(),
			"digitalocean_ssh_key":      dataSourceDigitalOceanSSHKey(),
			"digitalocean_tag":          dataSourceDigitalOceanTag(),
			"digitalocean_volume":       dataSourceDigitalOceanVolume(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"digitalocean_certificate":       resourceDigitalOceanCertificate(),
			"digitalocean_domain":            resourceDigitalOceanDomain(),
			"digitalocean_droplet":           resourceDigitalOceanDroplet(),
			"digitalocean_firewall":          resourceDigitalOceanFirewall(),
			"digitalocean_floating_ip":       resourceDigitalOceanFloatingIp(),
			"digitalocean_loadbalancer":      resourceDigitalOceanLoadbalancer(),
			"digitalocean_record":            resourceDigitalOceanRecord(),
			"digitalocean_snapshot":          resourceDigitalOceanSnapshot(),
			"digitalocean_ssh_key":           resourceDigitalOceanSSHKey(),
			"digitalocean_tag":               resourceDigitalOceanTag(),
			"digitalocean_volume":            resourceDigitalOceanVolume(),
			"digitalocean_volume_attachment": resourceDigitalOceanVolumeAttachment(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token: d.Get("token").(string),
	}

	return config.Client()
}
