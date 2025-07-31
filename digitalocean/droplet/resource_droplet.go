package droplet

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	errDropletBackupPolicy = errors.New("backup_policy can only be set when backups are enabled")
)

func ResourceDigitalOceanDroplet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDropletCreate,
		ReadContext:   resourceDigitalOceanDropletRead,
		UpdateContext: resourceDigitalOceanDropletUpdate,
		DeleteContext: resourceDigitalOceanDropletDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDropletImport,
		},
		MigrateState:  ResourceDigitalOceanDropletMigrateState,
		SchemaVersion: 1,

		// We are using these timeouts to be the minimum timeout for an operation.
		// This is how long an operation will wait for a state update, however
		// implementation of updates and deletes contain multiple instances of waiting for a state update
		// so the true timeout of an operation could be a multiple of the set value.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"image": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},

			"size": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 size slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},

			"graceful_shutdown": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"price_hourly": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"price_monthly": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"resize_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"backups": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"backup_policy": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				RequiredWith: []string{"backups"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"daily",
								"weekly",
							}, false),
						},
						"weekday": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT",
							}, false),
						},
						"hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 20),
						},
					},
				},
			},

			"ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ipv6_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"private_networking": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "This parameter has been deprecated. Use `vpc_uuid` instead to specify a VPC network for the Droplet. If no `vpc_uuid` is provided, the Droplet will be placed in your account's default VPC for the region.",
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_address_private": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ssh_keys": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},

			"user_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				StateFunc:    util.HashStringStateFunc(),
				// In order to support older statefiles with fully saved user data
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new != "" && old == d.Get("user_data")
				},
			},

			"volume_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},

			"monitoring": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"droplet_agent": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": tag.TagsSchema(),

			"vpc_uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},

		CustomizeDiff: customdiff.All(
			// If the `ipv6` attribute is changed to `true`, we need to mark the
			// `ipv6_address` attribute as changing in the plan. If not, the plan
			// will become inconsistent once the address is known when referenced
			// in another resource such as a domain record, e.g.:
			// https://github.com/digitalocean/terraform-provider-digitalocean/issues/981
			customdiff.IfValueChange("ipv6",
				func(ctx context.Context, old, new, meta interface{}) bool {
					return !old.(bool) && new.(bool)
				},
				customdiff.ComputedIf("ipv6_address", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.Get("ipv6").(bool)
				}),
			),
			// Forces replacement when IPv6 has attribute changes to `false`
			// https://github.com/digitalocean/terraform-provider-digitalocean/issues/1104
			customdiff.ForceNewIfChange("ipv6",
				func(ctx context.Context, old, new, meta interface{}) bool {
					return old.(bool) && !new.(bool)
				},
			),
		),
	}
}

func resourceDigitalOceanDropletCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	image := d.Get("image").(string)

	// Build up our creation options
	opts := &godo.DropletCreateRequest{
		Image:     godo.DropletCreateImage{},
		Name:      d.Get("name").(string),
		Region:    d.Get("region").(string),
		Size:      d.Get("size").(string),
		ProjectID: d.Get("project_id").(string),
		Tags:      tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
	}

	imageId, err := strconv.Atoi(image)
	if err == nil {
		// The image field is provided as an ID (number).
		opts.Image.ID = imageId
	} else {
		opts.Image.Slug = image
	}

	if attr, ok := d.GetOk("backups"); ok {
		_, exist := d.GetOk("backup_policy")
		if exist && !attr.(bool) { // Check there is no backup_policy specified when backups are disabled.
			return diag.FromErr(errDropletBackupPolicy)
		}
		opts.Backups = attr.(bool)
	}

	if attr, ok := d.GetOk("ipv6"); ok {
		opts.IPv6 = attr.(bool)
	}

	if attr, ok := d.GetOk("private_networking"); ok {
		opts.PrivateNetworking = attr.(bool)
	}

	if attr, ok := d.GetOk("user_data"); ok {
		opts.UserData = attr.(string)
	}

	if attr, ok := d.GetOk("volume_ids"); ok {
		for _, id := range attr.(*schema.Set).List() {
			if id == nil {
				continue
			}
			volumeId := id.(string)
			if volumeId == "" {
				continue
			}

			opts.Volumes = append(opts.Volumes, godo.DropletCreateVolume{
				ID: volumeId,
			})
		}
	}

	if attr, ok := d.GetOk("monitoring"); ok {
		opts.Monitoring = attr.(bool)
	}

	if attr, ok := d.GetOkExists("droplet_agent"); ok {
		opts.WithDropletAgent = godo.PtrTo(attr.(bool))
	}

	if attr, ok := d.GetOk("vpc_uuid"); ok {
		opts.VPCUUID = attr.(string)
	}

	// Get configured ssh_keys
	if v, ok := d.GetOk("ssh_keys"); ok {
		expandedSshKeys, err := expandSshKeys(v.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		opts.SSHKeys = expandedSshKeys
	}

	// Get configured backup_policy
	if policy, ok := d.GetOk("backup_policy"); ok {
		if !d.Get("backups").(bool) {
			return diag.FromErr(errDropletBackupPolicy)
		}

		backupPolicy, err := expandBackupPolicy(policy)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.BackupPolicy = backupPolicy
	}

	log.Printf("[DEBUG] Droplet create configuration: %#v", opts)

	droplet, _, err := client.Droplets.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating droplet: %s", err)
	}

	// Assign the droplets id
	d.SetId(strconv.Itoa(droplet.ID))
	log.Printf("[INFO] Droplet ID: %s", d.Id())

	// Ensure Droplet status has moved to "active."
	_, err = waitForDropletAttribute(ctx, d, "active", []string{"new"}, "status", schema.TimeoutCreate, meta)
	if err != nil {
		return diag.Errorf("Error waiting for droplet (%s) to become ready: %s", d.Id(), err)
	}

	// waitForDropletAttribute updates the Droplet's state and calls setDropletAttributes.
	// So there is no need to call resourceDigitalOceanDropletRead and add additional API calls.
	return nil
}

func resourceDigitalOceanDropletRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid droplet id: %v", err)
	}

	// Retrieve the droplet properties for updating the state
	droplet, resp, err := client.Droplets.Get(context.Background(), id)
	if err != nil {
		// check if the droplet no longer exists.
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Droplet (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving droplet: %s", err)
	}

	err = setDropletAttributes(d, droplet)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setDropletAttributes(d *schema.ResourceData, droplet *godo.Droplet) error {
	// Note that the image attribute is not set here. It is intentionally allowed
	// to drift once the Droplet has been created. This is to workaround the fact that
	// image slugs can move to point to images with a different ID. Image slugs are helpers
	// that always point to the most recent version of an image.
	// See: https://github.com/digitalocean/terraform-provider-digitalocean/issues/152
	d.Set("name", droplet.Name)
	d.Set("project_id", droplet.ProjectID)
	d.Set("urn", droplet.URN())
	d.Set("region", droplet.Region.Slug)
	d.Set("size", droplet.Size.Slug)
	d.Set("price_hourly", droplet.Size.PriceHourly)
	d.Set("price_monthly", droplet.Size.PriceMonthly)
	d.Set("disk", droplet.Disk)
	d.Set("vcpus", droplet.Vcpus)
	d.Set("memory", droplet.Memory)
	d.Set("status", droplet.Status)
	d.Set("locked", droplet.Locked)
	d.Set("created_at", droplet.Created)
	d.Set("vpc_uuid", droplet.VPCUUID)

	d.Set("ipv4_address", FindIPv4AddrByType(droplet, "public"))
	d.Set("ipv4_address_private", FindIPv4AddrByType(droplet, "private"))
	d.Set("ipv6_address", strings.ToLower(FindIPv6AddrByType(droplet, "public")))

	if features := droplet.Features; features != nil {
		d.Set("backups", containsDigitalOceanDropletFeature(features, "backups"))
		d.Set("ipv6", containsDigitalOceanDropletFeature(features, "ipv6"))
		d.Set("private_networking", containsDigitalOceanDropletFeature(features, "private_networking"))
		d.Set("monitoring", containsDigitalOceanDropletFeature(features, "monitoring"))
	}

	if err := d.Set("volume_ids", flattenDigitalOceanDropletVolumeIds(droplet.VolumeIDs)); err != nil {
		return fmt.Errorf("Error setting `volume_ids`: %+v", err)
	}

	if err := d.Set("tags", tag.FlattenTags(droplet.Tags)); err != nil {
		return fmt.Errorf("Error setting `tags`: %+v", err)
	}

	// Initialize the connection info
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": FindIPv4AddrByType(droplet, "public"),
	})

	return nil
}

func resourceDigitalOceanDropletImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Retrieve the image from API during import
	client := meta.(*config.CombinedConfig).GodoClient()
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("Invalid droplet id: %v", err)
	}

	droplet, _, err := client.Droplets.Get(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("Error importing droplet: %s", err)
	}

	if droplet.Image.Slug != "" {
		d.Set("image", droplet.Image.Slug)
	} else {
		d.Set("image", godo.Stringify(droplet.Image.ID))
	}

	// This is a non API attribute. So set to the default setting in the schema.
	d.Set("resize_disk", true)

	return []*schema.ResourceData{d}, nil
}

func FindIPv6AddrByType(d *godo.Droplet, addrType string) string {
	for _, addr := range d.Networks.V6 {
		if addr.Type == addrType {
			if ip := net.ParseIP(addr.IPAddress); ip != nil {
				return strings.ToLower(addr.IPAddress)
			}
		}
	}
	return ""
}

func FindIPv4AddrByType(d *godo.Droplet, addrType string) string {
	for _, addr := range d.Networks.V4 {
		if addr.Type == addrType {
			if ip := net.ParseIP(addr.IPAddress); ip != nil {
				return addr.IPAddress
			}
		}
	}
	return ""
}

func resourceDigitalOceanDropletUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	var warnings []diag.Diagnostic

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid droplet id: %v", err)
	}

	if d.HasChange("size") {
		newSize := d.Get("size")
		resizeDisk := d.Get("resize_disk").(bool)

		_, _, err = client.DropletActions.PowerOff(context.Background(), id)
		if err != nil && !strings.Contains(err.Error(), "Droplet is already powered off") {
			return diag.Errorf(
				"Error powering off droplet (%s): %s", d.Id(), err)
		}

		// Wait for power off
		_, err = waitForDropletAttribute(ctx, d, "off", []string{"active"}, "status", schema.TimeoutUpdate, meta)
		if err != nil {
			return diag.Errorf(
				"Error waiting for droplet (%s) to become powered off: %s", d.Id(), err)
		}

		// Resize the droplet
		var action *godo.Action
		action, _, err = client.DropletActions.Resize(context.Background(), id, newSize.(string), resizeDisk)
		if err != nil {
			newErr := powerOnAndWait(ctx, d, meta)
			if newErr != nil {
				return diag.Errorf(
					"Error powering on droplet (%s) after failed resize: %s", d.Id(), err)
			}
			return diag.Errorf(
				"Error resizing droplet (%s): %s", d.Id(), err)
		}

		// Wait for the resize action to complete.
		if err = util.WaitForAction(client, action); err != nil {
			newErr := powerOnAndWait(ctx, d, meta)
			if newErr != nil {
				return diag.Errorf(
					"Error powering on droplet (%s) after waiting for resize to finish: %s", d.Id(), err)
			}
			return diag.Errorf(
				"Error waiting for resize droplet (%s) to finish: %s", d.Id(), err)
		}

		_, _, err = client.DropletActions.PowerOn(context.Background(), id)

		if err != nil {
			return diag.Errorf(
				"Error powering on droplet (%s) after resize: %s", d.Id(), err)
		}

		// Wait for power on
		_, err = waitForDropletAttribute(ctx, d, "active", []string{"off"}, "status", schema.TimeoutUpdate, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")

		// Rename the droplet
		_, _, err = client.DropletActions.Rename(context.Background(), id, newName.(string))

		if err != nil {
			return diag.Errorf(
				"Error renaming droplet (%s): %s", d.Id(), err)
		}

		// Wait for the name to change
		_, err = waitForDropletAttribute(
			ctx, d, newName.(string), []string{"", oldName.(string)}, "name", schema.TimeoutUpdate, meta)

		if err != nil {
			return diag.Errorf(
				"Error waiting for rename droplet (%s) to finish: %s", d.Id(), err)
		}
	}

	if d.HasChange("backups") {
		if d.Get("backups").(bool) {
			// Enable backups on droplet
			var action *godo.Action
			// Apply backup_policy if specified, otherwise use the default policy
			policy, ok := d.GetOk("backup_policy")
			if ok {
				backupPolicy, err := expandBackupPolicy(policy)
				if err != nil {
					return diag.FromErr(err)
				}
				action, _, err = client.DropletActions.EnableBackupsWithPolicy(context.Background(), id, backupPolicy)
				if err != nil {
					return diag.Errorf(
						"Error enabling backups on droplet (%s): %s", d.Id(), err)
				}
			} else {
				action, _, err = client.DropletActions.EnableBackups(context.Background(), id)
				if err != nil {
					return diag.Errorf(
						"Error enabling backups on droplet (%s): %s", d.Id(), err)
				}
			}
			if err := util.WaitForAction(client, action); err != nil {
				return diag.Errorf("Error waiting for backups to be enabled for droplet (%s): %s", d.Id(), err)
			}
		} else {
			// Disable backups on droplet
			// Check there is no backup_policy specified
			_, ok := d.GetOk("backup_policy")
			if ok {
				return diag.FromErr(errDropletBackupPolicy)
			}
			action, _, err := client.DropletActions.DisableBackups(context.Background(), id)
			if err != nil {
				return diag.Errorf(
					"Error disabling backups on droplet (%s): %s", d.Id(), err)
			}

			if err := util.WaitForAction(client, action); err != nil {
				return diag.Errorf("Error waiting for backups to be disabled for droplet (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("backup_policy") {
		_, ok := d.GetOk("backup_policy")
		if ok {
			if !d.Get("backups").(bool) {
				return diag.FromErr(errDropletBackupPolicy)
			}

			_, new := d.GetChange("backup_policy")
			newPolicy, err := expandBackupPolicy(new)
			if err != nil {
				return diag.FromErr(err)
			}

			action, _, err := client.DropletActions.ChangeBackupPolicy(context.Background(), id, newPolicy)
			if err != nil {
				return diag.Errorf(
					"error changing backup policy on droplet (%s): %s", d.Id(), err)
			}

			if err := util.WaitForAction(client, action); err != nil {
				return diag.Errorf("error waiting for backup policy to be changed for droplet (%s): %s", d.Id(), err)
			}
		}
	}

	// As there is no way to disable private networking,
	// we only check if it needs to be enabled
	if d.HasChange("private_networking") && d.Get("private_networking").(bool) {
		_, _, err = client.DropletActions.EnablePrivateNetworking(context.Background(), id)

		if err != nil {
			return diag.Errorf(
				"Error enabling private networking for droplet (%s): %s", d.Id(), err)
		}

		// Wait for the private_networking to turn on
		_, err = waitForDropletAttribute(
			ctx, d, "true", []string{"", "false"}, "private_networking", schema.TimeoutUpdate, meta)

		if err != nil {
			return diag.Errorf(
				"Error waiting for private networking to be enabled on for droplet (%s): %s", d.Id(), err)
		}
	}

	// As there is no way to disable IPv6, we only check if it needs to be enabled
	if d.HasChange("ipv6") && d.Get("ipv6").(bool) {
		_, _, err = client.DropletActions.EnableIPv6(context.Background(), id)
		if err != nil {
			return diag.Errorf(
				"Error turning on ipv6 for droplet (%s): %s", d.Id(), err)
		}

		// Wait for ipv6 to turn on
		_, err = waitForDropletAttribute(
			ctx, d, "true", []string{"", "false"}, "ipv6", schema.TimeoutUpdate, meta)

		if err != nil {
			return diag.Errorf(
				"Error waiting for ipv6 to be turned on for droplet (%s): %s", d.Id(), err)
		}

		warnings = append(warnings, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Enabling IPv6 requires additional OS-level configuration",
			Detail:   "When enabling IPv6 on an existing Droplet, additional OS-level configuration is required. For more info, see: \nhttps://docs.digitalocean.com/products/networking/ipv6/how-to/enable/#on-existing-droplets",
		})
	}

	if d.HasChange("tags") {
		err = tag.SetTags(client, d, godo.DropletResourceType)
		if err != nil {
			return diag.Errorf("Error updating tags: %s", err)
		}
	}

	if d.HasChange("volume_ids") {
		oldIDs, newIDs := d.GetChange("volume_ids")
		newSet := func(ids []interface{}) map[string]struct{} {
			out := make(map[string]struct{}, len(ids))
			for _, id := range ids {
				out[id.(string)] = struct{}{}
			}
			return out
		}
		// leftDiff returns all elements in Left that are not in Right
		leftDiff := func(left, right map[string]struct{}) map[string]struct{} {
			out := make(map[string]struct{})
			for l := range left {
				if _, ok := right[l]; !ok {
					out[l] = struct{}{}
				}
			}
			return out
		}
		oldIDSet := newSet(oldIDs.(*schema.Set).List())
		newIDSet := newSet(newIDs.(*schema.Set).List())
		for volumeID := range leftDiff(newIDSet, oldIDSet) {
			action, _, err := client.StorageActions.Attach(context.Background(), volumeID, id)
			if err != nil {
				return diag.Errorf("Error attaching volume %q to droplet (%s): %s", volumeID, d.Id(), err)
			}
			// can't fire >1 action at a time, so waiting for each is OK
			if err := util.WaitForAction(client, action); err != nil {
				return diag.Errorf("Error waiting for volume %q to attach to droplet (%s): %s", volumeID, d.Id(), err)
			}
		}
		for volumeID := range leftDiff(oldIDSet, newIDSet) {
			err := detachVolumeIDOnDroplet(d, volumeID, meta)
			if err != nil {
				return diag.Errorf("Error detaching volume %q on droplet %s: %s", volumeID, d.Id(), err)

			}
		}
	}

	readErr := resourceDigitalOceanDropletRead(ctx, d, meta)
	if readErr != nil {
		warnings = append(warnings, readErr...)
	}

	return warnings
}

func resourceDigitalOceanDropletDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid droplet id: %v", err)
	}

	_, err = waitForDropletAttribute(
		ctx, d, "false", []string{"", "true"}, "locked", schema.TimeoutDelete, meta)

	if err != nil {
		return diag.Errorf(
			"Error waiting for droplet to be unlocked for destroy (%s): %s", d.Id(), err)
	}

	shutdown := d.Get("graceful_shutdown").(bool)
	if shutdown {
		log.Printf("[INFO] Shutting down droplet: %s", d.Id())

		// Shutdown the droplet
		// DO API doesn't return an error if we try to shutdown an already shutdown droplet
		_, _, err = client.DropletActions.Shutdown(context.Background(), id)
		if err != nil {
			return diag.Errorf(
				"Error shutting down the the droplet (%s): %s", d.Id(), err)
		}

		// Wait for shutdown
		_, err = waitForDropletAttribute(ctx, d, "off", []string{"active"}, "status", schema.TimeoutDelete, meta)
		if err != nil {
			return diag.Errorf("Error waiting for droplet (%s) to become off: %s", d.Id(), err)
		}
	}

	log.Printf("[INFO] Trying to Detach Storage Volumes (if any) from droplet: %s", d.Id())
	err = detachVolumesFromDroplet(d, meta)
	if err != nil {
		return diag.Errorf(
			"Error detaching the volumes from the droplet (%s): %s", d.Id(), err)
	}

	log.Printf("[INFO] Deleting droplet: %s", d.Id())

	// Destroy the droplet
	resp, err := client.Droplets.Delete(context.Background(), id)

	// Handle already destroyed droplets
	if err != nil && resp.StatusCode == 404 {
		return nil
	}

	_, err = waitForDropletDestroy(ctx, d, meta)
	if err != nil && strings.Contains(err.Error(), "404") {
		return nil
	} else if err != nil {
		return diag.Errorf("Error deleting droplet: %s", err)
	}

	return nil
}

