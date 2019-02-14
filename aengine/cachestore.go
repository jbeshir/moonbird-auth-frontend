package aengine

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/pkg/errors"
	"google.golang.org/appengine/memcache"
)

type CacheStore struct {
	Prefix string
	Codec  memcache.Codec
}

func (cs *CacheStore) Get(ctx context.Context, key string, v interface{}) error {
	_, err := cs.Codec.Get(ctx, cs.Prefix+key, v)
	return errors.Wrap(err, "")
}

func (cs *CacheStore) Set(ctx context.Context, key string, v interface{}) error {
	cacheItem := &memcache.Item{
		Key:    cs.Prefix + key,
		Object: v,
	}

	return errors.Wrap(cs.Codec.Set(ctx, cacheItem), "")
}

func (cs *CacheStore) Delete(ctx context.Context, key string) error {
	return memcache.Delete(ctx, cs.Prefix+key)
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
