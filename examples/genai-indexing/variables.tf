variable "knowledge_base_uuid" {
  description = "UUID of the knowledge base to query for indexing jobs. If empty, uses the first available knowledge base."
  type        = string
  default     = ""
}

variable "cancel_running_job" {
  description = "Whether to cancel a running job if one exists. Use with caution in production."
  type        = bool
  default     = false
}
