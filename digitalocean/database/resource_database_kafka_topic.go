package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseKafkaTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseKafkaTopicCreate,
		ReadContext:   resourceDigitalOceanDatabaseKafkaTopicRead,
		UpdateContext: resourceDigitalOceanDatabaseKafkaTopicUpdate,
		DeleteContext: resourceDigitalOceanDatabaseKafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseKafkaTopicImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"partition_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntBetween(3, 2048),
				Default:      3,
			},
			"replication_factor": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntAtLeast(2),
				Default:      2,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cleanup_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"delete",
								"compact",
								"compact_delete",
							}, false),
						},
						"compression_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"snappy",
								"gzip",
								"lz4",
								"producer",
								"uncompressed",
								"zstd",
							}, false),
						},
						"delete_retention_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"file_delete_delay_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"flush_messages": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"flush_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"index_interval_bytes": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"max_compaction_lag_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"max_message_bytes": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"message_down_conversion_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"message_format_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"0.8.0",
								"0.8.1",
								"0.8.2",
								"0.9.0",
								"0.10.0",
								"0.10.0-IV0",
								"0.10.0-IV1",
								"0.10.1",
								"0.10.1-IV0",
								"0.10.1-IV1",
								"0.10.1-IV2",
								"0.10.2",
								"0.10.2-IV0",
								"0.11.0",
								"0.11.0-IV0",
								"0.11.0-IV1",
								"0.11.0-IV2",
								"1.0",
								"1.0-IV0",
								"1.1",
								"1.1-IV0",
								"2.0",
								"2.0-IV0",
								"2.0-IV1",
								"2.1",
								"2.1-IV0",
								"2.1-IV1",
								"2.1-IV2",
								"2.2",
								"2.2-IV0",
								"2.2-IV1",
								"2.3",
								"2.3-IV0",
								"2.3-IV1",
								"2.4",
								"2.4-IV0",
								"2.4-IV1",
								"2.5",
								"2.5-IV0",
								"2.6",
								"2.6-IV0",
								"2.7",
								"2.7-IV0",
								"2.7-IV1",
								"2.7-IV2",
								"2.8",
								"2.8-IV0",
								"2.8-IV1",
								"3.0",
								"3.0-IV0",
								"3.0-IV1",
								"3.1",
								"3.1-IV0",
								"3.2",
								"3.2-IV0",
								"3.3",
								"3.3-IV0",
								"3.3-IV1",
								"3.3-IV2",
								"3.3-IV3",
								"3.4",
								"3.4-IV0",
								"3.5",
								"3.5-IV0",
								"3.5-IV1",
								"3.5-IV2",
								"3.6",
								"3.6-IV0",
								"3.6-IV1",
								"3.6-IV2",
							}, false),
						},
						"message_timestamp_difference_max_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateInt64(),
						},
						"message_timestamp_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"create_time",
								"log_append_time",
							}, false),
						},
						"min_cleanable_dirty_ratio": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.FloatBetween(0.0, 1.0),
						},
						"min_compaction_lag_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"min_insync_replicas": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"preallocate": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"retention_bytes": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateInt64(),
						},
						"retention_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateInt64(),
						},
						"segment_bytes": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"segment_index_bytes": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"segment_jitter_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"segment_ms": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateUint64(),
						},
						"unclean_leader_election_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDigitalOceanDatabaseKafkaTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	partition_count := uint32(d.Get("partition_count").(int))
	replication_factor := uint32(d.Get("replication_factor").(int))

	opts := &godo.DatabaseCreateTopicRequest{
		Name:              d.Get("name").(string),
		PartitionCount:    &partition_count,
		ReplicationFactor: &replication_factor,
	}

	if v, ok := d.GetOk("config"); ok {
		opts.Config = getTopicConfig(v.([]interface{}))
	}

	log.Printf("[DEBUG] Database kafka topic create configuration: %#v", opts)
	topic, _, err := client.Databases.CreateTopic(context.Background(), clusterID, opts)
	if err != nil {
		return diag.Errorf("Error creating database kafka topic: %s", err)
	}

	d.SetId(makeDatabaseTopicID(clusterID, topic.Name))
	log.Printf("[INFO] Database kafka topic name: %s", topic.Name)

	return resourceDigitalOceanDatabaseKafkaTopicRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseKafkaTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	topicName := d.Get("name").(string)
	partition_count := uint32(d.Get("partition_count").(int))
	replication_factor := uint32(d.Get("replication_factor").(int))

	opts := &godo.DatabaseUpdateTopicRequest{
		PartitionCount:    &partition_count,
		ReplicationFactor: &replication_factor,
	}

	if v, ok := d.GetOk("config"); ok {
		opts.Config = getTopicConfig(v.([]interface{}))
	}

	log.Printf("[DEBUG] Database kafka topic update configuration: %#v", opts)
	_, err := client.Databases.UpdateTopic(context.Background(), clusterID, topicName, opts)
	if err != nil {
		return diag.Errorf("Error updating database kafka topic: %s", err)
	}

	return resourceDigitalOceanDatabaseKafkaTopicRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseKafkaTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	topicName := d.Get("name").(string)

	topic, resp, err := client.Databases.GetTopic(ctx, clusterID, topicName)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving kafka topic: %s", err)
	}

	d.Set("state", topic.State)
	d.Set("replication_factor", topic.ReplicationFactor)
	// updating 'partition_count' is async, the number of partitions returned in the API will not be updated immeadiately in the response
	// setting this property to the current state rather than the number of `partitions` returned in the GetTopic response
	d.Set("partition_count", d.Get("partition_count").(int))

	if err := d.Set("config", flattenTopicConfig(topic.Config)); err != nil {
		return diag.Errorf("Error setting topic config: %#v", err)
	}

	return nil
}

func resourceDigitalOceanDatabaseKafkaTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	topicName := d.Get("name").(string)

	log.Printf("[INFO] Deleting kafka topic: %s", d.Id())
	_, err := client.Databases.DeleteTopic(ctx, clusterID, topicName)
	if err != nil {
		return diag.Errorf("Error deleting kafka topic: %s", err)
	}

	d.SetId("")
	return nil
}
func flattenTopicConfig(config *godo.TopicConfig) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	item := make(map[string]interface{})

	item["cleanup_policy"] = config.CleanupPolicy
	item["compression_type"] = config.CompressionType
	item["delete_retention_ms"] = strconv.FormatUint(*config.DeleteRetentionMS, 10)
	item["file_delete_delay_ms"] = strconv.FormatUint(*config.FileDeleteDelayMS, 10)
	item["flush_messages"] = strconv.FormatUint(*config.FlushMessages, 10)
	item["flush_ms"] = strconv.FormatUint(*config.FlushMS, 10)
	item["index_interval_bytes"] = strconv.FormatUint(*config.IndexIntervalBytes, 10)
	item["max_compaction_lag_ms"] = strconv.FormatUint(*config.MaxCompactionLagMS, 10)
	item["max_message_bytes"] = strconv.FormatUint(*config.MaxMessageBytes, 10)
	item["message_down_conversion_enable"] = *config.MessageDownConversionEnable
	item["message_format_version"] = config.MessageFormatVersion
	item["message_timestamp_difference_max_ms"] = strconv.FormatUint(*config.MessageTimestampDifferenceMaxMS, 10)
	item["message_timestamp_type"] = config.MessageTimestampType
	item["min_cleanable_dirty_ratio"] = *config.MinCleanableDirtyRatio
	item["min_compaction_lag_ms"] = strconv.FormatUint(*config.MinCompactionLagMS, 10)
	item["min_insync_replicas"] = int(*config.MinInsyncReplicas)
	item["retention_bytes"] = strconv.FormatInt(*config.RetentionBytes, 10)
	item["retention_ms"] = strconv.FormatInt(*config.RetentionMS, 10)
	item["segment_bytes"] = strconv.FormatUint(*config.SegmentBytes, 10)
	item["segment_index_bytes"] = strconv.FormatUint(*config.SegmentIndexBytes, 10)
	item["segment_jitter_ms"] = strconv.FormatUint(*config.SegmentJitterMS, 10)
	item["segment_ms"] = strconv.FormatUint(*config.SegmentMS, 10)
	item["unclean_leader_election_enable"] = *config.UncleanLeaderElectionEnable
	result = append(result, item)

	return result
}

