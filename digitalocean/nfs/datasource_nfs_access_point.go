package nfs

import (
	"context"
	"fmt"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanNfsAccessPoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanNfsAccessPointRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name", "share_id"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"share_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"anonuid":                      {Type: schema.TypeInt, Computed: true},
					"anongid":                      {Type: schema.TypeInt, Computed: true},
					"protocols":                    {Type: schema.TypeSet, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					"squash_config":                {Type: schema.TypeString, Computed: true},
					"identity_enforcement_enabled": {Type: schema.TypeBool, Computed: true},
				}},
			},
		},
	}
}

func dataSourceDigitalOceanNfsAccessPointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if id, ok := d.GetOk("id"); ok {
		accessPoint, _, err := getNfsAccessPoint(ctx, client, id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving NFS access point: %s", err)
		}
		setNfsAccessPointData(d, accessPoint)
		return nil
	}

	name, hasName := d.GetOk("name")
	shareID, hasShareID := d.GetOk("share_id")
	if !hasName || !hasShareID {
		return diag.Errorf("Either `id` or both `name` and `share_id` must be set")
	}

	vpcID := ""
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcID = v.(string)
	}

	accessPoints, _, err := listNfsAccessPoints(ctx, client, shareID.(string), vpcID)
	if err != nil {
		return diag.Errorf("Error listing NFS access points: %s", err)
	}

	var selected *nfsAccessPoint
	for _, ap := range accessPoints {
		if ap != nil && ap.Name == name.(string) {
			if vpcID != "" && ap.VpcID != vpcID {
				continue
			}
			if selected != nil {
				return diag.Errorf("too many access points found with name %s on share %s", name.(string), shareID.(string))
			}
			selected = ap
		}
	}
	if selected == nil {
		return diag.FromErr(fmt.Errorf("no access point found with name %s on share %s", name.(string), shareID.(string)))
	}

	setNfsAccessPointData(d, selected)
	return nil
}

func setNfsAccessPointData(d *schema.ResourceData, ap *nfsAccessPoint) {
	d.SetId(ap.ID)
	d.Set("name", ap.Name)
	d.Set("share_id", ap.ShareID)
	d.Set("status", ap.Status)
	d.Set("path", ap.Path)
	d.Set("is_default", ap.IsDefault)
	d.Set("created_at", ap.CreatedAt)
	d.Set("updated_at", ap.UpdatedAt)
	d.Set("vpc_id", ap.VpcID)
	d.Set("access_policy", flattenNfsAccessPointPolicy(ap.AccessPolicy))
}
