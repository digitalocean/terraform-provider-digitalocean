package genai_test

// import (
// 	"fmt"
// 	"strings"
// 	"testing"

// 	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// )

// const (
// 	initialInstr = "You are a helpful AI assistant."
// 	updatedInstr = "You are an even more helpful AI assistant (v2)."
// )

// func TestAccDataSourceDigitalOceanAgentVersions_Lifecycle(t *testing.T) {
// 	name := acceptance.RandomTestName() + "-agent"

// 	createCfg := testAccAgentConfig(name, initialInstr)
// 	updateCfg := testAccAgentConfig(name, updatedInstr)

// 	const dsCfg = `
// data "digitalocean_genai_agent_versions" "versions" {
//   agent_id = digitalocean_genai_agent.foo.id
// }`

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
// 		ProviderFactories: acceptance.TestAccProviderFactories,

// 		Steps: []resource.TestStep{
// 			{Config: createCfg},

// 			{
// 				Config: createCfg + dsCfg,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("data.digitalocean_genai_agent_versions.versions",
// 						"agent_versions.#", "1"),
// 					checkExactlyOneApplied("data.digitalocean_genai_agent_versions.versions"),
// 				),
// 			},
// 			{
// 				Config: updateCfg + dsCfg,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("data.digitalocean_genai_agent_versions.versions",
// 						"agent_versions.#", "2"),
// 					checkExactlyOneApplied("data.digitalocean_genai_agent_versions.versions"),
// 					resource.TestCheckResourceAttr(
// 						"data.digitalocean_genai_agent_versions.versions",
// 						"agent_versions.0.instruction", updatedInstr),
// 				),
// 			},
// 			{Config: updateCfg},
// 		},
// 	})
// }

// // Generates an agent resource with a paramaterised instruction string.
// func testAccAgentConfig(name, instruction string) string {
// 	return `
// resource "digitalocean_project" "test" {
//   name = "` + name + `-project"
// }

// resource "digitalocean_genai_agent" "foo" {
//   name        = "` + name + `"
//   instruction = "` + instruction + `"
//   model_uuid  = "` + defaultModelUUID + `"
//   project_id  = digitalocean_project.test.id
//   region      = "tor1"
// }`
// }

// // Custom check: ensure **exactly one** version has currently_applied=true.
// func checkExactlyOneApplied(resName string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[resName]
// 		if !ok {
// 			return fmt.Errorf("resource %s not found in state", resName)
// 		}
// 		applied := 0
// 		for k, v := range rs.Primary.Attributes {
// 			if strings.HasSuffix(k, ".currently_applied") && v == "true" {
// 				applied++
// 			}
// 		}
// 		if applied != 1 {
// 			return fmt.Errorf("expected exactly 1 currently_applied=true, got %d", applied)
// 		}
// 		return nil
// 	}
// }
