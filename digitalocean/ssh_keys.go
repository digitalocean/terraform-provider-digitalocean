package digitalocean

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func sshKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Description: "id of the ssh key",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "name of the ssh key",
		},
		"public_key": {
			Type:        schema.TypeString,
			Description: "public key part of the ssh key",
		},
		"fingerprint": {
			Type:        schema.TypeString,
			Description: "fingerprint of the ssh key",
		},
	}
}

func getDigitalOceanSshKeys(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*CombinedConfig).godoClient()

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var keyList []interface{}

	for {
		keys, resp, err := client.Keys.List(context.Background(), opts)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving ssh keys: %s", err)
		}

		for _, key := range keys {
			keyList = append(keyList, key)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving ssh keys: %s", err)
		}

		opts.Page = page + 1
	}

	return keyList, nil
}

func flattenDigitalOceanSshKey(rawSshKey, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	key := rawSshKey.(godo.Key)

	flattenedSshKey := map[string]interface{}{
		"id":          key.ID,
		"name":        key.Name,
		"fingerprint": key.Fingerprint,
		"public_key":  key.PublicKey,
	}

	return flattenedSshKey, nil
}
