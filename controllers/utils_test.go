package controllers

import (
	"context"
	"testing"
)

func newTestLimitedEndpointBiller(t *testing.T) *testLimitedEndpointBiller {
	return &testLimitedEndpointBiller{
		SetLimitFunc: func(ctx context.Context, token, endpoint string, limit int64) error {
			t.Error("SetLimit should not be called")
			return nil
		},
	}
}

type testLimitedEndpointBiller struct {
	SetLimitFunc func(ctx context.Context, token, endpoint string, limit int64) error
}

func (b *testLimitedEndpointBiller) SetLimit(ctx context.Context, token, endpoint string, limit int64) error {
	return b.SetLimitFunc(ctx, token, endpoint, limit)
}
