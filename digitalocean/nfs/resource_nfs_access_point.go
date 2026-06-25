package nfs

import (
	"context"
	"strings"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanNfsAccessPoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanNfsAccessPointCreate,
		ReadContext:   resourceDigitalOceanNfsAccessPointRead,
		DeleteContext: resourceDigitalOceanNfsAccessPointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"status": {
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
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"anonuid": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  65534,
					},
					"anongid": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  65534,
					},
					"protocols": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
						MinItems: 1,
					},
					"squash_config": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "ROOT_SQUASH",
						ValidateFunc: validation.StringInSlice([]string{"NO_SQUASH", "ROOT_SQUASH", "ALL_SQUASH"}, false),
					},
					"identity_enforcement_enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				}},
			},
		},
	}
}

func resourceDigitalOceanNfsAccessPointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	shareID := d.Get("share_id").(string)
	req := &nfsCreateAccessPointRequest{
		Name:         d.Get("name").(string),
		Path:         d.Get("path").(string),
		VpcID:        d.Get("vpc_id").(string),
		AccessPolicy: expandNfsAccessPointPolicy(d.Get("access_policy")),
	}

	accessPoint, _, err := createNfsAccessPoint(ctx, client, shareID, req)
	if err != nil {
		return diag.Errorf("Error creating NFS access point: %s", err)
	}

	d.SetId(accessPoint.ID)
	return resourceDigitalOceanNfsAccessPointRead(ctx, d, meta)
}

func resourceDigitalOceanNfsAccessPointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	accessPoint, resp, err := getNfsAccessPoint(ctx, client, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving NFS access point: %s", err)
	}

	d.Set("name", accessPoint.Name)
	d.Set("share_id", accessPoint.ShareID)
	d.Set("path", accessPoint.Path)
	d.Set("status", accessPoint.Status)
	d.Set("is_default", accessPoint.IsDefault)
	d.Set("created_at", accessPoint.CreatedAt)
	d.Set("updated_at", accessPoint.UpdatedAt)
	d.Set("vpc_id", accessPoint.VpcID)
	d.Set("access_policy", flattenNfsAccessPointPolicy(accessPoint.AccessPolicy))

	return nil
}

func resourceDigitalOceanNfsAccessPointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, _, err := deleteNfsAccessPoint(ctx, client, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting NFS access point: %s", err)
	}

	d.SetId("")
	return nil
}

func expandNfsAccessPointPolicy(v interface{}) nfsAccessPointPolicy {
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nfsAccessPointPolicy{Anonuid: 65534, Anongid: 65534, Protocols: []string{"NFS4"}, SquashConfig: "ROOT_SQUASH"}
	}

	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nfsAccessPointPolicy{Anonuid: 65534, Anongid: 65534, Protocols: []string{"NFS4"}, SquashConfig: "ROOT_SQUASH"}
	}

	policy := nfsAccessPointPolicy{
		Anonuid:                    uint64(m["anonuid"].(int)),
		Anongid:                    uint64(m["anongid"].(int)),
		SquashConfig:               strings.ToUpper(m["squash_config"].(string)),
		IdentityEnforcementEnabled: m["identity_enforcement_enabled"].(bool),
		Protocols:                  []string{"NFS4"},
	}

	if protocols, ok := m["protocols"].(*schema.Set); ok && protocols != nil && protocols.Len() > 0 {
		policy.Protocols = make([]string, 0, protocols.Len())
		for _, p := range protocols.List() {
			policy.Protocols = append(policy.Protocols, strings.ToUpper(p.(string)))
		}
	}

	return policy
}

func flattenNfsAccessPointPolicy(policy nfsAccessPointPolicy) []map[string]interface{} {
	protocols := make([]interface{}, 0, len(policy.Protocols))
	for _, p := range policy.Protocols {
		protocols = append(protocols, p)
	}

	return []map[string]interface{}{
		{
			"anonuid":                      int(policy.Anonuid),
			"anongid":                      int(policy.Anongid),
			"protocols":                    schema.NewSet(schema.HashString, protocols),
			"squash_config":                policy.SquashConfig,
			"identity_enforcement_enabled": policy.IdentityEnforcementEnabled,
		},
	}
}
