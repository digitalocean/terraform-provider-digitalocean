package database_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Database Engine Support Matrix for Logsinks:
// - PostgreSQL, MySQL, Kafka, Valkey: support rsyslog, opensearch logsinks
// - MongoDB: supports ONLY datadog logsinks (not opensearch or rsyslog)
//
// These tests cover rsyslog logsink functionality for supported engines only.

// TestAccDigitalOceanDatabaseLogsinkRsyslog_Basic tests creating a basic rsyslog logsink
// with default settings (TLS disabled, RFC5424 format). Expected: successful creation.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_Basic(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigBasic, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					testAccCheckDigitalOceanDatabaseLogsinkAttributes(&logsink, logsinkName, "rsyslog"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "name", logsinkName),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "server", "192.168.1.100"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "port", "514"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "format", "rfc5424"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_rsyslog.test", "cluster_id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_rsyslog.test", "logsink_id"),
					resource.TestCheckResourceAttrSet("digitalocean_database_logsink_rsyslog.test", "id"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkRsyslog_Update tests updating an rsyslog logsink
// configuration (port, TLS enabled, format change, structured data). Expected: successful update.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_Update(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigBasic, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "port", "514"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "format", "rfc5424"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigUpdated, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "port", "1514"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "tls", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "format", "rfc3164"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "structured_data", "[example@41058 iut=\"3\"]"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkRsyslog_CustomFormat tests creating an rsyslog logsink
// with custom format and logline template. Expected: successful creation with custom logline.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_CustomFormat(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigCustom, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "format", "custom"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "logline", "%timestamp% %HOSTNAME% %app-name% %procid% %msgid% %msg%"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkRsyslog_TLS tests creating an rsyslog logsink
// with TLS enabled and CA certificate. Expected: successful creation with TLS configuration.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_TLS(t *testing.T) {
	var logsink godo.DatabaseLogsink
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigTLS, clusterName, logsinkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseLogsinkExists("digitalocean_database_logsink_rsyslog.test", &logsink),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "tls", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_logsink_rsyslog.test", "ca_cert", "-----BEGIN CERTIFICATE-----\nMIIDCTCCAfGgAwIBAgIUdh0W7W79ns0Gc+6ZylC6JpCrF50wDQYJKoZIhvcNAQEL\nBQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDkxODE3NTMzNloXDTI2MDkx\nODE3NTMzNlowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF\nAAOCAQ8AMIIBCgKCAQEAsBxZyNgjCWqsDE6h5sfMZo1JfD3WFzGZN2XdPwaPDPH9\nGI6UokJbhdJXPhFPyKmXis8vRC7Dos434lCp6RuYEHYk27wBam2pZSAi/P+Be5EU\nbdJdRjikPtu31JVsbZ2ookIc9zfBxPbXd5F4wNlcUFRATv2LC2SFQ91l5fmuiThU\nXx8+0Prls1Jzuz3Ll/oLM+1vxQEZFWvZCcq4HPFyf0p5Y37alxyVGSQxOqnQW3Wu\nhxNVdMKbfhx50B9Kh62LZ4+Pcv06/ftReeIV7+lO+8/FQs1BsjbLlpsIsuXgueR5\nahfOMQ/3/Wu5sb7jN3o6DINjpBmGW8zItWnIiTm8CQIDAQABo1MwUTAdBgNVHQ4E\nFgQUmY7HILyhR4RiRKFkyDyT/7fXLRMwHwYDVR0jBBgwFoAUmY7HILyhR4RiRKFk\nyDyT/7fXLRMwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAP/wy\neDjrbAMgeuTUB0DisfkUZo2RKY/hJ9+9lH9VjTQ1foomWr7J8HUJHh7Co1n8Tnjd\n0dAl1agRY0o3VZrASj3gyYWFumbe6BBjhIynzZK3rsP9BzFvl8+xNUS9jkWiFhYU\n5x9f3YzMxXQsRf6sRSfS7/IIF8SCeOZTCJIVMB8l+8XbxsoYpTKz9sG+Opg7LD2K\nFbWGBKiSbxB6SKjax0Fk0MHO07ehjOqlxqns/a78w2AsBNKc2SDv73eXv24dRzJS\nlJu7YXccTSWs2/Y+wDxTMyp3DlJ9kzkgTveXhmKJdhKW8L8a+K1hzNGBrczJeHnm\nCwPzEPg7ca5lXYLDEA==\n-----END CERTIFICATE-----"),
				),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkRsyslog_MongoDB_ShouldFail tests that creating an rsyslog logsink
