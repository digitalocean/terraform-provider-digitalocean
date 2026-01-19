package triSHUIL

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/droplet"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VM API path (appended to base URL configured in config.go)
const vmAPIPath = "api/vm"

// vmRoot is the response structure for VM API
type vmRoot struct {
	Droplet *godo.Droplet `json:"droplet"`
}

func ResourceDigitalOceanVM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanVMCreate,
		ReadContext:   resourceDigitalOceanVMRead,
		UpdateContext: droplet.ResourceDigitalOceanDroplet().UpdateContext, // TODO: Implement for future
		DeleteContext: droplet.ResourceDigitalOceanDroplet().DeleteContext,

		SchemaVersion: 1,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Second),
		},

		Schema:        droplet.ResourceDigitalOceanDroplet().Schema,
		CustomizeDiff: droplet.ResourceDigitalOceanDroplet().CustomizeDiff,
	}
}

func resourceDigitalOceanVMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).VMClient()

	image := d.Get("image").(string)

	// Build up our creation options
	opts := &godo.DropletCreateRequest{
		Image:  godo.DropletCreateImage{},
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(string),
		Tags:   tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
	}

	imageId, err := strconv.Atoi(image)
	if err == nil {
		opts.Image.ID = imageId
	} else {
		opts.Image.Slug = image
	}

	if attr, ok := d.GetOk("backups"); ok {
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
		expandedSshKeys, err := expandVMSshKeys(v.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		opts.SSHKeys = expandedSshKeys
	}

	log.Printf("[DEBUG] VM create configuration: %#v", opts)

	// Create request using godo client (which is configured with VM API base URL)
	req, err := client.NewRequest(ctx, http.MethodPost, vmAPIPath, opts)
	if err != nil {
		return diag.Errorf("Error creating VM request: %s", err)
	}

	// Execute request
	root := new(vmRoot)
	_, err = client.Do(ctx, req, root)
	if err != nil {
		return diag.Errorf("Error creating VM: %s", err)
	}

	if root.Droplet == nil {
		return diag.Errorf("Error creating VM: empty response from API")
	}

	// Assign the VM id
	d.SetId(strconv.Itoa(root.Droplet.ID))
	log.Printf("[INFO] VM ID: %s", d.Id())

	// Wait for VM status to become "active"
	_, err = waitForVMAttribute(ctx, d, "active", []string{"new"}, "status", schema.TimeoutCreate, meta)
	if err != nil {
		return diag.Errorf("Error waiting for VM (%s) to become ready: %s", d.Id(), err)
	}

	return nil
}

func resourceDigitalOceanVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).VMClient()

	// Parse VM ID from resource state
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid VM id: %v", err)
	}

	// Retrieve the VM properties using VM API
	// API endpoint: http://167.71.171.33:8080/api/vm/{id}
	path := fmt.Sprintf("%s/%d", vmAPIPath, id)
	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return diag.Errorf("Error creating VM read request: %s", err)
	}

	// Execute request
	root := new(vmRoot)
	resp, err := client.Do(ctx, req, root)
	if err != nil {
		// Check if the VM no longer exists (404 = not found)
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean VM (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving VM: %s", err)
	}

	// Validate response
	if root.Droplet == nil {
		return diag.Errorf("Error retrieving VM: empty response from API")
	}

	// Update all attributes in Terraform state
	err = setVMAttributes(d, root.Droplet)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandVMSshKeys(sshKeys []interface{}) ([]godo.DropletCreateSSHKey, error) {
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

func setVMAttributes(d *schema.ResourceData, vm *godo.Droplet) error {
	d.Set("name", vm.Name)
	d.Set("urn", vm.URN())
	d.Set("region", vm.Region.Slug)
	d.Set("size", vm.Size.Slug)
	d.Set("price_hourly", vm.Size.PriceHourly)
	d.Set("price_monthly", vm.Size.PriceMonthly)
	d.Set("disk", vm.Disk)
	d.Set("vcpus", vm.Vcpus)
	d.Set("memory", vm.Memory)
	d.Set("status", vm.Status)
	d.Set("locked", vm.Locked)
	d.Set("created_at", vm.Created)
	d.Set("vpc_uuid", vm.VPCUUID)

	d.Set("ipv4_address", findVMIPv4AddrByType(vm, "public"))
	d.Set("ipv4_address_private", findVMIPv4AddrByType(vm, "private"))
	d.Set("ipv6_address", strings.ToLower(findVMIPv6AddrByType(vm, "public")))

	if features := vm.Features; features != nil {
		d.Set("backups", slices.Contains(features, "backups"))
		d.Set("ipv6", slices.Contains(features, "ipv6"))
		d.Set("private_networking", slices.Contains(features, "private_networking"))
		d.Set("monitoring", slices.Contains(features, "monitoring"))
	}

	if err := d.Set("volume_ids", flattenVMVolumeIds(vm.VolumeIDs)); err != nil {
		return fmt.Errorf("Error setting `volume_ids`: %+v", err)
	}

	if err := d.Set("tags", tag.FlattenTags(vm.Tags)); err != nil {
		return fmt.Errorf("Error setting `tags`: %+v", err)
	}

	// Initialize the connection info
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": findVMIPv4AddrByType(vm, "public"),
	})

	return nil
}

func findVMIPv6AddrByType(d *godo.Droplet, addrType string) string {
	for _, addr := range d.Networks.V6 {
		if addr.Type == addrType {
			if ip := net.ParseIP(addr.IPAddress); ip != nil {
				return strings.ToLower(addr.IPAddress)
			}
		}
	}
	return ""
}

func findVMIPv4AddrByType(d *godo.Droplet, addrType string) string {
	for _, addr := range d.Networks.V4 {
		if addr.Type == addrType {
			if ip := net.ParseIP(addr.IPAddress); ip != nil {
				return addr.IPAddress
			}
		}
	}
	return ""
}

func flattenVMVolumeIds(volumeids []string) *schema.Set {
	flattenedVolumes := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range volumeids {
		flattenedVolumes.Add(v)
	}

	return flattenedVolumes
}

func waitForVMAttribute(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, timeoutKey string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for VM (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &retry.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        vmStateRefreshFunc(ctx, d, attribute, meta),
		Timeout:        d.Timeout(timeoutKey),
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func vmStateRefreshFunc(
	ctx context.Context, d *schema.ResourceData, attribute string, meta interface{}) retry.StateRefreshFunc {
	client := meta.(*config.CombinedConfig).VMClient()

	return func() (interface{}, string, error) {
		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, "", err
		}

		// Retrieve the VM properties using VM client
		path := fmt.Sprintf("%s/%d", vmAPIPath, id)
		req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, "", fmt.Errorf("Error creating VM get request: %s", err)
		}

		root := new(vmRoot)
		_, err = client.Do(ctx, req, root)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving VM: %s", err)
		}

		if root.Droplet == nil {
			return nil, "", fmt.Errorf("Error retrieving VM: empty response")
		}

		vm := root.Droplet

		err = setVMAttributes(d, vm)
		if err != nil {
			return nil, "", err
		}

		// If the VM is locked, continue waiting
		if d.Get("locked").(bool) {
			log.Println("[DEBUG] VM is locked, skipping status check and retrying")
			return nil, "", nil
		}

		// See if we can access our attribute
		if attr, ok := d.GetOkExists(attribute); ok {
			switch attr := attr.(type) {
			case bool:
				return vm, strconv.FormatBool(attr), nil
			default:
				return vm, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}
