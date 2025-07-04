package genai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	defaultProjectID          = "84e1e297-ee40-41ac-95ff-1067cf2206e9"
	defaultEmbeddingModelUUID = "22653204-79ed-11ef-bf8f-4e013e2ddde4"
	defaultRegion             = "tor1"
)

func TestAccDataSourceDigitalOceanKnowledgeBases_Basic(t *testing.T) {
	kbName1 := acceptance.RandomTestName() + "-kb1"
	kbName2 := acceptance.RandomTestName() + "-kb2"

	resourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_knowledge_base" "test1" {
  name                  = "%s"
  project_id            = "%s"
  region                = "%s"
  embedding_model_uuid  = "%s"
  tags                  = ["terraform-test", "datasource-test"]
  is_public             = false

  datasources {
    web_crawler_data_source {
      base_url         = "https://docs.digitalocean.com/products/kubernetes/"
      crawling_option  = "SCOPED"
      embed_media      = true
    }
  }
}

resource "digitalocean_genai_knowledge_base" "test2" {
  name                  = "%s"
  project_id            = "%s"
  region                = "%s"
  embedding_model_uuid  = "%s"
  tags                  = ["terraform-test", "datasource-test"]
  is_public             = true

  datasources {
    web_crawler_data_source {
      base_url         = "https://docs.digitalocean.com/products/app-platform/"
      crawling_option  = "SCOPED"
      embed_media      = false
    }
  }
}

data "digitalocean_genai_knowledge_bases" "all" {}
`, kbName1, defaultProjectID, defaultRegion, defaultEmbeddingModelUUID,
		kbName2, defaultProjectID, defaultRegion, defaultEmbeddingModelUUID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_knowledge_bases.all", "knowledge_bases.#"),
				),
			},
		},
	})
}
