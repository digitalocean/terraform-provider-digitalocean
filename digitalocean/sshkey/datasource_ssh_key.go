package sshkey

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanSSHKey() *schema.Resource {
	recordSchema := sshKeySchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["name"].Required = true
	recordSchema["name"].Computed = false

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanSSHKeyRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keyList, err := getDigitalOceanSshKeys(meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	key, err := findSshKeyByName(keyList, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedKey, err := flattenDigitalOceanSshKey(*key, meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := util.SetResourceDataFromMap(d, flattenedKey); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(key.ID))

	return nil
}

func findSshKeyByName(keys []interface{}, name string) (*godo.Key, error) {
	results := make([]godo.Key, 0)
	for _, v := range keys {
		key := v.(godo.Key)
		if key.Name == name {
			results = append(results, key)
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
