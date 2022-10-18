package helpers

import (
	"crypto/sha256"
	"fmt"
)

func Hash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	token := hash.Sum(nil)

	return fmt.Sprintf("%x", token) // sha256 hex
}
