package vpcnatgateway_test

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanVPCNATGateway(t *testing.T) {
	var gateway godo.VPCNATGateway
	name := acceptance.RandomTestName()

	createConfig := testAccCheckDigitalOceanVPCNATGatewayConfig(name, "PUBLIC", 1)
	dataSourceIDConfig := `
data "digitalocean_vpc_nat_gateway" "foo" {
  id = digitalocean_vpc_nat_gateway.foobar.id
}`
	dataSourceNameConfig := `
data "digitalocean_vpc_nat_gateway" "foo" {
  name = digitalocean_vpc_nat_gateway.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVPCNATGatewayDestroy,
		Steps: []resource.TestStep{
			{
				// Test create
				Config: createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanVPCNATGatewayExists("digitalocean_vpc_nat_gateway.foobar", &gateway),
				),
			},
			{
				// Import by id
				Config: createConfig + dataSourceIDConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanVPCNATGatewayExists("digitalocean_vpc_nat_gateway.foobar", &gateway),
					resource.TestCheckResourceAttrSet("data.digitalocean_vpc_nat_gateway.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "type", "PUBLIC"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "size", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "vpcs.0.vpc_uuid"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "vpcs.0.gateway_ip"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "vpcs.0.default_gateway", "false"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "egresses.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "egresses.0.public_gateways.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "egresses.0.public_gateways.0.ipv4"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "udp_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "icmp_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "tcp_timeout_seconds", "30"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "updated_at"),
				),
			},
			{
				// Import by name
				Config: createConfig + dataSourceNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanVPCNATGatewayExists("digitalocean_vpc_nat_gateway.foobar", &gateway),
					resource.TestCheckResourceAttrSet("data.digitalocean_vpc_nat_gateway.foo", "id"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "type", "PUBLIC"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "size", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "vpcs.0.vpc_uuid"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "vpcs.0.gateway_ip"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "vpcs.0.default_gateway", "false"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "egresses.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "egresses.0.public_gateways.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "egresses.0.public_gateways.0.ipv4"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "udp_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "icmp_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_vpc_nat_gateway.foo", "tcp_timeout_seconds", "30"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_vpc_nat_gateway.foo", "updated_at"),
				),
			},
		},
	})
}