// with a MongoDB cluster fails due to engine incompatibility. Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_MongoDB_ShouldFail(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigMongoDB, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("log sink not supported"),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkRsyslog_InvalidPort tests validation for invalid port values.
// Uses port 0 which is outside the valid range (1-65535). Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_InvalidPort(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigInvalidPort, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("must be between 1 and 65535"),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkRsyslog_CustomFormatRequiresLogline tests validation for custom format.
// Uses format "custom" without providing logline field. Expected: validation error requiring logline.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_CustomFormatRequiresLogline(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigCustomNoLogline, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("logline is required when format is 'custom'"),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkRsyslog_InvalidFormat tests validation for invalid format values.
// Uses "invalid_format" which is not in allowed values (rfc5424, rfc3164, custom). Expected: validation error.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_InvalidFormat(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigInvalidFormat, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("must be one of: rfc5424, rfc3164, custom"),
			},
		},
	})
}

// TestAccDigitalOceanDatabaseLogsinkRsyslog_CertWithoutTLS tests validation for certificate fields without TLS.
// Provides ca_cert but leaves TLS disabled (default false). Expected: validation error requiring TLS.
func TestAccDigitalOceanDatabaseLogsinkRsyslog_CertWithoutTLS(t *testing.T) {
	clusterName := acceptance.RandomTestName()
	logsinkName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseLogsinkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigCertWithoutTLS, clusterName, logsinkName),
				ExpectError: regexp.MustCompile("tls must be true when ca_cert, client_cert, or client_key is set"),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseLogsinkExists(resource string, logsink *godo.DatabaseLogsink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		// Parse composite ID
		clusterID := rs.Primary.Attributes["cluster_id"]
		logsinkID := rs.Primary.Attributes["logsink_id"]

		foundLogsink, _, err := client.Databases.GetLogsink(context.Background(), clusterID, logsinkID)
		if err != nil {
			return err
		}

		if foundLogsink.ID != logsinkID {
			return fmt.Errorf("Record not found")
		}

		*logsink = *foundLogsink

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseLogsinkAttributes(logsink *godo.DatabaseLogsink, name, sinkType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if logsink.Name != name {
			return fmt.Errorf("Bad name: %s", logsink.Name)
		}

		if logsink.Type != sinkType {
			return fmt.Errorf("Bad sink type: %s", logsink.Type)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseLogsinkDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_logsink_rsyslog" && rs.Type != "digitalocean_database_logsink_opensearch" {
			continue
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		logsinkID := rs.Primary.Attributes["logsink_id"]

		_, _, err := client.Databases.GetLogsink(context.Background(), clusterID, logsinkID)
		if err == nil {
			return fmt.Errorf("DatabaseLogsink still exists")
		}
	}

	return nil
}

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigBasic = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  tls        = false
  format     = "rfc5424"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigUpdated = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id      = digitalocean_database_cluster.test.id
  name            = "%s"
  server          = "192.168.1.100"
  port            = 1514
  tls             = true
  format          = "rfc3164"
  structured_data = "[example@41058 iut=\"3\"]"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigCustom = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  tls        = false
  format     = "custom"
  logline    = "%%timestamp%% %%HOSTNAME%% %%app-name%% %%procid%% %%msgid%% %%msg%%"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigTLS = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  tls        = true
  ca_cert    = "-----BEGIN CERTIFICATE-----\nMIIDCTCCAfGgAwIBAgIUdh0W7W79ns0Gc+6ZylC6JpCrF50wDQYJKoZIhvcNAQEL\nBQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDkxODE3NTMzNloXDTI2MDkx\nODE3NTMzNlowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF\nAAOCAQ8AMIIBCgKCAQEAsBxZyNgjCWqsDE6h5sfMZo1JfD3WFzGZN2XdPwaPDPH9\nGI6UokJbhdJXPhFPyKmXis8vRC7Dos434lCp6RuYEHYk27wBam2pZSAi/P+Be5EU\nbdJdRjikPtu31JVsbZ2ookIc9zfBxPbXd5F4wNlcUFRATv2LC2SFQ91l5fmuiThU\nXx8+0Prls1Jzuz3Ll/oLM+1vxQEZFWvZCcq4HPFyf0p5Y37alxyVGSQxOqnQW3Wu\nhxNVdMKbfhx50B9Kh62LZ4+Pcv06/ftReeIV7+lO+8/FQs1BsjbLlpsIsuXgueR5\nahfOMQ/3/Wu5sb7jN3o6DINjpBmGW8zItWnIiTm8CQIDAQABo1MwUTAdBgNVHQ4E\nFgQUmY7HILyhR4RiRKFkyDyT/7fXLRMwHwYDVR0jBBgwFoAUmY7HILyhR4RiRKFk\nyDyT/7fXLRMwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAP/wy\neDjrbAMgeuTUB0DisfkUZo2RKY/hJ9+9lH9VjTQ1foomWr7J8HUJHh7Co1n8Tnjd\n0dAl1agRY0o3VZrASj3gyYWFumbe6BBjhIynzZK3rsP9BzFvl8+xNUS9jkWiFhYU\n5x9f3YzMxXQsRf6sRSfS7/IIF8SCeOZTCJIVMB8l+8XbxsoYpTKz9sG+Opg7LD2K\nFbWGBKiSbxB6SKjax0Fk0MHO07ehjOqlxqns/a78w2AsBNKc2SDv73eXv24dRzJS\nlJu7YXccTSWs2/Y+wDxTMyp3DlJ9kzkgTveXhmKJdhKW8L8a+K1hzNGBrczJeHnm\nCwPzEPg7ca5lXYLDEA==\n-----END CERTIFICATE-----"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigMongoDB = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "mongodb"
  version    = "7"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  format     = "rfc5424"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigInvalidPort = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 0
  tls        = false
  format     = "rfc5424"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigCustomNoLogline = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  tls        = false
  format     = "custom"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigInvalidFormat = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  tls        = false
  format     = "invalid_format"
}`

const testAccCheckDigitalOceanDatabaseLogsinkRsyslogConfigCertWithoutTLS = `
resource "digitalocean_database_cluster" "test" {
  name       = "%s"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_logsink_rsyslog" "test" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "%s"
  server     = "192.168.1.100"
  port       = 514
  format     = "rfc5424"
  ca_cert    = "-----BEGIN CERTIFICATE-----\nMIIDCTCCAfGgAwIBAgIUdh0W7W79ns0Gc+6ZylC6JpCrF50wDQYJKoZIhvcNAQEL\nBQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDkxODE3NTMzNloXDTI2MDkx\nODE3NTMzNlowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF\nAAOCAQ8AMIIBCgKCAQEAsBxZyNgjCWqsDE6h5sfMZo1JfD3WFzGZN2XdPwaPDPH9\nGI6UokJbhdJXPhFPyKmXis8vRC7Dos434lCp6RuYEHYk27wBam2pZSAi/P+Be5EU\nbdJdRjikPtu31JVsbZ2ookIc9zfBxPbXd5F4wNlcUFRATv2LC2SFQ91l5fmuiThU\nXx8+0Prls1Jzuz3Ll/oLM+1vxQEZFWvZCcq4HPFyf0p5Y37alxyVGSQxOqnQW3Wu\nhxNVdMKbfhx50B9Kh62LZ4+Pcv06/ftReeIV7+lO+8/FQs1BsjbLlpsIsuXgueR5\nahfOMQ/3/Wu5sb7jN3o6DINjpBmGW8zItWnIiTm8CQIDAQABo1MwUTAdBgNVHQ4E\nFgQUmY7HILyhR4RiRKFkyDyT/7fXLRMwHwYDVR0jBBgwFoAUmY7HILyhR4RiRKFk\nyDyT/7fXLRMwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAP/wy\neDjrbAMgeuTUB0DisfkUZo2RKY/hJ9+9lH9VjTQ1foomWr7J8HUJHh7Co1n8Tnjd\n0dAl1agRY0o3VZrASj3gyYWFumbe6BBjhIynzZK3rsP9BzFvl8+xNUS9jkWiFhYU\n5x9f3YzMxXQsRf6sRSfS7/IIF8SCeOZTCJIVMB8l+8XbxsoYpTKz9sG+Opg7LD2K\nFbWGBKiSbxB6SKjax0Fk0MHO07ehjOqlxqns/a78w2AsBNKc2SDv73eXv24dRzJS\nlJu7YXccTSWs2/Y+wDxTMyp3DlJ9kzkgTveXhmKJdhKW8L8a+K1hzNGBrczJeHnm\nCwPzEPg7ca5lXYLDEA==\n-----END CERTIFICATE-----"
}`
