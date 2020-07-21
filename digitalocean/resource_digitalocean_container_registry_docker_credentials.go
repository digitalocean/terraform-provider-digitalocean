package digitalocean

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const expirySecondsDefault = 1576800000 // Max allowed by the API, roughly 50 years

func resourceDigitalOceanContainerRegistryDockerCredentials() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanContainerRegistryDockerCredentialsCreate,
		Read:   resourceDigitalOceanContainerRegistryDockerCredentialsRead,
		Update: resourceDigitalOceanContainerRegistryDockerCredentialsUpdate,
		Delete: resourceDigitalOceanContainerRegistryDockerCredentialsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"registry_name": {
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
				Default:  expirySecondsDefault,
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

func resourceDigitalOceanContainerRegistryDockerCredentialsCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceDigitalOceanContainerRegistryDockerCredentialsRead(d, meta)
}

func resourceDigitalOceanContainerRegistryDockerCredentialsRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("registry_name", reg.Name)
	d.Set("write", write)

	err = updateExpiredDockerCredentials(d, write, client)
	if err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanContainerRegistryDockerCredentialsUpdate(d *schema.ResourceData, meta interface{}) error {
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

func resourceDigitalOceanContainerRegistryDockerCredentialsDelete(d *schema.ResourceData, meta interface{}) error {
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

func updateExpiredDockerCredentials(d *schema.ResourceData, readWrite bool, client *godo.Client) error {
	expirySeconds := d.Get("expiry_seconds").(int)
	expirationTime := d.Get("credential_expiration_time").(string)

	if (expirySeconds > expirySecondsDefault) || (expirySeconds <= 0) {
		return fmt.Errorf("expiry_seconds outside acceptable range")
	}

	d.Set("expiry_seconds", expirySeconds)

	currentTime := time.Now().UTC()
	if expirationTime != "" {
		expirationTime, err := time.Parse(time.RFC3339, expirationTime)
		if err != nil {
			return err
		}

		if expirationTime.Before(currentTime) {
			dockerConfigJSON, err := generateDockerCredentials(readWrite, expirySeconds, client)
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
		dockerConfigJSON, err := generateDockerCredentials(readWrite, expirySeconds, client)
		if err != nil {
			return err
		}
		d.Set("docker_credentials", dockerConfigJSON)
	}
	return nil
}
