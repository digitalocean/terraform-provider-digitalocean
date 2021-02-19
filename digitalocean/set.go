package digitalocean

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Helper function for sets of strings that are case insensitive
func HashStringIgnoreCase(v interface{}) int {
	return SDKHashString(strings.ToLower(v.(string)))
}

// getSetChanges compares two *schema.Set, "old" and "new." It returns one
// *schema.Set only containing items not found in the "new" set and another
// containing items not found in the "old" set.
//
// Originally written to update the resources in a project.
func getSetChanges(old *schema.Set, new *schema.Set) (remove *schema.Set, add *schema.Set) {
	remove = schema.NewSet(old.F, []interface{}{})
	for _, x := range old.List() {
		if !new.Contains(x) {
			remove.Add(x)
		}
	}

	add = schema.NewSet(new.F, []interface{}{})
	for _, x := range new.List() {
		if !old.Contains(x) {
			add.Add(x)
		}
	}

	return remove, add
}
