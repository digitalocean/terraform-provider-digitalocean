package genai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanKnowledgeBase_Basic(t *testing.T) {
	kbName := acceptance.RandomTestName() + "-kb"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_knowledge_base" "test" {
  name                 = "%s"
  project_id           = "%s"
  region               = "%s"
  embedding_model_uuid = "%s"
  tags                 = ["terraform-test", "acceptance"]
  is_public            = false

  datasources {
    web_crawler_data_source {
      base_url        = "https://docs.digitalocean.com/products/kubernetes/"
      crawling_option = "SCOPED"
      embed_media     = true
    }
  }
}
`, kbName, defaultProjectID, defaultRegion, defaultEmbeddingModelUUID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "name", kbName),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "project_id", defaultProjectID),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "region", defaultRegion),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "embedding_model_uuid", defaultEmbeddingModelUUID),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "is_public", "false"),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "datasources.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "datasources.0.web_crawler_data_source.0.base_url", "https://docs.digitalocean.com/products/kubernetes/"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_knowledge_base.test", "created_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKnowledgeBase_Update(t *testing.T) {
	kbName := acceptance.RandomTestName() + "-kb"
	updatedKbName := kbName + "-up" // don't exceed 32 characters

	resourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_knowledge_base" "test" {
  name                 = "%s"
  project_id           = "%s"
  region               = "%s"
  embedding_model_uuid = "%s"
  tags                 = ["terraform-test", "update-test"]
  is_public            = false

  datasources {
    web_crawler_data_source {
      base_url        = "https://docs.digitalocean.com/products/kubernetes/"
      crawling_option = "SCOPED"
      embed_media     = true
    }
  }
}
`, kbName, defaultProjectID, defaultRegion, defaultEmbeddingModelUUID)

	updatedResourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_knowledge_base" "test" {
  name                 = "%s"
  project_id           = "%s"
  region               = "%s"
  embedding_model_uuid = "%s"
  tags                 = ["terraform-test", "update-test", "updated"]
  is_public            = false

  datasources {
    web_crawler_data_source {
      base_url        = "https://docs.digitalocean.com/products/kubernetes/"
      crawling_option = "SCOPED"
      embed_media     = true
    }
  }
}
`, updatedKbName, defaultProjectID, defaultRegion, defaultEmbeddingModelUUID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "name", kbName),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "is_public", "false"),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "tags.#", "2"),
				),
			},
			{
				Config: updatedResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "name", updatedKbName),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "is_public", "false"),
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "tags.#", "3"),
				),
			},
		},
	})
}

func TestAccDigitalOceanKnowledgeBase_Delete(t *testing.T) {
	kbName := acceptance.RandomTestName() + "-kb"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_knowledge_base" "test" {
  name                 = "%s"
  project_id           = "%s"
  region               = "%s"
  embedding_model_uuid = "%s"
  tags                 = ["terraform-test", "delete-test"]
  is_public            = false

  datasources {
    web_crawler_data_source {
      base_url        = "https://docs.digitalocean.com/products/kubernetes/"
      crawling_option = "SCOPED"
      embed_media     = true
    }
  }
}
`, kbName, defaultProjectID, defaultRegion, defaultEmbeddingModelUUID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_genai_knowledge_base.test", "name", kbName),
				),
			},
			{
				ResourceName:      "digitalocean_genai_knowledge_base.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
