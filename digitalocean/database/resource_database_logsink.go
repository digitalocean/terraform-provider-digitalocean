package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseLogsink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseLogsinkCreate,
		ReadContext:   resourceDigitalOceanDatabaseLogsinkRead,
		UpdateContext: resourceDigitalOceanDatabaseLogsinkUpdate,
		DeleteContext: resourceDigitalOceanDatabaseLogsinkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseLogsinkImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"sink_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"sink_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sink_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"rsyslog",
					"elasticsearch",
					"opensearch",
				}, false),
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "DNS name or IPv4 address of the rsyslog server. Required for rsyslog.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The internal port on which the rsyslog server is listening. Required for rsyslog",
						},
						"tls": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Use TLS (as the messages are not filtered and may contain sensitive information, it is highly recommended to set this to true if the remote server supports it). Required for rsyslog.",
						},
						"format": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"rfc5424",
								"rfc3164",
								"custom",
							}, false),
							Description: "Message format used by the server, this can be either rfc3164 (the old BSD style message format), rfc5424 (current syslog message format) or custom. Required for rsyslog.",
						},
						"logline": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Conditional (required if format == custom). Syslog log line template for a custom format, supporting limited rsyslog style templating (using %tag%). Supported tags are: HOSTNAME, app-name, msg, msgid, pri, procid, structured-data, timestamp and timestamp:::date-rfc3339.",
						},
						"sd": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Content of the structured data block of rfc5424 message",
						},
						"ca": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PEM encoded CA certificate",
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(PEM format) client key if the server requires client authentication",
						},
						"cert": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(PEM format) client cert to use",
						},
						"index_days_max": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Default 7 days. Maximum number of days of logs to keep",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Connection URL. Required for Elasticsearch and Opensearch.",
						},
						"index_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Index prefix. Required for Opensearch and Elasticsearch.",
						},
						"timeout": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Description:  "Default 10 days. Elasticsearch/Opensearch request timeout limit",
							ValidateFunc: validation.FloatBetween(10, 120),
						},
					},
				},
			},
		},
	}
}

func resourceDigitalOceanDatabaseLogsinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.DatabaseCreateLogsinkRequest{
		Name: d.Get("sink_name").(string),
		Type: d.Get("sink_type").(string),
	}

	if v, ok := d.GetOk("config"); ok {
		opts.Config = expandLogsinkConfig(v.([]interface{}))
	}

	logsink, _, err := client.Databases.CreateLogsink(context.Background(), clusterID, opts)
	if err != nil {
		return diag.Errorf("Error creating database logsink: %s", err)
	}

	time.Sleep(30 * time.Second)

	log.Printf("[DEBUGGGG] Database LOGSINK NAMEE: %#v", logsink.Name)

	logsinkIDFormat := makeDatabaseLogsinkID(clusterID, logsink.ID)
	log.Printf("[DEBUGGGG] Database logsink create configuration: %#v", logsinkIDFormat)
	d.SetId(logsinkIDFormat)
	d.Set("sink_id", logsink.ID)

	log.Printf("[DEBUGGGG] Database LOGSINK - logsink.ID: %#v", logsink.ID)
	log.Printf("[DEBUGGGG] Database LOGSINK - d sink_id: %#v", d.Get("sink_id").(string))

	return resourceDigitalOceanDatabaseLogsinkRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	opts := &godo.DatabaseUpdateLogsinkRequest{}

	if v, ok := d.GetOk("config"); ok {
		opts.Config = expandLogsinkConfig(v.([]interface{}))
	}

	log.Printf("[DEBUG] Database logsink update configuration: %#v", opts)
	_, err := client.Databases.UpdateLogsink(context.Background(), clusterID, d.Id(), opts)
	if err != nil {
		return diag.Errorf("Error updating database logsink: %s", err)
	}

	return resourceDigitalOceanDatabaseLogsinkRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	log.Printf("[INFO] Deleting logsink: %s", d.Id())
	_, err := client.Databases.DeleteLogsink(ctx, clusterID, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting logsink topic: %s", err)
	}

	d.SetId("")
	return nil
}

