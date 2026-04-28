package xtypes

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
)

func TestFileModeAndJSONFileMode(t *testing.T) {
	t.Parallel()

	mode := FileMode(0o644)
	if got := (&mode).String(); got == "" {
		t.Fatalf("empty mode string")
	}
	if got := (*FileMode)(nil).String(); !strings.HasPrefix(got, "0") {
		t.Fatalf("nil String got %q", got)
	}

	if (&mode).OSFileMode() == 0 {
		t.Fatalf("expected non-zero os file mode")
	}
	neg := FileMode(-1)
	if (&neg).OSFileMode() != 0 {
		t.Fatalf("negative should map to 0")
	}
	huge := FileMode(math.MaxInt64)
	if (&huge).OSFileMode() != 0 {
		t.Fatalf("huge mode should map to 0")
	}

	j, err := (&mode).MarshalJSON()
	if err != nil {
		t.Fatalf("marshal mode: %v", err)
	}
	var parsed FileMode
	if err := (&parsed).UnmarshalJSON(j); err != nil {
		t.Fatalf("unmarshal mode from string json: %v", err)
	}
	if parsed != mode {
		t.Fatalf("mode mismatch")
	}
	if err := (&parsed).UnmarshalJSON([]byte(`420`)); err != nil {
		t.Fatalf("unmarshal numeric mode: %v", err)
	}
	if err := (&parsed).UnmarshalText([]byte("0644")); err != nil {
		t.Fatalf("unmarshal text mode: %v", err)
	}
	if err := (&parsed).UnmarshalText([]byte("bad")); err == nil {
		t.Fatalf("expected mode text error")
	}
	text, err := (&mode).MarshalText()
	if err != nil {
		t.Fatalf("marshal text mode: %v", err)
	}
	if len(text) == 0 {
		t.Fatalf("expected mode text")
	}

	var jm JSONFileMode = JSONFileMode(mode)
	blob, err := json.Marshal(jm)
	if err != nil {
		t.Fatalf("marshal JSONFileMode: %v", err)
	}
	var jm2 JSONFileMode
	if err := json.Unmarshal(blob, &jm2); err != nil {
		t.Fatalf("unmarshal JSONFileMode: %v", err)
	}
	if FileMode(jm2) != mode {
		t.Fatalf("json file mode mismatch")
	}
	text, err = jm.MarshalText()
	if err != nil {
		t.Fatalf("marshal text JSONFileMode: %v", err)
	}
	if err := (&jm2).UnmarshalText(text); err != nil {
		t.Fatalf("unmarshal text JSONFileMode: %v", err)
	}
	if err := (&jm2).UnmarshalJSON([]byte(`{}`)); err == nil {
		t.Fatalf("expected JSONFileMode error")
	}
}

