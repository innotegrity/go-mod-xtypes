package xtypes

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// uuidByteLen is the length of a UUID in bytes (RFC 4122).
const uuidByteLen = 16

// NewUUID generates a new UUID.
//
// This function first attempts to generate a v7 UUID. If that fails, a v8 UUID is generated instead.
func NewUUID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return generateUUIDv8()
	}

	return strings.ToUpper(id.String()), nil
}

// generateUUIDv8 generates a v8 UUID when v7 generation is unavailable.
func generateUUIDv8() (string, error) {
	vals := make([]byte, uuidByteLen)

	_, err := rand.Read(vals)
	if err != nil {
		return "", fmt.Errorf("generate random bytes for UUID v8: %w", err)
	}

	// replace bits 48-51 with the version (8)
	vals[6] = (((vals[6] << 4) & 255) >> 4) | 128 //nolint:mnd // fixed layout for UUID version field

	// replace bits 64 and 65 with the variant (2)
	vals[8] = (((vals[8] << 2) & 255) >> 2) | 128 //nolint:mnd // fixed layout for RFC 4122 variant

	return strings.ToUpper(fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(vals[0:4]),
		hex.EncodeToString(vals[4:6]),
		hex.EncodeToString(vals[6:8]),
		hex.EncodeToString(vals[8:10]),
		hex.EncodeToString(vals[10:]))), nil
}
