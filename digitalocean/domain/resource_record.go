package domain

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanRecordCreate,
		ReadContext:   resourceDigitalOceanRecordRead,
		UpdateContext: resourceDigitalOceanRecordUpdate,
		DeleteContext: resourceDigitalOceanRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A",
					"AAAA",
					"CAA",
					"CNAME",
					"MX",
					"NS",
					"TXT",
					"SRV",
					"SOA",
				}, false),
			},

			"domain": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					domain := d.Get("domain").(string) + "."

					return (old == "@" && new == domain) || (old+"."+domain == new)
				},
			},

			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},

			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},

			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					domain := d.Get("domain").(string) + "."

					return (old == "@" && new == domain) || (old == new+"."+domain)
				},
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"flags": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 255),
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"issue",
					"issuewild",
					"iodef",
				}, false),
			},
		},

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			recordType := diff.Get("type").(string)

			_, hasPriority := diff.GetOkExists("priority")
			if recordType == "MX" {
				if !hasPriority {
					return fmt.Errorf("`priority` is required for when type is `MX`")
				}
			}

			_, hasWeight := diff.GetOkExists("weight")
			if recordType == "SRV" {
				if !hasPriority {
					return fmt.Errorf("`priority` is required for when type is `SRV`")
				}
				if !hasWeight {
					return fmt.Errorf("`weight` is required for when type is `SRV`")
				}
			}

			_, hasFlags := diff.GetOkExists("flags")
			_, hasTag := diff.GetOk("tag")
			if recordType == "CAA" {
				if !hasFlags {
					return fmt.Errorf("`flags` is required for when type is `CAA`")
				}
				if !hasTag {
					return fmt.Errorf("`tag` is required for when type is `CAA`")
				}
			}

			return nil
		},
	}
}

func resourceDigitalOceanRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	newRecord, err := expandDigitalOceanRecordResource(d)
	if err != nil {
		return diag.Errorf("Error in constructing record request: %s", err)
	}

	newRecord.Type = d.Get("type").(string)

	_, hasPort := d.GetOkExists("port")
	if newRecord.Type == "SRV" && !hasPort {
		return diag.Errorf("`port` is required for when type is `SRV`")
	}

	log.Printf("[DEBUG] record create configuration: %#v", newRecord)
	rec, _, err := client.Domains.CreateRecord(context.Background(), d.Get("domain").(string), newRecord)
	if err != nil {
		return diag.Errorf("Failed to create record: %s", err)
	}

	d.SetId(strconv.Itoa(rec.ID))
	log.Printf("[INFO] Record ID: %s", d.Id())

	return resourceDigitalOceanRecordRead(ctx, d, meta)
}

func resourceDigitalOceanRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	domain := d.Get("domain").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid record ID: %v", err)
	}

	rec, resp, err := client.Domains.Record(context.Background(), domain, id)
	if err != nil {
		// If the record is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	if t := rec.Type; t == "CNAME" || t == "MX" || t == "NS" || t == "SRV" || t == "CAA" {
		if rec.Data != "@" && rec.Tag != "iodef" {
			rec.Data += "."
		}
	}

	d.Set("name", rec.Name)
	d.Set("type", rec.Type)
	d.Set("value", rec.Data)
	d.Set("port", rec.Port)
	d.Set("priority", rec.Priority)
	d.Set("ttl", rec.TTL)
	d.Set("weight", rec.Weight)
	d.Set("flags", rec.Flags)
	d.Set("tag", rec.Tag)

	en := ConstructFqdn(rec.Name, d.Get("domain").(string))
	log.Printf("[DEBUG] Constructed FQDN: %s", en)
	d.Set("fqdn", en)

	return nil
}

func resourceDigitalOceanRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		// Validate that this is an ID by making sure it can be converted into an int
		_, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("invalid record ID: %v", err)
		}

		d.SetId(s[1])
		d.Set("domain", s[0])
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDigitalOceanRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	domain := d.Get("domain").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid record ID: %v", err)
	}

	editRecord, err := expandDigitalOceanRecordResource(d)
	if err != nil {
		return diag.Errorf("Error in constructing record request: %s", err)
	}

	log.Printf("[DEBUG] record update configuration: %#v", editRecord)
	_, _, err = client.Domains.EditRecord(context.Background(), domain, id, editRecord)
	if err != nil {
		return diag.Errorf("Failed to update record: %s", err)
	}

	return resourceDigitalOceanRecordRead(ctx, d, meta)
}

func resourceDigitalOceanRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	domain := d.Get("domain").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid record ID: %v", err)
	}

	log.Printf("[INFO] Deleting record: %s, %d", domain, id)

	resp, delErr := client.Domains.DeleteRecord(context.Background(), domain, id)
	if delErr != nil {
		// If the record is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			return nil
		}

		return diag.Errorf("Error deleting record: %s", delErr)
	}

	return nil
}

func expandDigitalOceanRecordResource(d *schema.ResourceData) (*godo.DomainRecordEditRequest, error) {
	record := &godo.DomainRecordEditRequest{
		Name: d.Get("name").(string),
		Data: d.Get("value").(string),
	}

	if v, ok := d.GetOk("port"); ok {
		record.Port = v.(int)
	}
	if v, ok := d.GetOk("priority"); ok {
		record.Priority = v.(int)
	}
	if v, ok := d.GetOk("ttl"); ok {
		record.TTL = v.(int)
	}
	if v, ok := d.GetOk("weight"); ok {
		record.Weight = v.(int)
	}
	if v, ok := d.GetOk("flags"); ok {
		record.Flags = v.(int)
	}
	if v, ok := d.GetOk("tag"); ok {
		record.Tag = v.(string)
	}

	return record, nil
}

func ConstructFqdn(name, domain string) string {
	if name == "@" {
		return domain
	}

	rn := strings.ToLower(name)
	domainSuffix := domain + "."
	if strings.HasSuffix(rn, domainSuffix) {
		rn = strings.TrimSuffix(rn, ".")
	} else {
		rn = strings.Join([]string{name, domain}, ".")
	}
	return rn
}
