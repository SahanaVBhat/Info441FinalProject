package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, ErrNoSessionID
	}

	Sessionid, err1 := NewSessionID(signingKey)
	if err1 != nil {
		return Sessionid, err1
	}
	err2 := store.Save(Sessionid, sessionState)
	if err2 != nil {
		return Sessionid, err2
	}
	w.Header().Add(headerAuthorization, schemeBearer+Sessionid.String())

	return Sessionid, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	authToken := r.Header.Get(headerAuthorization)
	if len(authToken) == 0 {
		authToken = r.URL.Query().Get(paramAuthorization)
	}
	if len(authToken) == 0 {
		return InvalidSessionID, ErrNoSessionID
	}
	if strings.HasPrefix(authToken, schemeBearer) {
		authToken = strings.TrimPrefix(authToken, schemeBearer)
	} else {
		return InvalidSessionID, ErrInvalidScheme
	}
	sid, err := ValidateID(authToken, signingKey)
	return sid, err
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	err2 := store.Get(sid, sessionState)
	if err2 != nil {
		return InvalidSessionID, err2
	}
	return sid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	err2 := store.Delete(sid)
	if err2 != nil {
		return InvalidSessionID, err
	}

	return sid, nil
}
