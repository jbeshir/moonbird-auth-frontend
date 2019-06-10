package aengine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jbeshir/moonbird-predictor-frontend/ctxlogrus"
	"github.com/jbeshir/moonbird-predictor-frontend/data"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine/datastore"
)

type PersistentStore struct {
	Prefix string
}

func (ps *PersistentStore) Get(ctx context.Context, kind, key string, content interface{}) ([]data.Property, error) {
	l := ctxlogrus.Get(ctx)
	l.WithFields(logrus.Fields{"prefix": ps.Prefix, "kind": kind, "key": key}).Debug("datastore get")

	var aeProperties datastore.PropertyList
	k := ps.makeKey(ctx, kind, key)
	err := datastore.Get(ctx, k, &aeProperties)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	foundContent := false
	for i := len(aeProperties) - 1; i >= 0; i-- {
		if aeProperties[i].Name == "Content" {
			foundContent = true
			if content == nil {
				return nil, errors.New("entity contained content to deserialize, but content param was not set")
			}

			contentBytes, ok := aeProperties[i].Value.([]byte)
			if !ok {
				return nil, errors.New("entity contained content property with incorrect type")
			}
			err = json.Unmarshal(contentBytes, content)
			if err != nil {
				return nil, errors.Wrap(err, "unable to deserialize entity content")
			}

			// Splice this property out
			aeProperties = append(aeProperties[:i], aeProperties[i+1:]...)
			break
		}
	}
	if !foundContent && content != nil {
		return nil, errors.New("entity did not contain content to deserialize, but content param was set")
	}

	return propertiesFromAppEngine(aeProperties), nil
}

func (ps *PersistentStore) Set(ctx context.Context, kind, key string, properties []data.Property, content interface{}) error {
	l := ctxlogrus.Get(ctx)
	l.WithFields(logrus.Fields{"prefix": ps.Prefix, "kind": kind, "key": key}).Debug("datastore set")

	aeProperties, err := propertiesToAppEngine(properties)
	if err != nil {
		return err
	}
	if content != nil {
		opaque := &opaqueContent{}
		err := opaque.Marshal(content)
		if err != nil {
			return errors.Wrap(err, "")
		}

		aeProperties = append(aeProperties, datastore.Property{
			Name:    "Content",
			Value:   opaque.Content,
			NoIndex: true,
		})
	}

	k := ps.makeKey(ctx, kind, key)
	_, err = datastore.Put(ctx, k, &aeProperties)
	return errors.Wrap(err, "")
}

func (ps *PersistentStore) Transact(ctx context.Context, f func(ctx context.Context) error) error {
	l := ctxlogrus.Get(ctx)
	l.Debug("datastore transaction start")

	err := datastore.RunInTransaction(ctx, f, nil)

	l.Debug("datastore transaction end")

	return errors.Wrap(err, "")
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

func propertiesFromAppEngine(from datastore.PropertyList) (to []data.Property) {
	for _, v := range from {
		to = append(to, data.Property{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return
}

func propertiesToAppEngine(from []data.Property) (to datastore.PropertyList, err error) {
	for _, v := range from {
		switch v.Value.(type) {
		case int64:
		case bool:
		case string:
		case float64:
		default:
			return nil, fmt.Errorf("property '%s' had invalid type: %T", v.Name, v.Value)
		}

		if v.Name == "Content" {
			return nil, fmt.Errorf("property '%s' had reserved name", v.Name)
		}

		to = append(to, datastore.Property{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return
}
