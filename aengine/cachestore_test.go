package aengine

import (
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/memcache"
	"reflect"
	"testing"
)

func TestCacheStore_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
		return
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	cacheItem := &memcache.Item{
		Key:    "FooBar",
		Object: "bluh",
	}
	err = memcache.JSON.Set(ctx, cacheItem)
	if err != nil {
		t.Fatalf("Unexpected error adding to memcache: %s", err)
	}

	cs := &CacheStore{
		Prefix: "Foo",
		Codec:  memcache.JSON,
	}
	err = cs.Delete(ctx, "Bar")
	if err != nil {
		t.Errorf("Unexpected error from Delete: %s", err)
	}

	var data string
	_, err = memcache.JSON.Get(ctx, "FooBar", &data)
	if err != memcache.ErrCacheMiss {
		if err == nil {
			t.Errorf("Found memcache data still present, expected deleted")
		} else {
			t.Errorf("Expected cache miss error, got: %s", err)
		}
	}
}

func TestCacheStore_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
		return
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	cacheItem := &memcache.Item{
		Key:    "FooBar",
		Object: "bluh",
	}
	err = memcache.JSON.Set(ctx, cacheItem)
	if err != nil {
		t.Fatalf("Unexpected error adding to memcache: %s", err)
	}

	cs := &CacheStore{
		Prefix: "Foo",
		Codec:  memcache.JSON,
	}

	var data string
	err = cs.Get(ctx, "Bar", &data)
	if err != nil {
		t.Errorf("Unexpected error from Get: %s", err)
	}
	if data != "bluh" {
		t.Errorf("Data read from memcache was incorrect; expected %s, was %s", "bluh", data)
	}
}

func TestCacheStore_Get_CacheMiss(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
		return
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	cs := &CacheStore{
		Prefix: "Foo",
		Codec:  memcache.JSON,
	}

	var data string
	err = cs.Get(ctx, "Bar", &data)
	if err != memcache.ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss from Get, got: %s", err)
	}
}

func TestCacheStore_Set(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
		return
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	cs := &CacheStore{
		Prefix: "Foo",
		Codec:  memcache.JSON,
	}

	setStr := "bluh"
	err = cs.Set(ctx, "Bar", &setStr)
	if err != nil {
		t.Errorf("Unexpected error from Set: %s", err)
	}

	var data string
	_, err = memcache.JSON.Get(ctx, "FooBar", &data)
	if err != nil {
		t.Errorf("Error reading data written to memcache: %s", err)
	}
	if data != "bluh" {
		t.Errorf("Data written to memcache was incorrect; expected %s, was %s", "bluh", data)
	}
}

func TestBinaryMemcacheCodec_Marshal(t *testing.T) {
	key, err := BinaryMemcacheCodec.Marshal(&[2]float64{0.4, 0.1})
	if err != nil {
		t.Errorf("Unexpected error from Marshal: %s", err)
	}
	expectedKey := []byte{
		0x3F, 0xD9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A, // Big-endian 0.4
		0x3F, 0xB9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A, // Big-endian 0.1
	}
	if !reflect.DeepEqual(key, expectedKey) {
		t.Errorf("Incorrect generated marshalled value; expected %x, was %x", expectedKey, key)
	}
}

func TestBinaryMemcacheCodec_Marshal_Illegal(t *testing.T) {
	data := "bluh"
	_, err := BinaryMemcacheCodec.Marshal(&data)
	if err == nil {
		t.Errorf("Expected error from Marshal, got nil error")
	}
}

func TestBinaryMemcacheCodec_Unmarshal(t *testing.T) {
	input := []byte{
		0x3F, 0xD9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A, // Big-endian 0.4
		0x3F, 0xB9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A, // Big-endian 0.1
	}
	var data [2]float64
	err := BinaryMemcacheCodec.Unmarshal(input, &data)
	if err != nil {
		t.Errorf("Unexpected error from Unmarshal: %s", err)
	}

	expectedData := [2]float64{0.4, 0.1}
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Incorrect generated marshalled value; expected %x, was %x", expectedData, data)
	}
}

func TestBinaryMemcacheCodec_Unmarshal_Illegal(t *testing.T) {
	input := []byte{
		0x3F, 0xD9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A, // Big-endian 0.4
		0x3F, 0xB9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A, // Big-endian 0.1
	}
	var data string
	err := BinaryMemcacheCodec.Unmarshal(input, &data)
	if err == nil {
		t.Errorf("Expected error from Unmarshal, got nil error")
	}
}
