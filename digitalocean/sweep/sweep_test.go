package sweep_test

import (
	"testing"

	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/app"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/cdn"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/certificate"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/database"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/domain"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/droplet"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/firewall"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/image"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/kubernetes"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/loadbalancer"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/monitoring"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/project"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/reservedip"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/snapshot"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/spaces"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sshkey"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/uptime"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/volume"
	_ "github.com/digitalocean/terraform-provider-digitalocean/digitalocean/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
