package config

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"golang.org/x/oauth2"
)

// VM API base URL (hardcoded for now)
const vmAPIBaseURL = "http://167.71.171.33:8080"

type Config struct {
	Token             string
	APIEndpoint       string
	SpacesAPIEndpoint string
	AccessID          string
	SecretKey         string
	RequestsPerSecond float64
	TerraformVersion  string
	HTTPRetryMax      int
	HTTPRetryWaitMax  float64
	HTTPRetryWaitMin  float64
}

type CombinedConfig struct {
	client                 *godo.Client
	vmClient               *godo.Client // Client for VM API with different base URL
	spacesEndpointTemplate *template.Template
	accessID               string
	secretKey              string
}

func (c *CombinedConfig) GodoClient() *godo.Client { return c.client }

func (c *CombinedConfig) VMClient() *godo.Client { return c.vmClient }

func (c *CombinedConfig) SpacesClient(region string) (*session.Session, error) {
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
	var client *http.Client
	var godoOpts []godo.ClientOpt

	client = oauth2.NewClient(context.Background(), tokenSrc)

	if c.HTTPRetryMax > 0 {
		retryConfig := godo.RetryConfig{
			RetryMax:     c.HTTPRetryMax,
			RetryWaitMin: godo.PtrTo(c.HTTPRetryWaitMin),
			RetryWaitMax: godo.PtrTo(c.HTTPRetryWaitMax),
			Logger:       log.Default(),
		}

		godoOpts = []godo.ClientOpt{godo.WithRetryAndBackoffs(retryConfig)}
	}

	godoOpts = append(godoOpts, godo.SetUserAgent(userAgent))

	if c.RequestsPerSecond > 0.0 {
		godoOpts = append(godoOpts, godo.SetStaticRateLimit(c.RequestsPerSecond))
	}

	godoClient, err := godo.New(client, godoOpts...)

	// TODO: logging.NewTransport is deprecated and should be replaced with
	// logging.NewTransportWithRequestLogging.
	//
	//nolint:staticcheck
	clientTransport := logging.NewTransport("DigitalOcean", godoClient.HTTPClient.Transport)

	godoClient.HTTPClient.Transport = clientTransport

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

	// Create VM API client with same auth but different base URL
	vmGodoClient, err := godo.New(client, godoOpts...)
	if err != nil {
		return nil, err
	}
	vmAPIURL, err := url.Parse(vmAPIBaseURL)
	if err != nil {
		return nil, err
	}
	vmGodoClient.BaseURL = vmAPIURL

	// Add transport logging for VM client
	//nolint:staticcheck
	vmClientTransport := logging.NewTransport("DigitalOceanVM", vmGodoClient.HTTPClient.Transport)
	vmGodoClient.HTTPClient.Transport = vmClientTransport

	log.Printf("[INFO] DigitalOcean VM Client configured for URL: %s", vmGodoClient.BaseURL.String())

	return &CombinedConfig{
		client:                 godoClient,
		vmClient:               vmGodoClient,
		spacesEndpointTemplate: spacesEndpointTemplate,
		accessID:               c.AccessID,
		secretKey:              c.SecretKey,
	}, nil
}
