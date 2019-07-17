package api

import (
	"context"
	"github.com/jbeshir/moonbird-auth-frontend/data"
	"net/http"
	"net/url"
)

type ContextMaker interface {
	MakeContext(r *http.Request) (context.Context, error)
}

type UserService interface {
	ContextUser(ctx context.Context) string
}

type PersistentStore interface {
	Get(ctx context.Context, kind, key string, v interface{}) ([]data.Property, error)
	Set(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error
	Transact(ctx context.Context, f func(ctx context.Context) error) error
}

type TokenBiller interface {
	Bill(ctx context.Context, token string, url *url.URL) error
}
