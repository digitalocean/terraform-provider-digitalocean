package digitalocean

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

// Helper function for sets of strings that are case insensitive
func HashStringIgnoreCase(v interface{}) int {
	return hashcode.String(strings.ToLower(v.(string)))
}
