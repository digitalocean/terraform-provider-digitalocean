package database_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Database Engine Support Matrix for Logsinks:
// - PostgreSQL, MySQL, Kafka, Valkey: support rsyslog, opensearch logsinks
// - MongoDB: supports ONLY datadog logsinks (not opensearch or rsyslog)
//
// These tests cover opensearch logsink functionality for supported engines only.

// TestAccDigitalOceanDatabaseLogsinkOpensearch_Basic tests creating a basic opensearch logsink
// with required fields (endpoint, index_prefix, index_days_max). Expected: successful creation.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_Basic(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigBasic, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					testAccCheckDigitalOceanDatabaseLogsinkAttributes(&logsink, logsinkName, "opensearch"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "name", logsinkName),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "endpoint", "https://opensearch.example.com:9200"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_prefix", "db-logs"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_days_max", "7"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "timeout_seconds", "10"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_opensearch.test", "cluster_id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_opensearch.test", "logsink_id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_opensearch.test", "id"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_Update tests updating an opensearch logsink
// configuration (index_prefix, index_days_max, timeout_seconds). Expected: successful update.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_Update(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigBasic, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_prefix", "db-logs"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_days_max", "7"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigUpdated, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_prefix", "updated-logs"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "index_days_max", "14"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "timeout_seconds", "30"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_WithCA tests creating an opensearch logsink
// with CA certificate for TLS verification. Expected: successful creation with certificate.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_WithCA(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigWithCA, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_opensearch.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_opensearch.test", "ca_cert", "-----BEGIN CERTIFICATE-----\nMIIDCTCCAfGgAwIBAgIUdh0W7W79ns0Gc+6ZylC6JpCrF50wDQYJKoZIhvcNAQEL\nBQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDkxODE3NTMzNloXDTI2MDkx\nODE3NTMzNlowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF\nAAOCAQ8AMIIBCgKCAQEAsBxZyNgjCWqsDE6h5sfMZo1JfD3WFzGZN2XdPwaPDPH9\nGI6UokJbhdJXPhFPyKmXis8vRC7Dos434lCp6RuYEHYk27wBam2pZSAi/P+Be5EU\nbdJdRjikPtu31JVsbZ2ookIc9zfBxPbXd5F4wNlcUFRATv2LC2SFQ91l5fmuiThU\nXx8+0Prls1Jzuz3Ll/oLM+1vxQEZFWvZCcq4HPFyf0p5Y37alxyVGSQxOqnQW3Wu\nhxNVdMKbfhx50B9Kh62LZ4+Pcv06/ftReeIV7+lO+8/FQs1BsjbLlpsIsuXgueR5\nahfOMQ/3/Wu5sb7jN3o6DINjpBmGW8zItWnIiTm8CQIDAQABo1MwUTAdBgNVHQ4E\nFgQUmY7HILyhR4RiRKFkyDyT/7fXLRMwHwYDVR0jBBgwFoAUmY7HILyhR4RiRKFk\nyDyT/7fXLRMwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAP/wy\neDjrbAMgeuTUB0DisfkUZo2RKY/hJ9+9lH9VjTQ1foomWr7J8HUJHh7Co1n8Tnjd\n0dAl1agRY0o3VZrASj3gyYWFumbe6BBjhIynzZK3rsP9BzFvl8+xNUS9jkWiFhYU\n5x9f3YzMxXQsRf6sRSfS7/IIF8SCeOZTCJIVMB8l+8XbxsoYpTKz9sG+Opg7LD2K\nFbWGBKiSbxB6SKjax0Fk0MHO07ehjOqlxqns/a78w2AsBNKc2SDv73eXv24dRzJS\nlJu7YXccTSWs2/Y+wDxTMyp3DlJ9kzkgTveXhmKJdhKW8L8a+K1hzNGBrczJeHnm\nCwPzEPg7ca5lXYLDEA==\n-----END CERTIFICATE-----"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_InvalidIndexDaysMax tests validation for invalid index_days_max.
// Uses value 0 which is below minimum (must be >= 1). Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_InvalidIndexDaysMax(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigInvalidIndexDays, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("must be >= 1"),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_EmptyIndexPrefix tests validation for empty index_prefix.
// Uses empty string for index_prefix which is not allowed. Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_EmptyIndexPrefix(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigEmptyIndexPrefix, clusterName, logsinkName),
				ExpectError: regexp.MustCompile(`"index_prefix" cannot be empty`),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkOpensearch_MalformedEndpoint tests validation for malformed endpoint URLs.
