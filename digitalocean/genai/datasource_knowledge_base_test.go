package genai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanKnowledgeBase_BasicByID(t *testing.T) {
	kbName := acceptance.RandomTestName() + "-kb"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_knowledge_base" "test" {
  name                 = "%s"
  project_id           = "%s"
  region               = "%s"
  embedding_model_uuid = "%s"
  tags                 = ["terraform-test", "datasource-test"]
  is_public            = false

  datasources {
    web_crawler_data_source {
      base_url        = "https://docs.digitalocean.com/products/kubernetes/"
      crawling_option = "SCOPED"
      embed_media     = true
    }
  }
}

data "digitalocean_genai_knowledge_base" "byid" {
  uuid = digitalocean_genai_knowledge_base.test.id
}
`, kbName, defaultProjectID, defaultRegion, defaultEmbeddingModelUUID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_genai_knowledge_base.byid", "id",
						"digitalocean_genai_knowledge_base.test", "id",
					),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_knowledge_base.byid", "name", kbName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_knowledge_base.byid", "project_id", defaultProjectID),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_knowledge_base.byid", "region", defaultRegion),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_knowledge_base.byid", "embedding_model_uuid", defaultEmbeddingModelUUID),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_knowledge_base.byid", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_knowledge_base.byid", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_knowledge_base.byid", "datasources.#", "1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_knowledge_base.byid", "datasources.0.web_crawler_data_source.0.base_url", "https://docs.digitalocean.com/products/kubernetes/"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_genai_knowledge_base.byid", "created_at"),
				),
			},
		},
	})
}
