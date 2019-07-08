package aengine

import (
	"context"
	"github.com/jbeshir/moonbird-auth-frontend/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine"
	"io/ioutil"
	"net/http"
)

type ContextMaker struct {
	Namespace string
}

var aeLogger *logrus.Logger

func init() {
	aeLogger = logrus.New()
	aeLogger.Hooks.Add(&logHook{})
	aeLogger.Out = ioutil.Discard
}

func (cm *ContextMaker) MakeContext(r *http.Request) (context.Context, error) {
	ctx := appengine.NewContext(r)

	ctx = ctxlogrus.WithLogger(ctx, logrus.NewEntry(aeLogger))
	ctx = ctxlogrus.WithFields(ctx, logrus.Fields{"aengine-ctx": ctx})

	if cm.Namespace != "" {
		var err error
		ctx, err = appengine.Namespace(ctx, cm.Namespace)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		ctx = ctxlogrus.WithFields(ctx, logrus.Fields{"aengine-namespace": cm.Namespace})
	}

	return ctx, nil
}
