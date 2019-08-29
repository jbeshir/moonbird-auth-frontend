package main

import (
	"github.com/jbeshir/moonbird-auth-frontend/aengine"
	"github.com/jbeshir/moonbird-auth-frontend/api"
	"github.com/jbeshir/moonbird-auth-frontend/controllers"
	"github.com/jbeshir/moonbird-auth-frontend/responders"
	"google.golang.org/appengine"
	"net/http"
)

func main() {
	authContext := &aengine.ContextMaker{
		Namespace: "moonbird-auth",
	}

	admApiSetLimit := &controllers.AdminApiSetLimit{
		Biller: &api.EndpointBiller{
			PersistentStore: &aengine.PersistentStore{},
		},
	}
	http.HandleFunc("/admin/api/set-limit", admApiSetLimit.HandleFunc(authContext, &responders.WebApi{}))

	admApiCreateToken := &controllers.AdminApiCreateToken{
		ProjectTokenLister: &api.ProjectPermissionChecker{
			PersistentStore: &aengine.PersistentStore{},
		},
	}
	http.HandleFunc("/admin/api/create-token", admApiCreateToken.HandleFunc(authContext, &responders.WebApi{}))

	appengine.Main()
}
