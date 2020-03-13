---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_project_resources"
sidebar_current: "docs-do-resource-project-resources"
description: |-
  Assign resources to a DigitalOcean Project.
---

# digitalocean\_project\_resources

Assign resources to a DigitalOcean Project. This is useful if you need to assign resources
managed in Terraform to a DigitalOcean Project managed outside of Terraform.

The following resource types can be associated with a project:

* Database Clusters
* Domains
* Droplets
* Floating IP
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
  size   = "512mb"
  image  = "centos-7-x64"
  region = "nyc3"
}

resource "digitalocean_project_resources" "barfoo" {
  project = data.digitalocean_project.foo.id
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