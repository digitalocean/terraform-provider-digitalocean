package main

import (
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: digitalocean.Provider})
}

//Holaxd