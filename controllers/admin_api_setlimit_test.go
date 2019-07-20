package controllers

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-auth-frontend/testhelpers"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestExamplesUpdate_HandleFunc_Success(t *testing.T) {
	t.Parallel()

	var createdContext context.Context

	expectedToken := "bluh"
	expectedEndpoint := "bar"
	expectedLimit := int64(86400)

	calledSetLimit := false
	b := newTestLimitedEndpointBiller(t)
	b.SetLimitFunc = func(ctx context.Context, token, endpoint string, limit int64) error {
		calledSetLimit = true

		if token != expectedToken {
			t.Errorf("Expected token '%s', got token '%s'", expectedToken, token)
		}
		if endpoint != expectedEndpoint {
			t.Errorf("Expected endpoint '%s', got endpoint '%s'", expectedEndpoint, endpoint)
		}
		if limit != expectedLimit {
			t.Errorf("Expected limit %d, got limit %d", expectedLimit, limit)
		}
		return nil
	}

	calledOnSuccess := false
	r := testhelpers.NewWebApiResponder(t)
	r.OnSuccessFunc = func(w http.ResponseWriter, v interface{}) {
		calledOnSuccess = true
	}

	cm := testhelpers.NewContextMaker(t)
	cm.MakeContextFunc = func(r *http.Request) (i context.Context, e error) {
		createdContext = context.Background()
		return createdContext, nil
	}

	c := &AdminApiSetLimit{
		Biller: b,
	}
	handler := c.HandleFunc(cm, r)
	handler(nil, &http.Request{
		Form: url.Values{
			"token": []string{expectedToken},
			"endpoint": []string{expectedEndpoint},
			"limit": []string{strconv.FormatInt(expectedLimit, 10)},
		},
	})

	if !calledSetLimit {
		t.Error("Expected SetLimit to be called, was not called")
	}
	if !calledOnSuccess {
		t.Error("Expected responder's OnSuccess method to be called, was not called")
	}
}

func TestExamplesUpdate_HandleFunc_SetLimitErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("bluh")
	b := newTestLimitedEndpointBiller(t)
	b.SetLimitFunc = func(ctx context.Context, token, endpoint string, limit int64) error {
		return expectedErr
	}

	calledOnError := false
	r := testhelpers.NewWebApiResponder(t)
	r.OnErrorFunc = func(ctx context.Context, w http.ResponseWriter, err error) {
		calledOnError = true
	}

	cm := testhelpers.NewContextMaker(t)
	cm.MakeContextFunc = func(r *http.Request) (i context.Context, e error) {
		return context.Background(), nil
	}

	c := &AdminApiSetLimit{
		Biller: b,
	}
	handler := c.HandleFunc(cm, r)
	handler(nil, &http.Request{
		Form: url.Values{
			"token": []string{"bluh"},
			"endpoint": []string{"bar"},
			"limit": []string{strconv.FormatInt(1, 10)},
		},
	})

	if !calledOnError {
		t.Error("Expected responder's OnError method to be called, was not called")
	}
}