package api

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-predictor-frontend/testhelpers"
	"net/http"
	"net/url"
	"testing"
)

func TestTokenAuthenticator_MakeContext(t *testing.T) {
	t.Parallel()

	expectedToken := "bluh"

	formValues := make(url.Values)
	formValues.Add("apitoken", expectedToken)
	r := &http.Request{Form: formValues}

	a := &TokenAuthenticator{}
	c, err := a.MakeContext(r)
	if err != nil {
		t.Errorf("Expected nil error, got '%s'", err)
	}

	token := c.Value("apitoken").(string)
	if token != expectedToken {
		t.Errorf("Expected token '%s', got '%s'", expectedToken, token)
	}
}

func TestTokenAuthenticator_MakeContext_NoToken(t *testing.T) {
	t.Parallel()

	r := &http.Request{}
	a := &TokenAuthenticator{}
	_, err := a.MakeContext(r)
	if err == nil {
		t.Error("Expected non-nil error, got nil error")
	}
}

func TestTokenAuthenticator_MakeContext_MultipleTokens(t *testing.T) {
	t.Parallel()

	formValues := make(url.Values)
	formValues.Add("apitoken", "bluh")
	formValues.Add("apitoken", "bluh")

	r := &http.Request{Form: formValues}
	a := &TokenAuthenticator{}
	_, err := a.MakeContext(r)
	if err == nil {
		t.Error("Expected non-nil error, got nil error")
	}
}

func TestTokenAuthenticator_MakeContext_Wrapped(t *testing.T) {
	t.Parallel()

	expectedToken := "bluh"

	formValues := make(url.Values)
	formValues.Add("apitoken", expectedToken)
	r := &http.Request{Form: formValues}

	cm := testhelpers.NewContextMaker(t)
	cm.MakeContextFunc = func(req *http.Request) (i context.Context, e error) {
		if req != r {
			t.Error("Request object was not expected request object")
		}
		return context.WithValue(context.Background(), "foo", "bar"), nil
	}

	a := &TokenAuthenticator{
		Wrapped: cm,
	}
	c, err := a.MakeContext(r)
	if err != nil {
		t.Errorf("Expected nil error, got '%s'", err)
	}

	if s, ok := c.Value("foo").(string); !ok || s != "bar" {
		t.Error("MakeContext did not preserve existing context values")
	}
}

func TestTokenAuthenticator_MakeContext_WrappedErr(t *testing.T) {
	t.Parallel()

	expectedToken := "bluh"

	formValues := make(url.Values)
	formValues.Add("apitoken", expectedToken)
	r := &http.Request{Form: formValues}

	cm := testhelpers.NewContextMaker(t)
	cm.MakeContextFunc = func(req *http.Request) (i context.Context, e error) {
		return nil, errors.New("bluh")
	}

	a := &TokenAuthenticator{
		Wrapped: cm,
	}
	_, err := a.MakeContext(r)
	if err == nil {
		t.Error("Expected non-nil error, got nil error")
	}
}

func TestTokenAuthenticator_GetToken(t *testing.T) {
	t.Parallel()

	expectedToken := "bluh"

	a := &TokenAuthenticator{}
	c := context.WithValue(context.Background(), "apitoken", expectedToken)
	token := a.GetToken(c)
	if token != expectedToken {
		t.Errorf("Expected token '%s', got '%s'", expectedToken, token)
	}
}

func TestTokenAuthenticator_GetToken_None(t *testing.T) {
	t.Parallel()

	expectedToken := ""

	a := &TokenAuthenticator{}
	token := a.GetToken(context.Background())
	if token != expectedToken {
		t.Errorf("Expected token '%s', got '%s'", expectedToken, token)
	}
}
