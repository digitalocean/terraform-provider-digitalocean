package digitalocean

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (interface{}, error) {
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

	config := Config{
		Token:             os.Getenv("DIGITALOCEAN_TOKEN"),
		APIEndpoint:       apiEndpoint,
		SpacesAPIEndpoint: spacesEndpoint,
	}

	// configures a default client for the region, using the above env vars
	client, err := config.Client()
	if err != nil {
		return nil, fmt.Errorf("error getting DigitalOcean client")
	}

	return client, nil
}
