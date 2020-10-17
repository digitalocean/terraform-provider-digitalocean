package digitalocean

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDigitalOceanKubernetesVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanKubernetesVersionsRead,
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

func dataSourceDigitalOceanKubernetesVersionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	k8sOptions, _, err := client.Kubernetes.GetOptions(context.Background())
	if err != nil {
		return diag.Errorf("Error retrieving Kubernetes options: %s", err)
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
