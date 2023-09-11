# insert explanation and instructions here

terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = ">= 2.8.0"
    }
  }
}

provider "digitalocean" {
  # You need to set this in your .bashrc
  # export DIGITALOCEAN_TOKEN="Your API TOKEN"
  #
}

resource "digitalocean_app" "golang-sample" {
  spec {
    name   = "golang-sample"
    region = "ams"

    service {
      name               = "go-service"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "professional-xs"

        env {
            key   = "APP_ENV"
            value = "production"
            scope = "RUN_AND_BUILD_TIME"
            type  = "GENERAL"
        }

        env {
            key   = "SECRET_TOKEN"
            value = (var.encrypted_secret_token != "" ? var.encrypted_secret_token : var.secret_token)
            scope = "RUN_TIME"
            type  = "SECRET"
        }

        git {
            repo_clone_url = "https://github.com/digitalocean/sample-golang.git"
            branch         = "main"
        }
    }

    database {
      name       = "db"
      engine     = "PG"
      production = false
    }
  }
}