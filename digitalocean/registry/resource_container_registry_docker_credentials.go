package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	expirySecondsDefault     = 1576800000 // Max allowed by the API, roughly 50 years
	oauthTokenRevokeEndpoint = "https://cloud.digitalocean.com/v1/oauth/revoke"
)

func ResourceDigitalOceanContainerRegistryDockerCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanContainerRegistryDockerCredentialsCreate,
		ReadContext:   resourceDigitalOceanContainerRegistryDockerCredentialsRead,
		UpdateContext: resourceDigitalOceanContainerRegistryDockerCredentialsUpdate,
		DeleteContext: resourceDigitalOceanContainerRegistryDockerCredentialsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      expirySecondsDefault,
				ValidateFunc: validation.IntBetween(0, expirySecondsDefault),
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

func resourceDigitalOceanContainerRegistryDockerCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDigitalOceanContainerRegistryDockerCredentialsRead(ctx, d, meta)
}

func resourceDigitalOceanContainerRegistryDockerCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	reg, response, err := client.Registry.Get(context.Background())

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return diag.Errorf("registry not found: %s", err)
		}
		return diag.Errorf("Error retrieving registry: %s", err)
	}

	write := d.Get("write").(bool)
	d.SetId(reg.Name)
	d.Set("registry_name", reg.Name)
	d.Set("write", write)

	err = updateExpiredDockerCredentials(d, write, client)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDigitalOceanContainerRegistryDockerCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("expiry_seconds") {
		write := d.Get("write").(bool)
		expirySeconds := d.Get("expiry_seconds").(int)
		client := meta.(*config.CombinedConfig).GodoClient()
		currentTime := time.Now().UTC()
		expirationTime := currentTime.Add(time.Second * time.Duration(expirySeconds))
		d.Set("credential_expiration_time", expirationTime.Format(time.RFC3339))
		dockerConfigJSON, err := generateDockerCredentials(write, expirySeconds, client)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("write", write)
		d.Set("docker_credentials", dockerConfigJSON)
	} else {
		if d.HasChange("write") {
			write := d.Get("write").(bool)
			expirySeconds := d.Get("expiry_seconds").(int)
			client := meta.(*config.CombinedConfig).GodoClient()
			dockerConfigJSON, err := generateDockerCredentials(write, expirySeconds, client)
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set("write", write)
			d.Set("docker_credentials", dockerConfigJSON)
		}
	}

	return nil
}

type dockerConfig struct {
	Auths struct {
		Registry struct {
			Auth string `json:"auth"`
		} `json:"registry.digitalocean.com"`
	} `json:"auths"`
}

func resourceDigitalOceanContainerRegistryDockerCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configJSON := d.Get("docker_credentials")
	var config dockerConfig
	err := json.Unmarshal([]byte(configJSON.(string)), &config)
	if err != nil {
		return diag.FromErr(err)
	}

	// The OAuth token is used for both the username and password
	// and stored as a base64 encoded string.
	decoded, err := base64.StdEncoding.DecodeString(config.Auths.Registry.Auth)
	if err != nil {
		return diag.FromErr(err)
	}
	tokens := strings.Split(string(decoded), ":")
	if len(tokens) != 2 {
		return diag.FromErr(errors.New("unable to find OAuth token"))
	}

	err = RevokeOAuthToken(tokens[0], oauthTokenRevokeEndpoint)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func RevokeOAuthToken(token string, endpoint string) error {
	data := url.Values{}
	data.Set("token", token)
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}

	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return errors.New("error revoking token: " + http.StatusText(resp.StatusCode))
	}

	return err
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
