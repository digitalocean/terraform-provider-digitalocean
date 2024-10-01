package util

import (
	"crypto/sha1"
	"encoding/hex"
	"hash/crc32"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HashString produces a hash of a string.
func HashString(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// HashStringIgnoreCase is a helper function for sets of strings that are case insensitive
func HashStringIgnoreCase(v interface{}) int {
	return SDKHashString(strings.ToLower(v.(string)))
}

// HashStringStateFunc implements a schema.SchemaStateFunc with HashString
func HashStringStateFunc() schema.SchemaStateFunc {
	return func(v interface{}) string {
		switch v := v.(type) {
		case string:
			return HashString(v)
		default:
			return ""
		}
	}
}

// SDKHashString implements hashcode.String from the terraform-plugin-sdk which
// was made internal to the SDK in v2. Embed the implementation here to allow
// same hash function to continue to be used by the code in this provider that
// used it for hash computation.
func SDKHashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
