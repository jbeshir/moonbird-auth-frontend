package storeutil

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-auth-frontend/data"
	"github.com/jbeshir/moonbird-auth-frontend/testhelpers"
	"strings"
	"testing"
)

func TestHelper_EnsureExists_GetErr(t *testing.T) {
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

	err := h.EnsureExists(context.Background(), "bluh", "bar/baz", true)
	if !getCalled {
		t.Error("Expected EnsureExists to call Get, not called")
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected EnsureExists to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestHelper_EnsureExists_SetErr(t *testing.T) {
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

		return nil, data.ErrNoSuchEntity
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

	err := h.EnsureExists(context.Background(), "bluh", "bar/baz", true)
	if !getCalled {
		t.Error("Expected EnsureExists to call Get, not called")
	}
	if !setCalled {
		t.Error("Expected EnsureExists to call Set, not called")
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected EnsureExists to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestHelper_EnsureExists_AlreadyExists(t *testing.T) {
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

	err := h.EnsureExists(context.Background(), "bluh", "bar/baz", true)
	if !getCalled {
		t.Error("Expected EnsureExists to call Get, not called")
	}
	if err != nil {
		t.Errorf("Expected EnsureExists to return nil error, got '%s'", err)
	}
}

func TestHelper_EnsureExists_EntityAbsent(t *testing.T) {
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

		if properties != nil {
			t.Error("Expected empty property list")
		}

		return nil
	}

	h := &Helper{
		Store: ps,
	}

	err := h.EnsureExists(context.Background(), "bluh", "bar/baz", true)
	if !getCalled {
		t.Error("Expected EnsureExists to call Get, not called")
	}
	if !setCalled {
		t.Error("Expected EnsureExists to call Set, not called")
	}
	if err != nil {
		t.Errorf("Expected EnsureExists to return nil error, got '%s'", err)
	}
}

func TestHelper_EnsureExists_EntityAbsent_NoTransact(t *testing.T) {
	t.Parallel()

	getCalled := false
	setCalled := false
	ps := testhelpers.NewPersistentStore(t)

	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true

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

		expectedKind := "bluh"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bar/baz"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		if properties != nil {
			t.Error("Expected empty property list")
		}

		return nil
	}

	h := &Helper{
		Store: ps,
	}

	err := h.EnsureExists(context.Background(), "bluh", "bar/baz", false)
	if !getCalled {
		t.Error("Expected EnsureExists to call Get, not called")
	}
	if !setCalled {
		t.Error("Expected EnsureExists to call Set, not called")
	}
	if err != nil {
		t.Errorf("Expected EnsureExists to return nil error, got '%s'", err)
	}
}
