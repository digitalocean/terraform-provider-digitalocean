package database

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

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
			"rsyslog_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DNS name or IPv4 address of the rsyslog server",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The internal port on which the rsyslog server is listening",
						},
						"tls": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Use TLS (as the messages are not filtered and may contain sensitive information, it is highly recommended to set this to true if the remote server supports it)",
						},
						"format": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"rfc5424",
								"rfc3164",
								"custom",
							}, false),
							Description: "Message format used by the server, this can be either rfc3164 (the old BSD style message format), rfc5424 (current syslog message format) or custom",
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
					},
				},
			},
			"elasticsearch_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Connection URL. Required for Elasticsearch",
						},
						"index_prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Index prefix. Required for Elasticsearch",
						},
						"index_days_max": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Default 7 days. Maximum number of days of logs to keep",
						},
						"timeout": {
							Type:         schema.TypeFloat,
							Required:     true,
							Description:  "Default 10 days. Required for Elasticsearch",
							ValidateFunc: validation.FloatBetween(10, 120),
						},
						"ca": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PEM encoded CA certificate",
						},
					},
				},
			},
			"opensearch_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Connection URL. Required for Opensearch",
						},
						"index_prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Index prefix. Required for Opensearch",
						},
						"index_days_max": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Default 7 days. Maximum number of days of logs to keep",
						},
						"timeout": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Description:  "Default 10 days",
							ValidateFunc: validation.FloatBetween(10, 120),
						},
						"ca": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PEM encoded CA certificate",
						},
					},
				},
			},
		},
	}
}

func expandLogsinkRsyslogConfig(config []interface{}) *godo.RsyslogLogsinkConfig {
	logsinkConfigOpts := &godo.RsyslogLogsinkConfig{}
	if len(config) == 0 || config[0] == nil {
		return logsinkConfigOpts
	}
	configMap := config[0].(map[string]interface{})
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

	return logsinkConfigOpts
}

func expandLogsinkElasticsearchConfig(config []interface{}) *godo.ElasticsearchLogsinkConfig {
	logsinkConfigOpts := &godo.ElasticsearchLogsinkConfig{}
	if len(config) == 0 || config[0] == nil {
		return logsinkConfigOpts
	}
	configMap := config[0].(map[string]interface{})
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
	if v, ok := configMap["ca"]; ok {
		logsinkConfigOpts.CA = v.(string)
	}

	return logsinkConfigOpts
}

func expandLogsinkOpensearchConfig(config []interface{}) *godo.OpensearchLogsinkConfig {
	logsinkConfigOpts := &godo.OpensearchLogsinkConfig{}
	if len(config) == 0 || config[0] == nil {
		return logsinkConfigOpts
	}
	configMap := config[0].(map[string]interface{})
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
	if v, ok := configMap["ca"]; ok {
		logsinkConfigOpts.CA = v.(string)
	}

	return logsinkConfigOpts
}

func resourceDigitalOceanDatabaseLogsinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	sinkType := d.Get("sink_type").(string)
	opts := &godo.DatabaseCreateLogsinkRequest{
		Name: d.Get("sink_name").(string),
		Type: d.Get("sink_type").(string),
	}

	var iCfg interface{}

	switch sinkType {
	case "rsyslog":
		if v, ok := d.GetOk("rsyslog_config"); ok {
			iCfg = expandLogsinkRsyslogConfig(v.([]interface{}))
		} else {
			return diag.Errorf("Error creating database logsink: rsyslog_config is required when type is rsyslog")
		}
	case "elasticsearch":
		if v, ok := d.GetOk("elasticsearch_config"); ok {
			iCfg = expandLogsinkElasticsearchConfig(v.([]interface{}))
		} else {
			return diag.Errorf("Error creating database logsink: elasticsearch_config is required when type is elasticsearch")
		}
	case "opensearch":
		if v, ok := d.GetOk("opensearch_config"); ok {
			iCfg = expandLogsinkOpensearchConfig(v.([]interface{}))
		} else {
			return diag.Errorf("Error creating database logsink: opensearch_config is required when type is opensearch")
		}
	}

	opts.Config = &iCfg
	if opts.Config == nil {
		return diag.Errorf("Error creating database logsink: config is required")
	}

	logsink, _, err := client.Databases.CreateLogsink(context.Background(), clusterID, opts)
	if err != nil {
		return diag.Errorf("Error creating database logsink: %s", err)
	}

	logsinkIDFormat := makeDatabaseLogsinkID(clusterID, logsink.ID)
	d.SetId(logsinkIDFormat)
	d.Set("sink_id", logsink.ID)

	return resourceDigitalOceanDatabaseLogsinkRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	opts := &godo.DatabaseUpdateLogsinkRequest{}

	sinkType := d.Get("sink_type").(string)

	var iCfg interface{}

	switch sinkType {
	case "rsyslog":
		if v, ok := d.GetOk("rsyslog_config"); ok {
			iCfg = expandLogsinkRsyslogConfig(v.([]interface{}))
		} else {
			return diag.Errorf("Error updating database logsink: rsyslog_config is required when type is rsyslog")
		}
	case "elasticsearch":
		if v, ok := d.GetOk("elasticsearch_config"); ok {
			iCfg = expandLogsinkElasticsearchConfig(v.([]interface{}))
		} else {
			return diag.Errorf("Error updating database logsink: elasticsearch_config is required when type is elasticsearch")
		}
	case "opensearch":
		if v, ok := d.GetOk("opensearch_config"); ok {
			iCfg = expandLogsinkOpensearchConfig(v.([]interface{}))
		} else {
			return diag.Errorf("Error updating database logsink: opensearch_config is required when type is opensearch")
		}
	}

	opts.Config = &iCfg

	_, err := client.Databases.UpdateLogsink(context.Background(), clusterID, d.Id(), opts)
	if err != nil {
		return diag.Errorf("Error updating database logsink: %s", err)
	}

	return resourceDigitalOceanDatabaseLogsinkRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseLogsinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	logsinkID := d.Get("sink_id").(string)

	_, err := client.Databases.DeleteLogsink(ctx, clusterID, logsinkID)
	if err != nil {
		return diag.Errorf("Error deleting logsink topic: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseLogsinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	logsinkID := d.Get("sink_id").(string)

	logsink, resp, err := client.Databases.GetLogsink(ctx, clusterID, logsinkID)
	if err != nil {
		// If the logsink is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving logsink: %s", err)
	}

	d.Set("sink_name", logsink.Name)
	d.Set("sink_type", logsink.Type)

	switch logsink.Type {
	case "rsyslog":
		if cfg, ok := (*logsink.Config).(*godo.RsyslogLogsinkConfig); ok {
			if err := d.Set("config", flattenLogsinkRsyslogConfig(cfg)); err != nil {
				return diag.Errorf("Error setting logsink config: %#v", err)
			}
		} else {
			return diag.Errorf("Error asserting logsink config to RsyslogLogsinkConfig")
		}
	case "elasticsearch":
		if cfg, ok := (*logsink.Config).(*godo.ElasticsearchLogsinkConfig); ok {
			if err := d.Set("config", flattenLogsinkElasticsearchConfig(cfg)); err != nil {
				return diag.Errorf("Error setting logsink config: %#v", err)
			}
		} else {
			return diag.Errorf("Error asserting logsink config to ElasticsearchLogsinkConfig")
		}
	case "opensearch":
		if cfg, ok := (*logsink.Config).(*godo.OpensearchLogsinkConfig); ok {
			if err := d.Set("config", flattenLogsinkOpensearchConfig(cfg)); err != nil {
				return diag.Errorf("Error setting logsink config: %#v", err)
			}
		} else {
			return diag.Errorf("Error asserting logsink config to OpensearchLogsinkConfig")
		}
	}

	return nil
}

func resourceDigitalOceanDatabaseLogsinkImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
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

func flattenLogsinkRsyslogConfig(config *godo.RsyslogLogsinkConfig) map[string]interface{} {
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

		return r
	}

	return nil
}

func flattenLogsinkElasticsearchConfig(config *godo.ElasticsearchLogsinkConfig) map[string]interface{} {
	if config != nil {
		r := make(map[string]interface{})
		r["ca"] = (*config).CA
		r["url"] = (*config).URL
		r["index_prefix"] = (*config).IndexPrefix
		r["index_days_max"] = (*config).IndexDaysMax
		r["timeout"] = (*config).Timeout

		return r
	}

	return nil
}

func flattenLogsinkOpensearchConfig(config *godo.OpensearchLogsinkConfig) map[string]interface{} {
	if config != nil {
		r := make(map[string]interface{})
		r["ca"] = (*config).CA
		r["url"] = (*config).URL
		r["index_prefix"] = (*config).IndexPrefix
		r["index_days_max"] = (*config).IndexDaysMax
		r["timeout"] = (*config).Timeout

		return r
	}

	return nil
}
