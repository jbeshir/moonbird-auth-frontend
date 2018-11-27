package aengine

import (
	"context"
	"encoding/json"
	"google.golang.org/appengine/datastore"
)

type PersistentStore struct {
	Prefix string
}

func (ps *PersistentStore) GetOpaque(ctx context.Context, kind, key string, v interface{}) error {

	opaque := &opaqueContent{}
	k := ps.makeKey(ctx, kind, key)
	err := datastore.Get(ctx, k, opaque)
	if err != nil {
		return err
	}

	err = opaque.Unmarshal(v)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PersistentStore) SetOpaque(ctx context.Context, kind, key string, v interface{}) error {

	opaque := &opaqueContent{}
	err := opaque.Marshal(v)
	if err != nil {
		return err
	}

	k := ps.makeKey(ctx, kind, key)
	_, err = datastore.Put(ctx, k, opaque)
	return err
}

func (ps *PersistentStore) makeKey(ctx context.Context, kind, key string) *datastore.Key {
	return datastore.NewKey(ctx, kind, ps.Prefix+key, 0, nil)
}

type opaqueContent struct {
	Content []byte
}

func (o *opaqueContent) Marshal(v interface{}) (err error) {
	o.Content, err = json.Marshal(v)
	return
}

func (o *opaqueContent) Unmarshal(v interface{}) (err error) {
	err = json.Unmarshal(o.Content, v)
	return
}
