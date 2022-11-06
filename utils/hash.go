package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// SHA256Sum returns the SHA256 sum of the given string keys.
func SHA256Sum(keys ...string) string {
	s := strings.Join(keys, "-")

	h := sha256.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}
