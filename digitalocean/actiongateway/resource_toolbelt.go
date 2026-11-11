package actiongateway

import (
	"context"
	"log"
	"sort"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceDigitalOceanActionGatewayToolbelt returns the
// digitalocean_action_gateway_toolbelt resource. Changing the tools creates a
// new version of the toolbelt; the display name and description cannot be
// changed after creation.
func ResourceDigitalOceanActionGatewayToolbelt() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanActionGatewayToolbeltCreate,
		ReadContext:   resourceDigitalOceanActionGatewayToolbeltRead,
		UpdateContext: resourceDigitalOceanActionGatewayToolbeltUpdate,
		DeleteContext: resourceDigitalOceanActionGatewayToolbeltDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf("version", toolsChanged),
			customdiff.ComputedIf("reference", toolsChanged),
			customdiff.ComputedIf("pinned_tools", toolsChanged),
			customdiff.ComputedIf("tool_count", toolsChanged),
			customdiff.ComputedIf("updated_at", toolsChanged),
		),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the toolbelt, unique among your team's active toolbelts",
				ValidateFunc: validation.All(
					validation.StringMatch(toolbeltNameRegex, "must start with a lowercase letter and contain at most 64 characters from [a-z0-9_-]"),
					validation.StringNotInSlice([]string{"search"}, false),
				),
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Human-readable label for the toolbelt",
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Description of the toolbelt",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"tools": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Catalog tool slugs in the toolbelt, each optionally pinned as <tool_slug>@<version>. Unpinned tools are pinned to their current released version.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},
			"pinned_tools": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Members of the current version, each pinned as <tool_slug>@<version>",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current version of the toolbelt",
			},
			"reference": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version-pinned reference to the toolbelt, such as my-toolbelt@2",
			},
			"reference_latest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to the latest version of the toolbelt",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the toolbelt",
			},
			"tool_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tools in the toolbelt",
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func toolsChanged(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
	return d.Id() != "" && d.HasChange("tools")
}

func resourceDigitalOceanActionGatewayToolbeltCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := &godo.ActionGatewayToolbeltCreateRequest{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
		Description: d.Get("description").(string),
		Tools:       expandStringSet(d.Get("tools")),
	}

	created, _, err := client.ActionGateway.Toolbelts.Create(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating Action Gateway toolbelt: %s", err)
	}
	if created.Toolbelt == nil || created.Toolbelt.Name == "" {
		return diag.Errorf("Error creating Action Gateway toolbelt: response contained no toolbelt")
	}

	d.SetId(created.Toolbelt.Name)
	log.Printf("[INFO] Action Gateway toolbelt created: %s", d.Id())

	return resourceDigitalOceanActionGatewayToolbeltRead(ctx, d, meta)
}

func resourceDigitalOceanActionGatewayToolbeltRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	found, resp, err := client.ActionGateway.Toolbelts.Get(ctx, d.Id(), &godo.ActionGatewayGetToolbeltOptions{PageSize: 1})
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Action Gateway toolbelt (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Action Gateway toolbelt: %s", err)
	}
	if found.Toolbelt == nil {
		d.SetId("")
		return nil
	}

	values := flattenToolbelt(found.Toolbelt)
	values["pinned_tools"] = values["tools"]
	values["tools"] = reconcileToolbeltTools(expandStringSet(d.Get("tools")), found.Toolbelt.Tools)

	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	return nil
}

func resourceDigitalOceanActionGatewayToolbeltUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	name := d.Id()

	if d.HasChange("tools") {
		pinned, _ := d.GetChange("pinned_tools")
		remove, add := toolbeltMembershipChanges(expandStringSet(pinned), expandStringSet(d.Get("tools")))

		if len(remove) > 0 {
			if _, _, err := client.ActionGateway.Toolbelts.RemoveTools(ctx, name, &godo.ActionGatewayToolbeltToolsRequest{Tools: remove}); err != nil {
				return diag.Errorf("Error removing tools from Action Gateway toolbelt (%s): %s", name, err)
			}
		}
		if len(add) > 0 {
			if _, _, err := client.ActionGateway.Toolbelts.AddTools(ctx, name, &godo.ActionGatewayToolbeltToolsRequest{Tools: add}); err != nil {
				return diag.Errorf("Error adding tools to Action Gateway toolbelt (%s): %s", name, err)
			}
		}
	}

	return resourceDigitalOceanActionGatewayToolbeltRead(ctx, d, meta)
}

func resourceDigitalOceanActionGatewayToolbeltDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, resp, err := client.ActionGateway.Toolbelts.Delete(ctx, d.Id())
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting Action Gateway toolbelt: %s", err)
	}

	log.Printf("[INFO] Action Gateway toolbelt deleted: %s", d.Id())
	d.SetId("")
	return nil
}

func flattenToolbelt(tb *godo.ActionGatewayToolbelt) map[string]interface{} {
	return map[string]interface{}{
		"name":             tb.Name,
		"display_name":     tb.DisplayName,
		"description":      tb.Description,
		"tools":            tb.Tools,
		"version":          tb.Version,
		"reference":        tb.Reference,
		"reference_latest": tb.ReferenceLatest,
		"status":           tb.Status,
		"tool_count":       tb.ToolCount,
		"created_at":       formatTimestamp(tb.CreatedAt),
		"updated_at":       formatTimestamp(tb.UpdatedAt),
	}
}

// reconcileToolbeltTools reports the toolbelt's pinned members in the form
// they were configured, so an unpinned tool in the configuration does not
// show a diff against the version the API pinned it to.
func reconcileToolbeltTools(configured, pinned []string) []string {
	exact := make(map[string]bool, len(configured))
	for _, ref := range configured {
		exact[ref] = true
	}

	out := make([]string, 0, len(pinned))
	for _, ref := range pinned {
		slug, _ := splitToolPin(ref)
		if !exact[ref] && exact[slug] {
			out = append(out, slug)
		} else {
			out = append(out, ref)
		}
	}
	sort.Strings(out)
	return out
}

// toolbeltMembershipChanges compares the toolbelt's current pinned members
// with the desired tools and returns the slugs to remove and the references
// to add. A desired tool without a pin is satisfied by any pinned version of
// that tool, and a tool whose pin changes is removed and added again.
func toolbeltMembershipChanges(currentPinned, desired []string) (remove, add []string) {
	current := make(map[string]string, len(currentPinned))
	for _, ref := range currentPinned {
		slug, version := splitToolPin(ref)
		current[slug] = version
	}

	wanted := make(map[string]bool, len(desired))
	for _, ref := range desired {
		slug, version := splitToolPin(ref)
		wanted[slug] = true

		currentVersion, isMember := current[slug]
		switch {
		case !isMember:
			add = append(add, ref)
		case version != "" && version != currentVersion:
			remove = append(remove, slug)
			add = append(add, ref)
		}
	}

	for slug := range current {
		if !wanted[slug] {
			remove = append(remove, slug)
		}
	}

	sort.Strings(remove)
	sort.Strings(add)
	return remove, add
}
