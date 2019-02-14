package aengine

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine/log"
)

type logHook struct{}

func (lh *logHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
	}
}

func (lh *logHook) Fire(e *logrus.Entry) error {
	appCtx := e.Data["aengine-ctx"].(context.Context)

	// Make the entry's data separate from its template, so our change to exclude data doesn't permanently edit the map.
	data := make(logrus.Fields, len(e.Data)-1)
	for k, v := range e.Data {
		if k == "aengine-ctx" {
			continue
		}
		data[k] = v
	}
	e.Data = data

	msg, err := e.String()
	if err != nil {
		return errors.Wrap(err, "")
	}

	switch e.Level {
	case logrus.DebugLevel:
		log.Debugf(appCtx, "%s", msg)
	case logrus.InfoLevel:
		log.Infof(appCtx, "%s", msg)
	case logrus.WarnLevel:
		log.Warningf(appCtx, "%s", msg)
	case logrus.ErrorLevel:
		log.Errorf(appCtx, "%s", msg)
	default:
		log.Criticalf(appCtx, "Unknown level '%s' (%d): %s", e.Level, e.Level, msg)
	}
	return nil
}
