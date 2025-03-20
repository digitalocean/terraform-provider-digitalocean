---
page_title: "DigitalOcean: digitalocean_database_opensearch_config"
subcategory: "Databases"
---

# digitalocean\_database\_opensearch\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed Opensearch database cluster.

-> **Note** Opensearch configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_opensearch_config" "example" {
  cluster_id                                            = digitalocean_database_cluster.example.id
  ism_enabled                                           = true
  ism_history_enabled                                   = true
  ism_history_max_age_hours                             = 24
  ism_history_max_docs                                  = 2500000
  ism_history_rollover_check_period_hours               = 8
  ism_history_rollover_retention_period_days            = 30
  http_max_content_length_bytes                         = 100000000
  http_max_header_size_bytes                            = 8192
  http_max_initial_line_length_bytes                    = 4096
  indices_query_bool_max_clause_count                   = 1024
  search_max_buckets                                    = 10000
  indices_fielddata_cache_size_percentage               = 3
  indices_memory_index_buffer_size_percentage           = 10
  indices_memory_min_index_buffer_size_mb               = 48
  indices_memory_max_index_buffer_size_mb               = 3
  indices_queries_cache_size_percentage                 = 10
  indices_recovery_max_mb_per_sec                       = 40
  indices_recovery_max_concurrent_file_chunks           = 2
  action_auto_create_index_enabled                      = true
  action_destructive_requires_name                      = false
  enable_security_audit                                 = false
  thread_pool_search_size                               = 1
  thread_pool_search_throttled_size                     = 1
  thread_pool_search_throttled_queue_size               = 10
  thread_pool_search_queue_size                         = 10
  thread_pool_get_size                                  = 1
  thread_pool_get_queue_size                            = 10
  thread_pool_analyze_size                              = 1
  thread_pool_analyze_queue_size                        = 10
  thread_pool_write_size                                = 1
  thread_pool_write_queue_size                          = 10
  thread_pool_force_merge_size                          = 1
  override_main_response_version                        = false
  script_max_compilations_rate                          = "use-context"
  cluster_max_shards_per_node                           = 100
  cluster_routing_allocation_node_concurrent_recoveries = 2
  plugins_alerting_filter_by_backend_roles_enabled      = false
  reindex_remote_whitelist                              = ["cloud.digitalocean.com:8080"]
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-opensearch-cluster"
  engine     = "opensearch"
  version    = "2"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc3"
  node_count = 1
}
```


## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/digitalocean/#tag/Databases/operation/databases_patch_config)
for additional details on each option.

* `cluster_id` - (Required) The ID of the target Opensearch cluster.
* `ism_enabled` - (Optional) Specifies whether ISM is enabled or not. Default: `true`
* `ism_history_enabled` - (Optional) Specifies whether audit history is enabled or not. The logs from ISM are automatically indexed to a logs document. Default: `true`
* `ism_history_max_age_hours` - (Optional) Maximum age before rolling over the audit history index, in hours. Default: `24`
* `ism_history_max_docs` - (Optional) Maximum number of documents before rolling over the audit history index. Default: `2500000`
* `ism_history_rollover_check_period_hours` - (Optional) The time between rollover checks for the audit history index, in hours. Default: `8`
* `ism_history_rollover_retention_period_days` - (Optional) Length of time long audit history indices are kept, in days. Default: `30`
* `http_max_content_length_bytes` - (Optional) Maximum content length for HTTP requests to the OpenSearch HTTP API, in bytes. Default: `100000000`
* `http_max_header_size_bytes` - (Optional) Maximum size of allowed headers, in bytes. Default: `8192`
* `http_max_initial_line_length_bytes` - (Optional) Maximum length of an HTTP URL, in bytes. Default: `4096`
* `indices_query_bool_max_clause_count` - (Optional) Maximum number of clauses Lucene BooleanQuery can have. Only increase it if necessary, as it may cause performance issues. Default: `1024`
* `search_max_buckets` - (Optional) Maximum number of aggregation buckets allowed in a single response. Default: `10000`
* `indices_fielddata_cache_size_percentage` - (Optional) Maximum amount of heap memory used for field data cache, expressed as a percentage. Decreasing the value too much will increase overhead of loading field data. Increasing the value too much will decrease amount of heap available for other operations.
* `indices_memory_index_buffer_size_percentage` - (Optional) Total amount of heap used for indexing buffer before writing segments to disk, expressed as a percentage. Too low value will slow down indexing; too high value will increase indexing performance but causes performance issues for query performance. Default: `10`
* `indices_memory_min_index_buffer_size_mb` - (Optional) Minimum amount of heap used for indexing buffer before writing segments to disk, in mb. Works in conjunction with indices_memory_index_buffer_size_percentage, each being enforced. Default: `48`
* `indices_memory_max_index_buffer_size_mb` - (Optional) Maximum amount of heap used for indexing buffer before writing segments to disk, in mb. Works in conjunction with indices_memory_index_buffer_size_percentage, each being enforced. The default is unbounded.
* `indices_queries_cache_size_percentage` - (Optional) Maximum amount of heap used for query cache. Too low value will decrease query performance and increase performance for other operations; too high value will cause issues with other functionality. Default: `10`
* `indices_recovery_max_mb_per_sec` - (Optional) Limits total inbound and outbound recovery traffic for each node, expressed in mb per second. Applies to both peer recoveries as well as snapshot recoveries (i.e., restores from a snapshot). Default: `40`
* `indices_recovery_max_concurrent_file_chunks` - (Optional) Maximum number of file chunks sent in parallel for each recovery. Default: `2`
* `action_auto_create_index_enabled` - (Optional) Specifices whether to allow automatic creation of indices. Default: `true`
* `action_destructive_requires_name` - (Optional) Specifies whether to require explicit index names when deleting indices.
* `enable_security_audit` - (Optional) Specifies whether to allow security audit logging. Default: `false`
* `thread_pool_search_size` - (Optional) Number of workers in the search operation thread pool. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.
* `thread_pool_search_throttled_size` - (Optional) Number of workers in the search throttled operation thread pool. This pool is used for searching frozen indices. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.
* `thread_pool_search_throttled_queue_size` - (Optional) Size of queue for operations in the search throttled thread pool.
* `thread_pool_search_queue_size` - (Optional) Size of queue for operations in the search thread pool.
* `thread_pool_get_size` - (Optional) Number of workers in the get operation thread pool. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.
* `thread_pool_get_queue_size` - (Optional) Size of queue for operations in the get thread pool.
* `thread_pool_analyze_size` - (Optional) Number of workers in the analyze operation thread pool. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.
* `thread_pool_analyze_queue_size` - (Optional) Size of queue for operations in the analyze thread pool.
* `thread_pool_write_size` - (Optional) Number of workers in the write operation thread pool. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.
* `thread_pool_write_queue_size` - (Optional) Size of queue for operations in the write thread pool.
* `thread_pool_force_merge_size` - (Optional) Number of workers in the force merge operation thread pool. This pool is used for forcing a merge between shards of one or more indices. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.
* `override_main_response_version` - (Optional) Compatibility mode sets OpenSearch to report its version as 7.10 so clients continue to work. Default: `false`
* `script_max_compilations_rate` - (Optional) Limits the number of inline script compilations within a period of time. Default is `use-context`
* `cluster_max_shards_per_node` - (Optional) Maximum number of shards allowed per data node.
* `cluster_routing_allocation_node_concurrent_recoveries` - (Optional) Maximum concurrent incoming/outgoing shard recoveries (normally replicas) are allowed to happen per node. Default: `2`
* `plugins_alerting_filter_by_backend_roles_enabled` - (Optional) Enable or disable filtering of alerting by backend roles. Default: `false`
* `reindex_remote_whitelist` - (Optional) Allowlist of remote IP addresses for reindexing. Changing this value will cause all OpenSearch instances to restart.

## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A Opensearch database cluster's configuration can be imported using the `id` the parent cluster, e.g.

```
terraform import digitalocean_database_opensearch_config.example 4b62829a-9c42-465b-aaa3-84051048e712
```
