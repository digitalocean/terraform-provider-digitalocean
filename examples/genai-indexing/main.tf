terraform {
  required_version = "~> 1"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.44.1"
    }
  }
}

provider "digitalocean" {
  # You need to set this in your .bashrc
  # export DIGITALOCEAN_TOKEN="Your API TOKEN"
}

# Get all available knowledge bases
data "digitalocean_gradientai_knowledge_bases" "example" {}

# Use a specific knowledge base or the first available one
locals {
  knowledge_base_uuid = var.knowledge_base_uuid != "" ? var.knowledge_base_uuid : data.digitalocean_gradientai_knowledge_bases.example.knowledge_bases[0].uuid
}

# ========================================
# API 1: List Indexing Jobs for a Knowledge Base
# ========================================

data "digitalocean_gradientai_knowledge_base_indexing_jobs" "example" {
  knowledge_base_uuid = local.knowledge_base_uuid
}

# ========================================
# API 2: List Data Sources for Indexing Job
# ========================================

# Get data sources for the first job (if any jobs exist)
data "digitalocean_gradientai_indexing_job_data_sources" "example" {
  count = length(data.digitalocean_gradientai_knowledge_base_indexing_jobs.example.jobs) > 0 ? 1 : 0

  indexing_job_uuid = data.digitalocean_gradientai_knowledge_base_indexing_jobs.example.jobs[0].uuid
}

# ========================================
# API 3: Retrieve Status of Indexing Job
# ========================================

# Get detailed status of the first job (if any jobs exist)
data "digitalocean_gradientai_indexing_job" "example" {
  count = length(data.digitalocean_gradientai_knowledge_base_indexing_jobs.example.jobs) > 0 ? 1 : 0

  uuid = data.digitalocean_gradientai_knowledge_base_indexing_jobs.example.jobs[0].uuid
}

# ========================================
# API 4: Cancel Indexing Job (Optional)
# ========================================

# Only create the cancel resource if explicitly requested
resource "digitalocean_gradientai_indexing_job_cancel" "example" {
  count = var.cancel_running_job && length(data.digitalocean_gradientai_knowledge_base_indexing_jobs.example.jobs) > 0 ? 1 : 0

  uuid = data.digitalocean_gradientai_knowledge_base_indexing_jobs.example.jobs[0].uuid
}