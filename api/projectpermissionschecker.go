package api

import (
	"context"
	"github.com/jbeshir/moonbird-predictor-frontend/data"
	"net/url"
	"strings"
)

type ProjectPermissionChecker struct {
	PersistentStore    PersistentStore
	UserService        UserService
	TokenAuthenticator *TokenAuthenticator
}

func (pc *ProjectPermissionChecker) CheckRead(ctx context.Context, kind, key string) (bool, error) {
	escapedProject := strings.Split(key, "/")[0]

	if pc.UserService != nil {
		user := pc.UserService.ContextUser(ctx)
		if user != "" {
			_, err := pc.PersistentStore.Get(ctx, "ProjectAuth", escapedProject+"/user/"+url.PathEscape(user), nil)
			if err == nil {
				return true, nil
			}
			if err != data.ErrNoSuchEntity {
				return false, err
			}
		}
	}

	if pc.TokenAuthenticator != nil {
		token := pc.TokenAuthenticator.GetToken(ctx)
		if token != "" {
			_, err := pc.PersistentStore.Get(ctx, "ProjectAuth", escapedProject+"/token/"+url.PathEscape(token), nil)
			if err == nil {
				return true, nil
			}
			if err != data.ErrNoSuchEntity {
				return false, err
			}
		}
	}

	return false, nil
}

func (pc *ProjectPermissionChecker) CheckWrite(ctx context.Context, kind, key string) (bool, error) {
	return pc.CheckRead(ctx, kind, key)
}
