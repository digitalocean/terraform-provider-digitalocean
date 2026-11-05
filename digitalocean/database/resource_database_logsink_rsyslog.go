package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseLogsinkRsyslog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseLogsinkRsyslogCreate,
		ReadContext:   resourceDigitalOceanDatabaseLogsinkRsyslogRead,
		UpdateContext: resourceDigitalOceanDatabaseLogsinkRsyslogUpdate,
		DeleteContext: resourceDigitalOceanDatabaseLogsinkRsyslogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanDatabaseLogsinkRsyslogImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "UUID of the source database cluster that will forward logs",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Display name for the logsink",
			},
			"server": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Hostname or IP address of the rsyslog server",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateLogsinkPort,
				Description:  "Port number for the rsyslog server (1-65535)",
			},
			"tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable TLS encryption for rsyslog connection",
			},
			"format": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "rfc5424",
				ValidateFunc: validateRsyslogFormat,
				Description:  "Log format: rfc5424, rfc3164, or custom",
			},
			"logline": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom logline template (required when format is 'custom')",
			},
			"structured_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Structured data for rsyslog",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "CA certificate for TLS verification (PEM format)",
			},
			"client_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client certificate for mTLS (PEM format)",
			},
			"client_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Client private key for mTLS (PEM format)",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Composite ID of the logsink resource",
			},
			"logsink_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API sink_id returned by DigitalOcean",
			},
		},

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("name", func(_ context.Context, old, new, meta interface{}) bool {
				// Force recreation if name changes
				return old.(string) != new.(string)
			}),
			func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return validateLogsinkCustomDiff(diff, "rsyslog")
			},
		),
	}
}

func resourceDigitalOceanDatabaseLogsinkRsyslogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	req := &godo.DatabaseCreateLogsinkRequest{
		Name:   d.Get("name").(string),
		Type:   "rsyslog",
		Config: expandLogsinkConfigRsyslog(d),
	}

	log.Printf("[DEBUG] Database logsink rsyslog create configuration: %#v", req)
	logsink, _, err := client.Databases.CreateLogsink(ctx, clusterID, req)
	if err != nil {
		return diag.Errorf("Error creating database logsink rsyslog: %s", err)
	}

	log.Printf("[DEBUG] API Response logsink: %#v", logsink)
	log.Printf("[DEBUG] Logsink ID: '%s'", logsink.ID)
	log.Printf("[DEBUG] Logsink Name: '%s'", logsink.Name)
	log.Printf("[DEBUG] Logsink Type: '%s'", logsink.Type)

	d.SetId(createLogsinkID(clusterID, logsink.ID))
	log.Printf("[INFO] Database logsink rsyslog ID: %s", logsink.ID)

	// Post-create read for consistency
	return resourceDigitalOceanDatabaseLogsinkRsyslogRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkRsyslogRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	logsink, resp, err := client.Databases.GetLogsink(ctx, clusterID, logsinkID)
	if err != nil {
		// If the logsink is somehow already destroyed, mark as successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving database logsink rsyslog: %s", err)
	}

	if logsink == nil {
		return diag.Errorf("Error retrieving database logsink rsyslog: logsink is nil")
	}

	d.Set("cluster_id", clusterID)
	d.Set("name", logsink.Name)
	d.Set("logsink_id", logsink.ID)

	if err := flattenLogsinkConfigRsyslog(d, logsink.Config); err != nil {
		return diag.Errorf("Error setting logsink resource data: %s", err)
	}

	return nil
}

func resourceDigitalOceanDatabaseLogsinkRsyslogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	req := &godo.DatabaseUpdateLogsinkRequest{
		Config: expandLogsinkConfigRsyslog(d),
	}

	log.Printf("[DEBUG] Database logsink rsyslog update configuration: %#v", req)
	_, err := client.Databases.UpdateLogsink(ctx, clusterID, logsinkID, req)
	if err != nil {
		return diag.Errorf("Error updating database logsink rsyslog: %s", err)
	}

	// Re-read the resource to refresh state
	return resourceDigitalOceanDatabaseLogsinkRsyslogRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkRsyslogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	log.Printf("[INFO] Deleting database logsink rsyslog: %s", d.Id())
	_, err := client.Databases.DeleteLogsink(ctx, clusterID, logsinkID)
	if err != nil {
		// Treat 404 as success (already removed)
		if godoErr, ok := err.(*godo.ErrorResponse); ok && godoErr.Response.StatusCode == 404 {
			log.Printf("[INFO] Database logsink rsyslog %s was already deleted", d.Id())
		} else {
			return diag.Errorf("Error deleting database logsink rsyslog: %s", err)
		}
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseLogsinkRsyslogImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Validate the import ID format
	clusterID, logsinkID := splitLogsinkID(d.Id())
	if clusterID == "" || logsinkID == "" {
		return nil, fmt.Errorf("must use the format 'cluster_id,logsink_id' for import (e.g. 'deadbeef-dead-4aa5-beef-deadbeef347d,01234567-89ab-cdef-0123-456789abcdef')")
	}

	// The Read function will handle populating all fields from the API
	return []*schema.ResourceData{d}, nil
}

