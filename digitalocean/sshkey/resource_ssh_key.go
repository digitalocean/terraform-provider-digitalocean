package sshkey

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanSSHKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanSSHKeyCreate,
		ReadContext:   resourceDigitalOceanSSHKeyRead,
		UpdateContext: resourceDigitalOceanSSHKeyUpdate,
		DeleteContext: resourceDigitalOceanSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"public_key": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: resourceDigitalOceanSSHKeyPublicKeyDiffSuppress,
				ValidateFunc:     validation.NoZeroValues,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanSSHKeyPublicKeyDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func resourceDigitalOceanSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	// Build up our creation options
	opts := &godo.KeyCreateRequest{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("public_key").(string),
	}

	log.Printf("[DEBUG] SSH Key create configuration: %#v", opts)
	key, _, err := client.Keys.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating SSH Key: %s", err)
	}

	d.SetId(strconv.Itoa(key.ID))
	log.Printf("[INFO] SSH Key: %d", key.ID)

	return resourceDigitalOceanSSHKeyRead(ctx, d, meta)
}

func resourceDigitalOceanSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid SSH key id: %v", err)
	}

	key, resp, err := client.Keys.GetByID(context.Background(), id)
	if err != nil {
		// If the key is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving SSH key: %s", err)
	}

	d.Set("name", key.Name)
	d.Set("fingerprint", key.Fingerprint)
	d.Set("public_key", key.PublicKey)

	return nil
}

func resourceDigitalOceanSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid SSH key id: %v", err)
	}

	var newName string
	if v, ok := d.GetOk("name"); ok {
		newName = v.(string)
	}

	log.Printf("[DEBUG] SSH key update name: %#v", newName)
	opts := &godo.KeyUpdateRequest{
		Name: newName,
	}
	_, _, err = client.Keys.UpdateByID(context.Background(), id, opts)
	if err != nil {
		return diag.Errorf("Failed to update SSH key: %s", err)
	}

	return resourceDigitalOceanSSHKeyRead(ctx, d, meta)
}

func resourceDigitalOceanSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid SSH key id: %v", err)
	}

	log.Printf("[INFO] Deleting SSH key: %d", id)
	_, err = client.Keys.DeleteByID(context.Background(), id)
	if err != nil {
		return diag.Errorf("Error deleting SSH key: %s", err)
	}

	d.SetId("")
	return nil
}
