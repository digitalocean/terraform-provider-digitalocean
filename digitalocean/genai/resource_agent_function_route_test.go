package genai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanAgentFunctionRoute_Basic(t *testing.T) {
	var agent godo.Agent
	function_name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentFunctionRouteConfig_basic(function_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_function.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_function.test", "function_name", function_name),
					resource.TestCheckResourceAttr("digitalocean_genai_function.test", "description", "Adding a function route and this will also tell temperature"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "faas_name"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "faas_namespace"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "input_schema"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgentFunctionRoute_WithOptionalFields(t *testing.T) {
	var agent godo.Agent
	function_name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentFunctionRouteConfig_withOptionalFields(function_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_function.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_function.test", "function_name", function_name),
					resource.TestCheckResourceAttr("digitalocean_genai_function.test", "description", "Adding a function route with optional fields"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "faas_name"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "faas_namespace"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "input_schema"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "output_schema"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgentFunctionRoute_Update(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()
	updatedName := agentName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentFunctionRouteConfig_basic(agentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_agent.test", "name", agentName),
					resource.TestCheckResourceAttr("digitalocean_agent.test", "instruction", "You are a helpful AI assistant."),
				),
			},
			{
				Config: testAccCheckDigitalOceanAgentFunctionRouteConfig_updated(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_agent.test", "name", updatedName),
					resource.TestCheckResourceAttr("digitalocean_agent.test", "instruction", "You are an updated AI assistant with new capabilities."),
					resource.TestCheckResourceAttr("digitalocean_agent.test", "description", "Updated test agent"),
					resource.TestCheckResourceAttr("digitalocean_agent.test", "temperature", "0.8"),
					resource.TestCheckResourceAttr("digitalocean_agent.test", "max_tokens", "2000"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanAgentFunctionRouteConfig_basic(function_name string) string {
	agentId := "00000000-0000-0000-0000-000000000000"
	description := "Adding a function route and this will also tell temperature"
	faasName := "default/testing"
	faasNamespace := "fn-b90faf52-2b42-49c2-9792-75edfbb6f397"
	functionName := "terraform-tf-complete"
	inputSchema := `{
		"parameters": [
			{
				"in": "query",
				"name": "zipCode",
				"schema": {
					"type": "string"
				},
				"required": false,
				"description": "The ZIP code for which to fetch the weather"
			},
			{
				"name": "measurement",
				"schema": {
					"enum": [
						"F",
						"C"
					],
					"type": "string"
				},
				"required": false,
				"description": "The measurement unit for temperature (F or C)",
				"in": "query"
			}
		]
	}`
	return fmt.Sprintf(`
resource "digitalocean_genai_function" "check" {
  agent_id       = "%s"
  description    = "%s"
  faas_name      = "%s"
  faas_namespace = "%s"
  function_name  = "%s"
  input_schema   = <<EOF
%s
	EOF
}
`, agentId, description, faasName, faasNamespace, functionName, inputSchema)
}

func testAccCheckDigitalOceanAgentFunctionRouteConfig_withOptionalFields(function_name string) string {
	agentId := "00000000-0000-0000-0000-000000000000"
	description := "Adding a function route and this will also tell temperature"
	faasName := "default/testing"
	faasNamespace := "fn-b90faf52-2b42-49c2-9792-75edfbb6f397"
	functionName := "terraform-tf-complete"
	inputSchema := `{
		"parameters": [
			{
				"in": "query",
				"name": "zipCode",
				"schema": {
					"type": "string"
				},
				"required": false,
				"description": "The ZIP code for which to fetch the weather"
			},
			{
				"name": "measurement",
				"schema": {
					"enum": [
						"F",
						"C"
					],
					"type": "string"
				},
				"required": false,
				"description": "The measurement unit for temperature (F or C)",
				"in": "query"
			}
		]
	}`

	outputSchema := `{
		"properties": [
			{
				"name": "temperature",
				"type": "number",
				"description": "The temperature for the specified location"
			},
			{
				"name": "measurement",
				"type": "string",
				"description": "The measurement unit used for the temperature (F or C)"
			},
			{
				"name": "conditions",
				"type": "string",
				"description": "A description of the current weather conditions (Sunny, Cloudy, etc)"
			}
		]
	}`

	return fmt.Sprintf(`
resource "digitalocean_genai_function" "check" {
  agent_id       = "%s"
  description    = "%s"
  faas_name      = "%s"
  faas_namespace = "%s"
  function_name  = "%s"
  input_schema   = <<EOF
%s
	EOF
  output_schema  = <<EOF
%s
	EOF
}
`, agentId, description, faasName, faasNamespace, functionName, inputSchema, outputSchema)
}

func testAccCheckDigitalOceanAgentFunctionRouteConfig_updated(function_name string) string {
	agentId := "00000000-0000-0000-0000-000000000000"
	description := "Adding a function route and this will also tell temperature"
	faasName := "default/testing"
	faasNamespace := "fn-b90faf52-2b42-49c2-9792-75edfbb6f397"
	functionName := "terraform-tf"
	inputSchema := `{
		"parameters": [
			{
				"in": "query",
				"name": "zipCode",
				"schema": {
					"type": "string"
				},
				"required": false,
				"description": "The ZIP code for which to fetch the weather"
			},
			{
				"name": "measurement",
				"schema": {
					"enum": [
						"F",
						"C"
					],
					"type": "string"
				},
				"required": false,
				"description": "The measurement unit for temperature (F or C)",
				"in": "query"
			}
		]
	}`

	outputSchema := `{
		"properties": [
			{
				"name": "temperature",
				"type": "number",
				"description": "The temperature for the specified location"
			},
			{
				"name": "measurement",
				"type": "string",
				"description": "The measurement unit used for the temperature (F or C)"
			},
			{
				"name": "conditions",
				"type": "string",
				"description": "A description of the current weather conditions (Sunny, Cloudy, etc)"
			}
		]
	}`

	return fmt.Sprintf(`
resource "digitalocean_genai_function" "check" {
  agent_id       = "%s"
  description    = "%s"
  faas_name      = "%s"
  faas_namespace = "%s"
  function_name  = "%s"
  input_schema   = <<EOF
%s
	EOF

  output_schema = <<EOF
%s
	EOF
}
`, agentId, description, faasName, faasNamespace, functionName, inputSchema, outputSchema)
}
