package jwtsession

import (
	"crypto/rand"
	"fmt"
)

// GenerateRandToken generates a random token
func GenerateRandToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}
