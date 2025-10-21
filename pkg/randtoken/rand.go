package randtoken

import (
	"crypto/rand"
	"fmt"
)

// New generates a random token
func New(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}