// Uses invalid URL format that fails scheme validation. Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkOpensearch_MalformedEndpoint(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigMalformedEndpoint, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("must use HTTPS scheme"),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigBasic = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id      = digitalocean_database_cluster.test.id
  name            = "%s"
  endpoint        = "https://opensearch.example.com:9200"
  index_prefix    = "db-logs"
  index_days_max  = 7
  timeout_seconds = 10
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigUpdated = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id      = digitalocean_database_cluster.test.id
  name            = "%s"
  endpoint        = "https://opensearch.example.com:9200"
  index_prefix    = "updated-logs"
  index_days_max  = 14
  timeout_seconds = 30
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigWithCA = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id      = digitalocean_database_cluster.test.id
  name            = "%s"
  endpoint        = "https://opensearch.example.com:9200"
  index_prefix    = "db-logs"
  index_days_max  = 7
  timeout_seconds = 10
  ca_cert         = "-----BEGIN CERTIFICATE-----\nMIIDCTCCAfGgAwIBAgIUdh0W7W79ns0Gc+6ZylC6JpCrF50wDQYJKoZIhvcNAQEL\nBQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDkxODE3NTMzNloXDTI2MDkx\nODE3NTMzNlowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF\nAAOCAQ8AMIIBCgKCAQEAsBxZyNgjCWqsDE6h5sfMZo1JfD3WFzGZN2XdPwaPDPH9\nGI6UokJbhdJXPhFPyKmXis8vRC7Dos434lCp6RuYEHYk27wBam2pZSAi/P+Be5EU\nbdJdRjikPtu31JVsbZ2ookIc9zfBxPbXd5F4wNlcUFRATv2LC2SFQ91l5fmuiThU\nXx8+0Prls1Jzuz3Ll/oLM+1vxQEZFWvZCcq4HPFyf0p5Y37alxyVGSQxOqnQW3Wu\nhxNVdMKbfhx50B9Kh62LZ4+Pcv06/ftReeIV7+lO+8/FQs1BsjbLlpsIsuXgueR5\nahfOMQ/3/Wu5sb7jN3o6DINjpBmGW8zItWnIiTm8CQIDAQABo1MwUTAdBgNVHQ4E\nFgQUmY7HILyhR4RiRKFkyDyT/7fXLRMwHwYDVR0jBBgwFoAUmY7HILyhR4RiRKFk\nyDyT/7fXLRMwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAP/wy\neDjrbAMgeuTUB0DisfkUZo2RKY/hJ9+9lH9VjTQ1foomWr7J8HUJHh7Co1n8Tnjd\n0dAl1agRY0o3VZrASj3gyYWFumbe6BBjhIynzZK3rsP9BzFvl8+xNUS9jkWiFhYU\n5x9f3YzMxXQsRf6sRSfS7/IIF8SCeOZTCJIVMB8l+8XbxsoYpTKz9sG+Opg7LD2K\nFbWGBKiSbxB6SKjax0Fk0MHO07ehjOqlxqns/a78w2AsBNKc2SDv73eXv24dRzJS\nlJu7YXccTSWs2/Y+wDxTMyp3DlJ9kzkgTveXhmKJdhKW8L8a+K1hzNGBrczJeHnm\nCwPzEPg7ca5lXYLDEA==\n-----END CERTIFICATE-----"
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigInvalidIndexDays = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id     = digitalocean_database_cluster.test.id
  name           = "%s"
  endpoint       = "https://opensearch.example.com:9200"
  index_prefix   = "db-logs"
  index_days_max = 0
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigEmptyIndexPrefix = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id     = digitalocean_database_cluster.test.id
  name           = "%s"
  endpoint       = "https://opensearch.example.com:9200"
  index_prefix   = ""
  index_days_max = 7
}`

const testAccCheckDigitalOceanDatabaseLogsinkOpensearchConfigMalformedEndpoint = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_opensearch" "test" {
  cluster_id     = digitalocean_database_cluster.test.id
  name           = "%s"
  endpoint       = "not-a-valid-url"
  index_prefix   = "db-logs"
  index_days_max = 7
}`
