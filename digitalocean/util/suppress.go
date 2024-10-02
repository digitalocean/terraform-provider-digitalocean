package util

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CaseSensitive implements a schema.SchemaDiffSuppressFunc that ignores case
func CaseSensitive(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
