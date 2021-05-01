output "IPv4" {
  value = digitalocean_droplet.mywebserver.ipv4_address
}

output "IPv6" {
  value = digitalocean_droplet.mywebserver.ipv6_address
}

output "Name" {
  value = digitalocean_droplet.mywebserver.name
}
