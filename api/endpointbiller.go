package api

import (
	"context"
	"github.com/jbeshir/moonbird-auth-frontend/data"
	"net/url"
)

type EndpointBiller struct {
	PersistentStore PersistentStore
	UrlEndpoints    map[string]string
}

type tokenUsage struct {
	Count int64
}

func (b *EndpointBiller) Bill(ctx context.Context, token string, url *url.URL) error {
	endpoint := b.UrlEndpoints[url.Path]
	if endpoint == "" {
		return data.ErrOutOfCredit
	}

	properties, err := b.PersistentStore.Get(ctx, "TokenLimit", tokenEndpointKey(token, endpoint), nil)
	if err != nil {
		if err == data.ErrNoSuchEntity {
			return data.ErrOutOfCredit
		}
		return err
	}

	var limit int64
	for _, v := range properties {
		if v.Name == "Limit" {
			limit = v.Value.(int64)
		}
	}
	if limit == 0 {
		return data.ErrOutOfCredit
	}

	estUsage, err := b.estimateUsage(ctx, token, endpoint)
	if err != nil {
		return err
	}
	if estUsage >= limit {
		return data.ErrOutOfCredit
	}

	return b.incrementUsage(ctx, token, endpoint)
}

func (b *EndpointBiller) SetLimit(ctx context.Context, token, endpoint string, limit int64) error {
	key := tokenEndpointKey(token, endpoint)
	properties := []data.Property{
		{
			Name:  "Limit",
			Value: limit,
		},
	}
	return b.PersistentStore.Set(ctx, "TokenLimit", key, properties, nil)
}

// Permitted to be moderately out of date for performance.
func (b *EndpointBiller) estimateUsage(ctx context.Context, token string, endpoint string) (int64, error) {
	var usage tokenUsage
	key := tokenEndpointKey(token, endpoint) + "/1"
	_, err := b.PersistentStore.Get(ctx, "TokenUsage", key, &usage)
	if err != nil {
		if err == data.ErrNoSuchEntity {
			return 0, nil
		}
		return 0, err
	}
	return usage.Count, nil
}

func (b *EndpointBiller) incrementUsage(ctx context.Context, token string, endpoint string) error {
	return b.PersistentStore.Transact(ctx, func(ctx context.Context) error {
		var usage tokenUsage
		key := tokenEndpointKey(token, endpoint) + "/1"
		_, err := b.PersistentStore.Get(ctx, "TokenUsage", key, &usage)
		if err != nil {
			if err == data.ErrNoSuchEntity {
				usage.Count = 0
			} else {
				return err
			}
		}

		usage.Count++

		return b.PersistentStore.Set(ctx, "TokenUsage", key, nil, &usage)
	})
}

func tokenEndpointKey(token, endpoint string) string {
	return token + "/" + endpoint
}
