package digitalocean

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const expirySecondsDefault = 2147483647 // Max value of signed 32 bit integer

func resourceDigitalOceanContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanContainerRegistryCreate,
		Read:   resourceDigitalOceanContainerRegistryRead,
		Update: resourceDigitalOceanContainerRegistryUpdate,
		Delete: resourceDigitalOceanContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanContainerRegistryImport,
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
			"expiry_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  expirySecondsDefault, // Relatively close to max value of Duration
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
			"credential_expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
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

	expirySeconds := d.Get("expiry_seconds").(int)

	if (expirySeconds > expirySecondsDefault) || (expirySeconds <= 0) {
		return fmt.Errorf("expiry_seconds outside acceptable range")
	}

	d.Set("expiry_seconds", expirySeconds)

	expirationTime := d.Get("credential_expiration_time").(string)
	currentTime := time.Now().UTC()
	if expirationTime != "" {
		expirationTime, err := time.Parse(time.RFC3339, expirationTime)
		if err != nil {
			return err
		}

		if expirationTime.Before(currentTime) {
			dockerConfigJSON, err := generateDockerCredentials(write, expirySeconds, client)
			if err != nil {
				return err
			}
			d.Set("docker_credentials", dockerConfigJSON)
			expirationTime := currentTime.Add(time.Second * time.Duration(expirySeconds))
			d.Set("credential_expiration_time", expirationTime.Format(time.RFC3339))
		}

	} else {
		expirationTime := currentTime.Add(time.Second * time.Duration(expirySeconds))
		d.Set("credential_expiration_time", expirationTime.Format(time.RFC3339))
		dockerConfigJSON, err := generateDockerCredentials(write, expirySeconds, client)
		if err != nil {
			return err
		}
		d.Set("docker_credentials", dockerConfigJSON)
	}

	d.Set("server_url", RegistryHostname)

	return nil
}

func resourceDigitalOceanContainerRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("expiry_seconds") {
		write := d.Get("write").(bool)
		expirySeconds := d.Get("expiry_seconds").(int)
		client := meta.(*CombinedConfig).godoClient()
		currentTime := time.Now().UTC()
		expirationTime := currentTime.Add(time.Second * time.Duration(expirySeconds))
		d.Set("credential_expiration_time", expirationTime.Format(time.RFC3339))
		dockerConfigJSON, err := generateDockerCredentials(write, expirySeconds, client)
		if err != nil {
			return err
		}
		d.Set("write", write)
		d.Set("docker_credentials", dockerConfigJSON)
	} else {
		if d.HasChange("write") {
			write := d.Get("write").(bool)
			expirySeconds := d.Get("expiry_seconds").(int)
			client := meta.(*CombinedConfig).godoClient()
			dockerConfigJSON, err := generateDockerCredentials(write, expirySeconds, client)
			if err != nil {
				return err
			}
			d.Set("write", write)
			d.Set("docker_credentials", dockerConfigJSON)
		}
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

func generateDockerCredentials(readWrite bool, expirySeconds int, client *godo.Client) (string, error) {
	dockerCreds, response, err := client.Registry.DockerCredentials(context.Background(), &godo.RegistryDockerCredentialsRequest{ReadWrite: readWrite, ExpirySeconds: &expirySeconds})
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return "", fmt.Errorf("docker credentials not found: %s", err)
		}
		return "", fmt.Errorf("Error retrieving docker credentials: %s", err)
	}
	dockerConfigJSON := string(dockerCreds.DockerConfigJSON)
	return dockerConfigJSON, nil
}

func resourceDigitalOceanContainerRegistryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if v, ok := d.GetOk("expiry_seconds"); ok {
		d.Set("expiry_seconds", v.(int))
	} else {
		d.Set("expiry_seconds", expirySecondsDefault)
	}

	return []*schema.ResourceData{d}, nil
}
