---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_app"
sidebar_current: "docs-do-datasource-app"
description: |-
  Get information on a DigitalOcean App.
---

# digitalocean_app

Get information on a DigitalOcean App.

## Example Usage

Get the account:

```hcl
data "digitalocean_app" "example" {
  app_id = "e665d18d-7b56-44a9-92ce-31979174d544"
}

output "default_ingress" {
  value = data.digitalocean_app.example.default_ingress
}
```

## Argument Reference

* `app_id` - The ID of the app to retrieve information about.

## Attributes Reference

The following attributes are exported:

* `default_ingress` - The default URL to access the app.
* `live_url` - The live URL of the app.
* `active_deployment_id` - The ID the app's currently active deployment.
* `updated_at` - The date and time of when the app was last updated.
* `created_at` - The date and time of when the app was created.
* `spec` - A DigitalOcean App spec describing the app.

A spec can contain multiple components.

A `service` can contain:

* `name` - The name of the component
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `run_command` - An optional run command to override the component's default.
* `environment_slug` - An environment slug describing the type of this app.
* `instance_size_slug` - The instance size to use for this component.
* `instance_count` - The amount of instances that this component should be scaled to.
* `http_port` - The internal port on which this service's run command will listen.
* `git` - A Git repo to use as component's source. Only one of `git` and `github` may be set.
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as component's source. Only one of `git` and `github` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  -   - (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.
* `route` - An HTTP paths that should be routed to this component.
  - `path` - Paths must start with `/` and must be unique within the app.
* `health_check` - A health check to determine the availability of this component.
  - `path` - The route path used for the HTTP health check ping.

A `worker` can contain:

* `name` - The name of the component
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `run_command` - An optional run command to override the component's default.
* `environment_slug` - An environment slug describing the type of this app.
* `instance_size_slug` - The instance size to use for this component.
* `instance_count` - The amount of instances that this component should be scaled to.
* `git` - A Git repo to use as component's source. Only one of `git` and `github` may be set.
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as component's source. Only one of `git` and `github` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  -   - (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.

A `static_site` can contain:

* `name` - The name of the component
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `environment_slug` - An environment slug describing the type of this app.
* `output_dir` - An optional path to where the built assets will be located, relative to the build context. If not set, App Platform will automatically scan for these directory names: `_static`, `dist`, `public`.
* `index_document` - The name of the index document to use when serving this static site.
* `error_document` - The name of the error document to use when serving this static site*
* `git` - A Git repo to use as component's source. Only one of `git` and `github` may be set.
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as component's source. Only one of `git` and `github` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  -   - (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.
* `route` - An HTTP paths that should be routed to this component.
  - `path` - Paths must start with `/` and must be unique within the app.