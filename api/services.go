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
}

type TokenBiller interface {
	Bill(token string, url *url.URL) error
}
