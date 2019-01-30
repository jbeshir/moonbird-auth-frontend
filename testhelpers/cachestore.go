package testhelpers

import (
	"context"
	"testing"
)

type CacheStore struct {
	GetFunc    func(ctx context.Context, key string, v interface{}) error
	SetFunc    func(ctx context.Context, key string, v interface{}) error
	DeleteFunc func(ctx context.Context, key string) error
}

func NewCacheStore(t *testing.T) *CacheStore {
	return &CacheStore{
		GetFunc: func(ctx context.Context, key string, v interface{}) error {
			t.Error("Get should not be called")
			return nil
		},
		SetFunc: func(ctx context.Context, key string, v interface{}) error {
			t.Error("Set should not be called")
			return nil
		},
		DeleteFunc: func(ctx context.Context, key string) error {
			t.Error("Delete should not be called")
			return nil
		},
	}
}

func (cs *CacheStore) Get(ctx context.Context, key string, v interface{}) error {
	return cs.GetFunc(ctx, key, v)
}

func (cs *CacheStore) Set(ctx context.Context, key string, v interface{}) error {
	return cs.SetFunc(ctx, key, v)
}

func (cs *CacheStore) Delete(ctx context.Context, key string) error {
	return cs.DeleteFunc(ctx, key)
}
