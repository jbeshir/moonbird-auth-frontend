package aengine

import (
	"context"
	"google.golang.org/appengine/user"
)

type UserService struct{}

func (us *UserService) ContextUser(ctx context.Context) string {
	u := user.Current(ctx)
	if u == nil {
		return ""
	} else {
		return u.ID
	}
}