func makeDatabaseTopicID(clusterID string, name string) string {
	return fmt.Sprintf("%s/topic/%s", clusterID, name)
}

func validateInt64() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		_, err := strconv.ParseInt(i.(string), 10, 64)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected type of %s to be int64", k))
			return warnings, errors
		}
		return warnings, errors
	}
}

func validateUint64() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		_, err := strconv.ParseUint(i.(string), 10, 64)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected type of %s to be uint64", k))
			return warnings, errors
		}
		return warnings, errors
	}
}

func getTopicConfig(raw []interface{}) *godo.TopicConfig {
	res := &godo.TopicConfig{}
	res.CleanupPolicy = "compact_delete"
	for _, kv := range raw {
		cfg := kv.(map[string]interface{})

		if v, ok := cfg["cleanup_policy"]; ok {
			res.CleanupPolicy = v.(string)
		}
		if v, ok := cfg["compression_type"]; ok {
			res.CompressionType = v.(string)
		}
		if v, ok := cfg["delete_retention_ms"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.DeleteRetentionMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["file_delete_delay_ms"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.FileDeleteDelayMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["flush_messages"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.FlushMessages = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["flush_ms"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.FlushMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["index_interval_bytes"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.IndexIntervalBytes = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["max_compaction_lag_ms"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.MaxCompactionLagMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["max_message_bytes"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.MaxMessageBytes = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["message_down_conversion_enable"]; ok {
			res.MessageDownConversionEnable = godo.PtrTo(v.(bool))
		}
		if v, ok := cfg["message_format_version"]; ok {
			res.MessageFormatVersion = v.(string)
		}
		if v, ok := cfg["message_timestamp_difference_max_ms"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.MessageTimestampDifferenceMaxMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["message_timestamp_type"]; ok {
			res.MessageTimestampType = v.(string)
		}
		if v, ok := cfg["min_cleanable_dirty_ratio"]; ok {
			res.MinCleanableDirtyRatio = godo.PtrTo(float32(v.(float64)))
		}
		if v, ok := cfg["min_compaction_lag_ms"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.MinCompactionLagMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["min_insync_replicas"]; ok {
			res.MinInsyncReplicas = godo.PtrTo(uint32(v.(int)))
		}
		if v, ok := cfg["preallocate"]; ok {
			res.Preallocate = godo.PtrTo(v.(bool))
		}
		if v, ok := cfg["retention_bytes"]; ok {
			v, err := strconv.ParseInt(v.(string), 10, 64)
			if err == nil {
				res.RetentionBytes = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["retention_ms"]; ok {
			v, err := strconv.ParseInt(v.(string), 10, 64)
			if err == nil {
				res.RetentionMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["segment_bytes"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.SegmentBytes = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["segment_index_bytes"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.SegmentIndexBytes = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["segment_jitter_ms"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.SegmentJitterMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["segment_ms"]; ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				res.SegmentMS = godo.PtrTo(v)
			}
		}
		if v, ok := cfg["unclean_leader_election_enable"]; ok {
			res.UncleanLeaderElectionEnable = godo.PtrTo(v.(bool))
		}
	}

	return res
}

func resourceDigitalOceanDatabaseKafkaTopicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(makeDatabaseTopicID(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	} else {
		return nil, errors.New("must use the ID of the source database cluster and the name of the topic joined with a comma (e.g. `id,name`)")
	}

	return []*schema.ResourceData{d}, nil
}
