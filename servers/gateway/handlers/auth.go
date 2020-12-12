package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/info441-au20/assignments-melodc/servers/gateway/models/users"
	"github.com/info441-au20/assignments-melodc/servers/gateway/sessions"
)

//handles requests for the "users" resource.
func (hc *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		ctype := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ctype, "application/json") {
			http.Error(w, "The request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}

		// extract response body JSON to New User object
		usr := &users.NewUser{}
		jsonErr := json.NewDecoder(r.Body).Decode(usr)
		if jsonErr != nil {
			http.Error(w, "Unable to unpack json into new user", http.StatusBadRequest)
			return
		}

		//convert new user to user with validation
		convertedUser, conversionError := usr.ToUser()
		if conversionError != nil {
			http.Error(w, conversionError.Error(), http.StatusBadRequest)
			return
		}

		//create user account in database
		insertedUser, insertErr := hc.UserStore.Insert(convertedUser)
		if insertErr != nil {
			http.Error(w, insertErr.Error(), http.StatusBadRequest)
			return
		}

		//begin session
		sessions.BeginSession(hc.SigningKey, hc.SessionStore, &SessionState{time.Now(), insertedUser}, w)

		//Respond to the client
		w.Header().Set("Content-Type", "application/json")

		userJSON, _ := json.Marshal(insertedUser)
		w.WriteHeader(http.StatusCreated)
		w.Write(userJSON)
	} else {
		http.Error(w, "Expecting only POST method", http.StatusMethodNotAllowed)
		return
	}

}

//get authenticated user
// func (hc *HandlerContext) GetAuthenticatedUser(r *http.Request) (*users.User, error) {
// 	//fetch the session state from the session store,
// 		if err2 != nil {
// 		return nil, err2
// 	}
// 	//and return the authenticated user
// 	//or an error if the user is not authenticated
// 	return sessState.AuthenticatedUser, nil
// }

//handles requests for a specific user based on user Id
func (hc *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	//get authenticated user
	sessState := &SessionState{}
	_, err := sessions.GetState(r, hc.SigningKey, hc.SessionStore, sessState)
	if err != nil {
		http.Error(w, "please sign-in", http.StatusUnauthorized)
		return
	}

	userID := int64(-1)
	if r.Method == "GET" {
		if path.Base(r.URL.Path) == "me" {
			userID = sessState.AuthenticatedUser.ID
		} else {
			convertedID, idErr := strconv.Atoi(path.Base(r.URL.Path))
			if idErr != nil {
				http.Error(w, "Passed ID was not valid", http.StatusBadRequest)
				return
			}
			userID = int64(convertedID)
		}
		// Get the user profile associated for the requested user ID
		userProfile, Err := hc.UserStore.GetByID(sessState.AuthenticatedUser.ID)
		if Err != nil {
			http.Error(w, "No user is found with this ID", http.StatusNotFound)
			return
		}
		//Respond to the client
		w.Header().Set("Content-Type", "application/json")

		userJSON, _ := json.Marshal(userProfile)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userJSON))

	} else if r.Method == "PATCH" {
		if path.Base(r.URL.Path) == "me" {
			userID = sessState.AuthenticatedUser.ID
		} else {
			convertedID, idErr := strconv.Atoi(path.Base(r.URL.Path))
			if idErr != nil {
				http.Error(w, "Passed ID was not valid", http.StatusBadRequest)
				return
			}
			userID = int64(convertedID)
			if userID != sessState.AuthenticatedUser.ID {
				http.Error(w, "The user is not the currently authenticated user", http.StatusForbidden)
				return
			}
		}

		ctype := r.Header.Get("Content-Type")

		if !strings.HasPrefix(ctype, "application/json") {
			http.Error(w, "The request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		// extract response body JSON to updates
		userUpdates := &users.Updates{}
		jsonErr := json.NewDecoder(r.Body).Decode(userUpdates)
		if jsonErr != nil {
			http.Error(w, "Unable to unpack json into new user", http.StatusBadRequest)
			return
		}

		// update user based on Updates
		updatedUser, updateErr := hc.UserStore.Update(sessState.AuthenticatedUser.ID, userUpdates)
		if updateErr != nil {
			http.Error(w, updateErr.Error(), http.StatusUnsupportedMediaType)
			return
		}
		//Respond to the client
		w.Header().Set("Content-Type", "application/json")

		userJSON, _ := json.Marshal(updatedUser)
		w.WriteHeader(http.StatusOK)
		w.Write(userJSON)

	} else {
		http.Error(w, "Expecting GET or PATCH method", http.StatusMethodNotAllowed)
		return
	}
}

//handles requests for the "sessions" resource,
//and allows clients to begin a new session using an existing user's credentials.
func (hc *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ctype := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ctype, "application/json") {
			http.Error(w, "The request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		// extract response body JSON to credentials
		userCred := &users.Credentials{}
		jsonErr := json.NewDecoder(r.Body).Decode(userCred)
		if jsonErr != nil {
			http.Error(w, "Unable to unpack json into new user", http.StatusBadRequest)
			return
		}
		// get user by email
		userProfile, getErr := hc.UserStore.GetByEmail(userCred.Email)
		if getErr != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		//authenticate user
		authErr := userProfile.Authenticate(userCred.Password)
		if authErr != nil {
			http.Error(w, authErr.Error(), http.StatusUnauthorized)
			return
		}

		//begin session
		sessions.BeginSession(hc.SigningKey, hc.SessionStore, &SessionState{time.Now(), userProfile}, w)

		//Get client IP address
		clientIP := r.RemoteAddr
		if len(clientIP) == 0 {
			clientIP = r.Header.Get("X-Forwarded-For")
		}

		//Log user Sign Ins
		DbError := hc.UserStore.LogSignIn(userProfile, time.Now(), clientIP)
		if DbError != nil {
			http.Error(w, DbError.Error(), http.StatusBadRequest)
			return
		}

		//Respond to the client
		w.Header().Set("Content-Type", "application/json")

		userJSON, _ := json.Marshal(userProfile)
		w.WriteHeader(http.StatusCreated)
		w.Write(userJSON)

	} else {
		http.Error(w, "Expecting POST method", http.StatusMethodNotAllowed)
		return
	}
}

//handles requests related to a specific authenticated session
func (hc *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		if !strings.HasSuffix(r.URL.Path, "mine") {
			http.Error(w, "Request path doesnt end with mine", http.StatusForbidden)
			return
		}
		// end current session
		_, err := sessions.EndSession(r, hc.SigningKey, hc.SessionStore)
		if err != nil {
			http.Error(w, "Unable to end Session", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Signed out"))

	} else {
		http.Error(w, "Expecting DELETE method", http.StatusMethodNotAllowed)
		return
	}
}
