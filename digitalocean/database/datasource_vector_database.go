package database

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanVectorDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanVectorDatabaseRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "name"},
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"owner_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tag.TagsDataSourceSchema(),

			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_quantization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_auto_schema": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"weaviate_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grpc": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

func dataSourceDigitalOceanVectorDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var vectorDB *godo.VectorDB

	if id, ok := d.GetOk("id"); ok {
		found, _, err := client.VectorDBs.Get(ctx, id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving vector database: %s", err)
		}
		vectorDB = found
	} else {
		name := d.Get("name").(string)

		opts := &godo.ListOptions{
			Page:    1,
			PerPage: 200,
		}

		var vectorDBList []godo.VectorDB

		for {
			vectorDBs, resp, err := client.VectorDBs.List(ctx, opts)
			if err != nil {
				return diag.Errorf("Error retrieving vector databases: %s", err)
			}

			vectorDBList = append(vectorDBList, vectorDBs...)

			if resp.Links == nil || resp.Links.IsLastPage() {
				break
			}

			page, err := resp.Links.CurrentPage()
			if err != nil {
				return diag.Errorf("Error retrieving vector databases: %s", err)
			}

			opts.Page = page + 1
		}

		for i := range vectorDBList {
			if vectorDBList[i].Name == name {
				vectorDB = &vectorDBList[i]
				break
			}
		}

		if vectorDB == nil {
			return diag.Errorf("Unable to find vector database with name: %s", name)
		}
	}

	d.SetId(vectorDB.ID)
	d.Set("name", vectorDB.Name)
	d.Set("region", vectorDB.Region)
	d.Set("size", vectorDB.Size)
	d.Set("status", vectorDB.Status)
	d.Set("owner_uuid", vectorDB.OwnerUUID)
	d.Set("created_at", vectorDB.CreatedAt.UTC().String())
	d.Set("updated_at", vectorDB.UpdatedAt.UTC().String())
	d.Set("tags", tag.FlattenTags(vectorDB.Tags))

	if err := d.Set("config", flattenVectorDBConfig(vectorDB.Config)); err != nil {
		return diag.Errorf("Error setting vector database config: %s", err)
	}

	if err := d.Set("endpoints", flattenVectorDBEndpoints(vectorDB.Endpoints)); err != nil {
		return diag.Errorf("Error setting vector database endpoints: %s", err)
	}

	return nil
}
