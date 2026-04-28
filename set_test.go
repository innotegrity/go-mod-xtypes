package xtypes

import "testing"

func TestSetOperations(t *testing.T) {
	t.Parallel()

	s1 := NewSet("a", "b")
	s1.Add("b", "c")
	if !s1.Contains("a") || s1.Contains("z") {
		t.Fatalf("contains mismatch")
	}

	s2 := NewSet("b", "d")
	if len(s1.Intersection(s2)) != 1 {
		t.Fatalf("intersection mismatch")
	}
	if len(s1.Union(s2)) != 4 {
		t.Fatalf("union mismatch")
	}
	if s1.String() == "" || len(s1.Members()) == 0 {
		t.Fatalf("set output mismatch")
	}
}

