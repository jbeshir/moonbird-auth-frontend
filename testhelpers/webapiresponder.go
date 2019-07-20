package testhelpers

import (
	"context"
	"net/http"
	"testing"
)

func NewWebApiResponder(t *testing.T) *WebApiResponder {
	return &WebApiResponder{
		OnContextErrorFunc: func(w http.ResponseWriter, err error) {
			t.Error("OnContextErrorFunc should not be called")
		},
		OnErrorFunc: func(ctx context.Context, w http.ResponseWriter, err error) {
			t.Error("OnErrorFunc should not be called")
		},
		OnSuccessFunc: func(w http.ResponseWriter, v interface{}) {
			t.Error("OnSuccessFunc should not be called")
		},
	}
}

type WebApiResponder struct {
	OnContextErrorFunc func(w http.ResponseWriter, err error)
	OnErrorFunc        func(ctx context.Context, w http.ResponseWriter, err error)
	OnSuccessFunc      func(w http.ResponseWriter, v interface{})
}

func (r *WebApiResponder) OnContextError(w http.ResponseWriter, err error) {
	r.OnContextErrorFunc(w, err)
}

func (r *WebApiResponder) OnError(ctx context.Context, w http.ResponseWriter, err error) {
	r.OnErrorFunc(ctx, w, err)
}

func (r *WebApiResponder) OnSuccess(w http.ResponseWriter, v interface{}) {
	r.OnSuccessFunc(w, v)
}
