package aengine

import (
	"context"
	"google.golang.org/appengine"
	"net/http"
)

type ContextMaker struct {
	Namespace string
}

func (cm *ContextMaker) MakeContext(r *http.Request) (context.Context, error) {
	ctx := appengine.NewContext(r)
	if cm.Namespace != "" {
		return appengine.Namespace(ctx, cm.Namespace)
	} else {
		return ctx, nil
	}
}
