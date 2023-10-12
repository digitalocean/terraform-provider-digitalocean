---
page_title: "DigitalOcean: digitalocean_database_firewall"
---

# digitalocean\_database\_firewall

Provides a DigitalOcean database firewall resource allowing you to restrict
connections to your database to trusted sources. You may limit connections to
specific Droplets, Kubernetes clusters, or IP addresses.

## Example Usage

### Create a new database firewall allowing multiple IP addresses

```hcl
resource "digitalocean_database_firewall" "example-fw" {
  cluster_id = digitalocean_database_cluster.postgres-example.id

  rule {
    type  = "ip_addr"
    value = "192.168.1.1"
  }

  rule {
    type  = "ip_addr"
    value = "192.0.2.0"
  }
}

resource "digitalocean_database_cluster" "postgres-example" {
  name       = "example-postgres-cluster"
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

### Create a new database firewall allowing a Droplet

```hcl
resource "digitalocean_database_firewall" "example-fw" {
  cluster_id = digitalocean_database_cluster.postgres-example.id

  rule {
    type  = "droplet"
    value = digitalocean_droplet.web.id
  }
}

resource "digitalocean_droplet" "web" {
  name   = "web-01"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_database_cluster" "postgres-example" {
  name       = "example-postgres-cluster"
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

### Create a new database firewall for a database replica

```hcl
resource "digitalocean_database_cluster" "postgres-example" {
  name       = "example-postgres-cluster"
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_replica" "replica-example" {
  cluster_id = digitalocean_database_cluster.postgres-example.id
  name       = "replica-example"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
}

# Create firewall rule for database replica
resource "digitalocean_database_firewall" "example-fw" {
  cluster_id = digitalocean_database_replica.replica-example.uuid

  rule {
    type  = "ip_addr"
    value = "192.168.1.1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the target database cluster.
* `rule` - (Required) A rule specifying a resource allowed to access the database cluster. The following arguments must be specified:
  - `type` - (Required) The type of resource that the firewall rule allows to access the database cluster. The possible values are: `droplet`, `k8s`, `ip_addr`, `tag`, or `app`.
  - `value` - (Required) The ID of the specific resource, the name of a tag applied to a group of resources, or the IP address that the firewall rule allows to access the database cluster.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `uuid` - A unique identifier for the firewall rule.
* `created_at` - The date and time when the firewall rule was created.

## Import

Database firewalls can be imported using the `id` of the target database cluster
For example:

```
terraform import digitalocean_database_firewall.example-fw 5f55c6cd-863b-4907-99b8-7e09b0275d54
```
