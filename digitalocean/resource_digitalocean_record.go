package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDigitalOceanRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanRecordCreate,
		Read:   resourceDigitalOceanRecordRead,
		Update: resourceDigitalOceanRecordUpdate,
		Delete: resourceDigitalOceanRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"weight": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					domain := d.Get("domain").(string) + "."

					return (old == "@" && new == domain) || (old == new+domain)
				},
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"flags": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDigitalOceanRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	newRecord := godo.DomainRecordEditRequest{
		Type: d.Get("type").(string),
		Name: d.Get("name").(string),
		Data: d.Get("value").(string),
		Tag:  d.Get("tag").(string),
	}

	var err error
	if priority := d.Get("priority").(string); priority != "" {
		newRecord.Priority, err = strconv.Atoi(priority)
		if err != nil {
			return fmt.Errorf("Failed to parse priority as an integer: %v", err)
		}
	}
	if port := d.Get("port").(string); port != "" {
		newRecord.Port, err = strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("Failed to parse port as an integer: %v", err)
		}
	}
	if ttl := d.Get("ttl").(string); ttl != "" {
		newRecord.TTL, err = strconv.Atoi(ttl)
		if err != nil {
			return fmt.Errorf("Failed to parse ttl as an integer: %v", err)
		}
	}
	if weight := d.Get("weight").(string); weight != "" {
		newRecord.Weight, err = strconv.Atoi(weight)
		if err != nil {
			return fmt.Errorf("Failed to parse weight as an integer: %v", err)
		}
	}
	if flags := d.Get("flags").(string); flags != "" {
		newRecord.Flags, err = strconv.Atoi(flags)
		if err != nil {
			return fmt.Errorf("Failed to parse flags as an integer: %v", err)
		}
	}

	log.Printf("[DEBUG] record create configuration: %#v", newRecord)
	rec, _, err := client.Domains.CreateRecord(context.Background(), d.Get("domain").(string), &newRecord)
	if err != nil {
		return fmt.Errorf("Failed to create record: %s", err)
	}

	d.SetId(strconv.Itoa(rec.ID))
	log.Printf("[INFO] Record ID: %s", d.Id())

	return resourceDigitalOceanRecordRead(d, meta)
}

func resourceDigitalOceanRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)
	domain := d.Get("domain").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid record ID: %v", err)
	}

	rec, resp, err := client.Domains.Record(context.Background(), domain, id)
	if err != nil && resp != nil {
		// If the record is somehow already destroyed, mark as
		// successfully gone
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return err
	} else if err != nil {
		return err
	}

	if t := rec.Type; t == "CNAME" || t == "MX" || t == "NS" || t == "SRV" || t == "CAA" {
		if rec.Data != "@" {
			rec.Data += "."
		}
	}

	d.Set("name", rec.Name)
	d.Set("type", rec.Type)
	d.Set("value", rec.Data)
	d.Set("weight", strconv.Itoa(rec.Weight))
	d.Set("priority", strconv.Itoa(rec.Priority))
	d.Set("port", strconv.Itoa(rec.Port))
	d.Set("ttl", strconv.Itoa(rec.TTL))
	d.Set("flags", strconv.Itoa(rec.Flags))
	d.Set("tag", rec.Tag)

	en := constructFqdn(rec.Name, d.Get("domain").(string))
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

	err := resourceDigitalOceanRecordRead(d, meta)
	if err != nil {
		return nil, fmt.Errorf("unable to import record: %v", err)
	}

	results := make([]*schema.ResourceData, 0)
	results = append(results, d)

	return results, nil
}

func resourceDigitalOceanRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	domain := d.Get("domain").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid record ID: %v", err)
	}

	var editRecord godo.DomainRecordEditRequest
	if v, ok := d.GetOk("name"); ok {
		editRecord.Name = v.(string)
	}
	if v, ok := d.GetOk("value"); ok {
		editRecord.Data = v.(string)
	}
	if v, ok := d.GetOk("tag"); ok {
		editRecord.Tag = v.(string)
	}

	if d.HasChange("priority") {
		newPriority := d.Get("priority").(string)
		editRecord.Priority, err = strconv.Atoi(newPriority)
		if err != nil {
			return fmt.Errorf("Failed to parse priority as an integer: %v", err)
		}
	}
	if d.HasChange("port") {
		newPort := d.Get("port").(string)
		editRecord.Port, err = strconv.Atoi(newPort)
		if err != nil {
			return fmt.Errorf("Failed to parse port as an integer: %v", err)
		}
	}
	if d.HasChange("ttl") {
		newTTL := d.Get("ttl").(string)
		editRecord.TTL, err = strconv.Atoi(newTTL)
		if err != nil {
			return fmt.Errorf("Failed to parse ttl as an integer: %v", err)
		}
	}
	if d.HasChange("weight") {
		newWeight := d.Get("weight").(string)
		editRecord.Weight, err = strconv.Atoi(newWeight)
		if err != nil {
			return fmt.Errorf("Failed to parse weight as an integer: %v", err)
		}
	}
	if d.HasChange("flags") {
		newFlags := d.Get("flags").(string)
		editRecord.Flags, err = strconv.Atoi(newFlags)
		if err != nil {
			return fmt.Errorf("Failed to parse flags as an integer: %v", err)
		}
	}

	log.Printf("[DEBUG] record update configuration: %#v", editRecord)
	_, _, err = client.Domains.EditRecord(context.Background(), domain, id, &editRecord)
	if err != nil {
		return fmt.Errorf("Failed to update record: %s", err)
	}

	return resourceDigitalOceanRecordRead(d, meta)
}

func resourceDigitalOceanRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	domain := d.Get("domain").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid record ID: %v", err)
	}

	log.Printf("[INFO] Deleting record: %s, %d", domain, id)

	resp, delErr := client.Domains.DeleteRecord(context.Background(), domain, id)
	if delErr != nil {
		// If the record is somehow already destroyed, mark as
		// successfully gone
		if resp.StatusCode == 404 {
			return nil
		}

		return fmt.Errorf("Error deleting record: %s", delErr)
	}

	return nil
}

func constructFqdn(name, domain string) string {
	rn := strings.ToLower(strings.TrimSuffix(name, "."))
	domain = strings.TrimSuffix(domain, ".")
	if !strings.HasSuffix(rn, domain) {
		rn = strings.Join([]string{name, domain}, ".")
	}
	return rn
}
