output "endpoint_id" {
  description = "The ID of the dedicated inference endpoint."
  value       = digitalocean_dedicated_inference.example.id
}

output "endpoint_status" {
  description = "The current status of the dedicated inference endpoint."
  value       = digitalocean_dedicated_inference.example.status
}

output "private_endpoint_fqdn" {
  description = "The FQDN of the private endpoint."
  value       = digitalocean_dedicated_inference.example.private_endpoint_fqdn
}

output "public_endpoint_fqdn" {
  description = "The FQDN of the public endpoint (empty if not enabled)."
  value       = digitalocean_dedicated_inference.example.public_endpoint_fqdn
}

output "api_token" {
  description = "The API token for accessing the inference endpoint."
  value       = digitalocean_dedicated_inference_token.example.token
  sensitive   = true
}
