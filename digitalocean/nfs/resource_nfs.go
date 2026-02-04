package nfs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanNfs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanNfsCreate,
		ReadContext:   resourceDigitalOceanNfsRead,
		UpdateContext: resourceDigitalOceanNfsUpdate,
		DeleteContext: resourceDigitalOceanNfsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanNfsImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(50),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": { // API requires at least one VPC ID
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host IP of the NFS server accessible from the associated VPC",
			},
			"mount_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mount path for accessing the NFS share",
			},

			"tags": tag.TagsSchema(),
		},

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {

			// if the new size of the share is smaller than the old one return an error since
			// only expanding the share is allowed
			oldSize, newSize := diff.GetChange("size")
			if newSize.(int) < oldSize.(int) {
				return fmt.Errorf("share `size` can only be expanded and not shrunk")
			}

			return nil
		},
	}
}

func resourceDigitalOceanNfsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.NfsCreateRequest{
		Name:    d.Get("name").(string),
		SizeGib: d.Get("size").(int),
		Region:  d.Get("region").(string),
	}
	if v, ok := d.GetOk("region"); ok {
		opts.Region = strings.ToLower(v.(string))
	}
	if v, ok := d.GetOk("size"); ok {
		opts.SizeGib = v.(int)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		opts.VpcIDs = []string{v.(string)}
	}

	log.Printf("[DEBUG] Nfs create configuration: %#v", opts)
	share, _, err := client.Nfs.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating Share: %s", err)
	}

	d.SetId(share.ID)
	log.Printf("[INFO] Share name: %s", share.Name)

	// Wait for share to become ACTIVE so host and mount_path are populated
	err = waitForNfsActive(ctx, client, share.ID, opts.Region)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDigitalOceanNfsRead(ctx, d, meta)
}

func resourceDigitalOceanNfsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	region := strings.ToLower(d.Get("region").(string))

	if d.HasChange("size") {
		size := d.Get("size").(int)

		_, _, err := client.NfsActions.Resize(ctx, d.Id(), uint64(size), region)
		if err != nil {
			return diag.FromErr(err)
		}

		// Wait for resize to complete
		err = waitForNfsResize(ctx, client, d.Id(), d.Get("region").(string), d.Get("size").(int))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDigitalOceanNfsRead(ctx, d, meta)
}

func waitForNfsResize(ctx context.Context, client *godo.Client, id, region string, expectedSize int) error {
	for i := 0; i < 60; i++ {
		share, _, err := client.Nfs.Get(ctx, id, region)
		if err != nil {
			return err
		}

		if share.SizeGib == expectedSize && share.Status == "ACTIVE" {
			return nil
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for NFS resize to complete")
}

func waitForNfsActive(ctx context.Context, client *godo.Client, id, region string) error {
	for i := 0; i < 60; i++ {
		share, _, err := client.Nfs.Get(ctx, id, region)
		if err != nil {
			return err
		}

		if share.Status == "ACTIVE" {
			return nil
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for NFS share to become active")
}

func resourceDigitalOceanNfsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	share, resp, err := client.Nfs.Get(context.Background(), d.Id(), d.Get("region").(string))
	if err != nil {
		// If the share is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving share: %s", err)
	}

	d.Set("name", share.Name)
	d.Set("region", share.Region)
	d.Set("size", share.SizeGib)
	d.Set("status", share.Status)
	d.Set("vpc_id", share.VpcIDs[0])
	d.Set("host", share.Host)
	d.Set("mount_path", share.MountPath)

	if err = d.Set("vpc_ids", flattenDigitalOceanShareVpcIds(share.VpcIDs)); err != nil {
		return diag.Errorf("[DEBUG] Error setting vpc_ids: %#v", err)
	}

	return nil
}

func resourceDigitalOceanNfsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Deleting share: %s", d.Id())
	_, err := client.Nfs.Delete(context.Background(), d.Id(), d.Get("region").(string))
	if err != nil {
		return diag.Errorf("Error deleting share: %s", err)
	}

	d.SetId("")
	return nil
}

func flattenDigitalOceanShareVpcIds(vpcs []string) *schema.Set {
	flattenedVpcs := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range vpcs {
		flattenedVpcs.Add(v)
	}

	return flattenedVpcs
}

func resourceDigitalOceanNfsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		d.Set("region", s[1])
	}

	// Verify the resource exists before calling Read
	client := meta.(*config.CombinedConfig).GodoClient()
	region := d.Get("region").(string)

	_, resp, err := client.Nfs.Get(ctx, d.Id(), region)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil, fmt.Errorf("NFS share %s not found in region %s. Please verify the ID is correct", d.Id(), region)
		}
		return nil, fmt.Errorf("error retrieving NFS share: %s", err)
	}

	// Trigger a Read to populate all attributes including vpc_id
	diags := resourceDigitalOceanNfsRead(ctx, d, meta)
	if diags.HasError() {
		// Convert diag.Diagnostics to error
		return nil, fmt.Errorf("error reading NFS during import: %s", diags[0].Summary)
	}
	return []*schema.ResourceData{d}, nil
}