// createLogsinkID creates a composite ID for logsink resources
// Format: <cluster_id>,<logsink_id>
func createLogsinkID(clusterID string, logsinkID string) string {
	return fmt.Sprintf("%s,%s", clusterID, logsinkID)
}

// splitLogsinkID splits a composite logsink ID into cluster ID and logsink ID
func splitLogsinkID(id string) (string, string) {
	parts := strings.SplitN(id, ",", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// expandLogsinkConfigRsyslog converts Terraform schema data to godo.DatabaseLogsinkConfig for rsyslog
func expandLogsinkConfigRsyslog(d *schema.ResourceData) *godo.DatabaseLogsinkConfig {
	config := &godo.DatabaseLogsinkConfig{}

	config.Server = d.Get("server").(string)
	config.Port = d.Get("port").(int)
	config.TLS = d.Get("tls").(bool)
	config.Format = d.Get("format").(string)
	if v, ok := d.GetOk("logline"); ok {
		config.Logline = v.(string)
	}
	if v, ok := d.GetOk("structured_data"); ok {
		config.SD = v.(string)
	}
	if v, ok := d.GetOk("ca_cert"); ok {
		config.CA = strings.TrimSpace(v.(string))
	}
	if v, ok := d.GetOk("client_cert"); ok {
		config.Cert = strings.TrimSpace(v.(string))
	}
	if v, ok := d.GetOk("client_key"); ok {
		config.Key = strings.TrimSpace(v.(string))
	}

	return config
}

// flattenLogsinkConfigRsyslog converts godo.DatabaseLogsinkConfig to Terraform schema data for rsyslog
func flattenLogsinkConfigRsyslog(d *schema.ResourceData, config *godo.DatabaseLogsinkConfig) error {
	if config == nil {
		return nil
	}

	if config.Server != "" {
		d.Set("server", config.Server)
	}
	if config.Port != 0 {
		d.Set("port", config.Port)
	}
	d.Set("tls", config.TLS)
	if config.Format != "" {
		d.Set("format", config.Format)
	}
	if config.Logline != "" {
		d.Set("logline", config.Logline)
	}
	if config.SD != "" {
		d.Set("structured_data", config.SD)
	}
	if config.CA != "" {
		d.Set("ca_cert", strings.TrimSpace(config.CA))
	}
	if config.Cert != "" {
		d.Set("client_cert", strings.TrimSpace(config.Cert))
	}
	if config.Key != "" {
		d.Set("client_key", strings.TrimSpace(config.Key))
	}

	return nil
}

// validateLogsinkPort validates port is in range 1-65535
func validateLogsinkPort(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(int)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be an integer", key))
		return
	}

	if v < 1 || v > 65535 {
		errs = append(errs, fmt.Errorf("%q must be between 1 and 65535", key))
	}

	return
}

// validateRsyslogFormat validates format is one of the allowed values
func validateRsyslogFormat(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be a string", key))
		return
	}

	validFormats := []string{"rfc5424", "rfc3164", "custom"}
	for _, format := range validFormats {
		if v == format {
			return
		}
	}

	errs = append(errs, fmt.Errorf("%q must be one of: %s", key, strings.Join(validFormats, ", ")))
	return
}

// validateLogsinkCustomDiff validates cross-field dependencies for rsyslog logsink resources
func validateLogsinkCustomDiff(diff *schema.ResourceDiff, sinkType string) error {
	if sinkType != "rsyslog" {
		return nil
	}

	// If format is custom, require logline
	format := diff.Get("format").(string)
	logline := diff.Get("logline").(string)

	if format == "custom" && strings.TrimSpace(logline) == "" {
		return fmt.Errorf("logline is required when format is 'custom'")
	}

	// If any TLS cert fields are set, require tls = true
	tls := diff.Get("tls").(bool)
	caCert := diff.Get("ca_cert").(string)
	clientCert := diff.Get("client_cert").(string)
	clientKey := diff.Get("client_key").(string)

	if !tls && (caCert != "" || clientCert != "" || clientKey != "") {
		return fmt.Errorf("tls must be true when ca_cert, client_cert, or client_key is set")
	}

	// If client_cert or client_key is set, require both
	if (clientCert != "" || clientKey != "") && (clientCert == "" || clientKey == "") {
		return fmt.Errorf("both client_cert and client_key must be set for mTLS")
	}

	return nil
}
