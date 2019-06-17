package api

import (
	"context"
	"net/http"
)

type ContextMaker interface {
	MakeContext(r *http.Request) (context.Context, error)
}
