package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateShortExternalID creates a short and unique external ID
func GenerateShortExternalID() string {
	// Generate random bytes for uniqueness
	randomBytes := make([]byte, 3) // 3 bytes = 6 characters in hex
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // Handle error as appropriate in production
	}

	// Format ID: TX-{timestamp}-{random}
	timestamp := time.Now().Unix()
	return fmt.Sprintf("TX-%d-%s", timestamp, hex.EncodeToString(randomBytes))
}
