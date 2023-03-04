package utils

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

// generates a random token for user auth.
// we do not use JWT
func GenToken() string {
	randomBytes := make([]byte, 20)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(base32.StdEncoding.EncodeToString(randomBytes))
}
