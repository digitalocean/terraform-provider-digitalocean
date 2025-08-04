package vpcnatgateway_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDigitalOceanVPCNATGateway(t *testing.T) {
	var gateway godo.VPCNATGateway
	name := acceptance.RandomTestName()

	createConfig := testAccCheckDigitalOceanVPCNATGatewayConfig(name, "PUBLIC", 1)
	// Update name and timeouts
	updateConfig := strings.ReplaceAll(createConfig, name, fmt.Sprintf("%s-updated", name))
	updateConfig = strings.ReplaceAll(updateConfig, "udp_timeout_seconds = 30", "udp_timeout_seconds = 60")
	updateConfig = strings.ReplaceAll(updateConfig, "icmp_timeout_seconds = 30", "icmp_timeout_seconds = 60")
	updateConfig = strings.ReplaceAll(updateConfig, "tcp_timeout_seconds = 30", "tcp_timeout_seconds = 60")

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
					resource.TestCheckResourceAttrSet("digitalocean_vpc_nat_gateway.foobar", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "type", "PUBLIC"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "size", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "vpcs.0.vpc_uuid"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "vpcs.0.gateway_ip"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "vpcs.0.default_gateway", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "egresses.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "egresses.0.public_gateways.#", "1"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "egresses.0.public_gateways.0.ipv4"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "udp_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "icmp_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "tcp_timeout_seconds", "30"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "updated_at"),
				),
			},
			{
				// Test update (name and timeout values)
				Config: updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanVPCNATGatewayExists("digitalocean_vpc_nat_gateway.foobar", &gateway),
					resource.TestCheckResourceAttrSet("digitalocean_vpc_nat_gateway.foobar", "id"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "name", fmt.Sprintf("%s-updated", name)),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "type", "PUBLIC"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "size", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "vpcs.0.vpc_uuid"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "vpcs.0.gateway_ip"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "vpcs.0.default_gateway", "false"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "egresses.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "egresses.0.public_gateways.#", "1"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "egresses.0.public_gateways.0.ipv4"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "udp_timeout_seconds", "60"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "icmp_timeout_seconds", "60"),
					resource.TestCheckResourceAttr(
						"digitalocean_vpc_nat_gateway.foobar", "tcp_timeout_seconds", "60"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_vpc_nat_gateway.foobar", "updated_at"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVPCNATGatewayConfig(name, gatewayType string, size int) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "foo" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_vpc_nat_gateway" "foobar" {
  name   = "%s"
  type   = "%s"
  region = "nyc3"
  size   = "%d"
  vpcs {
    vpc_uuid = digitalocean_vpc.foo.id
  }
  udp_timeout_seconds  = 30
  icmp_timeout_seconds = 30
  tcp_timeout_seconds  = 30
}`,
		fmt.Sprintf("test-%s-vpc", name),
		name,
		gatewayType,
		size)
}

func testAccCheckDigitalOceanVPCNATGatewayDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_vpc_nat_gateway" {
			continue
		}
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		_, _, err := client.VPCNATGateways.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("nat gateway with id %s not found", rs.Primary.ID)) {
				return nil
			}
			return fmt.Errorf("VPC NAT Gateway still exists")
		}
	}
	return nil
}

func testAccCheckDigitalOceanVPCNATGatewayExists(resourceName string, gateway *godo.VPCNATGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %v", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set")
		}
		// Check for valid ID response to validate that the resource has been created
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		gotGateway, _, err := client.VPCNATGateways.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if gotGateway.ID != rs.Primary.ID {
			return fmt.Errorf("VPC NAT gateway not found")
		}
		*gateway = *gotGateway
		return nil
	}
}
