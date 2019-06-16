package aengine

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-predictor-frontend/data"
	"github.com/jbeshir/moonbird-predictor-frontend/testhelpers"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"math"
	"reflect"
	"strings"
	"testing"
)

func makeTestProperties() (properties []data.Property) {
	properties = append(properties, data.Property{
		Name:  "Foo1",
		Value: "Bar",
	})
	properties = append(properties, data.Property{
		Name:  "Foo2",
		Value: int64(7),
	})
	properties = append(properties, data.Property{
		Name:  "Foo3",
		Value: true,
	})
	properties = append(properties, data.Property{
		Name:  "Foo4",
		Value: float64(0.3),
	})
	return
}

func TestPropertiesToAppEngine(t *testing.T) {
	from := makeTestProperties()
	to, err := propertiesToAppEngine(from)
	if err != nil {
		t.Errorf("Unexpected error converting properties to appengine format: %s", err)
	}
	if len(to) != len(from) {
		t.Errorf("Made a list of %d properties, expected %d", len(to), len(from))
	}

	for i := 0; i < len(to); i++ {
		if to[i].Name != from[i].Name {
			t.Errorf("Property %d had name '%s', expected '%s'", i, to[i].Name, from[i].Name)
		}
		if to[i].Value != from[i].Value {
			t.Errorf("Property %d had value '%v', expected '%v'", i, to[i].Value, from[i].Value)
		}
		if to[i].NoIndex {
			t.Errorf("Property %d had no index set, this is incorrect", i)
		}
		if to[i].Multiple {
			t.Errorf("Property %d had multiple set, this is incorrect", i)
		}

	}
}

func TestPropertiesToAppEngine_InvalidValue(t *testing.T) {
	from := makeTestProperties()
	from[1].Value = 7

	to, err := propertiesToAppEngine(from)
	if err == nil || !strings.Contains(err.Error(), "property 'Foo2' had invalid type: int") {
		t.Errorf("Did not receive expected error from conversion of properties to appengine format")
	}
	if len(to) != 0 {
		t.Errorf("Expected zero-length properties output with error, got non-zero-length properties list")
	}
}

func TestPropertiesToAppEngine_InvalidName(t *testing.T) {
	from := makeTestProperties()
	from[2].Name = "Content"

	to, err := propertiesToAppEngine(from)
	if err == nil || !strings.Contains(err.Error(), "property 'Content' had reserved name") {
		t.Errorf("Did not receive expected error from conversion of properties to appengine format")
	}
	if len(to) != 0 {
		t.Errorf("Expected zero-length properties output with error, got non-zero-length properties list")
	}
}

func TestPersistentStore_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)
	aeProperties = append(aeProperties, datastore.Property{
		Name:    "Content",
		Value:   []byte(`{"Foo":"Bar"}`),
		NoIndex: true,
	})

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	expectedData := map[string]interface{}{
		"Foo": "Bar",
	}
	var d map[string]interface{}
	properties, err := ps.Get(ctx, "Baz", "Bar", &d)
	if err != nil {
		t.Errorf("Unexpected error from Get: %s", err)
	}
	if !reflect.DeepEqual(properties, expectedProperties) {
		t.Errorf("Unmarshalled properties did not equal expected properties")
	}
	if !reflect.DeepEqual(d, expectedData) {
		t.Errorf("Unmarshalled d did not equal expected d")
	}
}

func TestPersistentStore_Get_NoContent(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	properties, err := ps.Get(ctx, "Baz", "Bar", nil)
	if err != nil {
		t.Errorf("Unexpected error from Get: %s", err)
	}
	if !reflect.DeepEqual(properties, expectedProperties) {
		t.Errorf("Unmarshalled properties did not equal expected properties")
	}
}

func TestPersistentStore_Get_ContentParamWithNoContent(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	var d map[string]interface{}
	_, err = ps.Get(ctx, "Baz", "Bar", &d)
	if err == nil {
		t.Errorf("Expected error from Get, got nil error")
	}
}

func TestPersistentStore_Get_ContentWithNoContentParam(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)
	aeProperties = append(aeProperties, datastore.Property{
		Name:    "Content",
		Value:   []byte(`{"Foo":"Bar"}`),
		NoIndex: true,
	})

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	_, err = ps.Get(ctx, "Baz", "Bar", nil)
	if err == nil {
		t.Errorf("Expected error from Get, got nil error")
	}
}

func TestPersistentStore_Get_ContentNotBytes(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)
	aeProperties = append(aeProperties, datastore.Property{
		Name:    "Content",
		Value:   true,
		NoIndex: true,
	})

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	var d map[string]interface{}
	_, err = ps.Get(ctx, "Baz", "Bar", &d)
	if err == nil {
		t.Errorf("Expected error from Get, got nil error")
	}
}

