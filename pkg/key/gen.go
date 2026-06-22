package key

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GenerateRandomKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("failed to generate random key: " + err.Error())
	}
	return base64.StdEncoding.EncodeToString(key)
}

func GenerateDBPassword() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Errorf("failed to generate db password: %w", err))
	}
	return hex.EncodeToString(b)
}
