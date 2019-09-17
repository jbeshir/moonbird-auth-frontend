package storeutil

import (
	"context"
	"github.com/jbeshir/moonbird-auth-frontend/data"
)

type PersistentStore interface {
	Get(ctx context.Context, kind, key string, v interface{}) ([]data.Property, error)
	Set(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error
	Transact(ctx context.Context, f func(ctx context.Context) error) error
}
