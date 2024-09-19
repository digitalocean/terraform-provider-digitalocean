package database

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseMongoDBConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseMongoDBConfigCreate,
		ReadContext:   resourceDigitalOceanDatabaseMongoDBConfigRead,
		UpdateContext: resourceDigitalOceanDatabaseMongoDBConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabaseMongoDBConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseMongoDBConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"default_read_concern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"local",
						"available",
						"majority",
					},
					true,
				),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"default_write_concern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"transaction_lifetime_limit_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"slow_op_threshold_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"verbosity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDatabaseMongoDBConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	if err := updateMongoDBConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating MongoDB configuration: %s", err)
	}

	d.SetId(makeDatabaseMongoDBConfigID(clusterID))

	return resourceDigitalOceanDatabaseMongoDBConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseMongoDBConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updateMongoDBConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating MongoDB configuration: %s", err)
	}

	return resourceDigitalOceanDatabaseMongoDBConfigRead(ctx, d, meta)
}

func updateMongoDBConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.MongoDBConfig{}

	if v, ok := d.GetOk("default_read_concern"); ok {
		opts.DefaultReadConcern = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("default_write_concern"); ok {
		opts.DefaultWriteConcern = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("transaction_lifetime_limit_seconds"); ok {
		opts.TransactionLifetimeLimitSeconds = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("slow_op_threshold_ms"); ok {
		opts.SlowOpThresholdMs = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("verbosity"); ok {
		opts.Verbosity = godo.PtrTo(v.(int))
	}

	log.Printf("[DEBUG] MongoDB configuration: %s", godo.Stringify(opts))

	if _, err := client.Databases.UpdateMongoDBConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabaseMongoDBConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetMongoDBConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving MongoDB configuration: %s", err)
	}

	d.Set("default_read_concern", config.DefaultReadConcern)
	d.Set("default_write_concern", config.DefaultWriteConcern)
	d.Set("transaction_lifetime_limit_seconds", config.TransactionLifetimeLimitSeconds)
	d.Set("slow_op_threshold_ms", config.SlowOpThresholdMs)
	d.Set("verbosity", config.Verbosity)

	return nil
}

func resourceDigitalOceanDatabaseMongoDBConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_mongodb_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceDigitalOceanDatabaseMongoDBConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseMongoDBConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseMongoDBConfigID(clusterID string) string {
	return fmt.Sprintf("%s/mongodb-config", clusterID)
}
