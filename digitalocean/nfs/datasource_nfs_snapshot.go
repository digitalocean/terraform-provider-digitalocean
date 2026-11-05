package nfs

import (
	"context"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanNfsSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanNfsSnapshotRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name"},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tag.TagsDataSourceSchema(),
		},
	}
}

// dataSourceDoSnapshotRead performs the Snapshot lookup.
func dataSourceDigitalOceanNfsSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	name, hasName := d.GetOk("name")
	nameRegex, hasNameRegex := d.GetOk("name_regex")
	region := ""
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}
	shareID := d.Get("share_id").(string)

	if !hasName && !hasNameRegex {
		return diag.Errorf("One of `name` or `name_regex` must be assigned")
	}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var snapshotList []godo.NfsSnapshot

	for {
		snapshots, resp, err := client.Nfs.ListSnapshots(context.Background(), opts, shareID, region)

		if err != nil {
			return diag.Errorf("Error retrieving share snapshots: %s", err)
		}

		// Dereference pointers and append
		for _, s := range snapshots {
			snapshotList = append(snapshotList, *s)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return diag.Errorf("Error retrieving share snapshots: %s", err)
		}

		opts.Page = page + 1
	}

	// Go through all the possible filters
	if hasName {
		snapshotList = filterSnapshotsByName(snapshotList, name.(string))
	} else {
		snapshotList = filterSnapshotsByNameRegex(snapshotList, nameRegex.(string))
	}

	// Get the queried snapshot or fail if it can't be determined
	var snapshot *godo.NfsSnapshot
	if len(snapshotList) == 0 {
		return diag.Errorf("no share snapshot found with name %s", name)
	}
	if len(snapshotList) > 1 {
		recent := d.Get("most_recent").(bool)
		if recent {
			snapshot = findMostRecentSnapshot(snapshotList)
		} else {
			return diag.Errorf("too many share snapshots found with name %s (found %d, expected 1)", name, len(snapshotList))
		}
	} else {
		snapshot = &snapshotList[0]
	}

	log.Printf("[DEBUG] do_snapshot - Single Share Snapshot found: %s", snapshot.ID)

	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("created_at", snapshot.CreatedAt)
	d.Set("region", snapshot.Region)
	d.Set("share_id", snapshot.ShareID)
	d.Set("size", snapshot.SizeGib)

	return nil
}

func filterSnapshotsByName(snapshots []godo.NfsSnapshot, name string) []godo.NfsSnapshot {
	result := make([]godo.NfsSnapshot, 0)
	for _, s := range snapshots {
		if s.Name == name {
			result = append(result, s)
		}
	}
	return result
}

func filterSnapshotsByNameRegex(snapshots []godo.NfsSnapshot, name string) []godo.NfsSnapshot {
	r := regexp.MustCompile(name)
	result := make([]godo.NfsSnapshot, 0)
	for _, s := range snapshots {
		if r.MatchString(s.Name) {
			result = append(result, s)
		}
	}
	return result
}

// Returns the most recent Snapshot out of a slice of Snapshots.
func findMostRecentSnapshot(snapshots []godo.NfsSnapshot) *godo.NfsSnapshot {
	sort.Slice(snapshots, func(i, j int) bool {
		itime, _ := time.Parse(time.RFC3339, snapshots[i].CreatedAt)
		jtime, _ := time.Parse(time.RFC3339, snapshots[j].CreatedAt)
		return itime.Unix() > jtime.Unix()
	})

	return &snapshots[0]
}
