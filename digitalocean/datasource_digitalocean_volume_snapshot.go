package digitalocean

import (
	"context"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanVolumeSnapshotRead,
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
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Computed values.
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"tags": tagsDataSourceSchema(),
		},
	}
}

// dataSourceDoSnapshotRead performs the Snapshot lookup.
func dataSourceDigitalOceanVolumeSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	name, hasName := d.GetOk("name")
	nameRegex, hasNameRegex := d.GetOk("name_regex")
	region, hasRegion := d.GetOk("region")

	if !hasName && !hasNameRegex {
		return diag.Errorf("One of `name` or `name_regex` must be assigned")
	}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var snapshotList []godo.Snapshot

	for {
		snapshots, resp, err := client.Snapshots.ListVolume(context.Background(), opts)

		if err != nil {
			return diag.Errorf("Error retrieving volume snapshots: %s", err)
		}

		for _, s := range snapshots {
			snapshotList = append(snapshotList, s)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return diag.Errorf("Error retrieving volume snapshots: %s", err)
		}

		opts.Page = page + 1
	}

	// Go through all the possible filters
	if hasName {
		snapshotList = filterSnapshotsByName(snapshotList, name.(string))
	} else {
		snapshotList = filterSnapshotsByNameRegex(snapshotList, nameRegex.(string))
	}
	if hasRegion {
		snapshotList = filterSnapshotsByRegion(snapshotList, region.(string))
	}

	// Get the queried snapshot or fail if it can't be determined
	var snapshot *godo.Snapshot
	if len(snapshotList) == 0 {
		return diag.Errorf("no volume snapshot found with name %s", name)
	}
	if len(snapshotList) > 1 {
		recent := d.Get("most_recent").(bool)
		if recent {
			snapshot = findMostRecentSnapshot(snapshotList)
		} else {
			return diag.Errorf("too many volume snapshots found with name %s (found %d, expected 1)", name, len(snapshotList))
		}
	} else {
		snapshot = &snapshotList[0]
	}

	log.Printf("[DEBUG] do_snapshot - Single Volume Snapshot found: %s", snapshot.ID)

	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("created_at", snapshot.Created)
	d.Set("min_disk_size", snapshot.MinDiskSize)
	d.Set("regions", snapshot.Regions)
	d.Set("volume_id", snapshot.ResourceID)
	d.Set("size", snapshot.SizeGigaBytes)
	d.Set("tags", flattenTags(snapshot.Tags))

	return nil
}

func filterSnapshotsByName(snapshots []godo.Snapshot, name string) []godo.Snapshot {
	result := make([]godo.Snapshot, 0)
	for _, s := range snapshots {
		if s.Name == name {
			result = append(result, s)
		}
	}
	return result
}

func filterSnapshotsByNameRegex(snapshots []godo.Snapshot, name string) []godo.Snapshot {
	r := regexp.MustCompile(name)
	result := make([]godo.Snapshot, 0)
	for _, s := range snapshots {
		if r.MatchString(s.Name) {
			result = append(result, s)
		}
	}
	return result
}

func filterSnapshotsByRegion(snapshots []godo.Snapshot, region string) []godo.Snapshot {
	result := make([]godo.Snapshot, 0)
	for _, s := range snapshots {
		for _, r := range s.Regions {
			if r == region {
				result = append(result, s)
				break
			}
		}
	}
	return result
}

// Returns the most recent Snapshot out of a slice of Snapshots.
func findMostRecentSnapshot(snapshots []godo.Snapshot) *godo.Snapshot {
	sort.Slice(snapshots, func(i, j int) bool {
		itime, _ := time.Parse(time.RFC3339, snapshots[i].Created)
		jtime, _ := time.Parse(time.RFC3339, snapshots[j].Created)
		return itime.Unix() > jtime.Unix()
	})

	return &snapshots[0]
}
