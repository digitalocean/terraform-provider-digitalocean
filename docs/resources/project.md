---
page_title: "DigitalOcean: digitalocean_project"
---

# digitalocean\_project

Provides a DigitalOcean Project resource.

Projects allow you to organize your resources into groups that fit the way you work.
You can group resources (like Droplets, Spaces, Load Balancers, domains, and Floating IPs)
in ways that align with the applications you host on DigitalOcean.

The following resource types can be associated with a project:

* App Platform Apps
* Database Clusters
* Domains
* Droplets
* Floating IPs
* Kubernetes Clusters
* Load Balancers
* Spaces Buckets
* Volumes

**Note:** A Terraform managed project cannot be set as a default project.

## Example Usage

The following example demonstrates the creation of an empty project:

```hcl
resource "digitalocean_project" "playground" {
  name        = "playground"
  description = "A project to represent development resources."
  purpose     = "Web Application"
  environment = "Development"
}
```

The following example demonstrates the creation of a project with a Droplet resource:

```hcl
resource "digitalocean_droplet" "foobar" {
  name   = "example"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_project" "playground" {
  name        = "playground"
  description = "A project to represent development resources."
  purpose     = "Web Application"
  environment = "Development"
  resources   = [digitalocean_droplet.foobar.urn]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Project
* `description` - (Optional) the description of the project
* `purpose` - (Optional) the purpose of the project, (Default: "Web Application")
* `environment` - (Optional) the environment of the project's resources. The possible values are: `Development`, `Staging`, `Production`)
* `resources` - a list of uniform resource names (URNs) for the resources associated with the project
* `is_default` - (Optional) a boolean indicating whether or not the project is the default project. (Default: "false")

## Attributes Reference

The following attributes are exported:

* `id` - The id of the project
* `owner_uuid` - the unique universal identifier of the project owner.
* `owner_id` - the id of the project owner.
* `created_at` - the date and time when the project was created, (ISO8601)
* `updated_at` - the date and time when the project was last updated, (ISO8601)

## Import

Projects can be imported using the `id` returned from DigitalOcean, e.g.

```
terraform import digitalocean_project.myproject 245bcfd0-7f31-4ce6-a2bc-475a116cca97
```
