variable do_token {}
provider digitalocean {
    token = var.do_token
}

data "digitalocean_container_registry" "registry"{
    name = "megger"
    write = true
}

output "image_sha" {
    value = data.digitalocean_container_registry.registry.docker_credentials
}
