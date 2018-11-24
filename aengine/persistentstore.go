package aengine

import (
	"context"
	"encoding/json"
	"google.golang.org/appengine/datastore"
)

type PersistentStore struct {
	Prefix string
}

type opaqueContent struct {
	Content []byte
}

func (ps *PersistentStore) GetOpaque(ctx context.Context, kind, key string, v interface{}) error {
	opaque := &opaqueContent{}

	k := datastore.NewKey(ctx, kind, ps.Prefix+key, 0, nil)
	err := datastore.Get(ctx, k, opaque)
	if err != nil {
		return err
	}

	err = json.Unmarshal(opaque.Content, v)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PersistentStore) SetOpaque(ctx context.Context, kind, key string, v interface{}) error {
	opaque := &opaqueContent{}

	var err error
	opaque.Content, err = json.Marshal(v)
	if err != nil {
		return err
	}

	k := datastore.NewKey(ctx, kind, ps.Prefix+key, 0, nil)
	_, err = datastore.Put(ctx, k, opaque)
	return err
}
