package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanContainerRegistryCreate,
		Read:   resourceDigitalOceanContainerRegistryRead,
		Update: resourceDigitalOceanContainerRegistryUpdate,
		Delete: resourceDigitalOceanContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"write": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_url": {
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

func resourceDigitalOceanContainerRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	// Build up our creation options
	opts := &godo.RegistryCreateRequest{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Container Registry create configuration: %#v", opts)
	reg, _, err := client.Registry.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating container registry: %s", err)
	}

	d.SetId(reg.Name)
	log.Printf("[INFO] Container Registry: %s", reg.Name)

	return resourceDigitalOceanContainerRegistryRead(d, meta)
}

func resourceDigitalOceanContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	reg, resp, err := client.Registry.Get(context.Background())
	if err != nil {
		// If the registry is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving container registry: %s", err)
	}

	write := d.Get("write").(bool)
	d.SetId(reg.Name)
	d.Set("name", reg.Name)
	d.Set("write", write)
	d.Set("endpoint", fmt.Sprintf("%s/%s", RegistryHostname, reg.Name))
	dockerConfigJSON, err := generateDockerCredentials(write, client)
	if err != nil {
		return err
	}
	d.Set("docker_credentials", dockerConfigJSON)
	d.Set("server_url", RegistryHostname)

	return nil
}

func resourceDigitalOceanContainerRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("write") {
		write := d.Get("write").(bool)
		client := meta.(*CombinedConfig).godoClient()
		dockerConfigJSON, err := generateDockerCredentials(write, client)
		if err != nil {
			return err
		}
		d.Set("write", write)
		d.Set("docker_credentials", dockerConfigJSON)
	}
	return nil
}

func resourceDigitalOceanContainerRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting container registry: %s", d.Id())
	_, err := client.Registry.Delete(context.Background())
	if err != nil {
		return fmt.Errorf("Error deleting container registry: %s", err)
	}

	d.SetId("")
	return nil
}

func generateDockerCredentials(readWrite bool, client *godo.Client) (string, error) {
	dockerCreds, response, err := client.Registry.DockerCredentials(context.Background(), &godo.RegistryDockerCredentialsRequest{ReadWrite: readWrite})
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return "", fmt.Errorf("docker credentials not found: %s", err)
		}
		return "", fmt.Errorf("Error retrieving docker credentials: %s", err)
	}
	dockerConfigJSON := string(dockerCreds.DockerConfigJSON)
	return dockerConfigJSON, nil
}