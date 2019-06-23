package testhelpers

import (
	"context"
	"testing"
)

func NewUserService(t *testing.T) *UserService {
	return &UserService{
		ContextUserFunc: func(ctx context.Context) string {
			t.Error("ContextUser should not be called")
			return ""
		},
	}
}

type UserService struct {
	ContextUserFunc func(ctx context.Context) string
}

func (us *UserService) ContextUser(ctx context.Context) string {
	return us.ContextUserFunc(ctx)
}
