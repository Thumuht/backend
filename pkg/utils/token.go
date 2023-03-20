package utils

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

func GenRandStr(n int) string {
	randomBytes := make([]byte, n)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)
}

// generates a random token for user auth.
// we do not use JWT
func GenToken() string {
	return strings.ToLower(GenRandStr(20))
}
