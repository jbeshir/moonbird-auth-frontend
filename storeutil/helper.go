package storeutil

import (
	"context"
	"github.com/jbeshir/moonbird-auth-frontend/data"
)

type Helper struct {
	Store PersistentStore
}

func (h *Helper) EnsureProperty(ctx context.Context, kind, key, name, value string, transact bool) error {
	return h.maybeTransact(ctx, func(ctx context.Context) error {
		properties, err := h.Store.Get(ctx, kind, key, nil)
		if err != nil && err != data.ErrNoSuchEntity {
			return err
		}

		found := false
		for i, _ := range properties {
			if properties[i].Name == name {
				if properties[i].Value == value {
					return nil
				}

				found = true
				properties[i].Value = value
				break
			}
		}
		if !found {
			properties = append(properties, data.Property{
				Name:  name,
				Value: value,
			})
		}

		return h.Store.Set(ctx, kind, key, properties, nil)
	}, transact)
}

func (h *Helper) EnsureExists(ctx context.Context, kind, key string, transact bool) error {
	return h.maybeTransact(ctx, func(ctx context.Context) error {
		_, err := h.Store.Get(ctx, kind, key, nil)
		if err != data.ErrNoSuchEntity {
			return err
		}

		return h.Store.Set(ctx, kind, key, nil, nil)
	}, transact)
}

func (h *Helper) maybeTransact(ctx context.Context, f func(context.Context) error, transact bool) error {
	if transact {
		return h.Store.Transact(ctx, f)
	} else {
		return f(ctx)
	}
}
