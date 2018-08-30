package digitalocean

import (
	"crypto/sha1"
	"encoding/hex"
)

func HashString(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
