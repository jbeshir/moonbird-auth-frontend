package api

import (
	"context"
	"errors"
	"github.com/jbeshir/moonbird-auth-frontend/data"
	"github.com/jbeshir/moonbird-auth-frontend/testhelpers"
	"strings"
	"testing"
)

func TestProjectPermissionsChecker_CheckRead_NoAuth(t *testing.T) {
	t.Parallel()

	pc := &ProjectPermissionChecker{}
	ok, err := pc.CheckRead(context.Background(), "foo", "bar")

	expectedOk := false
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}
}

func TestProjectPermissionsChecker_CheckRead_User_None(t *testing.T) {
	t.Parallel()

	expectedContext := context.Background()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true
		return nil, nil
	}

	contextUserCalled := false
	us := testhelpers.NewUserService(t)
	us.ContextUserFunc = func(ctx context.Context) string {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		contextUserCalled = true
		return ""
	}

	pc := &ProjectPermissionChecker{
		PersistentStore: ps,
		UserService:     us,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := false
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if getCalled {
		t.Error("Expected get function to not be called, was called")
	}

	if !contextUserCalled {
		t.Error("Expected context user function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckRead_User_Ok(t *testing.T) {
	t.Parallel()

	expectedContext := context.Background()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store get kind %v, was %v", expectedKind, kind)
		}

		expectedKey := "bar/user/superman"
		if key != expectedKey {
			t.Errorf("Expected store get key %v, was %v", expectedKey, key)
		}

		if v != nil {
			t.Error("Expected nil content, got non-nil content")
		}

		getCalled = true
		return nil, nil
	}

	contextUserCalled := false
	us := testhelpers.NewUserService(t)
	us.ContextUserFunc = func(ctx context.Context) string {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		contextUserCalled = true
		return "superman"
	}

	pc := &ProjectPermissionChecker{
		PersistentStore: ps,
		UserService:     us,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := true
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if !getCalled {
		t.Error("Expected get function to be called, was not called")
	}

	if !contextUserCalled {
		t.Error("Expected context user function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckRead_User_Escaped(t *testing.T) {
	t.Parallel()

	expectedContext := context.Background()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store get kind %v, was %v", expectedKind, kind)
		}

		expectedKey := "bar/user/%2Fname%2F"
		if key != expectedKey {
			t.Errorf("Expected store get key %v, was %v", expectedKey, key)
		}

		if v != nil {
			t.Error("Expected nil content, got non-nil content")
		}

		getCalled = true
		return nil, nil
	}

	contextUserCalled := false
	us := testhelpers.NewUserService(t)
	us.ContextUserFunc = func(ctx context.Context) string {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		contextUserCalled = true
		return "/name/"
	}

	pc := &ProjectPermissionChecker{
		PersistentStore: ps,
		UserService:     us,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := true
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if !getCalled {
		t.Error("Expected get function to be called, was not called")
	}

	if !contextUserCalled {
		t.Error("Expected context user function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckRead_User_NotOk(t *testing.T) {
	t.Parallel()

	expectedContext := context.Background()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store get kind %v, was %v", expectedKind, kind)
		}

		expectedKey := "bar/user/superman"
		if key != expectedKey {
			t.Errorf("Expected store get key %v, was %v", expectedKey, key)
		}

		if v != nil {
			t.Error("Expected nil content, got non-nil content")
		}

		getCalled = true
		return nil, data.ErrNoSuchEntity
	}

	contextUserCalled := false
	us := testhelpers.NewUserService(t)
	us.ContextUserFunc = func(ctx context.Context) string {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		contextUserCalled = true
		return "superman"
	}

	pc := &ProjectPermissionChecker{
		PersistentStore: ps,
		UserService:     us,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := false
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if !getCalled {
		t.Error("Expected get function to be called, was not called")
	}

	if !contextUserCalled {
		t.Error("Expected context user function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckRead_User_Err(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("bluh")
	expectedContext := context.Background()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store get kind %v, was %v", expectedKind, kind)
		}

		expectedKey := "bar/user/superman"
		if key != expectedKey {
			t.Errorf("Expected store get key %v, was %v", expectedKey, key)
		}

		if v != nil {
			t.Error("Expected nil content, got non-nil content")
		}

		getCalled = true
		return nil, expectedError
	}

	contextUserCalled := false
	us := testhelpers.NewUserService(t)
	us.ContextUserFunc = func(ctx context.Context) string {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		contextUserCalled = true
		return "superman"
	}

	pc := &ProjectPermissionChecker{
		PersistentStore: ps,
		UserService:     us,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := false
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != expectedError {
		t.Errorf("Expected err from check '%v', got '%v'", expectedError, err)
	}

	if !getCalled {
		t.Error("Expected get function to be called, was not called")
	}

	if !contextUserCalled {
		t.Error("Expected context user function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckRead_Token_Ok(t *testing.T) {
	t.Parallel()

	expectedToken := "bluh"
	expectedContext := context.WithValue(context.Background(), "apitoken", expectedToken)

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store get kind %v, was %v", expectedKind, kind)
		}

		expectedKey := "bar/token/bluh"
		if key != expectedKey {
			t.Errorf("Expected store get key %v, was %v", expectedKey, key)
		}

		if v != nil {
			t.Error("Expected nil content, got non-nil content")
		}

		getCalled = true
		return nil, nil
	}

	a := &TokenAuthenticator{}

	pc := &ProjectPermissionChecker{
		PersistentStore:    ps,
		TokenAuthenticator: a,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := true
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if !getCalled {
		t.Error("Expected get function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckRead_Token_Escaped(t *testing.T) {
	t.Parallel()

	expectedToken := "/trickytoken/"
	expectedContext := context.WithValue(context.Background(), "apitoken", expectedToken)

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store get kind %v, was %v", expectedKind, kind)
		}

		expectedKey := "bar/token/%2Ftrickytoken%2F"
		if key != expectedKey {
			t.Errorf("Expected store get key %v, was %v", expectedKey, key)
		}

		if v != nil {
			t.Error("Expected nil content, got non-nil content")
		}

		getCalled = true
		return nil, nil
	}

	a := &TokenAuthenticator{}

	pc := &ProjectPermissionChecker{
		PersistentStore:    ps,
		TokenAuthenticator: a,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := true
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if !getCalled {
		t.Error("Expected get function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckRead_Token_None(t *testing.T) {
	t.Parallel()

	expectedContext := context.Background()

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		getCalled = true
		return nil, nil
	}

	a := &TokenAuthenticator{}

	pc := &ProjectPermissionChecker{
		PersistentStore:    ps,
		TokenAuthenticator: a,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := false
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if getCalled {
		t.Error("Expected get function to not be called, was called")
	}
}

func TestProjectPermissionsChecker_CheckRead_Token_NotOk(t *testing.T) {
	t.Parallel()

	expectedToken := "bluh"
	expectedContext := context.WithValue(context.Background(), "apitoken", expectedToken)

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store get kind %v, was %v", expectedKind, kind)
		}

		expectedKey := "bar/token/bluh"
		if key != expectedKey {
			t.Errorf("Expected store get key %v, was %v", expectedKey, key)
		}

		if v != nil {
			t.Error("Expected nil content, got non-nil content")
		}

		getCalled = true
		return nil, data.ErrNoSuchEntity
	}

	a := &TokenAuthenticator{}

	pc := &ProjectPermissionChecker{
		PersistentStore:    ps,
		TokenAuthenticator: a,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := false
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if !getCalled {
		t.Error("Expected get function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckRead_Token_Err(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("blah")
	expectedToken := "bluh"
	expectedContext := context.WithValue(context.Background(), "apitoken", expectedToken)

	getCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.GetFunc = func(ctx context.Context, kind, key string, v interface{}) (properties []data.Property, e error) {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store get kind %v, was %v", expectedKind, kind)
		}

		expectedKey := "bar/token/bluh"
		if key != expectedKey {
			t.Errorf("Expected store get key %v, was %v", expectedKey, key)
		}

		if v != nil {
			t.Error("Expected nil content, got non-nil content")
		}

		getCalled = true
		return nil, expectedError
	}

	a := &TokenAuthenticator{}

	pc := &ProjectPermissionChecker{
		PersistentStore:    ps,
		TokenAuthenticator: a,
	}
	ok, err := pc.CheckRead(expectedContext, "foo", "bar/baz")

	expectedOk := false
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != expectedError {
		t.Errorf("Expected err from check '%v' got '%v'", expectedError, err)
	}

	if !getCalled {
		t.Error("Expected get function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CheckWrite_NoAuth(t *testing.T) {
	t.Parallel()

	pc := &ProjectPermissionChecker{}
	ok, err := pc.CheckWrite(context.Background(), "foo", "bar")

	expectedOk := false
	if ok != expectedOk {
		t.Errorf("Permission check response was expected to be %v, was %v", expectedOk, ok)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}
}

func TestProjectPermissionsChecker_CreateToken_Ok(t *testing.T) {
	t.Parallel()

	setToken := ""
	expectedContext := context.Background()

	setCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKind := "ProjectAuth"
		if kind != expectedKind {
			t.Errorf("Expected store set kind %v, was %v", expectedKind, kind)
		}

		expectedKeyPrefix := "bar/token/"
		if !strings.HasPrefix(key, expectedKeyPrefix) {
			t.Errorf("Expected store set key prefix %v, key was %v", expectedKeyPrefix, key)
		}

		setToken = key[len(expectedKeyPrefix):]
		if len(setToken) != 44 {
			t.Errorf("Expected store set key token to be random 40 hex chars, token was %v", setToken)
		}

		setCalled = true
		return nil
	}

	pc := &ProjectPermissionChecker{
		PersistentStore: ps,
	}
	token, err := pc.CreateToken(expectedContext, "bar")

	if token != setToken {
		t.Errorf("CreateToken did not return the same token '%s' it placed into the datastore, returned '%s'", setToken, token)
	}

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if !setCalled {
		t.Error("Expected set function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CreateToken_Escaped(t *testing.T) {
	t.Parallel()

	expectedContext := context.Background()

	setCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		if ctx != expectedContext {
			t.Error("Context was not expected context")
		}

		expectedKeyPrefix := "a%2Fb/token/"
		if !strings.HasPrefix(key, expectedKeyPrefix) {
			t.Errorf("Expected store set key prefix %v, key was %v", expectedKeyPrefix, key)
		}

		setCalled = true
		return nil
	}

	pc := &ProjectPermissionChecker{
		PersistentStore: ps,
	}
	_, err := pc.CreateToken(expectedContext, "a/b")

	if err != nil {
		t.Errorf("Unexpected non-nil err from check: %v", err)
	}

	if !setCalled {
		t.Error("Expected set function to be called, was not called")
	}
}

func TestProjectPermissionsChecker_CreateToken_Err(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("blah")
	expectedContext := context.Background()

	setCalled := false
	ps := testhelpers.NewPersistentStore(t)
	ps.SetFunc = func(ctx context.Context, kind, key string, properties []data.Property, v interface{}) error {
		setCalled = true
		return expectedError
	}

	pc := &ProjectPermissionChecker{
		PersistentStore: ps,
	}
	_, err := pc.CreateToken(expectedContext, "bar")

	if err != expectedError {
		t.Errorf("Expected err from CreateToken '%v' got '%v'", expectedError, err)
	}

	if !setCalled {
		t.Error("Expected set function to be called, was not called")
	}
}
