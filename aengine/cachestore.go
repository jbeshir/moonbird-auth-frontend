package aengine

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/jbeshir/moonbird-auth-frontend/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine/memcache"
)

type CacheStore struct {
	Prefix string
	Codec  memcache.Codec
}

func (cs *CacheStore) Get(ctx context.Context, key string, v interface{}) error {
	l := ctxlogrus.Get(ctx)
	l.WithFields(logrus.Fields{"prefix": cs.Prefix, "key": key}).Debug("cache get")

	_, err := cs.Codec.Get(ctx, cs.Prefix+key, v)
	return errors.Wrap(err, "")
}

func (cs *CacheStore) Set(ctx context.Context, key string, v interface{}) error {
	l := ctxlogrus.Get(ctx)
	l.WithFields(logrus.Fields{"prefix": cs.Prefix, "key": key}).Debug("cache set")

	cacheItem := &memcache.Item{
		Key:    cs.Prefix + key,
		Object: v,
	}

	return errors.Wrap(cs.Codec.Set(ctx, cacheItem), "")
}

func (cs *CacheStore) Delete(ctx context.Context, key string) error {
	l := ctxlogrus.Get(ctx)
	l.WithFields(logrus.Fields{"prefix": cs.Prefix, "key": key}).Debug("cache delete")

	return memcache.Delete(ctx, cs.Prefix+key)
}

func (cs *CacheStore) Flush(ctx context.Context) error {
	l := ctxlogrus.Get(ctx)
	l.Debug("cache clear - full purge")

	return memcache.Flush(ctx)
}

// Can only marshal fixed-size data as defined by the encoding/binary package.
var BinaryMemcacheCodec = memcache.Codec{
	Marshal:   binaryMarshal,
	Unmarshal: binaryUnmarshal,
}

func binaryMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, v)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return buf.Bytes(), nil
}

func binaryUnmarshal(data []byte, v interface{}) error {
	return binary.Read(bytes.NewReader(data), binary.BigEndian, v)
}
