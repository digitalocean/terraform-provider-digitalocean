package sshkey

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanSSHKeys() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        sshKeySchema(),
		ResultAttributeName: "ssh_keys",
		GetRecords:          getDigitalOceanSshKeys,
		FlattenRecord:       flattenDigitalOceanSshKey,
	}

	return datalist.NewResource(dataListConfig)
}
