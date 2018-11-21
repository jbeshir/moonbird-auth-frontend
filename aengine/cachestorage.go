package aengine

import (
	"bytes"
	"context"
	"encoding/binary"
	"google.golang.org/appengine/memcache"
)

type CacheStorage struct {
	Prefix string
	Codec  memcache.Codec
}

func (cs *CacheStorage) Get(ctx context.Context, key string, v interface{}) error {
	_, err := cs.Codec.Get(ctx, key, v)
	return err
}

func (cs *CacheStorage) Set(ctx context.Context, key string, v interface{}) error {
	cacheItem := &memcache.Item{
		Key:    cs.Prefix + key,
		Object: v,
	}
	return cs.Codec.Set(ctx, cacheItem)
}

var BinaryMemcacheCodec = memcache.Codec{
	Marshal:   binaryMarshal,
	Unmarshal: binaryUnmarshal,
}

// Can only marshal fixed-size data as defined by the encoding/binary package.
func binaryMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Can only unmarshal fixed-size data as defined by the encoding/binary package.
func binaryUnmarshal(data []byte, v interface{}) error {
	return binary.Read(bytes.NewReader(data), binary.BigEndian, v)
}
