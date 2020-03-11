---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_project_resource"
sidebar_current: "docs-do-resource-project-resource"
description: |-
  Assigns a DigitalOcean to a Project.
---

# digitalocean\_project\_resource

Assign a single resource to a DigitalOcean Project resource. This is useful if you need to 
assign a resource managed in Terraform to a DigitalOcean Project managed outside of Terraform.

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
  name        = "playground"
}

resource "digitalocean_droplet" "foobar" {
  name   = "example"
  size   = "512mb"
  image  = "centos-7-x64"
  region = "nyc3"
}

resource "digitalocean_project_resource" "barfoo" {
  project = data.digitalocean_project.foo.id
  resource = digitalocean_droplet.foobar.urn
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) the ID of the project
* `resource` - (Required) the uniform resource name (URN) of the resource to assign to the project

## Attributes Reference

No additional attributes are exported.

## Import

Importing this resource is not currently supported.