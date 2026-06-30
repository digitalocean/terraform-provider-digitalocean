package nfs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
						DefaultFunc: func() (interface{}, error) {
							return schema.NewSet(schema.HashString, []interface{}{"NFS4"}), nil
						},
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

// GetNfsAccessPoint retrieves an NFS access point by ID for acceptance tests.
func GetNfsAccessPoint(ctx context.Context, client *godo.Client, accessPointID string) (*godo.NfsAccessPoint, *godo.Response, error) {
	return client.Nfs.GetAccessPoint(ctx, accessPointID)
}

func resourceDigitalOceanNfsAccessPointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	shareID := d.Get("share_id").(string)
	req := &godo.NfsCreateAccessPointRequest{
		Name:         d.Get("name").(string),
		Path:         d.Get("path").(string),
		VpcID:        d.Get("vpc_id").(string),
		AccessPolicy: expandNfsAccessPointPolicy(d.Get("access_policy")),
	}

	resp, _, err := client.Nfs.CreateAccessPoint(ctx, shareID, req)
	if err != nil {
		return diag.Errorf("Error creating NFS access point: %s", err)
	}
	if resp == nil || resp.AccessPoint == nil {
		return diag.Errorf("Error creating NFS access point: empty response")
	}

	d.SetId(resp.AccessPoint.ID)
	if err := waitForNfsAccessPointActive(ctx, client, d.Id()); err != nil {
		return diag.Errorf("Error waiting for NFS access point to become active: %s", err)
	}

	return resourceDigitalOceanNfsAccessPointRead(ctx, d, meta)
}

func resourceDigitalOceanNfsAccessPointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	accessPoint, resp, err := client.Nfs.GetAccessPoint(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving NFS access point: %s", err)
	}

	if nfsAccessPointRemoved(accessPoint) {
		d.SetId("")
		return nil
	}

	setNfsAccessPointResourceData(d, accessPoint)
	return nil
}

func resourceDigitalOceanNfsAccessPointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, _, err := client.Nfs.DeleteAccessPoint(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting NFS access point: %s", err)
	}

	if err := waitForNfsAccessPointDeleted(ctx, client, d.Id()); err != nil {
		return diag.Errorf("Error waiting for NFS access point to be deleted: %s", err)
	}

	d.SetId("")
	return nil
}

func nfsAccessPointRemoved(accessPoint *godo.NfsAccessPoint) bool {
	return accessPoint == nil || accessPoint.Status == godo.NfsAccessPointDeleted
}

func nfsAccessPointVpcID(accessPoint *godo.NfsAccessPoint) string {
	if accessPoint == nil || accessPoint.VpcID == nil {
		return ""
	}
	return *accessPoint.VpcID
}

func waitForNfsAccessPointActive(ctx context.Context, client *godo.Client, accessPointID string) error {
	conf := &retry.StateChangeConf{
		Pending: []string{string(godo.NfsAccessPointCreating)},
		Target:  []string{string(godo.NfsAccessPointActive)},
		Refresh: func() (interface{}, string, error) {
			accessPoint, resp, err := client.Nfs.GetAccessPoint(ctx, accessPointID)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil, "", fmt.Errorf("NFS access point %s not found while waiting to become active", accessPointID)
				}
				return nil, "", err
			}
			return accessPoint, string(accessPoint.Status), nil
		},
		Timeout:    15 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 15 * time.Second,
	}
	_, err := conf.WaitForStateContext(ctx)
	return err
}

func waitForNfsAccessPointDeleted(ctx context.Context, client *godo.Client, accessPointID string) error {
	conf := &retry.StateChangeConf{
		Pending: []string{
			string(godo.NfsAccessPointActive),
			string(godo.NfsAccessPointCreating),
		},
		Target: []string{string(godo.NfsAccessPointDeleted), "not-found"},
		Refresh: func() (interface{}, string, error) {
			accessPoint, resp, err := client.Nfs.GetAccessPoint(ctx, accessPointID)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return accessPoint, "not-found", nil
				}
				return nil, "", err
			}
			return accessPoint, string(accessPoint.Status), nil
		},
		Timeout:    15 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 15 * time.Second,
	}
	_, err := conf.WaitForStateContext(ctx)
	return err
}

func setNfsAccessPointResourceData(d *schema.ResourceData, accessPoint *godo.NfsAccessPoint) {
	d.Set("name", accessPoint.Name)
	d.Set("share_id", accessPoint.ShareID)
	d.Set("path", accessPoint.Path)
	d.Set("status", string(accessPoint.Status))
	d.Set("is_default", accessPoint.IsDefault)
	d.Set("created_at", accessPoint.CreatedAt)
	d.Set("updated_at", accessPoint.UpdatedAt)
	d.Set("vpc_id", nfsAccessPointVpcID(accessPoint))
	d.Set("access_policy", flattenNfsAccessPointPolicy(accessPoint.AccessPolicy))
}

func defaultNfsAccessPointProtocols() []godo.NfsAccessPolicyProtocol {
	return []godo.NfsAccessPolicyProtocol{godo.NfsAccessPolicyProtocolNFS4}
}

func expandNfsAccessPointPolicy(v interface{}) godo.NfsAccessPolicy {
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return godo.NfsAccessPolicy{
			Anonuid:                    65534,
			Anongid:                    65534,
			Protocols:                  defaultNfsAccessPointProtocols(),
			SquashConfig:               godo.NfsSquashConfigRootSquash,
			IdentityEnforcementEnabled: false,
		}
	}

	m, ok := list[0].(map[string]interface{})
	if !ok {
		return godo.NfsAccessPolicy{
			Anonuid:                    65534,
			Anongid:                    65534,
			Protocols:                  defaultNfsAccessPointProtocols(),
			SquashConfig:               godo.NfsSquashConfigRootSquash,
			IdentityEnforcementEnabled: false,
		}
	}

	policy := godo.NfsAccessPolicy{
		Anonuid:                    uint64(m["anonuid"].(int)),
		Anongid:                    uint64(m["anongid"].(int)),
		SquashConfig:               godo.NfsSquashConfig(strings.ToUpper(m["squash_config"].(string))),
		IdentityEnforcementEnabled: m["identity_enforcement_enabled"].(bool),
		Protocols:                  defaultNfsAccessPointProtocols(),
	}

	if protocols, ok := m["protocols"].(*schema.Set); ok && protocols != nil && protocols.Len() > 0 {
		policy.Protocols = make([]godo.NfsAccessPolicyProtocol, 0, protocols.Len())
		for _, p := range protocols.List() {
			policy.Protocols = append(policy.Protocols, godo.NfsAccessPolicyProtocol(strings.ToUpper(p.(string))))
		}
	}

	return policy
}

func flattenNfsAccessPointPolicy(policy godo.NfsAccessPolicy) []map[string]interface{} {
	protocolsList := policy.Protocols
	if len(protocolsList) == 0 {
		protocolsList = defaultNfsAccessPointProtocols()
	}

	protocols := make([]interface{}, 0, len(protocolsList))
	for _, p := range protocolsList {
		protocols = append(protocols, string(p))
	}

	return []map[string]interface{}{
		{
			"anonuid":                      int(policy.Anonuid),
			"anongid":                      int(policy.Anongid),
			"protocols":                    schema.NewSet(schema.HashString, protocols),
			"squash_config":                string(policy.SquashConfig),
			"identity_enforcement_enabled": policy.IdentityEnforcementEnabled,
		},
	}
}
