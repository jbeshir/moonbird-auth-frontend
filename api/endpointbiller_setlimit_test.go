package api

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-auth-frontend/data"
	"github.com/jbeshir/moonbird-auth-frontend/testhelpers"
	"reflect"
	"strings"
	"testing"
)

func TestEndpointBiller_SetLimit(t *testing.T) {
	t.Parallel()

	ps := testhelpers.NewPersistentStore(t)

	setCallCount := 0
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) (e error) {
		setCallCount++

		expectedKind := "TokenLimit"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bluh/bar"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKind, kind)
		}

		reflect.DeepEqual(properties, []data.Property{
			{
				Name:  "Limit",
				Value: int64(86400),
			},
		})

		return nil
	}

	b := &EndpointBiller{
		PersistentStore: ps,
	}
	err := b.SetLimit(context.Background(), "bluh", "bar", 86400)
	if err != nil {
		t.Errorf("Expected nil err from SetLimit, got '%s'", err)
	}
	if setCallCount != 1 {
		t.Errorf("Expected Set to be called %d times, was called %d times", 1, setCallCount)
	}
}

func TestEndpointBiller_SetLimit_Err(t *testing.T) {
	t.Parallel()

	ps := testhelpers.NewPersistentStore(t)

	setCallCount := 0
	expectedErr := errors.New("bluh")
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) (e error) {
		setCallCount++
		return expectedErr
	}

	b := &EndpointBiller{
		PersistentStore: ps,
	}
	err := b.SetLimit(context.Background(), "bluh", "bar", 86400)
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected err '%s' from SetLimit, got '%s'", expectedErr, err)
	}
	if setCallCount != 1 {
		t.Errorf("Expected Set to be called %d times, was called %d times", 1, setCallCount)
	}
}
