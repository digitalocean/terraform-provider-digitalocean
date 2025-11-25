package byoipprefix

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBYOIPPrefix() *schema.Resource {
	return &schema.Resource{
		Create: resourceBYOIPPrefixCreate,
		Read:   resourceBYOIPPrefixRead,
		Update: resourceBYOIPPrefixUpdate,
		Delete: resourceBYOIPPrefixDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBYOIPPrefixImport,
		},
		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"signature": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"advertised": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failure_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getBYOIPService(meta interface{}) godo.BYOIPPrefixesService {
	client := meta.(*config.CombinedConfig).GodoClient()
	return client.BYOIPPrefixes
}

func resourceBYOIPPrefixCreate(d *schema.ResourceData, meta interface{}) error {
	service := getBYOIPService(meta)
	prefix := d.Get("prefix").(string)
	region := d.Get("region").(string)
	signature, sigOk := d.GetOk("signature")
	if !sigOk || signature.(string) == "" {
		return fmt.Errorf("signature must be provided for BYOIP prefix creation")
	}

	createReq := &godo.BYOIPPrefixCreateReq{
		Prefix:    prefix,
		Signature: signature.(string),
		Region:    region,
	}

	resp, _, err := service.Create(context.Background(), createReq)
	if err != nil {
		return err
	}

	d.SetId(resp.UUID)
	return resourceBYOIPPrefixRead(d, meta)
}

func resourceBYOIPPrefixRead(d *schema.ResourceData, meta interface{}) error {
	service := getBYOIPService(meta)
	uuid := d.Id()

	prefix, _, err := service.Get(context.Background(), uuid)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("prefix", prefix.Prefix)
	d.Set("region", prefix.Region)
	d.Set("status", prefix.Status)
	d.Set("uuid", prefix.UUID)
	d.Set("advertised", prefix.Advertised)
	d.Set("failure_reason", prefix.FailureReason)
	return nil
}

func resourceBYOIPPrefixUpdate(d *schema.ResourceData, meta interface{}) error {
	service := getBYOIPService(meta)
	uuid := d.Id()

	if d.HasChange("advertised") {
		advertised := d.Get("advertised").(bool)
		updateReq := &godo.BYOIPPrefixUpdateReq{
			Advertise: &advertised,
		}
		_, _, err := service.Update(context.Background(), uuid, updateReq)
		if err != nil {
			return err
		}
	}
	return resourceBYOIPPrefixRead(d, meta)
}

func resourceBYOIPPrefixDelete(d *schema.ResourceData, meta interface{}) error {
	service := getBYOIPService(meta)
	uuid := d.Id()
	_, err := service.Delete(context.Background(), uuid)
	if err != nil {
		return err
	}
	return nil
}

func resourceBYOIPPrefixImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceBYOIPPrefixRead(d, meta)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func isNotFound(err error) bool {
	errResp, ok := err.(*godo.ErrorResponse)
	if ok && errResp.Response.StatusCode == 404 {
		return true
	}

	return false
}
