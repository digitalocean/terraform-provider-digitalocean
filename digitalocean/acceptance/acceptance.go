package acceptance

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const TestNamePrefix = "tf-acc-test-"

var (
	TestAccProvider          *schema.Provider
	TestAccProviders         map[string]*schema.Provider
	TestAccProviderFactories map[string]func() (*schema.Provider, error)
)

func init() {
	TestAccProvider = digitalocean.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"digitalocean": TestAccProvider,
	}
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"digitalocean": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("DIGITALOCEAN_TOKEN"); v == "" {
		t.Fatal("DIGITALOCEAN_TOKEN must be set for acceptance tests")
	}

	err := TestAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}

func RandomTestName(additionalNames ...string) string {
	prefix := TestNamePrefix
	for _, n := range additionalNames {
		prefix += "-" + strings.Replace(n, " ", "_", -1)
	}
	return randomName(prefix, 10)
}

func randomName(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(length))
}
