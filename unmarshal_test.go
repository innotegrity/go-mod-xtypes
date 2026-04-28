package xtypes

import "testing"

type anyUnmarshalImpl struct{}

func (anyUnmarshalImpl) UnmarshalAny(data any) error {
	_ = data
	return nil
}

func TestAnyUnmarshalerCompileAssertion(t *testing.T) {
	t.Parallel()
	var _ AnyUnmarshaler = anyUnmarshalImpl{}
}

