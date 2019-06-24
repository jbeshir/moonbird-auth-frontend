package api

import (
	"context"
	"github.com/jbeshir/moonbird-predictor-frontend/data"
	"strings"
)

type ProjectPermissionChecker struct {
	PersistentStore    PersistentStore
	UserService        UserService
	TokenAuthenticator *TokenAuthenticator
}

func (pc *ProjectPermissionChecker) CheckRead(ctx context.Context, kind, key string) (bool, error) {
	project := strings.Split(key, "/")[0]

	if pc.UserService != nil {
		user := pc.UserService.ContextUser(ctx)
		if user != "" {
			_, err := pc.PersistentStore.Get(ctx, "ProjectAuth", project+"/user/"+user, nil)
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
			_, err := pc.PersistentStore.Get(ctx, "ProjectAuth", project+"/token/"+token, nil)
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
