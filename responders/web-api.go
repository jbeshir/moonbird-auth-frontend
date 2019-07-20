package responders

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jbeshir/moonbird-auth-frontend/ctxlogrus"
	"net/http"
)

type WebApi struct {
	ExposeErrors bool
}

func (r *WebApi) OnContextError(w http.ResponseWriter, err error) {
	if r.ExposeErrors {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), 500)
	} else {
		http.Error(w, "Internal Server Error", 500)
	}
}

func (r *WebApi) OnError(ctx context.Context, w http.ResponseWriter, err error) {
	l := ctxlogrus.Get(ctx)
	l.Error(err)

	if r.ExposeErrors {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), 500)
	} else {
		http.Error(w, "Internal Server Error", 500)
	}
}

func (r *WebApi) OnSuccess(w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}
