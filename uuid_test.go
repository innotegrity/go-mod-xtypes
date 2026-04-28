package xtypes

import (
	"strings"
	"testing"
)

func TestUUIDGeneration(t *testing.T) {
	t.Parallel()

	id, err := NewUUID()
	if err != nil {
		t.Fatalf("NewUUID: %v", err)
	}
	if len(id) != 36 {
		t.Fatalf("uuid len=%d", len(id))
	}
	if strings.ToUpper(id) != id {
		t.Fatalf("uuid should be upper-case: %q", id)
	}

	id2, err := generateUUIDv8()
	if err != nil {
		t.Fatalf("generateUUIDv8: %v", err)
	}
	if len(id2) != 36 {
		t.Fatalf("uuidv8 len=%d", len(id2))
	}
}