func TestPersistentStore_Get_ContentNotJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)
	aeProperties = append(aeProperties, datastore.Property{
		Name:    "Content",
		Value:   []byte(`bluh`),
		NoIndex: true,
	})

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	var d map[string]interface{}
	_, err = ps.Get(ctx, "Baz", "Bar", &d)
	if err == nil {
		t.Errorf("Expected error from Get, got nil error")
	}
}

func TestPersistentStore_Get_NoEntity(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	var d map[string]interface{}
	_, err = ps.Get(ctx, "Baz", "Bar", &d)
	if err != data.ErrNoSuchEntity {
		t.Errorf("Expected error '%s' from Get, got '%s'", data.ErrNoSuchEntity, err)
	}
}

func TestPersistentStore_Get_Permission(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	pc := testhelpers.NewPermissionChecker(t)
	pc.CheckReadFunc = func(ctx context.Context, kind, key string) (b bool, e error) {
		return true, nil
	}
	ps := &PersistentStore{
		Prefix:            "Foo",
		PermissionChecker: pc,
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	_, err = ps.Get(ctx, "Baz", "Bar", nil)
	if err != nil {
		t.Errorf("Expected nil error from Get, got '%s'", err)
	}
}

func TestPersistentStore_Get_NoPermission(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	pc := testhelpers.NewPermissionChecker(t)
	pc.CheckReadFunc = func(ctx context.Context, kind, key string) (b bool, e error) {
		return false, nil
	}
	ps := &PersistentStore{
		Prefix:            "Foo",
		PermissionChecker: pc,
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	_, err = ps.Get(ctx, "Baz", "Bar", nil)
	if err != data.ErrNoSuchEntity {
		t.Errorf("Expected error '%s' from Get, got '%s'", data.ErrNoSuchEntity, err)
	}
}

func TestPersistentStore_Get_PermissionError(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	expectedError := errors.New("bluh")
	pc := testhelpers.NewPermissionChecker(t)
	pc.CheckReadFunc = func(ctx context.Context, kind, key string) (b bool, e error) {
		return false, expectedError
	}
	ps := &PersistentStore{
		Prefix:            "Foo",
		PermissionChecker: pc,
	}

	expectedProperties := makeTestProperties()
	aeProperties, _ := propertiesToAppEngine(expectedProperties)

	k := ps.makeKey(ctx, "Baz", "Bar")
	_, err = datastore.Put(ctx, k, &aeProperties)
	if err != nil {
		t.Fatalf("Unexpected error writing data to datastore: %s", err)
	}

	_, err = ps.Get(ctx, "Baz", "Bar", nil)
	if err != expectedError {
		t.Errorf("Expected error '%s' from Get, got '%s'", expectedError, err)
	}
}

func TestPersistentStore_Set(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	err = ps.Set(ctx, "Baz", "Bar", makeTestProperties(), &map[string]interface{}{
		"Foo": "Bar",
	})
	if err != nil {
		t.Errorf("Unexpected error from Set: %s", err)
	}

	k := ps.makeKey(ctx, "Baz", "Bar")

	expectedProperties := makeTestProperties()
	expectedAEProperties, _ := propertiesToAppEngine(expectedProperties)
	expectedAEProperties = append(expectedAEProperties, datastore.Property{
		Name:    "Content",
		Value:   []byte(`{"Foo":"Bar"}`),
		NoIndex: true,
	})

	var aeProperties datastore.PropertyList
	err = datastore.Get(ctx, k, &aeProperties)
	if err != nil {
		t.Errorf("Unexpected error reading data from datastore: %s", err)
	}

	if !reflect.DeepEqual(aeProperties, expectedAEProperties) {
		t.Errorf("Set entity did not match expected data")
	}
}

func TestPersistentStore_Set_NoContent(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	err = ps.Set(ctx, "Baz", "Bar", makeTestProperties(), nil)
	if err != nil {
		t.Errorf("Unexpected error from Set: %s", err)
	}

	k := ps.makeKey(ctx, "Baz", "Bar")

	expectedProperties := makeTestProperties()
	expectedAEProperties, _ := propertiesToAppEngine(expectedProperties)

	var aeProperties datastore.PropertyList
	err = datastore.Get(ctx, k, &aeProperties)
	if err != nil {
		t.Errorf("Unexpected error reading data from datastore: %s", err)
	}

	if !reflect.DeepEqual(aeProperties, expectedAEProperties) {
		t.Errorf("Set entity did not match expected data")
	}
}

func TestPersistentStore_Set_InvalidProperty(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	err = ps.Set(ctx, "Baz", "Bar", []data.Property{
		{
			Name:  "Content",
			Value: true,
		},
	}, nil)
	if err == nil {
		t.Errorf("Expected error from Set, got nil error.")
	}
}

func TestPersistentStore_Set_InvalidContent(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	err = ps.Set(ctx, "Baz", "Bar", nil, &map[string]interface{}{
		"Foo": math.NaN(),
	})
	if err == nil {
		t.Errorf("Expected error from Set, got nil error.")
	}
}

func TestPersistentStore_Set_Permission(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	pc := testhelpers.NewPermissionChecker(t)
	pc.CheckWriteFunc = func(ctx context.Context, kind, key string) (b bool, e error) {
		return true, nil
	}
	ps := &PersistentStore{
		Prefix:            "Foo",
		PermissionChecker: pc,
	}

	err = ps.Set(ctx, "Baz", "Bar", makeTestProperties(), nil)
	if err != nil {
		t.Errorf("Expected nil error from Set, got '%s'", err)
	}
}

func TestPersistentStore_Set_NoPermission(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	pc := testhelpers.NewPermissionChecker(t)
	pc.CheckWriteFunc = func(ctx context.Context, kind, key string) (b bool, e error) {
		return false, nil
	}
	ps := &PersistentStore{
		Prefix:            "Foo",
		PermissionChecker: pc,
	}

	err = ps.Set(ctx, "Baz", "Bar", makeTestProperties(), nil)
	if err != data.ErrWriteAccessDenied {
		t.Errorf("Expected error '%s' from Set, got '%s'", data.ErrWriteAccessDenied, err)
	}
}

func TestPersistentStore_Set_PermissionError(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	expectedError := errors.New("bluh")
	pc := testhelpers.NewPermissionChecker(t)
	pc.CheckWriteFunc = func(ctx context.Context, kind, key string) (b bool, e error) {
		return false, expectedError
	}
	ps := &PersistentStore{
		Prefix:            "Foo",
		PermissionChecker: pc,
	}

	err = ps.Set(ctx, "Baz", "Bar", makeTestProperties(), nil)
	if err != expectedError {
		t.Errorf("Expected error '%s' from Set, got '%s'", expectedError, err)
	}
}

func TestPersistentStore_Transact(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	callCount := 0
	k := ps.makeKey(ctx, "Baz", "Bar")
	midTransCheck := make(chan bool)
	midTransCheckDone := make(chan bool)
	go func() {
		<-midTransCheck
		o := new(opaqueContent)
		err = datastore.Get(ctx, k, o)

		wantErr := datastore.ErrNoSuchEntity
		if err != wantErr {
			t.Errorf("Expected err %s, got %s", wantErr, err)
		}
		midTransCheckDone <- true
	}()

	err = ps.Transact(ctx, func(ctx context.Context) error {
		callCount++

		o := new(opaqueContent)
		o.Content = []byte("foo")
		_, _ = datastore.Put(ctx, k, o)
		midTransCheck <- true
		<-midTransCheckDone
		return nil
	})
	if err != nil {
		t.Errorf("Expected nil error from Transact, got %s", err)
	}
	wantCallCount := 1
	if callCount != wantCallCount {
		t.Errorf("Expected call count to be %d, was %d", wantCallCount, callCount)
	}
}

func TestPersistentStore_Transact_WithError(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	ps := &PersistentStore{
		Prefix: "Foo",
	}

	err = ps.Transact(ctx, func(ctx context.Context) error {
		return errors.New("bluh")
	})
	if err == nil {
		t.Errorf("Expected non-nil error from Transact, got nil error")
	}
}

func Test_PersistentStore_makeKey(t *testing.T) {
	if testing.Short() {
		t.Skip("AppEngine dev server testing is expensive")
	}

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
	t.Parallel()

	o := &opaqueContent{}

	d := map[string]interface{}{
		"Foo": "Bar",
	}

	err := o.Marshal(&d)
	if err != nil {
		t.Errorf("Unexpected error from Marshal: %s", err)
	}
	if string(o.Content) != `{"Foo":"Bar"}` {
		t.Errorf("Unexpected result of marshalling; expected `%s`, got `%s`", `{"Foo":"Bar"}`, o.Content)
	}
}

func TestOpaqueContent_Unmarshal(t *testing.T) {
	t.Parallel()

	o := &opaqueContent{
		Content: []byte(`{"Foo":"Bar"}`),
	}

	expectedData := map[string]interface{}{
		"Foo": "Bar",
	}
	var d map[string]interface{}
	err := o.Unmarshal(&d)
	if err != nil {
		t.Errorf("Unexpected error from Unmarshal: %s", err)
	}
	if !reflect.DeepEqual(d, expectedData) {
		t.Errorf("Unmarshalled data did not equal expected data")
	}
}
