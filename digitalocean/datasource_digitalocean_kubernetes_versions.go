package digitalocean

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDigitalOceanKubernetesVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanKubernetesVersionsRead,
		Schema: map[string]*schema.Schema{
			"version_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDigitalOceanKubernetesVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	k8sOptions, _, err := client.Kubernetes.GetOptions(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving Kubernetes options: %s", err)
	}

	d.SetId(resource.UniqueId())

	validVersions := make([]string, 0)
	for _, v := range k8sOptions.Versions {
		if strings.HasPrefix(v.Slug, d.Get("version_prefix").(string)) {
			validVersions = append(validVersions, v.Slug)
		}
	}
	d.Set("valid_versions", validVersions)

	if len(validVersions) > 0 {
		d.Set("latest_version", validVersions[0])
	}

	return nil
}
