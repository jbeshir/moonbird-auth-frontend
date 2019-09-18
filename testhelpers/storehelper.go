package testhelpers

import (
	"context"
	"testing"
)

type StoreHelper struct {
	EnsureExistsFunc   func(ctx context.Context, kind, key string, transact bool) error
	EnsurePropertyFunc func(ctx context.Context, kind, key, name, value string, transact bool) error
}

func NewStoreHelper(t *testing.T) *StoreHelper {
	return &StoreHelper{
		EnsureExistsFunc: func(ctx context.Context, kind, key string, transact bool) error {
			t.Error("EnsureExists should not be called")
			return nil
		},
		EnsurePropertyFunc: func(ctx context.Context, kind, key, name, value string, transact bool) error {
			t.Error("EnsureProperty should not be called")
			return nil
		},
	}
}

func (h *StoreHelper) EnsureExists(ctx context.Context, kind, key string, transact bool) error {
	return h.EnsureExistsFunc(ctx, kind, key, transact)
}

func (h *StoreHelper) EnsureProperty(ctx context.Context, kind, key, name, value string, transact bool) error {
	return h.EnsurePropertyFunc(ctx, kind, key, name, value, transact)
}
