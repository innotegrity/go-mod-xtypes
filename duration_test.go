package xtypes

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDurationAndJSONDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want time.Duration
	}{
		{"", 0},
		{"2h", 2 * time.Hour},
		{"1d", 24 * time.Hour},
		{"2w", 14 * 24 * time.Hour},
		{"1mo", 30 * 24 * time.Hour},
		{"1y", 365 * 24 * time.Hour},
	}
	for _, tc := range tests {
		got, err := ParseDuration(tc.in)
		if err != nil {
			t.Fatalf("ParseDuration(%q): %v", tc.in, err)
		}
		if time.Duration(got) != tc.want {
			t.Fatalf("ParseDuration(%q)=%v want %v", tc.in, time.Duration(got), tc.want)
		}
	}

	if _, err := ParseDuration("abc"); err == nil {
		t.Fatalf("expected parse error")
	}
	if _, err := ParseDuration("999999999999999999999999999y"); err == nil {
		t.Fatalf("expected overflow error")
	}

	d := Duration(2 * time.Hour)
	if got := (&d).String(); got != "2h0m0s" {
		t.Fatalf("String got %q", got)
	}
	if got := (*Duration)(nil).String(); got != "0s" {
		t.Fatalf("nil String got %q", got)
	}

	j, err := (&d).MarshalJSON()
	if err != nil {
		t.Fatalf("marshal duration: %v", err)
	}
	var parsed Duration
	if err := (&parsed).UnmarshalJSON(j); err != nil {
		t.Fatalf("unmarshal duration: %v", err)
	}
	if parsed != d {
		t.Fatalf("parsed duration mismatch")
	}
	if err := (&parsed).UnmarshalJSON([]byte(`123`)); err != nil {
		t.Fatalf("unmarshal numeric duration: %v", err)
	}
	if err := (&parsed).UnmarshalText([]byte("1h")); err != nil {
		t.Fatalf("unmarshal text duration: %v", err)
	}
	if err := (&parsed).UnmarshalJSON([]byte(`{}`)); err == nil {
		t.Fatalf("expected JSON error")
	}

	var jd JSONDuration = JSONDuration(d)
	blob, err := json.Marshal(jd)
	if err != nil {
		t.Fatalf("marshal JSONDuration: %v", err)
	}
	var jd2 JSONDuration
	if err := json.Unmarshal(blob, &jd2); err != nil {
		t.Fatalf("unmarshal JSONDuration: %v", err)
	}
	if Duration(jd2) != d {
		t.Fatalf("json duration mismatch")
	}
	txt, err := jd.MarshalText()
	if err != nil {
		t.Fatalf("marshal text JSONDuration: %v", err)
	}
	if err := (&jd2).UnmarshalText(txt); err != nil {
		t.Fatalf("unmarshal text JSONDuration: %v", err)
	}
	if err := (&jd2).UnmarshalJSON([]byte(`{}`)); err == nil {
		t.Fatalf("expected JSONDuration error")
	}
}

