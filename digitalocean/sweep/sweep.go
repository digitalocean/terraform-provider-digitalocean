package sweep

import (
	"fmt"
	"log"
	"os"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
)

const TestNamePrefix = "tf-acc-test-"

func SharedConfigForRegion(region string) (interface{}, error) {
	if os.Getenv("DIGITALOCEAN_TOKEN") == "" {
		return nil, fmt.Errorf("empty DIGITALOCEAN_TOKEN")
	}

	apiEndpoint := os.Getenv("DIGITALOCEAN_API_URL")
	if apiEndpoint == "" {
		apiEndpoint = "https://api.digitalocean.com"
	}

	spacesEndpoint := os.Getenv("SPACES_ENDPOINT_URL")
	if spacesEndpoint == "" {
		spacesEndpoint = "https://{{.Region}}.digitaloceanspaces.com"
	}

	spacesAccessKeyID := os.Getenv("SPACES_ACCESS_KEY_ID")
	if spacesAccessKeyID == "" {
		log.Println("[DEBUG] SPACES_ACCESS_KEY_ID not set")
	}

	spacesSecretAccessKey := os.Getenv("SPACES_SECRET_ACCESS_KEY")
	if spacesSecretAccessKey == "" {
		log.Println("[DEBUG] SPACES_SECRET_ACCESS_KEY not set")
	}

	conf := config.Config{
		Token:             os.Getenv("DIGITALOCEAN_TOKEN"),
		APIEndpoint:       apiEndpoint,
		SpacesAPIEndpoint: spacesEndpoint,
		AccessID:          spacesAccessKeyID,
		SecretKey:         spacesSecretAccessKey,
	}

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, fmt.Errorf("error getting DigitalOcean client")
	}

	return client, nil
}
