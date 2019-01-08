package digitalocean

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"golang.org/x/oauth2"
)

type Config struct {
	Token     string
	AccessID  string
	SecretKey string
}

type CombinedConfig struct {
	client    *godo.Client
	accessID  string
	secretKey string
}

func (c *CombinedConfig) godoClient() *godo.Client { return c.client }

func (c *CombinedConfig) spacesClient(region string) (*session.Session, error) {
	if c.accessID == "" || c.secretKey == "" {
		err := fmt.Errorf("Spaces credentials not configured")
		return &session.Session{}, err
	}

	endpoint := fmt.Sprintf("https://%s.digitaloceanspaces.com", region)
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

	userAgent := fmt.Sprintf("Terraform/%s", terraform.VersionString())

	do, err := godo.New(oauth2.NewClient(oauth2.NoContext, tokenSrc), godo.SetUserAgent(userAgent))
	if err != nil {
		return nil, err
	}

	if logging.IsDebugOrHigher() {
		do.OnRequestCompleted(logRequestAndResponse)
	}

	log.Printf("[INFO] DigitalOcean Client configured for URL: %s", do.BaseURL.String())

	return &CombinedConfig{
		client:    do,
		accessID:  c.AccessID,
		secretKey: c.SecretKey,
	}, nil
}

func logRequestAndResponse(req *http.Request, resp *http.Response) {
	reqData, err := httputil.DumpRequest(req, true)
	if err == nil {
		log.Printf("[DEBUG] "+logReqMsg, string(reqData))
	} else {
		log.Printf("[ERROR] DigitalOcean API Request error: %#v", err)
	}

	respData, err := httputil.DumpResponse(resp, true)
	if err == nil {
		log.Printf("[DEBUG] "+logRespMsg, string(respData))
	} else {
		log.Printf("[ERROR] DigitalOcean API Response error: %#v", err)
	}
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

const logReqMsg = `DigitalOcean API Request Details:
---[ REQUEST ]---------------------------------------
%s
-----------------------------------------------------`

const logRespMsg = `DigitalOcean API Response Details:
---[ RESPONSE ]--------------------------------------
%s
-----------------------------------------------------`
