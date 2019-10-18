package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-digitalocean/digitalocean"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: digitalocean.Provider})
}