func expandLogsinkConfig(config []interface{}) *godo.DatabaseLogsinkConfig {
	logsinkConfigOpts := &godo.DatabaseLogsinkConfig{}
	configMap := config[0].(map[string]interface{}) // TODO: check out expandAppSpecServices

	if v, ok := configMap["server"]; ok {
		logsinkConfigOpts.Server = v.(string)
	}

	if v, ok := configMap["port"]; ok {
		logsinkConfigOpts.Port = v.(int)
	}

	if v, ok := configMap["tls"]; ok {
		logsinkConfigOpts.TLS = v.(bool)
	}

	if v, ok := configMap["format"]; ok {
		logsinkConfigOpts.Format = v.(string)
	}

	if v, ok := configMap["logline"]; ok {
		logsinkConfigOpts.Logline = v.(string)
	}

	if v, ok := configMap["sd"]; ok {
		logsinkConfigOpts.SD = v.(string)
	}

	if v, ok := configMap["ca"]; ok {
		logsinkConfigOpts.CA = v.(string)
	}

	if v, ok := configMap["key"]; ok {
		logsinkConfigOpts.Key = v.(string)
	}

	if v, ok := configMap["cert"]; ok {
		logsinkConfigOpts.Cert = v.(string)
	}

	if v, ok := configMap["url"]; ok {
		logsinkConfigOpts.URL = v.(string)
	}

	if v, ok := configMap["index_prefix"]; ok {
		logsinkConfigOpts.IndexPrefix = v.(string)
	}

	if v, ok := configMap["index_days_max"]; ok {
		logsinkConfigOpts.IndexDaysMax = v.(int)
	}

	if v, ok := configMap["timeout"]; ok {
		if v.(float64) > float64(math.SmallestNonzeroFloat32) || v.(float64) < float64(math.MaxFloat32) {
			logsinkConfigOpts.Timeout = float32(v.(float64))
		}
	}

	return logsinkConfigOpts
}

func resourceDigitalOceanDatabaseLogsinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	logsinkID := d.Get("sink_id").(string)

	logsink, resp, err := client.Databases.GetLogsink(ctx, clusterID, logsinkID)
	log.Printf("[DEBUG] Read LOGSINK - logsink: %#v", logsink)
	log.Printf("[DEBUG] Read LOGSINK - resp: %#v", resp)
	log.Printf("[DEBUG] Read LOGSINK - err: %#v", err)
	if err != nil {
		// If the logsink is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving logsink: %s", err)
	}

	d.Set("sink_name", logsink.Name) // TODO: nil - fix
	d.Set("sink_type", logsink.Type)

	log.Printf("[DEBUG] TRACE 1")
	log.Printf("[DEBUG] logsink.Config: %v ", logsink.Config)
	if err := d.Set("config", flattenDatabaseLogsinkConfig(logsink.Config)); err != nil {
		log.Printf("[DEBUG] TRACE 2")
		return diag.Errorf("Error setting logsink config: %#v", err)
	}
	log.Printf("[DEBUG] TRACE 3")

	return nil
}

func resourceDigitalOceanDatabaseLogsinkImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG44] Database logsink create configuration: %#v", d.Id())

	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		log.Printf("[DEBUG33] Database logsink create configuration: %#v", s)

		d.SetId(makeDatabaseLogsinkID(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("sink_id", s[1])
	} else {
		return nil, errors.New("must use the ID of the source cluster and logsink id joined with a comma (e.g. `id,sink_id`)")
	}

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseLogsinkID(clusterID string, logsinkID string) string {
	return fmt.Sprintf("%s/logsink/%s", clusterID, logsinkID)
}

func flattenDatabaseLogsinkConfig(config *godo.DatabaseLogsinkConfig) map[string]interface{} {

	if config != nil {
		r := make(map[string]interface{})

		r["server"] = (*config).Server
		r["port"] = (*config).Port
		r["tls"] = (*config).TLS
		r["format"] = (*config).Format
		r["logline"] = (*config).Logline
		r["sd"] = (*config).SD
		r["ca"] = (*config).CA
		r["key"] = (*config).Key
		r["cert"] = (*config).Cert
		r["url"] = (*config).URL
		r["index_prefix"] = (*config).IndexPrefix
		r["index_days_max"] = (*config).IndexDaysMax
		r["timeout"] = (*config).Timeout

		return r
	}

	return nil
}
