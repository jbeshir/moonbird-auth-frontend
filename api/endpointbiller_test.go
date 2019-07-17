package api

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-auth-frontend/data"
	"github.com/jbeshir/moonbird-auth-frontend/testhelpers"
	"net/url"
	"strings"
	"testing"
)

func TestEndpointBiller_Bill_NoMatch(t *testing.T) {
	b := &EndpointBiller{
		UrlEndpoints: map[string]string{},
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	expectedErr := data.ErrOutOfCredit
	err = b.Bill(context.Background(), "bluh", u)
	if err != expectedErr {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill_LimitGetErr(t *testing.T) {
	getCalled := false
	ps := testhelpers.NewPersistentStore(t)

	expectedErr := errors.New("bluh")
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true

		expectedKind := "TokenLimit"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bluh/bar"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		return nil, expectedErr
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if !getCalled {
		t.Error("Expected bill to call Get, not called")
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill_NoLimitEntity(t *testing.T) {
	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true

		expectedKind := "TokenLimit"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bluh/bar"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		return nil, data.ErrNoSuchEntity
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if !getCalled {
		t.Error("Expected bill to call Get, not called")
	}
	expectedErr := data.ErrOutOfCredit
	if err != expectedErr {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill_ZeroLimit(t *testing.T) {
	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true

		expectedKind := "TokenLimit"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bluh/bar"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		return []data.Property{
			{
				Name:  "Limit",
				Value: int64(0),
			},
		}, nil
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if !getCalled {
		t.Error("Expected bill to call Get, not called")
	}
	expectedErr := data.ErrOutOfCredit
	if err != expectedErr {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill_EstUsageCheckErr(t *testing.T) {
	getCallCount := 0
	expectedErr := errors.New("bluh")
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCallCount++
		switch getCallCount {
		case 1:
			expectedKind := "TokenLimit"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			return []data.Property{
				{
					Name:  "Limit",
					Value: int64(86400),
				},
			}, nil
		case 2:
			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			return nil, expectedErr
		}
		t.Error("Unexpected third call to Get")
		return nil, nil
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if getCallCount != 2 {
		t.Error("Expected bill to call Get twice, not called")
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill_LimitReached(t *testing.T) {
	getCallCount := 0
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCallCount++
		switch getCallCount {
		case 1:
			expectedKind := "TokenLimit"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			return []data.Property{
				{
					Name:  "Limit",
					Value: int64(86400),
				},
			}, nil
		case 2:
			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			usage, ok := v.(*tokenUsage)
			if !ok {
				t.Error("Expected token usage struct to unpack into")
			} else {
				usage.Count = 86400
			}

			return nil, nil
		}
		t.Error("Unexpected third call to Get")
		return nil, nil
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if getCallCount != 2 {
		t.Error("Expected bill to call Get twice, not called")
	}
	expectedErr := data.ErrOutOfCredit
	if err != expectedErr {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill_IncrementTransactErr(t *testing.T) {
	ps := testhelpers.NewPersistentStore(t)

	expectedErr := errors.New("bluh")
	transactionCallCount := 0
	inTransaction := false
	ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
		if transactionCallCount != 0 {
			t.Error("Expected Transact to only be called once")
		}
		return expectedErr
	}

	getCallCount := 0
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCallCount++
		switch getCallCount {
		case 1:
			if inTransaction {
				t.Errorf("Expected get call %d to be outside transaction", getCallCount)
			}

			expectedKind := "TokenLimit"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKind, kind)
			}

			return []data.Property{
				{
					Name:  "Limit",
					Value: int64(86400),
				},
			}, nil
		case 2:
			if inTransaction {
				t.Errorf("Expected get call %d to be outside transaction", getCallCount)
			}

			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			usage, ok := v.(*tokenUsage)
			if !ok {
				t.Error("Expected token usage struct to unpack into")
			} else {
				usage.Count = 86399
			}

			return nil, nil

		case 3:
			return nil, nil
		}
		t.Error("Unexpected third call to Get")
		return nil, nil
	}

	setCallCount := 0
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		setCallCount++
		return nil
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if getCallCount != 2 {
		t.Errorf("Expected bill to call Get %d times, called %d times", 3, getCallCount)
	}
	if setCallCount != 0 {
		t.Errorf("Expected bill to call Set %d times, called %d times", 1, setCallCount)
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill_IncrementGetErr(t *testing.T) {
	ps := testhelpers.NewPersistentStore(t)

	expectedErr := errors.New("bluh")
	transactionCallCount := 0
	inTransaction := false
	ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
		if transactionCallCount != 0 {
			t.Error("Expected Transact to only be called once")
		}
		transactionCallCount++

		inTransaction = true
		err := f(ctx)
		inTransaction = false
		return err
	}

	getCallCount := 0
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCallCount++
		switch getCallCount {
		case 1:
			if inTransaction {
				t.Errorf("Expected get call %d to be outside transaction", getCallCount)
			}

			expectedKind := "TokenLimit"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKind, kind)
			}

			return []data.Property{
				{
					Name:  "Limit",
					Value: int64(86400),
				},
			}, nil
		case 2:
			if inTransaction {
				t.Errorf("Expected get call %d to be outside transaction", getCallCount)
			}

			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			usage, ok := v.(*tokenUsage)
			if !ok {
				t.Error("Expected token usage struct to unpack into")
			} else {
				usage.Count = 86399
			}

			return nil, nil

		case 3:
			if !inTransaction {
				t.Errorf("Expected get call %d to be inside transaction", getCallCount)
			}

			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}
			return nil, expectedErr
		}
		t.Error("Unexpected third call to Get")
		return nil, nil
	}

	setCallCount := 0
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		setCallCount++
		return nil
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if getCallCount != 3 {
		t.Errorf("Expected bill to call Get %d times, called %d times", 3, getCallCount)
	}
	if setCallCount != 0 {
		t.Errorf("Expected bill to call Set %d times, called %d times", 1, setCallCount)
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill_IncrementSetErr(t *testing.T) {
	ps := testhelpers.NewPersistentStore(t)

	expectedErr := errors.New("bluh")
	transactionCallCount := 0
	inTransaction := false
	ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
		if transactionCallCount != 0 {
			t.Error("Expected Transact to only be called once")
		}
		transactionCallCount++

		inTransaction = true
		err := f(ctx)
		inTransaction = false
		return err
	}

	getCallCount := 0
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCallCount++
		switch getCallCount {
		case 1:
			if inTransaction {
				t.Errorf("Expected get call %d to be outside transaction", getCallCount)
			}

			expectedKind := "TokenLimit"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKind, kind)
			}

			return []data.Property{
				{
					Name:  "Limit",
					Value: int64(86400),
				},
			}, nil
		case 2:
			if inTransaction {
				t.Errorf("Expected get call %d to be outside transaction", getCallCount)
			}

			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			usage, ok := v.(*tokenUsage)
			if !ok {
				t.Error("Expected token usage struct to unpack into")
			} else {
				usage.Count = 86399
			}

			return nil, nil

		case 3:
			if !inTransaction {
				t.Errorf("Expected get call %d to be inside transaction", getCallCount)
			}

			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			usage, ok := v.(*tokenUsage)
			if !ok {
				t.Error("Expected token usage struct to unpack into")
			} else {
				usage.Count = 86401
			}

			return nil, nil
		}
		t.Error("Unexpected third call to Get")
		return nil, nil
	}

	setCallCount := 0
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		setCallCount++
		return expectedErr
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if getCallCount != 3 {
		t.Errorf("Expected bill to call Get %d times, called %d times", 3, getCallCount)
	}
	if setCallCount != 1 {
		t.Errorf("Expected bill to call Set %d times, called %d times", 1, setCallCount)
	}
	if err == nil || !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Expected bill to return error '%s', got '%s'", expectedErr, err)
	}
}

func TestEndpointBiller_Bill(t *testing.T) {
	ps := testhelpers.NewPersistentStore(t)

	transactionCallCount := 0
	inTransaction := false
	ps.TransactFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
		if transactionCallCount != 0 {
			t.Error("Expected Transact to only be called once")
		}
		transactionCallCount++

		inTransaction = true
		err := f(ctx)
		inTransaction = false
		return err
	}

	getCallCount := 0
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCallCount++
		switch getCallCount {
		case 1:
			if inTransaction {
				t.Errorf("Expected get call %d to be outside transaction", getCallCount)
			}

			expectedKind := "TokenLimit"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKind, kind)
			}

			return []data.Property{
				{
					Name:  "Limit",
					Value: int64(86400),
				},
			}, nil
		case 2:
			if inTransaction {
				t.Errorf("Expected get call %d to be outside transaction", getCallCount)
			}

			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			usage, ok := v.(*tokenUsage)
			if !ok {
				t.Error("Expected token usage struct to unpack into")
			} else {
				usage.Count = 86399
			}

			return nil, nil

		case 3:
			if !inTransaction {
				t.Errorf("Expected get call %d to be inside transaction", getCallCount)
			}

			expectedKind := "TokenUsage"
			if kind != expectedKind {
				t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
			}

			expectedKey := "bluh/bar/1"
			if key != expectedKey {
				t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
			}

			usage, ok := v.(*tokenUsage)
			if !ok {
				t.Error("Expected token usage struct to unpack into")
			} else {
				usage.Count = 86401
			}

			return nil, nil
		}
		t.Error("Unexpected third call to Get")
		return nil, nil
	}

	setCallCount := 0
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		setCallCount++

		expectedKind := "TokenUsage"
		if kind != expectedKind {
			t.Errorf("Expected kind '%s', got '%s'", expectedKind, kind)
		}

		expectedKey := "bluh/bar/1"
		if key != expectedKey {
			t.Errorf("Expected key '%s', got '%s'", expectedKey, key)
		}

		usage, ok := v.(*tokenUsage)
		if !ok {
			t.Error("Expected token usage struct to save")
		}

		expectedCount := int64(86402)
		if usage.Count != expectedCount {
			t.Errorf("Expected new usage count to be %d, was %d", expectedCount, usage.Count)
		}

		return nil
	}

	u, err := url.Parse("https://example.com/api/foo")
	if err != nil {
		t.Fatal("Failed to parse URL")
	}

	b := &EndpointBiller{
		PersistentStore: ps,
		UrlEndpoints: map[string]string{
			"/api/foo": "bar",
		},
	}
	err = b.Bill(context.Background(), "bluh", u)
	if getCallCount != 3 {
		t.Errorf("Expected bill to call Get %d times, called %d times", 3, getCallCount)
	}
	if setCallCount != 1 {
		t.Errorf("Expected bill to call Set %d times, called %d times", 1, setCallCount)
	}
	if err != nil {
		t.Errorf("Expected bill to return nil error, got '%s'", err)
	}
}
