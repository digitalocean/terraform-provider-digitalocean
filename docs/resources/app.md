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
    domain {
      name = "foo.example.com"
    }

    alert {
      rule = "DEPLOYMENT_FAILED"
    }

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

      alert {
        value    = 75
        operator = "GREATER_THAN"
        window   = "TEN_MINUTES"
        rule     = "CPU_UTILIZATION"
      }

      log_destination {
        name = "MyLogs"
        papertrail {
          endpoint = "syslog+tls://example.com:12345"
        }
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
 - `domain` - Describes a domain where the application will be made available.
     * `name` - The hostname for the domain.
     * `type` - The domain type, which can be one of the following:
         - `DEFAULT`: The default .ondigitalocean.app domain assigned to this app.
         - `PRIMARY`: The primary domain for this app that is displayed as the default in the control panel, used in bindable environment variables, and any other places that reference an app's live URL. Only one domain may be set as primary.
         - `ALIAS`: A non-primary domain.
     * `wildcard` - A boolean indicating whether the domain includes all sub-domains, in addition to the given domain.
     * `zone` - If the domain uses DigitalOcean DNS and you would like App Platform to automatically manage it for you, set this to the name of the domain on your account.
 - `env` - Describes an app-wide environment variable made available to all components.
     * `key` - The name of the environment variable.
     * `value` - The value of the environment variable.
     * `scope` - The visibility scope of the environment variable. One of `RUN_TIME`, `BUILD_TIME`, or `RUN_AND_BUILD_TIME` (default).
     * `type` - The type of the environment variable, `GENERAL` or `SECRET`.
 - `alert` - Describes an alert policy for the app.
     * `rule` - The type of the alert to configure. Top-level app alert policies can be: `DEPLOYMENT_FAILED`, `DEPLOYMENT_LIVE`, `DOMAIN_FAILED`, or `DOMAIN_LIVE`.
     * `disabled` - Determines whether or not the alert is disabled (default: `false`).

A spec can contain multiple components.

A `service` can contain:

* `name` - The name of the component
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `run_command` - An optional run command to override the component's default.
* `environment_slug` - An environment slug describing the type of this app.
* `instance_size_slug` - The instance size to use for this component. This determines the plan (basic or professional) and the available CPU and memory. The list of available instance sizes can be [found with the API](https://docs.digitalocean.com/reference/api/api-reference/#operation/list_instance_sizes) or using the [doctl CLI](https://docs.digitalocean.com/reference/doctl/) (`doctl apps tier instance-size list`). Default: `basic-xxs`
* `instance_count` - The amount of instances that this component should be scaled to.
* `http_port` - The internal port on which this service's run command will listen.
* `internal_ports` - A list of ports on which this service will listen for internal traffic.
* `git` - A Git repo to use as the component's source. The repository must be able to be cloned without authentication.  Only one of `git`, `github` or `gitlab`  may be set
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/github/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `gitlab` - A Gitlab repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/gitlab/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `image` - An image to use as the component's source. Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `registry_type` - The registry type. One of `DOCR` (DigitalOcean container registry) or `DOCKER_HUB`.
  - `registry` - The registry name. Must be left empty for the `DOCR` registry type. Required for the `DOCKER_HUB` registry type.
  - `repository` - The repository name.
  - `tag` - The repository tag. Defaults to `latest` if not provided.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  - `scope` - The visibility scope of the environment variable. One of `RUN_TIME`, `BUILD_TIME`, or `RUN_AND_BUILD_TIME` (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.
* `route` - An HTTP paths that should be routed to this component.
  - `path` - Paths must start with `/` and must be unique within the app.
  - `preserve_path_prefix` -  An optional flag to preserve the path that is forwarded to the backend service.
* `health_check` - A health check to determine the availability of this component.
  - `http_path` - The route path used for the HTTP health check ping.
  - `initial_delay_seconds` - The number of seconds to wait before beginning health checks.
  - `period_seconds` - The number of seconds to wait between health checks.
  - `timeout_seconds` - The number of seconds after which the check times out.
  - `success_threshold` - The number of successful health checks before considered healthy.
  - `failure_threshold` - The number of failed health checks before considered unhealthy.
* `cors` - The [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) policies of the app.
	- `allow_origins` - The `Access-Control-Allow-Origin` can be
    - `exact` - The `Access-Control-Allow-Origin` header will be set to the client's origin only if the client's origin exactly matches the value you provide.
    - `prefix` -  The `Access-Control-Allow-Origin` header will be set to the client's origin if the beginning of the client's origin matches the value you provide.
    - `regex` - The `Access-Control-Allow-Origin` header will be set to the client's origin if the client’s origin matches the regex you provide, in [RE2 style syntax](https://github.com/google/re2/wiki/Syntax).
  - `allow_headers` - The set of allowed HTTP request headers. This configures the `Access-Control-Allow-Headers` header.
  - `max_age` - An optional duration specifying how long browsers can cache the results of a preflight request. This configures the Access-Control-Max-Age header. Example: `5h30m`.
  - `expose_headers` - The set of HTTP response headers that browsers are allowed to access. This configures the `Access-Control-Expose-Headers` header.
  - `allow_methods` - The set of allowed HTTP methods. This configures the `Access-Control-Allow-Methods` header.
  - `allow_credentials` - Whether browsers should expose the response to the client-side JavaScript code when the request's credentials mode is `include`. This configures the `Access-Control-Allow-Credentials` header.
* `alert` - Describes an alert policy for the component.
  - `rule` - The type of the alert to configure. Component app alert policies can be: `CPU_UTILIZATION`, `MEM_UTILIZATION`, or `RESTART_COUNT`.
  - `value` - The threshold for the type of the warning.
  - `operator` - The operator to use. This is either of `GREATER_THAN` or `LESS_THAN`.
  - `window` - The time before alerts should be triggered. This is may be one of: `FIVE_MINUTES`, `TEN_MINUTES`, `THIRTY_MINUTES`, `ONE_HOUR`.
  - `disabled` - Determines whether or not the alert is disabled (default: `false`).
* `log_destination` - Describes a log forwarding destination.
  - `name` - Name of the log destination. Minimum length: 2. Maximum length: 42.
  - `papertrail` - Papertrail configuration.
    - `endpoint` - Papertrail syslog endpoint.
  - `datadog` - Datadog configuration.
    - `endpoint` - Datadog HTTP log intake endpoint.
    - `api_key` - Datadog API key.
  - `logtail` - Logtail configuration.
    - `token` - Logtail token.

A `static_site` can contain:

* `name` - The name of the component.
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `environment_slug` - An environment slug describing the type of this app.
* `output_dir` - An optional path to where the built assets will be located, relative to the build context. If not set, App Platform will automatically scan for these directory names: `_static`, `dist`, `public`.
* `index_document` - The name of the index document to use when serving this static site.
* `error_document` - The name of the error document to use when serving this static site.
* `catchall_document` - The name of the document to use as the fallback for any requests to documents that are not found when serving this static site.
* `git` - A Git repo to use as the component's source. The repository must be able to be cloned without authentication.  Only one of `git`, `github` or `gitlab`  may be set.
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/github/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `gitlab` - A Gitlab repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/gitlab/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
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
  - `preserve_path_prefix` -  An optional flag to preserve the path that is forwarded to the backend service.
* `cors` - The [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) policies of the app.
	- `allow_origins` - The `Access-Control-Allow-Origin` can be
    - `exact` - The `Access-Control-Allow-Origin` header will be set to the client's origin only if the client's origin exactly matches the value you provide.
    - `prefix` -  The `Access-Control-Allow-Origin` header will be set to the client's origin if the beginning of the client's origin matches the value you provide.
    - `regex` - The `Access-Control-Allow-Origin` header will be set to the client's origin if the client’s origin matches the regex you provide, in [RE2 style syntax](https://github.com/google/re2/wiki/Syntax).
  - `allow_headers` - The set of allowed HTTP request headers. This configures the `Access-Control-Allow-Headers` header.
  - `max_age` - An optional duration specifying how long browsers can cache the results of a preflight request. This configures the Access-Control-Max-Age header. Example: `5h30m`.
  - `expose_headers` - The set of HTTP response headers that browsers are allowed to access. This configures the `Access-Control-Expose-Headers` header.
  - `allow_methods` - The set of allowed HTTP methods. This configures the `Access-Control-Allow-Methods` header.
  - `allow_credentials` - Whether browsers should expose the response to the client-side JavaScript code when the request's credentials mode is `include`. This configures the `Access-Control-Allow-Credentials` header.

A `worker` can contain:

* `name` - The name of the component
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `run_command` - An optional run command to override the component's default.
* `environment_slug` - An environment slug describing the type of this app.
* `instance_size_slug` - The instance size to use for this component. This determines the plan (basic or professional) and the available CPU and memory. The list of available instance sizes can be [found with the API](https://docs.digitalocean.com/reference/api/api-reference/#operation/list_instance_sizes) or using the [doctl CLI](https://docs.digitalocean.com/reference/doctl/) (`doctl apps tier instance-size list`). Default: `basic-xxs`
* `instance_count` - The amount of instances that this component should be scaled to.
* `git` - A Git repo to use as the component's source. The repository must be able to be cloned without authentication.  Only one of `git`, `github` or `gitlab`  may be set
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/github/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `gitlab` - A Gitlab repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/gitlab/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `image` - An image to use as the component's source. Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `registry_type` - The registry type. One of `DOCR` (DigitalOcean container registry) or `DOCKER_HUB`.
  - `registry` - The registry name. Must be left empty for the `DOCR` registry type. Required for the `DOCKER_HUB` registry type.
  - `repository` - The repository name.
  - `tag` - The repository tag. Defaults to `latest` if not provided.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  - `scope` - The visibility scope of the environment variable. One of `RUN_TIME`, `BUILD_TIME`, or `RUN_AND_BUILD_TIME` (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.
* `alert` - Describes an alert policy for the component.
  - `rule` - The type of the alert to configure. Component app alert policies can be: `CPU_UTILIZATION`, `MEM_UTILIZATION`, or `RESTART_COUNT`.
  - `value` - The threshold for the type of the warning.
  - `operator` - The operator to use. This is either of `GREATER_THAN` or `LESS_THAN`.
  - `window` - The time before alerts should be triggered. This is may be one of: `FIVE_MINUTES`, `TEN_MINUTES`, `THIRTY_MINUTES`, `ONE_HOUR`.
  - `disabled` - Determines whether or not the alert is disabled (default: `false`).
* `log_destination` - Describes a log forwarding destination.
  - `name` - Name of the log destination. Minimum length: 2. Maximum length: 42.
  - `papertrail` - Papertrail configuration.
    - `endpoint` - Papertrail syslog endpoint.
  - `datadog` - Datadog configuration.
    - `endpoint` - Datadog HTTP log intake endpoint.
    - `api_key` - Datadog API key.
  - `logtail` - Logtail configuration.
    - `token` - Logtail token.

A `job` can contain:

* `name` - The name of the component
* `kind` - The type of job and when it will be run during the deployment process. It may be one of:
  - `UNSPECIFIED`: Default job type, will auto-complete to POST_DEPLOY kind.
  - `PRE_DEPLOY`: Indicates a job that runs before an app deployment.
  - `POST_DEPLOY`: Indicates a job that runs after an app deployment.
  - `FAILED_DEPLOY`: Indicates a job that runs after a component fails to deploy.
* `build_command` - An optional build command to run while building this component from source.
* `dockerfile_path` - The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.
* `source_dir` - An optional path to the working directory to use for the build.
* `run_command` - An optional run command to override the component's default.
* `environment_slug` - An environment slug describing the type of this app.
* `instance_size_slug` - The instance size to use for this component. This determines the plan (basic or professional) and the available CPU and memory. The list of available instance sizes can be [found with the API](https://docs.digitalocean.com/reference/api/api-reference/#operation/list_instance_sizes) or using the [doctl CLI](https://docs.digitalocean.com/reference/doctl/) (`doctl apps tier instance-size list`). Default: `basic-xxs`
* `instance_count` - The amount of instances that this component should be scaled to.
* `git` - A Git repo to use as the component's source. The repository must be able to be cloned without authentication.  Only one of `git`, `github` or `gitlab`  may be set
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/github/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `gitlab` - A Gitlab repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/gitlab/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `image` - An image to use as the component's source. Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `registry_type` - The registry type. One of `DOCR` (DigitalOcean container registry) or `DOCKER_HUB`.
  - `registry` - The registry name. Must be left empty for the `DOCR` registry type. Required for the `DOCKER_HUB` registry type.
  - `repository` - The repository name.
  - `tag` - The repository tag. Defaults to `latest` if not provided.
* `env` - Describes an environment variable made available to an app competent.
  - `key` - The name of the environment variable.
  - `value` - The value of the environment variable.
  - `scope` - The visibility scope of the environment variable. One of `RUN_TIME`, `BUILD_TIME`, or `RUN_AND_BUILD_TIME` (default).
  - `type` - The type of the environment variable, `GENERAL` or `SECRET`.
* `alert` - Describes an alert policy for the component.
  - `rule` - The type of the alert to configure. Component app alert policies can be: `CPU_UTILIZATION`, `MEM_UTILIZATION`, or `RESTART_COUNT`.
  - `value` - The threshold for the type of the warning.
  - `operator` - The operator to use. This is either of `GREATER_THAN` or `LESS_THAN`.
  - `window` - The time before alerts should be triggered. This is may be one of: `FIVE_MINUTES`, `TEN_MINUTES`, `THIRTY_MINUTES`, `ONE_HOUR`.
  - `disabled` - Determines whether or not the alert is disabled (default: `false`).
* `log_destination` - Describes a log forwarding destination.
  - `name` - Name of the log destination. Minimum length: 2. Maximum length: 42.
  - `papertrail` - Papertrail configuration.
    - `endpoint` - Papertrail syslog endpoint.
  - `datadog` - Datadog configuration.
    - `endpoint` - Datadog HTTP log intake endpoint.
    - `api_key` - Datadog API key.
  - `logtail` - Logtail configuration.
    - `token` - Logtail token.


A `function` component can contain:

* `name` - The name of the component.
* `source_dir` - An optional path to the working directory to use for the build.
* `git` - A Git repo to use as the component's source. The repository must be able to be cloned without authentication.  Only one of `git`, `github` or `gitlab`  may be set.
  - `repo_clone_url` - The clone URL of the repo.
  - `branch` - The name of the branch to use.
* `github` - A GitHub repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/github/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
  - `repo` - The name of the repo in the format `owner/repo`.
  - `branch` - The name of the branch to use.
  - `deploy_on_push` - Whether to automatically deploy new commits made to the repo.
* `gitlab` - A Gitlab repo to use as the component's source. DigitalOcean App Platform must have [access to the repository](https://cloud.digitalocean.com/apps/gitlab/install). Only one of `git`, `github`, `gitlab`, or `image` may be set.
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
  - `preserve_path_prefix` -  An optional flag to preserve the path that is forwarded to the backend service.
* `cors` - The [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) policies of the app.
	- `allow_origins` - The `Access-Control-Allow-Origin` can be
    - `exact` - The `Access-Control-Allow-Origin` header will be set to the client's origin only if the client's origin exactly matches the value you provide.
    - `prefix` -  The `Access-Control-Allow-Origin` header will be set to the client's origin if the beginning of the client's origin matches the value you provide.
    - `regex` - The `Access-Control-Allow-Origin` header will be set to the client's origin if the client’s origin matches the regex you provide, in [RE2 style syntax](https://github.com/google/re2/wiki/Syntax).
  - `allow_headers` - The set of allowed HTTP request headers. This configures the `Access-Control-Allow-Headers` header.
  - `max_age` - An optional duration specifying how long browsers can cache the results of a preflight request. This configures the Access-Control-Max-Age header. Example: `5h30m`.
  - `expose_headers` - The set of HTTP response headers that browsers are allowed to access. This configures the `Access-Control-Expose-Headers` header.
  - `allow_methods` - The set of allowed HTTP methods. This configures the `Access-Control-Allow-Methods` header.
  - `allow_credentials` - Whether browsers should expose the response to the client-side JavaScript code when the request's credentials mode is `include`. This configures the `Access-Control-Allow-Credentials` header.
* `alert` - Describes an alert policy for the component.
  - `rule` - The type of the alert to configure. Component app alert policies can be: `CPU_UTILIZATION`, `MEM_UTILIZATION`, or `RESTART_COUNT`.
  - `value` - The threshold for the type of the warning.
  - `operator` - The operator to use. This is either of `GREATER_THAN` or `LESS_THAN`.
  - `window` - The time before alerts should be triggered. This is may be one of: `FIVE_MINUTES`, `TEN_MINUTES`, `THIRTY_MINUTES`, `ONE_HOUR`.
  - `disabled` - Determines whether or not the alert is disabled (default: `false`).
* `log_destination` - Describes a log forwarding destination.
  - `name` - Name of the log destination. Minimum length: 2. Maximum length: 42.
  - `papertrail` - Papertrail configuration.
    - `endpoint` - Papertrail syslog endpoint.
  - `datadog` - Datadog configuration.
    - `endpoint` - Datadog HTTP log intake endpoint.
    - `api_key` - Datadog API key.
  - `logtail` - Logtail configuration.
    - `token` - Logtail token.

A `database` can contain:

* `name` - The name of the component.
* `engine` - The database engine to use (`MYSQL`, `PG`, `REDIS`, or `MONGODB`).
* `version` - The version of the database engine.
* `production` - Whether this is a production or dev database.
* `cluster_name` - The name of the underlying DigitalOcean DBaaS cluster. This is required for production databases. For dev databases, if `cluster_name` is not set, a new cluster will be provisioned.
* `db_name` - The name of the MySQL or PostgreSQL database to configure.
* `db_user` - The name of the MySQL or PostgreSQL user to configure.

This resource supports [customized create timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts). The default timeout is 30 minutes.

## Attributes Reference

In addition to the above attributes, the following are exported:

* `id` - The ID of the app.
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
