package ctxlogrus

import "context"
import "github.com/sirupsen/logrus"

type key struct{}

func Get(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(key{})

	if logger == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return logger.(*logrus.Entry)
}

func WithLogger(ctx context.Context, l *logrus.Entry) context.Context {
	return context.WithValue(ctx, key{}, l)
}

func WithFields(ctx context.Context, fields logrus.Fields) context.Context {
	return WithLogger(ctx, Get(ctx).WithFields(fields))
}
