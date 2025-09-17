package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseLogsinkOpensearch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseLogsinkOpensearchCreate,
		ReadContext:   resourceDigitalOceanDatabaseLogsinkOpensearchRead,
		UpdateContext: resourceDigitalOceanDatabaseLogsinkOpensearchUpdate,
		DeleteContext: resourceDigitalOceanDatabaseLogsinkOpensearchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanDatabaseLogsinkOpensearchImport,
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
			"endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateHTTPSEndpoint,
				Description:  "HTTPS URL to OpenSearch (https://host:port)",
			},
			"index_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIndexPrefix,
				Description:  "Prefix for OpenSearch indices",
			},
			"index_days_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIndexDaysMax,
				Description:  "Maximum number of days to retain indices (>= 1)",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "CA certificate for TLS verification (PEM format)",
			},
			"timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateLogsinkTimeout,
				Description:  "Request timeout for log deliveries in seconds (>= 1)",
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
		),
	}
}

func resourceDigitalOceanDatabaseLogsinkOpensearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	req := &godo.DatabaseCreateLogsinkRequest{
		Name:   d.Get("name").(string),
		Type:   "opensearch",
		Config: expandLogsinkConfigOpensearch(d),
	}

	log.Printf("[DEBUG] Database logsink opensearch create configuration: %#v", req)
	logsink, _, err := client.Databases.CreateLogsink(ctx, clusterID, req)
	if err != nil {
		return diag.Errorf("Error creating database logsink opensearch: %s", err)
	}

	d.SetId(createLogsinkID(clusterID, logsink.ID))
	log.Printf("[INFO] Database logsink opensearch ID: %s", logsink.ID)

	// Post-create read for consistency
	return resourceDigitalOceanDatabaseLogsinkOpensearchRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkOpensearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("Error retrieving database logsink opensearch: %s", err)
	}

	if logsink == nil {
		return diag.Errorf("Error retrieving database logsink opensearch: logsink is nil")
	}

	d.Set("cluster_id", clusterID)
	d.Set("name", logsink.Name)
	d.Set("logsink_id", logsink.ID)

	if err := flattenLogsinkConfigOpensearch(d, logsink.Config); err != nil {
		return diag.Errorf("Error setting logsink resource data: %s", err)
	}

	return nil
}

func resourceDigitalOceanDatabaseLogsinkOpensearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	req := &godo.DatabaseUpdateLogsinkRequest{
		Config: expandLogsinkConfigOpensearch(d),
	}

	log.Printf("[DEBUG] Database logsink opensearch update configuration: %#v", req)
	_, err := client.Databases.UpdateLogsink(ctx, clusterID, logsinkID, req)
	if err != nil {
		return diag.Errorf("Error updating database logsink opensearch: %s", err)
	}

	// Re-read the resource to refresh state
	return resourceDigitalOceanDatabaseLogsinkOpensearchRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkOpensearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID, logsinkID := splitLogsinkID(d.Id())

	if clusterID == "" || logsinkID == "" {
		return diag.Errorf("Invalid logsink ID format: %s", d.Id())
	}

	log.Printf("[INFO] Deleting database logsink opensearch: %s", d.Id())
	_, err := client.Databases.DeleteLogsink(ctx, clusterID, logsinkID)
	if err != nil {
		// Treat 404 as success (already removed)
		if godoErr, ok := err.(*godo.ErrorResponse); ok && godoErr.Response.StatusCode == 404 {
			log.Printf("[INFO] Database logsink opensearch %s was already deleted", d.Id())
		} else {
			return diag.Errorf("Error deleting database logsink opensearch: %s", err)
		}
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseLogsinkOpensearchImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Validate the import ID format
	clusterID, logsinkID := splitLogsinkID(d.Id())
	if clusterID == "" || logsinkID == "" {
		return nil, fmt.Errorf("must use the format 'cluster_id,logsink_id' for import (e.g. 'deadbeef-dead-4aa5-beef-deadbeef347d,01234567-89ab-cdef-0123-456789abcdef')")
	}

	// The Read function will handle populating all fields from the API
	return []*schema.ResourceData{d}, nil
}

// expandLogsinkConfigOpensearch converts Terraform schema data to godo.DatabaseLogsinkConfig for opensearch
func expandLogsinkConfigOpensearch(d *schema.ResourceData) *godo.DatabaseLogsinkConfig {
	config := &godo.DatabaseLogsinkConfig{}

	if v, ok := d.GetOk("endpoint"); ok {
		config.URL = v.(string)
	}
	if v, ok := d.GetOk("index_prefix"); ok {
		config.IndexPrefix = v.(string)
	}
	if v, ok := d.GetOk("index_days_max"); ok {
		config.IndexDaysMax = v.(int)
	}
	if v, ok := d.GetOk("ca_cert"); ok {
		config.CA = strings.TrimSpace(v.(string))
	}
	if v, ok := d.GetOk("timeout_seconds"); ok {
		config.Timeout = float32(v.(int))
	}

	return config
}

// flattenLogsinkConfigOpensearch converts godo.DatabaseLogsinkConfig to Terraform schema data for opensearch
func flattenLogsinkConfigOpensearch(d *schema.ResourceData, config *godo.DatabaseLogsinkConfig) error {
	if config == nil {
		return nil
	}

	if config.URL != "" {
		d.Set("endpoint", config.URL)
	}
	if config.IndexPrefix != "" {
		d.Set("index_prefix", config.IndexPrefix)
	}
	if config.IndexDaysMax != 0 {
		d.Set("index_days_max", config.IndexDaysMax)
	}
	if config.CA != "" {
		d.Set("ca_cert", strings.TrimSpace(config.CA))
	}
	d.Set("timeout_seconds", int(config.Timeout))

	return nil
}

// validateHTTPSEndpoint validates that URL uses HTTPS scheme
func validateHTTPSEndpoint(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be a string", key))
		return
	}

	if strings.TrimSpace(v) == "" {
		errs = append(errs, fmt.Errorf("%q cannot be empty", key))
		return
	}

	u, err := url.Parse(v)
	if err != nil {
		errs = append(errs, fmt.Errorf("%q must be a valid URL: %s", key, err))
		return
	}

	if u.Scheme != "https" {
		errs = append(errs, fmt.Errorf("%q must use HTTPS scheme", key))
	}

	return
}

// validateIndexPrefix validates index_prefix is not empty
func validateIndexPrefix(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be a string", key))
		return
	}

	if strings.TrimSpace(v) == "" {
		errs = append(errs, fmt.Errorf("%q cannot be empty", key))
	}

	return
}

// validateIndexDaysMax validates index_days_max is >= 1
func validateIndexDaysMax(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(int)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be an integer", key))
		return
	}

	if v < 1 {
		errs = append(errs, fmt.Errorf("%q must be >= 1", key))
	}

	return
}

// validateLogsinkTimeout validates timeout is >= 1 second
func validateLogsinkTimeout(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(int)
	if !ok {
		errs = append(errs, fmt.Errorf("%q must be an integer", key))
		return
	}

	if v < 1 {
		errs = append(errs, fmt.Errorf("%q must be >= 1", key))
	}

	return
}
