package controllers

import (
	"context"
	"github.com/jbeshir/moonbird-auth-frontend/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type AdminApiSetLimit struct {
	Biller LimitedEndpointBiller
}

type AdminApiSetLimitInput struct {
	Token    string
	Endpoint string
	Limit    int64
}

func (c *AdminApiSetLimit) HandleFunc(cm ContextMaker, resp WebApiResponder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := cm.MakeContext(r)
		if err != nil {
			resp.OnContextError(w, err)
			return
		}

		limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 64)
		if err != nil {
			resp.OnError(ctx, w, err)
			return
		}

		input := AdminApiSetLimitInput{
			Token: r.FormValue("token"),
			Endpoint: r.FormValue("endpoint"),
			Limit: limit,
		}
		err = c.handle(ctx, input)
		if err != nil {
			resp.OnError(ctx, w, err)
		} else {
			resp.OnSuccess(w, true)
		}
	}
}

func (c *AdminApiSetLimit) handle(ctx context.Context, input AdminApiSetLimitInput) error {
	ctx = ctxlogrus.WithFields(ctx, logrus.Fields{
		"controller": "AdminApiSetLimit",
	})

	err := c.Biller.SetLimit(ctx, input.Token, input.Endpoint, input.Limit)
	return errors.Wrap(err, "")
}
