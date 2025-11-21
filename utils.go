package main

import (
	"crypto/rand"
	"encoding/hex"
)

// Stack overflow
func generateRandomToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
