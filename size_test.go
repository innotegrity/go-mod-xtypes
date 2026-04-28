package xtypes

import (
	"encoding/json"
	"math"
	"testing"
)

func TestSizeAndJSONSize(t *testing.T) {
	t.Parallel()

	cases := []string{"", "100", "1kb", "1kib", "1mb", "1mib", "1gb", "1gib", "1tb", "1tib", "1pb", "1pib"}
	for _, in := range cases {
		if _, err := ParseSize(in); err != nil {
			t.Fatalf("ParseSize(%q): %v", in, err)
		}
	}
	if _, err := ParseSize("x"); err == nil {
		t.Fatalf("expected ParseSize error")
	}
	if _, err := ParseSize("1e309kb"); err == nil {
		t.Fatalf("expected ParseSize overflow/parse error")
	}
	if _, err := scaleSize(math.MaxFloat64, 2); err == nil {
		t.Fatalf("expected scaleSize overflow")
	}

	s := Size(1024)
	if got := (&s).String(); got == "" {
		t.Fatalf("empty size string")
	}
	if got := (*Size)(nil).String(); got == "" {
		t.Fatalf("empty nil size string")
	}
	j, err := (&s).MarshalJSON()
	if err != nil {
		t.Fatalf("marshal size: %v", err)
	}
	var parsed Size
	if err := (&parsed).UnmarshalJSON(j); err != nil {
		t.Fatalf("unmarshal size: %v", err)
	}
	if err := (&parsed).UnmarshalText([]byte("1kb")); err != nil {
		t.Fatalf("unmarshal text size: %v", err)
	}
	if err := (&parsed).UnmarshalJSON([]byte(`{}`)); err == nil {
		t.Fatalf("expected size json error")
	}
	if _, err := (&s).MarshalText(); err != nil {
		t.Fatalf("marshal text size: %v", err)
	}

	thresholds := []*Size{
		func() *Size { v := Size(500); return &v }(),
		func() *Size { v := Size(5000); return &v }(),
		func() *Size { v := Size(5_000_000); return &v }(),
		func() *Size { v := Size(5_000_000_000); return &v }(),
		func() *Size { v := Size(5_000_000_000_000); return &v }(),
		func() *Size { v := Size(5_000_000_000_000_000); return &v }(),
	}
	for _, v := range thresholds {
		if got := v.String(); got == "" {
			t.Fatalf("expected non-empty String")
		}
	}

	var js JSONSize = JSONSize(s)
	blob, err := json.Marshal(js)
	if err != nil {
		t.Fatalf("marshal JSONSize: %v", err)
	}
	var js2 JSONSize
	if err := json.Unmarshal(blob, &js2); err != nil {
		t.Fatalf("unmarshal JSONSize: %v", err)
	}
	if Size(js2) == 0 {
		t.Fatalf("expected non-zero JSONSize after unmarshal")
	}
	text, err := js.MarshalText()
	if err != nil {
		t.Fatalf("marshal text JSONSize: %v", err)
	}
	if err := (&js2).UnmarshalText(text); err != nil {
		t.Fatalf("unmarshal text JSONSize: %v", err)
	}
	if err := (&js2).UnmarshalJSON([]byte(`{}`)); err == nil {
		t.Fatalf("expected JSONSize error")
	}
}

