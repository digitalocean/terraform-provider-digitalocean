package tag

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var tagNameRe = regexp.MustCompile("^[a-zA-Z0-9:\\-_]{1,255}$")

func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: ValidateTag,
		},
		Set: util.HashStringIgnoreCase,
	}
}

func TagsDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func ValidateTag(value interface{}, key string) ([]string, []error) {
	if !tagNameRe.MatchString(value.(string)) {
		return nil, []error{fmt.Errorf("tags may contain lowercase letters, numbers, colons, dashes, and underscores; there is a limit of 255 characters per tag")}
	}

	return nil, nil
}

// SetTags is a helper to set the tags for a resource. It expects the
// tags field to be named "tags"
func SetTags(conn *godo.Client, d *schema.ResourceData, resourceType godo.ResourceType) error {
	oraw, nraw := d.GetChange("tags")
	remove, create := DiffTags(TagsFromSchema(oraw), TagsFromSchema(nraw))

	log.Printf("[DEBUG] Removing tags: %#v from %s", remove, d.Id())
	for _, tag := range remove {
		_, err := conn.Tags.UntagResources(context.Background(), tag, &godo.UntagResourcesRequest{
			Resources: []godo.Resource{
				{
					ID:   d.Id(),
					Type: resourceType,
				},
			},
		})
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Creating tags: %s for %s", create, d.Id())
	for _, tag := range create {

		createdTag, _, err := conn.Tags.Create(context.Background(), &godo.TagCreateRequest{
			Name: tag,
		})
		if err != nil {
			return err
		}

		_, err = conn.Tags.TagResources(context.Background(), createdTag.Name, &godo.TagResourcesRequest{
			Resources: []godo.Resource{
				{
					ID:   d.Id(),
					Type: resourceType,
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// TagsFromSchema takes the raw schema tags and returns them as a
// properly asserted map[string]string
func TagsFromSchema(raw interface{}) map[string]string {
	result := make(map[string]string)
	for _, t := range raw.(*schema.Set).List() {
		result[t.(string)] = t.(string)
	}

	return result
}

// DiffTags takes the old and the new tag sets and returns the difference of
// both. The remaining tags are those that need to be removed and created
func DiffTags(oldTags, newTags map[string]string) (map[string]string, map[string]string) {
	for k := range oldTags {
		_, ok := newTags[k]
		if ok {
			delete(newTags, k)
			delete(oldTags, k)
		}
	}

	return oldTags, newTags
}

func ExpandTags(tags []interface{}) []string {
	expandedTags := make([]string, len(tags))
	for i, v := range tags {
		expandedTags[i] = v.(string)
	}

	return expandedTags
}

func FlattenTags(tags []string) *schema.Set {
	if tags == nil {
		return nil
	}

	flattenedTags := schema.NewSet(util.HashStringIgnoreCase, []interface{}{})
	for _, v := range tags {
		if v != "" {
			flattenedTags.Add(v)
		}
	}

	return flattenedTags
}
