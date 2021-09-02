package digitalocean

import "testing"

// maybe these can be functions instead
var test = `
resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "nyc2"
	size   = "s-1vcpu-1gb"
  }
  
  resource "digitalocean_monitoring" "cpu_alert" {
	alerts      = {
	  email = ["benny@digitalocean.com"]
	  slack = {
		"channel"   = "Production Alerts",
		"url"       = "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
	}
	window      = "5m"
	type        = "v1/insights/droplet/cpu"
	compare     = "GreaterThan"
	value       = 95
	enabled     = true
	entities    = [digitalocean_droplet.web.id]
	description = "Alert about CPU usage"
  }
`

func TestAccDigitalOceanMonitorAlert(t *testing.T) {

}