func waitForDropletDestroy(ctx context.Context, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for droplet (%s) to be destroyed", d.Id())

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active", "off"},
		Target:     []string{"archive"},
		Refresh:    dropletStateRefreshFunc(ctx, d, "status", meta),
		Timeout:    60 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func waitForDropletAttribute(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, timeoutKey string, meta interface{}) (interface{}, error) {
	// Wait for the droplet so we can get the networking attributes
	// that show up after a while
	log.Printf(
		"[INFO] Waiting for droplet (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    dropletStateRefreshFunc(ctx, d, attribute, meta),
		Timeout:    d.Timeout(timeoutKey),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		// This is a hack around DO API strangeness.
		// https://github.com/hashicorp/terraform/issues/481
		//
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

// TODO This function still needs a little more refactoring to make it
// cleaner and more efficient
func dropletStateRefreshFunc(
	ctx context.Context, d *schema.ResourceData, attribute string, meta interface{}) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).GodoClient()
	return func() (interface{}, string, error) {
		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, "", err
		}

		// Retrieve the droplet properties
		droplet, _, err := client.Droplets.Get(context.Background(), id)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving droplet: %s", err)
		}

		err = setDropletAttributes(d, droplet)
		if err != nil {
			return nil, "", err
		}

		// If the droplet is locked, continue waiting. We can
		// only perform actions on unlocked droplets, so it's
		// pointless to look at that status
		if d.Get("locked").(bool) {
			log.Println("[DEBUG] Droplet is locked, skipping status check and retrying")
			return nil, "", nil
		}

		// See if we can access our attribute
		if attr, ok := d.GetOkExists(attribute); ok {
			switch attr := attr.(type) {
			case bool:
				return &droplet, strconv.FormatBool(attr), nil
			default:
				return &droplet, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}

// Powers on the droplet and waits for it to be active
func powerOnAndWait(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid droplet id: %v", err)
	}

	client := meta.(*config.CombinedConfig).GodoClient()
	_, _, err = client.DropletActions.PowerOn(context.Background(), id)
	if err != nil {
		return err
	}
	// this method is only used for droplet updates so use that as the timeout parameter
	// Wait for power on
	_, err = waitForDropletAttribute(ctx, d, "active", []string{"off"}, "status", schema.TimeoutUpdate, meta)
	if err != nil {
		return err
	}

	return nil
}

// Detach volumes from droplet
func detachVolumesFromDroplet(d *schema.ResourceData, meta interface{}) error {
	var errors []error
	if attr, ok := d.GetOk("volume_ids"); ok {
		errors = make([]error, 0, attr.(*schema.Set).Len())
		for _, volumeID := range attr.(*schema.Set).List() {
			err := detachVolumeIDOnDroplet(d, volumeID.(string), meta)
			if err != nil {
				return err
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Error detaching one or more volumes: %v", errors)
	}

	return nil
}

func detachVolumeIDOnDroplet(d *schema.ResourceData, volumeID string, meta interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid droplet id: %v", err)
	}
	client := meta.(*config.CombinedConfig).GodoClient()
	action, _, err := client.StorageActions.DetachByDropletID(context.Background(), volumeID, id)
	if err != nil {
		return fmt.Errorf("Error detaching volume %q from droplet (%s): %s", volumeID, d.Id(), err)
	}
	// can't fire >1 action at a time, so waiting for each is OK
	if err := util.WaitForAction(client, action); err != nil {
		return fmt.Errorf("Error waiting for volume %q to detach from droplet (%s): %s", volumeID, d.Id(), err)
	}

	return nil
}

func containsDigitalOceanDropletFeature(features []string, name string) bool {
	for _, v := range features {
		if v == name {
			return true
		}
	}
	return false
}

func expandSshKeys(sshKeys []interface{}) ([]godo.DropletCreateSSHKey, error) {
	expandedSshKeys := make([]godo.DropletCreateSSHKey, len(sshKeys))
	for i, s := range sshKeys {
		sshKey := s.(string)

		var expandedSshKey godo.DropletCreateSSHKey
		if id, err := strconv.Atoi(sshKey); err == nil {
			expandedSshKey.ID = id
		} else {
			expandedSshKey.Fingerprint = sshKey
		}

		expandedSshKeys[i] = expandedSshKey
	}

	return expandedSshKeys, nil
}

func flattenDigitalOceanDropletVolumeIds(volumeids []string) *schema.Set {
	flattenedVolumes := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range volumeids {
		flattenedVolumes.Add(v)
	}

	return flattenedVolumes
}

func expandBackupPolicy(v interface{}) (*godo.DropletBackupPolicyRequest, error) {
	var policy godo.DropletBackupPolicyRequest
	policyList := v.([]interface{})

	for _, rawPolicy := range policyList {
		policyMap, ok := rawPolicy.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("droplet backup policy type assertion failed: expected map[string]interface{}, got %T", rawPolicy)
		}

		planVal, exists := policyMap["plan"]
		if !exists {
			return nil, errors.New("backup_policy plan key does not exist")
		}
		plan, ok := planVal.(string)
		if !ok {
			return nil, errors.New("backup_policy plan is not a string")
		}
		policy.Plan = plan

		weekdayVal, exists := policyMap["weekday"]
		if !exists {
			return nil, errors.New("backup_policy weekday key does not exist")
		}
		weekday, ok := weekdayVal.(string)
		if !ok {
			return nil, errors.New("backup_policy weekday is not a string")
		}
		policy.Weekday = weekday

		hourVal, exists := policyMap["hour"]
		if !exists {
			return nil, errors.New("backup_policy hour key does not exist")
		}
		hour, ok := hourVal.(int)
		if !ok {
			return nil, errors.New("backup_policy hour is not an int")
		}
		policy.Hour = &hour
	}

	return &policy, nil
}
