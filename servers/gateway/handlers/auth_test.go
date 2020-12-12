package handlers

import (
	//"encoding/json"
	//"io"
	//"io/ioutil"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/info441-au20/assignments-melodc/servers/gateway/models/users"
	"github.com/info441-au20/assignments-melodc/servers/gateway/sessions"

	//"reflect"
	//"testing"
	"time"
	//"regexp"
)

func TestUsersHandler(t *testing.T) {
	cases := []struct {
		name         string
		method       string
		apiCall      string
		jsonStr      string
		ctype        string
		expectedResp int
		expectError  bool
	}{
		{
			"User authenticated correctly",
			"POST",
			"/v1/users",
			`{"email":"test@test.com","password":"password123","passwordConf":"password123","userName":"username",
			"firstName":"firstname","lastName":"lastname"}`,
			"application/json",
			http.StatusCreated,
			false,
		},
		{
			"Wrong method",
			"GET",
			"/v1/users",
			"{}",
			"application/json",
			http.StatusMethodNotAllowed,
			true,
		},
		{
			"Non-JSON ctype",
			"POST",
			"/v1/users",
			"{{}",
			"",
			http.StatusUnsupportedMediaType,
			true,
		},
		{
			"Bad JSON formatted type",
			"POST",
			"/v1/users",
			"{{&!@[]}",
			"application/json",
			http.StatusBadRequest,
			true,
		},
		{
			"Unable to validate user",
			"POST",
			"/v1/users",
			`{"email":"test@test@test.com","password":"bad","passwordConf":"bad","userName":"user",
			"firstName":"firstname","lastName":"lastname"}`,
			"application/json",
			http.StatusBadRequest,
			true,
		},
		{
			"Unable to convert user",
			"POST",
			"/v1/users",
			`{"wrong":"wrong","wrong":"wrong","wrong":"wrong"}`,
			"application/json",
			http.StatusBadRequest,
			true,
		},
	}
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	userStore := &users.FakeMySQLStore{}
	signingkey := "test key"
	context := &HandlerContext{signingkey, sessionStore, userStore}

	for _, c := range cases {
		resp := httptest.NewRecorder()
		jsonStr := []byte(c.jsonStr)
		req, err := http.NewRequest(c.method, c.apiCall, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		handler := http.HandlerFunc(context.UsersHandler)
		req.Header.Set("Content-Type", c.ctype)

		if !c.expectError {
			handler.ServeHTTP(resp, req)

			if resp.Code != c.expectedResp {
				t.Errorf("incorrect response status code: expected %d but got %d", c.expectedResp, resp.Code)
			}
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if expectedctype != ctype {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		} else {
			handler.ServeHTTP(resp, req)

			if resp.Code != c.expectedResp {
				t.Errorf("expected error response status %d but got %d for case %v", c.expectedResp, resp.Code, c.name)
			}
		}

	}

}

func TestSpecificUserHandler(t *testing.T) {
	cases := []struct {
		name         string
		method       string
		apiCall      string
		jsonStr      string
		ctype        string
		expectedResp int
		signedIn     bool
		expectError  bool
	}{
		{
			"User authenticated correctly",
			"GET",
			"/v1/users/0",
			`{"email":"test@test.com","password":"password123","passwordConf":"password123","userName":"username",
			"firstName":"firstname","lastName":"lastname"}`,
			"application/json",
			http.StatusOK,
			true,
			false,
		},
		{
			"Wrong method",
			"POST",
			"/v1/users/0",
			`{"firstName":"newfirstname","lastName":"newlastname"}`,
			"application/json",
			http.StatusMethodNotAllowed,
			true,
			true,
		},
		{
			"Update user correctly",
			"PATCH",
			"/v1/users/0",
			`{"firstName":"newfirstname","lastName":"newlastname"}`,
			"application/json",
			http.StatusOK,
			true,
			false,
		},
		{
			"Update user incorrectly",
			"PATCH",
			"/v1/users/43928347",
			`{"firstName":"newfirstname","lastName":"newlastname"}`,
			"application/json",
			http.StatusForbidden,
			true,
			true,
		},
		{
			"Update user with bad JSON",
			"PATCH",
			"/v1/users/0",
			`{"{"}`,
			"",
			http.StatusUnsupportedMediaType,
			true,
			true,
		},
		{
			"Update user with wrong format",
			"PATCH",
			"/v1/users/0",
			`{"wrong":"newfirstname"`,
			"application/json",
			http.StatusUnsupportedMediaType,
			true,
			true,
		},
		{
			"User not signed in",
			"GET",
			"/v1/users/0",
			`{"email":"test@test.com","password":"password123","passwordConf":"password123","userName":"username",
			"firstName":"firstname","lastName":"lastname"}`,
			"application/json",
			http.StatusUnauthorized,
			false,
			true,
		},
	}
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	userStore := &users.FakeMySQLStore{}
	signingkey := "test key"
	context := &HandlerContext{signingkey, sessionStore, userStore}
	usr, _ := userStore.Insert(&users.User{0,
		"test@test.com",
		[]byte("password123"),
		"username",
		"firstname",
		"lastname",
		"photourl"})
	for _, c := range cases {
		resp := httptest.NewRecorder()
		token := "abc123"

		if c.signedIn {
			response, errBeginSession := sessions.BeginSession(signingkey, sessionStore, &SessionState{time.Now(), usr}, resp)
			if errBeginSession != nil {
				t.Errorf(errBeginSession.Error())
			}
			token = string(response)

		}
		jsonStr := []byte(c.jsonStr)

		req, err := http.NewRequest(c.method, c.apiCall, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		if !c.expectError {
			req.Header.Set("Content-Type", c.ctype)
			req.Header.Set("Authorization", "Bearer "+token)
			handler := http.HandlerFunc(context.SpecificUserHandler)
			handler.ServeHTTP(resp, req)

			if resp.Code != c.expectedResp {
				t.Errorf("incorrect response status code: expected %d but got %d for case %v", c.expectedResp, resp.Code, c.name)
			}
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if expectedctype != ctype {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		} else {
			jsonStr := []byte(c.jsonStr)
			req, err := http.NewRequest(c.method, c.apiCall, bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer "+token)
			handler := http.HandlerFunc(context.SpecificUserHandler)
			handler.ServeHTTP(resp, req)

			if resp.Code != c.expectedResp {
				t.Errorf("expected error response status %d but got %d for case %v", c.expectedResp, resp.Code, c.name)
			}
		}

	}
}

func TestSessionsHandler(t *testing.T) {
	cases := []struct {
		name         string
		method       string
		apiCall      string
		jsonStr      string
		ctype        string
		expectedResp int
		expectError  bool
	}{
		{
			"Session begins successfully",
			"POST",
			"/v1/sessions",
			`{"email":"test@test.com","password":""}`,
			"application/json",
			http.StatusCreated,
			false,
		},
		{
			"Wrong method",
			"GET",
			"/v1/sessions",
			"{}",
			"application/json",
			http.StatusMethodNotAllowed,
			true,
		},
		{
			"Non-JSON ctype",
			"POST",
			"/v1/sessions",
			"{{}",
			"",
			http.StatusUnsupportedMediaType,
			true,
		},
		{
			"Bad JSON formatted type",
			"POST",
			"/v1/sessions",
			"{{&!@[]}",
			"application/json",
			http.StatusBadRequest,
			true,
		},
		{
			"Invalid email credentials",
			"POST",
			"/v1/sessions",
			`{"email":"bad@bad.com","password":"bad"}`,
			"application/json",
			http.StatusUnauthorized,
			true,
		},
		{
			"Unable to authenticate user",
			"POST",
			"/v1/sessions",
			`{"wrong":"wrong","wrong":"wrong","wrong":"wrong"}`,
			"application/json",
			http.StatusUnauthorized,
			true,
		},
	}
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	userStore := &users.FakeMySQLStore{}
	signingkey := "test key"
	context := &HandlerContext{signingkey, sessionStore, userStore}
	usr, _ := userStore.Insert(&users.User{0,
		"test@test.com",
		[]byte("password123"),
		"username",
		"firstname",
		"lastname",
		"photourl"})

	for _, c := range cases {
		resp := httptest.NewRecorder()
		response, errBeginSession := sessions.BeginSession(signingkey, sessionStore, &SessionState{time.Now(), usr}, resp)
		if errBeginSession != nil {
			t.Errorf(errBeginSession.Error())
		}
		token := string(response)

		jsonStr := []byte(c.jsonStr)
		req, err := http.NewRequest(c.method, c.apiCall, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		handler := http.HandlerFunc(context.SessionsHandler)
		req.Header.Set("Content-Type", c.ctype)

		if !c.expectError {
			handler.ServeHTTP(resp, req)

			if resp.Code != c.expectedResp {
				t.Errorf("incorrect response status code: expected %d but got %d for case %v", c.expectedResp, resp.Code, c.name)
			}
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if expectedctype != ctype {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		} else {
			handler.ServeHTTP(resp, req)

			if resp.Code != c.expectedResp {
				t.Errorf("expected error response status %d but got %d for case %v", c.expectedResp, resp.Code, c.name)
			}
		}
	}
}

func TestSpecificSessionsHandler(t *testing.T) {
	cases := []struct {
		name         string
		method       string
		apiCall      string
		expectedResp int
		signedIn     bool
		expectError  bool
	}{
		{
			"Session ended successfully",
			"DELETE",
			"/v1/sessions/mine",
			http.StatusOK,
			true,
			false,
		},
		{
			"Wrong request path",
			"DELETE",
			"/v1/sessions/wrong",
			http.StatusForbidden,
			true,
			true,
		},
		{
			"Wrong method",
			"POST",
			"/v1/sessions/wrong",
			http.StatusMethodNotAllowed,
			true,
			true,
		},
		{
			"Unable to end session",
			"DELETE",
			"/v1/sessions/mine",
			http.StatusBadRequest,
			false,
			true,
		},
	}
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	userStore := &users.FakeMySQLStore{}
	signingkey := "test key"
	context := &HandlerContext{signingkey, sessionStore, userStore}
	usr := &users.User{0,
		"test@test.com",
		[]byte("password123"),
		"username",
		"firstname",
		"lastname",
		"photourl"}
	for _, c := range cases {

		resp := httptest.NewRecorder()
		token := "abc123"

		if c.signedIn {
			response, errBeginSession := sessions.BeginSession(signingkey, sessionStore, &SessionState{time.Now(), usr}, resp)
			if errBeginSession != nil {
				t.Errorf(errBeginSession.Error())
			}
			token = string(response)

		}
		req, err := http.NewRequest(c.method, c.apiCall, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		if err != nil {
			t.Fatal(err)
		}
		handler := http.HandlerFunc(context.SpecificSessionHandler)

		if !c.expectError {
			handler.ServeHTTP(resp, req)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			bodyString := string(bodyBytes)
			if "Signed out" != bodyString {
				t.Errorf("Not signed out successfully: %v", bodyString)
			}
		} else {
			handler.ServeHTTP(resp, req)

			if resp.Code != c.expectedResp {
				t.Errorf("expected error response status %d but got %d for case %v", c.expectedResp, resp.Code, c.name)
			}
		}
	}

}
