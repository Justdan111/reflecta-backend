package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken returns a SHA-256 hash so raw tokens are never stored in the database
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
