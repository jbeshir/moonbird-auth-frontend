package controllers

import (
	"context"
	"net/http"
)

type ContextMaker interface {
	MakeContext(r *http.Request) (context.Context, error)
}

type LimitedEndpointBiller interface {
	SetLimit(ctx context.Context, token, endpoint string, limit int64) error
}

type ProjectTokenLister interface {
	CreateToken(ctx context.Context, project string) (string, error)
}

type WebApiResponder interface {
	OnContextError(w http.ResponseWriter, err error)
	OnError(ctx context.Context, w http.ResponseWriter, err error)
	OnSuccess(w http.ResponseWriter, v interface{})
}
