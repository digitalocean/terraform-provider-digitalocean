package gradientai

import (
	"context"
	"fmt"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanCustomModel() *schema.Resource {
	recordSchema := CustomModelSchemaRead().Schema

	// The lookup key is the model UUID; remaining attributes are computed.
	recordSchema["uuid"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "UUID of the custom model to look up.",
		ValidateFunc: validation.NoZeroValues,
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanCustomModelRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanCustomModelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	uuid := d.Get("uuid").(string)
	model, _, err := client.GradientAI.GetCustomModel(ctx, uuid)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading custom model (%s): %w", uuid, err))
	}
	if model == nil {
		return diag.Errorf("custom model (%s) not found", uuid)
	}

	d.SetId(model.Uuid)

	flat, err := FlattenDigitalOceanCustomModel(model)
	if err != nil {
		return diag.FromErr(err)
	}
	for key, value := range flat {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(fmt.Errorf("error setting %s: %w", key, err))
		}
	}
	return nil
}
