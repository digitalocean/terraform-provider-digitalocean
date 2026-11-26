package main

import (
	"context"
	"log"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
)

func main() {
	providerServer, err := digitalocean.ProviderServer(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = tf5server.Serve("registry.terraform.io/digitalocean/digitalocean", providerServer)
	if err != nil {
		log.Fatal(err)
	}
}
