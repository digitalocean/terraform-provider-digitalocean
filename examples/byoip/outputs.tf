output "byoip_prefix_cidr" {
  description = "The CIDR of the BYOIP prefix"
  value       = data.digitalocean_byoip_prefix.example.prefix
}

output "byoip_prefix_region" {
  description = "The region of the BYOIP prefix"
  value       = data.digitalocean_byoip_prefix.example.region
}

output "byoip_prefix_status" {
  description = "The status of the BYOIP prefix"
  value       = data.digitalocean_byoip_prefix.example.status
}

output "available_ips_count" {
  description = "Number of IP addresses allocated from the BYOIP prefix"
  value       = length(data.digitalocean_byoip_addresses.example.addresses)
}

output "droplet_ip" {
  description = "The BYOIP IP address assigned to the Droplet"
  value       = digitalocean_reserved_ip.byoip_ip.ip_address
}

output "droplet_urn" {
  description = "The URN of the Droplet"
  value       = digitalocean_droplet.web.urn
}
