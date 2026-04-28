package xtypes

import (
	"errors"
	"testing"
)

func TestErrorConstructorsHaveCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		xerr interface{ Code() int }
		code int
	}{
		{"path", newPathError(errors.New("boom"), "a"), PathErrorCode},
		{"chmod", newPathChmodError(errors.New("boom"), "a"), PathChmodErrorCode},
		{"chown", newPathChownError(errors.New("boom"), "a"), PathChownErrorCode},
		{"create", newPathCreateError(errors.New("boom"), "a"), PathCreateErrorCode},
		{"open", newPathOpenFileError(errors.New("boom"), "a"), PathOpenFileErrorCode},
		{"write", newPathWriteError(errors.New("boom"), "a"), PathWriteErrorCode},
	}
	for _, tc := range tests {
		if tc.xerr.Code() != tc.code {
			t.Fatalf("%s code=%d want %d", tc.name, tc.xerr.Code(), tc.code)
		}
	}
}
