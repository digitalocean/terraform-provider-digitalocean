package microvm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceDigitalOceanMicroVMIngress returns the
// digitalocean_microvm_ingress resource schema. It models a single
// system-minted ingress URL on an existing MicroVM.
//
// URLs are minted by the platform (config cannot declare a hostname), so the
// resource exposes only the guest `port` as an input and surfaces the
// assigned `hostname`, `default`, and `status` as computed attributes.
// Creating the resource calls the add-URL API and deleting it calls the
// remove-URL API; both update the parent MicroVM in place and never force it
// to be replaced. Every input is ForceNew, so changing the port replaces the
// URL (remove + add) rather than the MicroVM.
func ResourceDigitalOceanMicroVMIngress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanMicroVMIngressCreate,
		ReadContext:   resourceDigitalOceanMicroVMIngressRead,
		DeleteContext: resourceDigitalOceanMicroVMIngressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanMicroVMIngressImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"microvm_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "ID of the MicroVM to add the ingress URL to",
				ValidateFunc: validation.NoZeroValues,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "Guest port the URL routes to. Must be one of the MicroVM's ports. Defaults to the MicroVM's http_port when omitted.",
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assigned hostname (no scheme). System-minted.",
			},
			"url_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stable identifier of the URL (same as the resource ID)",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is the MicroVM's system default URL",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL lifecycle status (PENDING or ACTIVE)",
			},
		},
	}
}

func resourceDigitalOceanMicroVMIngressCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	microVMID := d.Get("microvm_id").(string)
	req := &godo.MicroVMURLCreateRequest{}
	if v, ok := d.GetOk("port"); ok {
		req.Port = v.(int)
	}

	u, _, err := client.MicroVMs.AddURL(ctx, microVMID, req)
	if err != nil {
		return diag.Errorf("Error adding ingress URL to MicroVM (%s): %s", microVMID, err)
	}
	if u == nil || u.ID == "" {
		return diag.Errorf("Error adding ingress URL to MicroVM (%s): empty URL in response", microVMID)
	}

	d.SetId(u.ID)
	log.Printf("[INFO] MicroVM (%s) ingress URL added, ID: %s", microVMID, d.Id())

	// A freshly added URL may still be PENDING; the platform publishes the
	// route and flips it to ACTIVE without further input. Wait so apply only
	// returns once the URL is usable.
	stateConf := &retry.StateChangeConf{
		Pending:    []string{string(godo.MicroVMURLStatusPending)},
		Target:     []string{string(godo.MicroVMURLStatusActive)},
		Refresh:    microVMIngressStatusRefreshFunc(ctx, client, microVMID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for MicroVM (%s) ingress URL (%s) to become active: %s", microVMID, d.Id(), err)
	}

	return resourceDigitalOceanMicroVMIngressRead(ctx, d, meta)
}

func resourceDigitalOceanMicroVMIngressRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	microVMID := d.Get("microvm_id").(string)
	m, resp, err := client.MicroVMs.Get(ctx, microVMID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[WARN] MicroVM (%s) not found - removing ingress URL (%s) from state", microVMID, d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving MicroVM (%s) for ingress URL (%s): %s", microVMID, d.Id(), err)
	}

	u := findMicroVMURL(m, d.Id())
	if u == nil {
		// The URL was removed outside Terraform. Drop it from state so the
		// next plan proposes recreating it (drift detection).
		log.Printf("[WARN] ingress URL (%s) not found on MicroVM (%s) - removing from state", d.Id(), microVMID)
		d.SetId("")
		return nil
	}

	d.Set("microvm_id", microVMID)
	d.Set("url_id", u.ID)
	d.Set("hostname", u.Hostname)
	d.Set("port", u.Port)
	d.Set("default", u.Default)
	d.Set("status", string(u.Status))
	return nil
}

func resourceDigitalOceanMicroVMIngressDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	microVMID := d.Get("microvm_id").(string)
	resp, err := client.MicroVMs.RemoveURL(ctx, microVMID, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error removing ingress URL (%s) from MicroVM (%s): %s", d.Id(), microVMID, err)
	}

	log.Printf("[INFO] MicroVM (%s) ingress URL removed, ID: %s", microVMID, d.Id())
	d.SetId("")
	return nil
}

// resourceDigitalOceanMicroVMIngressImport accepts a "microvm_id,url_id"
// import ID because a URL is only addressable through its parent MicroVM.
func resourceDigitalOceanMicroVMIngressImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ",", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("import ID must be in the form 'microvm_id,url_id', got %q", d.Id())
	}
	d.Set("microvm_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}

// microVMIngressStatusRefreshFunc reads the parent MicroVM and reports the
// status of the one URL identified by urlID. A URL that has not yet appeared
// in the MicroVM's URL list is reported as PENDING so the waiter keeps
// polling rather than failing.
func microVMIngressStatusRefreshFunc(ctx context.Context, client *godo.Client, microVMID, urlID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		m, _, err := client.MicroVMs.Get(ctx, microVMID)
		if err != nil {
			return nil, "", err
		}
		u := findMicroVMURL(m, urlID)
		if u == nil {
			return m, string(godo.MicroVMURLStatusPending), nil
		}
		return u, string(u.Status), nil
	}
}

// findMicroVMURL returns the URL with the given ID from a MicroVM, or nil.
func findMicroVMURL(m *godo.MicroVM, urlID string) *godo.MicroVMURL {
	if m == nil {
		return nil
	}
	for i := range m.URLs {
		if m.URLs[i].ID == urlID {
			return &m.URLs[i]
		}
	}
	return nil
}
