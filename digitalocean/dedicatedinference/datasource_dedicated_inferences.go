package dedicatedinference

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanDedicatedInferences() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dedicatedInferenceListItemSchema(),
		ResultAttributeName: "dedicated_inferences",
		GetRecords:          getDigitalOceanDedicatedInferences,
		FlattenRecord:       flattenDigitalOceanDedicatedInferenceListItem,
	}

	return datalist.NewResource(dataListConfig)
}

func dedicatedInferenceListItemSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The unique ID of the dedicated inference endpoint.",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the dedicated inference endpoint.",
		},
		"region": {
			Type:        schema.TypeString,
			Description: "The region where the dedicated inference endpoint is deployed.",
		},
		"status": {
			Type:        schema.TypeString,
			Description: "The current status of the dedicated inference endpoint.",
		},
		"vpc_uuid": {
			Type:        schema.TypeString,
			Description: "The UUID of the VPC the dedicated inference endpoint is deployed in.",
		},
		"public_endpoint_fqdn": {
			Type:        schema.TypeString,
			Description: "The fully-qualified domain name of the public endpoint, if enabled.",
		},
		"private_endpoint_fqdn": {
			Type:        schema.TypeString,
			Description: "The fully-qualified domain name of the private endpoint.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "The date and time when the dedicated inference endpoint was created.",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Description: "The date and time when the dedicated inference endpoint was last updated.",
		},
		"provider_model_id": {
			Type:        schema.TypeList,
			Description: "The list of provider model IDs for the dedicated inference endpoint.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func getDigitalOceanDedicatedInferences(meta interface{}, _ map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	var allItems []godo.DedicatedInferenceListItem
	opts := &godo.DedicatedInferenceListOptions{
		ListOptions: godo.ListOptions{
			Page:    1,
			PerPage: 200,
		},
	}

	for {
		items, resp, err := client.DedicatedInference.List(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving dedicated inference endpoints: %s", err)
		}
		allItems = append(allItems, items...)
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error retrieving dedicated inference endpoints: %s", err)
		}
		opts.Page = page + 1
	}

	records := make([]interface{}, len(allItems))
	for i, item := range allItems {
		records[i] = item
	}
	return records, nil
}

func flattenDigitalOceanDedicatedInferenceListItem(record, _ interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
	item := record.(godo.DedicatedInferenceListItem)

	flat := map[string]interface{}{
		"id":                item.ID,
		"name":              item.Name,
		"region":            item.Region,
		"status":            item.Status,
		"vpc_uuid":          item.VPCUUID,
		"provider_model_id": item.ProviderModelID,
	}

	if item.Endpoints != nil {
		flat["public_endpoint_fqdn"] = item.Endpoints.PublicEndpointFQDN
		flat["private_endpoint_fqdn"] = item.Endpoints.PrivateEndpointFQDN
	} else {
		flat["public_endpoint_fqdn"] = ""
		flat["private_endpoint_fqdn"] = ""
	}

	if !item.CreatedAt.IsZero() {
		flat["created_at"] = item.CreatedAt.UTC().String()
	} else {
		flat["created_at"] = ""
	}
	if !item.UpdatedAt.IsZero() {
		flat["updated_at"] = item.UpdatedAt.UTC().String()
	} else {
		flat["updated_at"] = ""
	}

	return flat, nil
}
