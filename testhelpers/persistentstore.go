package testhelpers

import (
	"context"
	"github.com/jbeshir/moonbird-auth-frontend/data"
	"testing"
)

type PersistentStore struct {
	GetFunc      func(ctx context.Context, kind, key string, v interface{}) ([]data.Property, error)
	SetFunc      func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error
	TransactFunc func(ctx context.Context, f func(ctx context.Context) error) error
}

func NewPersistentStore(t *testing.T) *PersistentStore {
	return &PersistentStore{
		GetFunc: func(ctx context.Context, kind, key string, v interface{}) ([]data.Property, error) {
			t.Error("Get should not be called")
			return nil, nil
		},
		SetFunc: func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
			t.Error("Set should not be called")
			return nil
		},
		TransactFunc: func(ctx context.Context, f func(ctx context.Context) error) error {
			t.Error("Transact should not be called")
			return nil
		},
	}
}

func (ps *PersistentStore) Get(ctx context.Context, kind, key string, v interface{}) ([]data.Property, error) {
	return ps.GetFunc(ctx, kind, key, v)
}

func (ps *PersistentStore) Set(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
	return ps.SetFunc(ctx, kind, key, properties, v)
}

func (ps *PersistentStore) Transact(ctx context.Context, f func(ctx context.Context) error) error {
	return ps.TransactFunc(ctx, f)
}
