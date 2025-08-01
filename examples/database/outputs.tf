output "monitoring_endpoints" {
  description = "The metrics endpoints for the database cluster"
  value       = digitalocean_database_cluster.test.metrics_endpoints
}

output "monitoring_user" {
  description = "The username for accessing database metrics"
  value       = data.digitalocean_database_metrics_credentials.test.username
}

output "monitoring_password" {
  description = "The password for accessing database metrics"
  sensitive   = true
  value       = data.digitalocean_database_metrics_credentials.test.password
}
