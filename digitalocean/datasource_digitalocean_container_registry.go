package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const RegistryHostname = "registry.digitalocean.com"

func dataSourceDigitalOceanContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanContainerRegistryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the container registry",
				ValidateFunc: validation.NoZeroValues,
			},
			//TODO: Need better name?
			"write": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"docker_credentials": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceDigitalOceanContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	reg, response, err := client.Registry.Get(context.Background())

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return fmt.Errorf("registry not found: %s", err)
		}
		return fmt.Errorf("Error retrieving registry: %s", err)
	}

	write := d.Get("write").(bool)
	d.SetId(reg.Name)
	d.Set("name", reg.Name)
	d.Set("endpoint", fmt.Sprintf("%s/%s", RegistryHostname, reg.Name))
	dockerCreds, response, err := client.Registry.DockerCredentials(context.Background(), &godo.RegistryDockerCredentialsRequest{ReadWrite: write})
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return fmt.Errorf("docker credentials not found: %s", err)
		}
		return fmt.Errorf("Error retrieving docker credentials: %s", err)
	}
	dockerConfigJSON := string(dockerCreds.DockerConfigJSON)
	// TODO: Do we need this
	if dockerConfigJSON == "" {
		return fmt.Errorf("Empty docker credentials")
	}
	d.Set("docker_credentials", string(dockerCreds.DockerConfigJSON))
	return nil
}
