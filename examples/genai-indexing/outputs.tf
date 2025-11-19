# ========================================
# Knowledge Base Information
# ========================================

output "knowledge_base_uuid" {
  description = "UUID of the knowledge base being used"
  value       = local.knowledge_base_uuid
}

# ========================================
# API 1: List Indexing Jobs Output
# ========================================

output "indexing_jobs_count" {
  description = "Total number of indexing jobs"
  value       = length(data.digitalocean_genai_knowledge_base_indexing_jobs.example.jobs)
}

output "indexing_jobs" {
  description = "List of all indexing jobs"
  value = [
    for job in data.digitalocean_genai_knowledge_base_indexing_jobs.example.jobs : {
      uuid                  = job.uuid
      status                = job.status
      phase                 = job.phase
      tokens                = job.tokens
      completed_datasources = job.completed_datasources
      total_datasources     = job.total_datasources
      created_at            = job.created_at
      updated_at            = job.updated_at
    }
  ]
}

# ========================================
# API 2: Indexing Job Data Sources Output
# ========================================

output "job_data_sources_count" {
  description = "Number of data sources in the first indexing job"
  value       = length(data.digitalocean_genai_indexing_job_data_sources.example) > 0 ? length(data.digitalocean_genai_indexing_job_data_sources.example[0].indexed_data_sources) : 0
}

output "job_data_sources" {
  description = "Data sources for the first indexing job"
  value = length(data.digitalocean_genai_indexing_job_data_sources.example) > 0 ? [
    for ds in data.digitalocean_genai_indexing_job_data_sources.example[0].indexed_data_sources : {
      data_source_uuid    = ds.data_source_uuid
      status              = ds.status
      indexed_file_count  = ds.indexed_file_count
      indexed_item_count  = ds.indexed_item_count
      total_bytes         = ds.total_bytes
      total_bytes_indexed = ds.total_bytes_indexed
      started_at          = ds.started_at
      completed_at        = ds.completed_at
    }
  ] : []
}

# ========================================
# API 3: Detailed Job Status Output
# ========================================

output "job_status" {
  description = "Detailed status of the first indexing job"
  value = length(data.digitalocean_genai_indexing_job.example) > 0 ? {
    uuid                  = data.digitalocean_genai_indexing_job.example[0].uuid
    status                = data.digitalocean_genai_indexing_job.example[0].status
    phase                 = data.digitalocean_genai_indexing_job.example[0].phase
    knowledge_base_uuid   = data.digitalocean_genai_indexing_job.example[0].knowledge_base_uuid
    tokens                = data.digitalocean_genai_indexing_job.example[0].tokens
    completed_datasources = data.digitalocean_genai_indexing_job.example[0].completed_datasources
    total_datasources     = data.digitalocean_genai_indexing_job.example[0].total_datasources
    total_items_failed    = data.digitalocean_genai_indexing_job.example[0].total_items_failed
    total_items_indexed   = data.digitalocean_genai_indexing_job.example[0].total_items_indexed
    total_items_skipped   = data.digitalocean_genai_indexing_job.example[0].total_items_skipped
    data_source_uuids     = data.digitalocean_genai_indexing_job.example[0].data_source_uuids
    created_at            = data.digitalocean_genai_indexing_job.example[0].created_at
    started_at            = data.digitalocean_genai_indexing_job.example[0].started_at
    finished_at           = data.digitalocean_genai_indexing_job.example[0].finished_at
    updated_at            = data.digitalocean_genai_indexing_job.example[0].updated_at
  } : null
}

# ========================================
# API 4: Job Cancellation Output
# ========================================

output "cancelled_job_status" {
  description = "Status of the cancelled job (if cancellation was performed)"
  value = length(resource.digitalocean_genai_indexing_job_cancel.example) > 0 ? {
    uuid   = resource.digitalocean_genai_indexing_job_cancel.example[0].uuid
    status = resource.digitalocean_genai_indexing_job_cancel.example[0].status
    phase  = resource.digitalocean_genai_indexing_job_cancel.example[0].phase
  } : null
}