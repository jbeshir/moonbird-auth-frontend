package testhelpers

import (
	"context"
	"testing"
)

type PermissionChecker struct {
	CheckReadFunc  func(ctx context.Context, kind, key string) (bool, error)
	CheckWriteFunc func(ctx context.Context, kind, key string) (bool, error)
}

func NewPermissionChecker(t *testing.T) *PermissionChecker {
	return &PermissionChecker{
		CheckReadFunc: func(ctx context.Context, kind, key string) (b bool, e error) {
			t.Error("CheckRead should not be called")
			return false, nil
		},
		CheckWriteFunc: func(ctx context.Context, kind, key string) (b bool, e error) {
			t.Error("CheckWrite should not be called")
			return false, nil
		},
	}
}

func (pc *PermissionChecker) CheckRead(ctx context.Context, kind, key string) (b bool, e error) {
	return pc.CheckReadFunc(ctx, kind, key)
}

func (pc *PermissionChecker) CheckWrite(ctx context.Context, kind, key string) (b bool, e error) {
	return pc.CheckWriteFunc(ctx, kind, key)
}
