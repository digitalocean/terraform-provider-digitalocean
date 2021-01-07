---
page_title: "DigitalOcean: digitalocean_app"
---

# digitalocean\_app

Provides a DigitalOcean App resource.

## Example Usage

To create an app, provide a [DigitalOcean app spec](https://www.digitalocean.com/docs/app-platform/references/app-specification-reference/) specifying the app's components.

### Basic Example

```hcl
resource "digitalocean_app" "golang-sample" {
  spec {
    name   = "golang-sample"
    region = "ams"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "professional-xs"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
        branch         = "main"
      }
    }
  }
}
```

### Static Site Example

```hcl
resource "digitalocean_app" "static-ste-example" {
  spec {
    name   = "static-ste-example"
    region = "ams"

    static_site {
      name          = "sample-jekyll"
      build_command = "bundle exec jekyll build -d ./public"
      output_dir    = "/public"

      git {
        repo_clone_url = "https://github.com/digitalocean/sample-jekyll.git"
        branch         = "main"
      }
    }
  }
}
```

### Multiple Components Example

```hcl
resource "digitalocean_app" "mono-repo-example" {
  spec {
    name    = "mono-repo-example"
    region  = "ams"
    domains = ["foo.example.com"]

    # Build a Go project in the api/ directory that listens on port 3000
    # and serves it at https://foo.example.com/api
    service {
      name               = "api"
      environment_slug   = "go"
      instance_count     = 2
      instance_size_slug = "professional-xs"

      github {
        branch         = "main"
        deploy_on_push = true
        repo           = "username/repo"
      }

      source_dir = "api/"
      http_port  = 3000

      routes {
        path = "/api"
      }

      run_command = "bin/api"
    }

    # Builds a static site in the project's root directory
    # and serves it at https://foo.example.com/
    static_site {
      name          = "web"
      build_command = "npm run build"

      github {
        branch         = "main"
        deploy_on_push = true
        repo           = "username/repo"
      }

      routes {
        path = "/"
      }
    }

    database {
      name       = "starter-db"
      engine     = "PG"
      production = false
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `spec` - (Required) A DigitalOcean App spec describing the app.
 - `name` - (Required) The name of the app. Must be unique across all apps in the same account.
 - `region` - The slug for the DigitalOcean data center region hosting the app.
 - `domains` - A list of hostnames where the application will be available.

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
* `git` - A Git repo to use as the component's source. The repository must be able to be cloned without authentication. Only one of `git` and `github` may be set.
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/github/install). Only one of `git` and `github` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  - `scope` - The visibility scope of the environment variable. One of `RUN_TIME`, `BUILD_TIME`, or `RUN_AND_BUILD_TIME` (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.
* `route` - An HTTP paths that should be routed to this component.
  - `path` - Paths must start with `/` and must be unique within the app.
* `health_check` - A health check to determine the availability of this component.
  - `http_path` - The route path used for the HTTP health check ping.
  - `initial_delay_seconds` - The number of seconds to wait before beginning health checks.
  - `period_seconds` - The number of seconds to wait between health checks.
  - `timeout_seconds` - The number of seconds after which the check times out.
  - `success_threshold` - The number of successful health checks before considered healthy.
  - `failure_threshold` - The number of failed health checks before considered unhealthy.

A `worker` can contain:

* `name` - The name of the component
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `run_command` - An optional run command to override the component's default.
* `environment_slug` - An environment slug describing the type of this app.
* `instance_size_slug` - The instance size to use for this component.
* `instance_count` - The amount of instances that this component should be scaled to.
* `git` - A Git repo to use as the component's source. The repository must be able to be cloned without authentication. Only one of `git` and `github` may be set.
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/github/install). Only one of `git` and `github` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  - `scope` - The visibility scope of the environment variable. One of `RUN_TIME`, `BUILD_TIME`, or `RUN_AND_BUILD_TIME` (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.

A `static_site` can contain:

* `name` - The name of the component
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `environment_slug` - An environment slug describing the type of this app.
* `output_dir` - An optional path to where the built assets will be located, relative to the build context. If not set, App Platform will automatically scan for these directory names: `_static`, `dist`, `public`.
* `index_document` - The name of the index document to use when serving this static site.
* `error_document` - The name of the error document to use when serving this static site.
* `catchall_document` - The name of the document to use as the fallback for any requests to documents that are not found when serving this static site.
* `git` - A Git repo to use as the component's source. The repository must be able to be cloned without authentication. Only one of `git` and `github` may be set.
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/github/install). Only one of `git` and `github` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  - `scope` - The visibility scope of the environment variable. One of `RUN_TIME`, `BUILD_TIME`, or `RUN_AND_BUILD_TIME` (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.
* `route` - An HTTP paths that should be routed to this component.
  - `path` - Paths must start with `/` and must be unique within the app.

A `database` can contain:

* `name` - The name of the component.
* `engine` - The database engine to use (`MYSQL`, `PG`, or `REDIS`).
* `version` - The version of the database engine.
* `production` - Whether this is a production or dev database.
* `cluster_name` - The name of the underlying DigitalOcean DBaaS cluster. This is required for production databases. For dev databases, if `cluster_name` is not set, a new cluster will be provisioned.
* `db_name` - The name of the MySQL or PostgreSQL database to configure.
* `db_user` - The name of the MySQL or PostgreSQL user to configure.

## Attributes Reference

In addition to the above attributes, the following are exported:

* `default_ingress` - The default URL to access the app.
* `live_url` - The live URL of the app.
* `active_deployment_id` - The ID the app's currently active deployment.
* `updated_at` - The date and time of when the app was last updated.
* `created_at` - The date and time of when the app was created.

## Import

An app can be imported using its `id`, e.g.

```
terraform import digitalocean_app.myapp fb06ad00-351f-45c8-b5eb-13523c438661
```
