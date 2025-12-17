variable "region" {
  description = "DigitalOcean region where resources will be created"
  type        = string
  default     = "sfo3"
}

variable "db_name" {
  description = "Name for the database cluster"
  type        = string
  default     = "test-pg"
}

variable "db_engine" {
  description = "Database engine to use"
  type        = string
  default     = "pg"
}

variable "db_version" {
  description = "Database version to use"
  type        = string
  default     = "17"
}

variable "db_size" {
  description = "Database size slug"
  type        = string
  default     = "db-s-2vcpu-4gb"
}

variable "db_node_count" {
  description = "Number of nodes in the database cluster"
  type        = number
  default     = 2
}

variable "rsyslog_server" {
  description = "Hostname or IP address of the rsyslog server"
  type        = string
  default     = "logs.example.com"
}

variable "rsyslog_port" {
  description = "Port number for the rsyslog server"
  type        = number
  default     = 514
}

variable "rsyslog_format" {
  description = "Log format to use (rfc5424, rfc3164, or custom)"
  type        = string
  default     = "rfc5424"
}
