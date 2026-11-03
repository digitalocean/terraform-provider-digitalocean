package actiongateway

import (
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	// userIDRegexp matches connection user IDs and session actor IDs.
	userIDRegexp      = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)
	toolbeltNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
	serverRefRegexp   = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)

	validateUserID = validation.StringMatch(userIDRegexp, "must be 1 to 64 characters from [A-Za-z0-9._-]")
)

func isNotFound(resp *godo.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

func formatTimestamp(ts *godo.Timestamp) string {
	if ts == nil || ts.IsZero() {
		return ""
	}
	return ts.UTC().Format(time.RFC3339)
}

func expandStringSet(v interface{}) []string {
	set, ok := v.(*schema.Set)
	if !ok || set == nil {
		return nil
	}
	out := make([]string, 0, set.Len())
	for _, item := range set.List() {
		out = append(out, item.(string))
	}
	sort.Strings(out)
	return out
}

func expandStringList(v interface{}) []string {
	raw, _ := v.([]interface{})
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// flattenStringList returns in as the []interface{} that datalist filters
// and sorts expect for TypeList values.
func flattenStringList(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}

// computedOnly converts a resource schema, including nested blocks, into the
// computed-only form a data source exposes.
func computedOnly(in map[string]*schema.Schema) map[string]*schema.Schema {
	out := make(map[string]*schema.Schema, len(in))
	for k, v := range in {
		clone := *v
		clone.Required = false
		clone.Optional = false
		clone.ForceNew = false
		clone.Default = nil
		clone.ValidateFunc = nil
		clone.DiffSuppressFunc = nil
		clone.ConflictsWith = nil
		clone.MaxItems = 0
		clone.MinItems = 0
		clone.Computed = true
		if elem, ok := v.Elem.(*schema.Resource); ok {
			clone.Elem = &schema.Resource{Schema: computedOnly(elem.Schema)}
		} else if elem, ok := v.Elem.(*schema.Schema); ok {
			elemClone := *elem
			elemClone.ValidateFunc = nil
			clone.Elem = &elemClone
		}
		out[k] = &clone
	}
	return out
}

// splitToolPin splits a tool reference such as `exa_web_search@1` into its
// slug and pinned version. The version is empty for an unpinned reference.
func splitToolPin(ref string) (slug, version string) {
	if i := strings.LastIndex(ref, "@"); i >= 0 {
		return ref[:i], ref[i+1:]
	}
	return ref, ""
}
