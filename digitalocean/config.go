package digitalocean

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"golang.org/x/oauth2"
)

type Config struct {
	Token             string
	APIEndpoint       string
	SpacesAPIEndpoint string
	AccessID          string
	SecretKey         string
	TerraformVersion  string
}

type CombinedConfig struct {
	client                 *godo.Client
	spacesEndpointTemplate *template.Template
	accessID               string
	secretKey              string
}

func (c *CombinedConfig) godoClient() *godo.Client { return c.client }

func (c *CombinedConfig) spacesClient(region string) (*session.Session, error) {
	if c.accessID == "" || c.secretKey == "" {
		err := fmt.Errorf("Spaces credentials not configured")
		return &session.Session{}, err
	}

	endpointWriter := strings.Builder{}
	err := c.spacesEndpointTemplate.Execute(&endpointWriter, map[string]string{
		"Region": strings.ToLower(region),
	})
	if err != nil {
		return &session.Session{}, err
	}
	endpoint := endpointWriter.String()

	client, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(c.accessID, c.secretKey, ""),
		Endpoint:    aws.String(endpoint)},
	)
	if err != nil {
		return &session.Session{}, err
	}

	return client, nil
}

// Client() returns a new client for accessing digital ocean.
func (c *Config) Client() (*CombinedConfig, error) {
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.Token,
	})

	userAgent := fmt.Sprintf("Terraform/%s", c.TerraformVersion)
	client := oauth2.NewClient(oauth2.NoContext, tokenSrc)

	client.Transport = logging.NewTransport("DigitalOcean", client.Transport)

	godoClient, err := godo.New(client, godo.SetUserAgent(userAgent))
	if err != nil {
		return nil, err
	}

	apiURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}
	godoClient.BaseURL = apiURL

	spacesEndpointTemplate, err := template.New("spaces").Parse(c.SpacesAPIEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse spaces_endpoint '%s' as template: %s", c.SpacesAPIEndpoint, err)
	}

	log.Printf("[INFO] DigitalOcean Client configured for URL: %s", godoClient.BaseURL.String())

	return &CombinedConfig{
		client:                 godoClient,
		spacesEndpointTemplate: spacesEndpointTemplate,
		accessID:               c.AccessID,
		secretKey:              c.SecretKey,
	}, nil
}

// waitForAction waits for the action to finish using the resource.StateChangeConf.
func waitForAction(client *godo.Client, action *godo.Action) error {
	var (
		pending   = "in-progress"
		target    = "completed"
		refreshfn = func() (result interface{}, state string, err error) {
			a, _, err := client.Actions.Get(context.Background(), action.ID)
			if err != nil {
				return nil, "", err
			}
			if a.Status == "errored" {
				return a, "errored", nil
			}
			if a.CompletedAt != nil {
				return a, target, nil
			}
			return a, pending, nil
		}
	)
	_, err := (&resource.StateChangeConf{
		Pending: []string{pending},
		Refresh: refreshfn,
		Target:  []string{target},

		Delay:      10 * time.Second,
		Timeout:    60 * time.Minute,
		MinTimeout: 3 * time.Second,

		// This is a hack around DO API strangeness.
		// https://github.com/hashicorp/terraform/issues/481
		//
		NotFoundChecks: 60,
	}).WaitForState()
	return err
}

func isDigitalOceanError(err error, code int, message string) bool {
	if err, ok := err.(*godo.ErrorResponse); ok {
		return err.Response.StatusCode == code &&
			strings.Contains(strings.ToLower(err.Message), strings.ToLower(message))
	}
	return false
}
