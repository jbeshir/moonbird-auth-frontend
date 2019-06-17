package testhelpers

import (
	"context"
	"net/http"
	"testing"
)

func NewContextMaker(t *testing.T) *ContextMaker {
	return &ContextMaker{
		MakeContextFunc: func(r *http.Request) (context.Context, error) {
			t.Error("MakeContextFunc should not be called")
			return nil, nil
		},
	}
}

type ContextMaker struct {
	MakeContextFunc func(r *http.Request) (context.Context, error)
}

func (cm *ContextMaker) MakeContext(r *http.Request) (context.Context, error) {
	return cm.MakeContextFunc(r)
}
