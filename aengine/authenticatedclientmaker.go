package aengine

import (
	"context"
	"golang.org/x/oauth2/google"
	"net/http"
)

type AuthenticatedClientMaker struct {
	Scope string
}

func (cm *AuthenticatedClientMaker) MakeClient(ctx context.Context) (*http.Client, error) {
	return google.DefaultClient(ctx, cm.Scope)
}
