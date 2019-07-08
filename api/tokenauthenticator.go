package api

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-auth-frontend/ctxlogrus"
	"net/http"
)

type TokenAuthenticator struct {
	Wrapped ContextMaker
}

func (a *TokenAuthenticator) MakeContext(r *http.Request) (context.Context, error) {
	var wrappedCtx context.Context
	if a.Wrapped != nil {
		var err error
		wrappedCtx, err = a.Wrapped.MakeContext(r)
		if err != nil {
			return nil, err
		}
	} else {
		wrappedCtx = context.Background()
	}

	_ = r.ParseForm()
	if len(r.Form["apitoken"]) != 1 {
		ctxlogrus.Get(wrappedCtx).Error("expected exactly one api token for an API request")
		return nil, errors.New("expected exactly one api token for an API request")
	}
	token := r.Form["apitoken"][0]

	c := context.WithValue(wrappedCtx, "apitoken", token)
	return c, nil
}

func (a *TokenAuthenticator) GetToken(ctx context.Context) string {
	s, _ := ctx.Value("apitoken").(string)
	return s
}
