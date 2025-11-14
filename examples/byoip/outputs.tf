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

output "byoip_prefix_advertised" {
  description = "Whether the BYOIP prefix is advertised"
  value       = data.digitalocean_byoip_prefix.example.advertised
}

output "assigned_ips_count" {
  description = "Number of IP addresses currently assigned from the BYOIP prefix"
  value       = length(data.digitalocean_byoip_addresses.example.addresses)
}

output "assigned_ips" {
  description = "List of IP addresses assigned from the BYOIP prefix"
  value = [
    for addr in data.digitalocean_byoip_addresses.example.addresses : {
      id          = addr.id
      ip_address  = addr.ip_address
      region      = addr.region
      assigned_at = addr.assigned_at
    }
  ]
}
