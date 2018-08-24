package digitalocean

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func CaseSensitive(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}
