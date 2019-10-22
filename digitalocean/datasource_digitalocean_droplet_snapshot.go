package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceDigitalOceanDropletSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanDropletSnapshotRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.ValidateRegexp,
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
			"droplet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

// dataSourceDoSnapshotRead performs the Snapshot lookup.
func dataSourceDigitalOceanDropletSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	name, hasName := d.GetOk("name")
	nameRegex, hasNameRegex := d.GetOk("name_regex")
	region, hasRegion := d.GetOk("region")

	if !hasName && !hasNameRegex {
		return fmt.Errorf("One of `name` or `name_regex` must be assigned")
	}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var snapshotList []godo.Snapshot

	for {
		snapshots, resp, err := client.Snapshots.ListDroplet(context.Background(), opts)

		if err != nil {
			return fmt.Errorf("Error retrieving Droplet snapshots: %s", err)
		}

		for _, s := range snapshots {
			snapshotList = append(snapshotList, s)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return fmt.Errorf("Error retrieving Droplet snapshots: %s", err)
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
		return fmt.Errorf("No DROPLET snapshot found with name %s", name)
	}
	if len(snapshotList) > 1 {
		recent := d.Get("most_recent").(bool)
		if recent {
			snapshot = findMostRecentSnapshot(snapshotList)
		} else {
			return fmt.Errorf("too many Droplet snapshots found with name %s (found %d, expected 1)", name, len(snapshotList))
		}
	} else {
		snapshot = &snapshotList[0]
	}

	log.Printf("[DEBUG] do_snapshot - Single Droplet Snapshot found: %s", snapshot.ID)

	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("created_at", snapshot.Created)
	d.Set("min_disk_size", snapshot.MinDiskSize)
	d.Set("regions", snapshot.Regions)
	d.Set("droplet_id", snapshot.ResourceID)
	d.Set("size", snapshot.SizeGigaBytes)

	return nil
}
