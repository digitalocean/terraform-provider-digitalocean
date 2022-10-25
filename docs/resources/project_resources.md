---
page_title: "DigitalOcean: digitalocean_project_resources"
---

# digitalocean\_project\_resources

Assign resources to a DigitalOcean Project. This is useful if you need to assign resources
managed in Terraform to a DigitalOcean Project managed outside of Terraform.

The following resource types can be associated with a project:

* Database Clusters
* Domains
* Droplets
* Floating IP
* Kubernetes Cluster
* Load Balancers
* Spaces Bucket
* Volume

## Example Usage

The following example assigns a droplet to a Project managed outside of Terraform:

```hcl
data "digitalocean_project" "playground" {
  name = "playground"
}

resource "digitalocean_droplet" "foobar" {
  name   = "example"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_project_resources" "barfoo" {
  project = data.digitalocean_project.playground.id
  resources = [
    digitalocean_droplet.foobar.urn
  ]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) the ID of the project
* `resources` - (Required) a list of uniform resource names (URNs) for the resources associated with the project

## Attributes Reference

No additional attributes are exported.

## Import

Importing this resource is not supported.
