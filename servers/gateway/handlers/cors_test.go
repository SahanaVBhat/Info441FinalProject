package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/info441-au20/assignments-melodc/servers/gateway/models/users"
	"github.com/info441-au20/assignments-melodc/servers/gateway/sessions"
)

func TestCorsMiddlewareHandler(t *testing.T) {
	pass1 := true
	json := "{test:test}"

	sessionStore := &sessions.RedisStore{}
	userStore := &users.FakeMySQLStore{}
	signingkey := "test key"
	context := &HandlerContext{signingkey, sessionStore, userStore}

	req, err := http.NewRequest("POST", "/v1/sessions", strings.NewReader((json)))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(context.SessionsHandler)
	corsHandler := &CorsMiddleware{handler}
	resp := httptest.NewRecorder()
	corsHandler.ServeHTTP(resp, req)

	expectedcors := "*"
	ctype1 := resp.Header().Get("Access-Control-Allow-Origin")
	if ctype1 != expectedcors {
		t.Errorf("No `Access-Control-Allow-Origin` header found in the response: must be there start with `%s`", expectedcors)
		pass1 = false
	}

	expectedmethods := "GET, PUT, POST, PATCH, DELETE"
	ctype2 := resp.Header().Get("Access-Control-Allow-Methods")
	if ctype2 != expectedmethods {
		t.Errorf("No `Access-Control-Allow-Methods` header found in the response: must be there start with `%s`", expectedmethods)
		pass1 = false
	}

	expectedheaders := "Content-Type, Authorization"
	ctype3 := resp.Header().Get("Access-Control-Allow-Headers")
	if ctype3 != expectedheaders {
		t.Errorf("No `Access-Control-Allow-Headers` header found in the response: must be there start with `%s`", expectedheaders)
		pass1 = false
	}

	expectedMaxAge := "600"
	ctype4 := resp.Header().Get("Access-Control-Max-Age")
	if ctype4 != expectedMaxAge {
		t.Errorf("No `Access-Control-Max-Age` header found in the response: must be there start with `%s`", expectedMaxAge)
		pass1 = false
	}

	if pass1 {
		t.Log("TestCorsMiddlewareHandler passed")
	}

}
