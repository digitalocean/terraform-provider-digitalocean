package database

import (
	"context"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanDatabaseCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDatabaseClusterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"maintenance_window": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hour": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"private_uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"database": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_network_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tag.TagsSchema(),

			"storage_size_mib": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanDatabaseClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	name := d.Get("name").(string)

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var databaseList []godo.Database

	for {
		databases, resp, err := client.Databases.List(context.Background(), opts)
		if err != nil {
			return diag.Errorf("Error retrieving DatabaseClusters: %s", err)
		}

		databaseList = append(databaseList, databases...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return diag.Errorf("Error retrieving DatabaseClusters: %s", err)
		}

		opts.Page = page + 1
	}

	if len(databaseList) == 0 {
		return diag.Errorf("Unable to find any database clusters")
	}

	for _, db := range databaseList {
		if db.Name == name {
			d.SetId(db.ID)

			d.Set("name", db.Name)
			d.Set("engine", db.EngineSlug)
			d.Set("version", db.VersionSlug)
			d.Set("size", db.SizeSlug)
			d.Set("region", db.RegionSlug)
			d.Set("node_count", db.NumNodes)
			d.Set("tags", tag.FlattenTags(db.Tags))
			d.Set("storage_size_mib", strconv.FormatUint(db.StorageSizeMib, 10))

			if _, ok := d.GetOk("maintenance_window"); ok {
				if err := d.Set("maintenance_window", flattenMaintWindowOpts(*db.MaintenanceWindow)); err != nil {
					return diag.Errorf("[DEBUG] Error setting maintenance_window - error: %#v", err)
				}
			}

			err := setDatabaseConnectionInfo(&db, d)
			if err != nil {
				return diag.Errorf("Error setting connection info for database cluster: %s", err)
			}
			d.Set("urn", db.URN())
			d.Set("private_network_uuid", db.PrivateNetworkUUID)
			d.Set("project_id", db.ProjectID)

			break
		}
	}

	return nil
}
