package xtypes

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.innotegrity.dev/mod/xerrors"
)

func TestErrorConstructorsHaveCodes(t *testing.T) {
	t.Parallel()
	ctx := xerrors.ContextWithErrorOptions(context.Background(), xerrors.WithCaptureCaller())

	tests := []struct {
		name string
		xerr xerrors.Error
		code int
	}{
		{"path", newPathError(ctx, errors.New("boom"), "a"), PathErrorCode},
		{"chmod", newPathChmodError(ctx, errors.New("boom"), "a"), PathChmodErrorCode},
		{"chown", newPathChownError(ctx, errors.New("boom"), "a"), PathChownErrorCode},
		{"create", newPathCreateError(ctx, errors.New("boom"), "a"), PathCreateErrorCode},
		{"open", newPathOpenFileError(ctx, errors.New("boom"), "a"), PathOpenFileErrorCode},
		{"write", newPathWriteError(ctx, errors.New("boom"), "a"), PathWriteErrorCode},
	}
	for _, tc := range tests {
		if tc.xerr.Code() != tc.code {
			t.Fatalf("%s code=%d want %d", tc.name, tc.xerr.Code(), tc.code)
		}
		if reflect.DeepEqual(tc.xerr.Caller(), xerrors.CallerInfo{}) {
			t.Fatalf("%s expected caller info", tc.name)
		}
	}
}
