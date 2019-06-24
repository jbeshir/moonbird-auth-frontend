package api

import (
	"context"
	"github.com/jbeshir/moonbird-predictor-frontend/data"
	"net/http"
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
