package controllers

import (
	"context"
	"github.com/jbeshir/moonbird-auth-frontend/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AdminApiCreateToken struct {
	ProjectTokenLister ProjectTokenLister
}

type AdminApiCreateTokenInput struct {
	Project string
}

func (c *AdminApiCreateToken) HandleFunc(cm ContextMaker, resp WebApiResponder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := cm.MakeContext(r)
		if err != nil {
			resp.OnContextError(w, err)
			return
		}

		input := AdminApiCreateTokenInput{
			Project: r.FormValue("project"),
		}
		token, err := c.handle(ctx, input)
		if err != nil {
			resp.OnError(ctx, w, err)
		} else {
			resp.OnSuccess(w, token)
		}
	}
}

func (c *AdminApiCreateToken) handle(ctx context.Context, input AdminApiCreateTokenInput) (string, error) {
	ctx = ctxlogrus.WithFields(ctx, logrus.Fields{
		"controller": "AdminApiCreateToken",
	})

	token, err := c.ProjectTokenLister.CreateToken(ctx, input.Project)

	return token, errors.Wrap(err, "")
}
