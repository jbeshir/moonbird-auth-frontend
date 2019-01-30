package testhelpers

import (
	"context"
	"testing"
)

type PersistentStore struct {
	GetOpaqueFunc func(ctx context.Context, kind, key string, v interface{}) error
	SetOpaqueFunc func(ctx context.Context, kind, key string, v interface{}) error
}

func NewPersistentStore(t *testing.T) *PersistentStore {
	return &PersistentStore{
		GetOpaqueFunc: func(ctx context.Context, kind, key string, v interface{}) error {
			t.Error("GetOpaque should not be called")
			return nil
		},
		SetOpaqueFunc: func(ctx context.Context, kind, key string, v interface{}) error {
			t.Error("SetOpaque should not be called")
			return nil
		},
	}
}

func (ps *PersistentStore) GetOpaque(ctx context.Context, kind, key string, v interface{}) error {
	return ps.GetOpaqueFunc(ctx, kind, key, v)
}

func (ps *PersistentStore) SetOpaque(ctx context.Context, kind, key string, v interface{}) error {
	return ps.SetOpaqueFunc(ctx, kind, key, v)
}
