package digitalocean

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceDigitalOceanSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanSSHKeyRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the ssh key",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "public key part of the ssh key",
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "fingerprint of the ssh key",
			},
		},
	}
}

func dataSourceDigitalOceanSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	name := d.Get("name").(string)

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	keyList := []godo.Key{}

	for {
		keys, resp, err := client.Keys.List(context.Background(), opts)

		if err != nil {
			return fmt.Errorf("Error retrieving ssh keys: %s", err)
		}

		for _, key := range keys {
			keyList = append(keyList, key)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return fmt.Errorf("Error retrieving ssh keys: %s", err)
		}

		opts.Page = page + 1
	}

	key, err := findKeyByName(keyList, name)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(key.ID))
	d.Set("name", key.Name)
	d.Set("public_key", key.PublicKey)
	d.Set("fingerprint", key.Fingerprint)

	return nil
}

func findKeyByName(keys []godo.Key, name string) (*godo.Key, error) {
	results := make([]godo.Key, 0)
	for _, v := range keys {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no ssh key found with name %s", name)
	}
	return nil, fmt.Errorf("too many ssh keys found with name %s (found %d, expected 1)", name, len(results))
}
