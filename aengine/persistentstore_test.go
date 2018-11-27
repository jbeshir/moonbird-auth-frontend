package aengine

import (
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"math"
	"reflect"
	"testing"
)

func TestPersistentStore_GetOpaque(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	o := &opaqueContent{
		Content: []byte(`{"Foo":"Bar"}`),
	}

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, o)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	expectedData := map[string]interface{}{
		"Foo": "Bar",
	}
	var data map[string]interface{}
	err = ps.GetOpaque(ctx, "Baz", "Bar", &data)
	if err != nil {
		t.Errorf("Unexpected error from GetOpaque: %s", err)
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Unmarshalled data did not equal expected data")
	}
}

func TestPersistentStore_GetOpaque_NoEntity(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	var data map[string]interface{}
	err = ps.GetOpaque(ctx, "Baz", "Bar", &data)
	if err == nil {
		t.Errorf("Expected error from GetOpaque, got nil error")
	}
}

func TestPersistentStore_GetOpaque_InvalidEntity(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &struct {
		Foo int
	}{
		Foo: 1,
	})
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	var data map[string]interface{}
	err = ps.GetOpaque(ctx, "Baz", "Bar", &data)
	if err == nil {
		t.Errorf("Expected error from GetOpaque, got nil error")
	}
}

func TestPersistentStore_GetOpaque_InvalidEntityContent(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	o := &opaqueContent{
		Content: []byte(`nope`),
	}
	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, o)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	var data map[string]interface{}
	err = ps.GetOpaque(ctx, "Baz", "Bar", &data)
	if err == nil {
		t.Errorf("Expected error from GetOpaque, got nil error")
	}
}

func TestPersistentStore_SetOpaque(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	err = ps.SetOpaque(ctx, "Baz", "Bar", &map[string]interface{}{
		"Foo": "Bar",
	})
	if err != nil {
		t.Errorf("Unexpected error from SetOpaque: %s", err)
	}

	k := ps.makeKey(ctx, "Baz", "Bar")

	o := new(opaqueContent)
	err = datastore.Get(ctx, k, o)
	if err != nil {
		t.Errorf("Unexpected error reading data from datastore: %s", err)
	}

	if string(o.Content) != `{"Foo":"Bar"}` {
		t.Errorf("Data written to datastore did not match what was expected: expected `%s`, got `%s`", `{"Foo":"Bar"}`, string(o.Content))
	}
}

func TestPersistentStore_SetOpaque_Invalid(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	err = ps.SetOpaque(ctx, "Baz", "Bar", &map[string]interface{}{
		"Foo": math.NaN(),
	})
	if err == nil {
		t.Errorf("Expected error from SetOpaque, got nil error.")
	}
}

func Test_PersistentStore_makeKey(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	k := ps.makeKey(ctx, "Baz", "Bar")

	if k.Kind() != "Baz" {
		t.Errorf("Incorrect key kind, expected %s, was %s", "Baz", k.Kind())
	}
	if k.StringID() != "FooBar" {
		t.Errorf("Incorrect key ID, expected %s, was %s", "FooBar", k.StringID())
	}
}

func TestOpaqueContent_Marshal(t *testing.T) {
	o := &opaqueContent{}

	data := map[string]interface{}{
		"Foo": "Bar",
	}

	err := o.Marshal(&data)
	if err != nil {
		t.Errorf("Unexpected error from Marshal: %s", err)
	}
	if string(o.Content) != `{"Foo":"Bar"}` {
		t.Errorf("Unexpected result of marshalling; expected `%s`, got `%s`", `{"Foo":"Bar"}`, o.Content)
	}
}

func TestOpaqueContent_Unmarshal(t *testing.T) {
	o := &opaqueContent{
		Content: []byte(`{"Foo":"Bar"}`),
	}

	expectedData := map[string]interface{}{
		"Foo": "Bar",
	}
	var data map[string]interface{}
	err := o.Unmarshal(&data)
	if err != nil {
		t.Errorf("Unexpected error from Unmarshal: %s", err)
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Unmarshalled data did not equal expected data")
	}
}
