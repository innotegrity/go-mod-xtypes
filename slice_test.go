package xtypes

import "testing"

func TestAnySlice(t *testing.T) {
	t.Parallel()

	anyVals := AnySlice([]int{1, 2, 3})
	if len(anyVals) != 3 {
		t.Fatalf("AnySlice len mismatch")
	}
}

