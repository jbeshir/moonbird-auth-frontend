package storeutil

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-auth-frontend/data"
	"github.com/jbeshir/moonbird-auth-frontend/testhelpers"
	"reflect"
	"strings"
	"testing"
)

func TestHelper_EnsureProperty_GetErr(t *testing.T) {
	t.Parallel()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)

	inTransaction := false
	ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
		inTransaction = true
		defer func() {
			inTransaction = false
		}()
		return f(ctx)
	}

	expectedErr := errors.New("bluh")
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true

		if !inTransaction {
			t.Error("Expected Get to called in transaction, was not called in transaction")
		}

		expectedKind := "bluh"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bar/baz"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		return nil, expectedErr
	}

	h := &Helper{
		Store: ps,
	}

	err := h.EnsureProperty(context.Background(), "bluh", "bar/baz", "a", "b")
	if !getCalled {
		t.Error("Expected EnsureProperty to call Get, not called")
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected EnsureProperty to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestHelper_EnsureProperty_SetErr(t *testing.T) {
	t.Parallel()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)

	inTransaction := false
	ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
		inTransaction = true
		defer func() {
			inTransaction = false
		}()
		return f(ctx)
	}

	expectedErr := errors.New("bluh")
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true

		if !inTransaction {
			t.Error("Expected Get to called in transaction, was not called in transaction")
		}

		expectedKind := "bluh"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bar/baz"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		return []data.Property{}, nil
	}

	setCalled := false
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		setCalled = true

		if !inTransaction {
			t.Error("Expected Get to called in transaction, was not called in transaction")
		}

		expectedKind := "bluh"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bar/baz"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		return expectedErr
	}

	h := &Helper{
		Store: ps,
	}

	err := h.EnsureProperty(context.Background(), "bluh", "bar/baz", "a", "b")
	if !getCalled {
		t.Error("Expected EnsureProperty to call Get, not called")
	}
	if !setCalled {
		t.Error("Expected EnsureProperty to call Set, not called")
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected EnsureProperty to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestHelper_EnsureProperty_PropertyAlreadySet(t *testing.T) {
	t.Parallel()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)

	inTransaction := false
	ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
		inTransaction = true
		defer func() {
			inTransaction = false
		}()
		return f(ctx)
	}

	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true

		if !inTransaction {
			t.Error("Expected Get to called in transaction, was not called in transaction")
		}

		expectedKind := "bluh"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bar/baz"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		return []data.Property{
			{
				Name:  "a",
				Value: "b",
			},
		}, nil
	}

	h := &Helper{
		Store: ps,
	}

	err := h.EnsureProperty(context.Background(), "bluh", "bar/baz", "a", "b")
	if !getCalled {
		t.Error("Expected EnsureProperty to call Get, not called")
	}
	if err != nil {
		t.Errorf("Expected EnsureProperty to return nil error, got '%s'", err)
	}
}

func TestHelper_EnsureProperty(t *testing.T) {
	testCases := []struct {
		Label              string
		Properties         []data.Property
		ExpectedProperties []data.Property
	}{
		{
			Label: "PropertyAbsent",
			Properties: []data.Property{
				{
					Name:  "c",
					Value: "d",
				},
			},
			ExpectedProperties: []data.Property{
				{
					Name:  "c",
					Value: "d",
				},
				{
					Name:  "a",
					Value: "b",
				},
			},
		},
		{
			Label: "PropertyWrong",
			Properties: []data.Property{
				{
					Name:  "a",
					Value: "c",
				},
			},
			ExpectedProperties: []data.Property{
				{
					Name:  "a",
					Value: "b",
				},
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.Label, func(t *testing.T) {
			t.Parallel()

			getCalled := false
			setCalled := false
			ps := testhelpers.NewPersistentStore(t)

			inTransaction := false
			ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
				inTransaction = true
				defer func() {
					inTransaction = false
				}()
				return f(ctx)
			}

			ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
				getCalled = true

				if !inTransaction {
					t.Error("Expected Get to called in transaction, was not called in transaction")
				}

				expectedKind := "bluh"
				if kind != expectedKind {
					t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
				}

				expectedKey := "bar/baz"
				if key != expectedKey {
					t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
				}

				return testCase.Properties, nil
			}

			ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
				setCalled = true

				if !inTransaction {
					t.Error("Expected Get to called in transaction, was not called in transaction")
				}

				expectedKind := "bluh"
				if kind != expectedKind {
					t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
				}

				expectedKey := "bar/baz"
				if key != expectedKey {
					t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
				}

				if !reflect.DeepEqual(properties, testCase.ExpectedProperties) {
					t.Error("Properties did not match expected properties")
				}

				return nil
			}

			h := &Helper{
				Store: ps,
			}

			err := h.EnsureProperty(context.Background(), "bluh", "bar/baz", "a", "b")
			if !getCalled {
				t.Error("Expected EnsureProperty to call Get, not called")
			}
			if !setCalled {
				t.Error("Expected EnsureProperty to call Set, not called")
			}
			if err != nil {
				t.Errorf("Expected EnsureProperty to return nil error, got '%s'", err)
			}
		})
	}
}

func TestHelper_EnsureProperty_EntityAbsent(t *testing.T) {
	t.Parallel()

	getCalled := false
	setCalled := false
	ps := testhelpers.NewPersistentStore(t)

	inTransaction := false
	ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
		inTransaction = true
		defer func() {
			inTransaction = false
		}()
		return f(ctx)
	}

	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true

		if !inTransaction {
			t.Error("Expected Get to called in transaction, was not called in transaction")
		}

		expectedKind := "bluh"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bar/baz"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		return nil, data.ErrNoSuchEntity
	}

	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		setCalled = true

		if !inTransaction {
			t.Error("Expected Get to called in transaction, was not called in transaction")
		}

		expectedKind := "bluh"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bar/baz"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		expectedProperties := []data.Property{
			{
				Name:  "a",
				Value: "b",
			},
		}
		if !reflect.DeepEqual(properties, expectedProperties) {
			t.Error("Properties did not match expected properties")
		}

		return nil
	}

	h := &Helper{
		Store: ps,
	}

	err := h.EnsureProperty(context.Background(), "bluh", "bar/baz", "a", "b")
	if !getCalled {
		t.Error("Expected EnsureProperty to call Get, not called")
	}
	if !setCalled {
		t.Error("Expected EnsureProperty to call Set, not called")
	}
	if err != nil {
		t.Errorf("Expected EnsureProperty to return nil error, got '%s'", err)
	}
}
