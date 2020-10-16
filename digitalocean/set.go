package digitalocean

import (
	"strings"
)

// Helper function for sets of strings that are case insensitive
func HashStringIgnoreCase(v interface{}) int {
	return SDKHashString(strings.ToLower(v.(string)))
}
