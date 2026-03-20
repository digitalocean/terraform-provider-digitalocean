package dedicatedinference

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanDedicatedInferenceTokens() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dedicatedInferenceTokenSchema(),
		ResultAttributeName: "tokens",
		ExtraQuerySchema: map[string]*schema.Schema{
			"dedicated_inference_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of the dedicated inference endpoint to list tokens for.",
				ValidateFunc: validation.NoZeroValues,
			},
		},
		GetRecords:    getDigitalOceanDedicatedInferenceTokens,
		FlattenRecord: flattenDigitalOceanDedicatedInferenceToken,
	}

	return datalist.NewResource(dataListConfig)
}

func dedicatedInferenceTokenSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The unique ID of the token.",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the token.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "The date and time when the token was created.",
		},
	}
}

func getDigitalOceanDedicatedInferenceTokens(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	diID := extra["dedicated_inference_id"].(string)

	var allTokens []godo.DedicatedInferenceToken
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		tokens, resp, err := client.DedicatedInference.ListTokens(context.Background(), diID, opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving tokens for dedicated inference endpoint (%s): %s", diID, err)
		}
		allTokens = append(allTokens, tokens...)
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error retrieving tokens for dedicated inference endpoint (%s): %s", diID, err)
		}
		opts.Page = page + 1
	}

	records := make([]interface{}, len(allTokens))
	for i, token := range allTokens {
		records[i] = token
	}
	return records, nil
}

func flattenDigitalOceanDedicatedInferenceToken(record, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	token := record.(godo.DedicatedInferenceToken)

	flat := map[string]interface{}{
		"id":   token.ID,
		"name": token.Name,
	}

	if !token.CreatedAt.IsZero() {
		flat["created_at"] = token.CreatedAt.UTC().String()
	} else {
		flat["created_at"] = ""
	}

	return flat, nil
}
