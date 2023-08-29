terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.4.0"
    }
  }
}

resource "digitalocean_spaces_bucket" "foobar" {
	name   = "samiemadidriesbengals"
	region = "nyc3"
  }

resource "digitalocean_spaces_bucket_cors_configuration" "foo" {
  bucket = digitalocean_spaces_bucket.foobar.id
  region = "nyc3"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}