package aengine

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "")
	}

	return errors.Wrap(opaque.Unmarshal(v), "")
}

func (ps *PersistentStore) SetOpaque(ctx context.Context, kind, key string, v interface{}) error {

	opaque := &opaqueContent{}
	err := opaque.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "")
	}

	k := ps.makeKey(ctx, kind, key)
	_, err = datastore.Put(ctx, k, opaque)
	return errors.Wrap(err, "")
}

func (ps *PersistentStore) Transact(ctx context.Context, f func(ctx context.Context) error) error {
	return datastore.RunInTransaction(ctx, f, nil)
}

func (ps *PersistentStore) makeKey(ctx context.Context, kind, key string) *datastore.Key {
	return datastore.NewKey(ctx, kind, ps.Prefix+key, 0, nil)
}

type opaqueContent struct {
	Content []byte
}

func (o *opaqueContent) Marshal(v interface{}) error {
	var err error
	o.Content, err = json.Marshal(v)
	return errors.Wrap(err, "")
}

func (o *opaqueContent) Unmarshal(v interface{}) error {
	return errors.Wrap(json.Unmarshal(o.Content, v), "")
}
