terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = ">= 2.8.0"
    }
  }
}

resource "digitalocean_domain" "example_com" {
  name = "bengalswhodey1212.com"
}

resource "digitalocean_record" "test_txt_2" {
  domain   = digitalocean_domain.example_com.name
  type     = "TXT"
  name     = "@"
  value    = "2"
  ttl      = 1800
}

resource "digitalocean_record" "test_txt_1" {
  domain   = digitalocean_domain.example_com.name
  type     = "TXT"
  name     = "@"
  value    = "1"
  ttl      = 60
}