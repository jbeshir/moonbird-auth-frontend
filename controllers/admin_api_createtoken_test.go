package controllers

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-auth-frontend/testhelpers"
	"net/http"
	"net/url"
	"testing"
)

func TestAdminApiCreateToken_HandleFunc_Success(t *testing.T) {
	t.Parallel()

	var createdContext context.Context

	expectedProject := "bar"

	expectedToken := "foo"
	calledCreateToken := false
	ptl := newTestProjectTokenLister(t)
	ptl.CreateTokenFunc = func(ctx context.Context, project string) (s string, e error) {
		calledCreateToken = true

		if project != expectedProject {
			t.Errorf("Expected project '%s', got project '%s'", expectedProject, project)
		}
		return expectedToken, nil
	}

	calledOnSuccess := false
	r := testhelpers.NewWebApiResponder(t)
	r.OnSuccessFunc = func(w http.ResponseWriter, v interface{}) {
		if v != expectedToken {
			t.Errorf("Expected token '%s', got token '%s'", expectedToken, v)
		}
		calledOnSuccess = true
	}

	cm := testhelpers.NewContextMaker(t)
	cm.MakeContextFunc = func(r *http.Request) (i context.Context, e error) {
		createdContext = context.Background()
		return createdContext, nil
	}

	c := &AdminApiCreateToken{
		ProjectTokenLister: ptl,
	}
	handler := c.HandleFunc(cm, r)
	handler(nil, &http.Request{
		Form: url.Values{
			"project": []string{expectedProject},
		},
	})

	if !calledCreateToken {
		t.Error("Expected CreateToken to be called, was not called")
	}
	if !calledOnSuccess {
		t.Error("Expected responder's OnSuccess method to be called, was not called")
	}
}

func TestAdminApiCreateToken_HandleFunc_CreateTokenErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("bluh")
	ptl := newTestProjectTokenLister(t)
	ptl.CreateTokenFunc = func(ctx context.Context, project string) (s string, e error) {
		return "", expectedErr
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

	c := &AdminApiCreateToken{
		ProjectTokenLister: ptl,
	}
	handler := c.HandleFunc(cm, r)
	handler(nil, &http.Request{
		Form: url.Values{
			"project": []string{"bar"},
		},
	})

	if !calledOnError {
		t.Error("Expected responder's OnError method to be called, was not called")
	}
}
