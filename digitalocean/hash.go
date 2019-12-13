package digitalocean

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func HashString(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func HashStringStateFunc() schema.SchemaStateFunc {
	return func(v interface{}) string {
		switch v.(type) {
		case string:
			return HashString(v.(string))
		default:
			return ""
		}
	}
}
