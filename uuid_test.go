package xtypes

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
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

func TestUUIDFallbackAndReadErrorBranches(t *testing.T) {
	origNewV7 := uuidNewV7
	origRead := uuidRead
	t.Cleanup(func() {
		uuidNewV7 = origNewV7
		uuidRead = origRead
	})

	uuidNewV7 = func() (uuid.UUID, error) { return uuid.UUID{}, errors.New("boom") }
	uuidRead = func(p []byte) (int, error) {
		for i := range p {
			p[i] = byte(i)
		}
		return len(p), nil
	}

	id, err := NewUUID()
	if err != nil {
		t.Fatalf("fallback NewUUID: %v", err)
	}
	if len(id) != 36 {
		t.Fatalf("fallback uuid len=%d", len(id))
	}

	uuidRead = func([]byte) (int, error) { return 0, errors.New("rand fail") }
	if _, err := generateUUIDv8(); err == nil {
		t.Fatalf("expected generateUUIDv8 read error")
	}
}

