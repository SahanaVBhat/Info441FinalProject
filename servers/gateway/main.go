package main

import (
	"Info441FinalProject/servers/gateway/handlers"
	"Info441FinalProject/servers/gateway/models/users"
	"Info441FinalProject/servers/gateway/sessions"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

//main is the main entry point for the server
func main() {

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	tlscert := os.Getenv("TLSCERT")
	tlskey := os.Getenv("TLSKEY")

	if len(tlscert) == 0 || len(tlskey) == 0 {
		fmt.Fprintln(os.Stderr, "TLS certificate or key not found")
		os.Exit(1)
	}

	signingKey := os.Getenv("SESSIONKEY")
	reddisAddr := os.Getenv("REDISADDR")
	dsn := os.Getenv("DSN")

	// for microservices
	microserviceAddr := os.Getenv("MICROSERVICEADDR")
	if len(microserviceAddr) == 0 {
		fmt.Print("MICROSERVICEADDR not provided")
	}

	if len(reddisAddr) == 0 {
		reddisAddr = "127.0.0.1:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: reddisAddr,
	})
	redisStore := sessions.NewRedisStore(client, time.Hour)

	// Create a MYSQL local database with password
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Errorf("error opening database: %v\n", err)
		os.Exit(1)
	}
	sqlStore := users.NewMainSQLStore(db)
	// Create and initialize a new instance of your handler context
	context := handlers.HandlerContext{
		SigningKey:   signingKey,
		SessionStore: redisStore,
		UserStore:    sqlStore}

	mux := http.NewServeMux()

	// check for ccurrently authenticated user
	microserviceDirector := func(r *http.Request) {
		var microserviceURLs []*url.URL
		microserviceAddresses := strings.Split(microserviceAddr, ", ")
		for _, msaddr := range microserviceAddresses {
			url, err := url.Parse(msaddr)
			if err != nil {
				log.Fatalf("error parsing address: %v\n", err)
			}
			microserviceURLs = append(microserviceURLs, url)
		}
		serverName := microserviceURLs[0]
		if len(microserviceURLs) > 1 {
			rand.Seed(time.Now().UnixNano())
			serverNum := rand.Intn(len(microserviceURLs))
			serverName = microserviceURLs[serverNum]
		}

		r.Header.Add("X-Forwarded-Host", r.Host)

		r.Host = serverName.Host
		r.URL.Host = serverName.Host
		r.URL.Scheme = serverName.Scheme

		sessionState := &handlers.SessionState{}
		_, err := sessions.GetState(r, signingKey, redisStore, sessionState)

		if err == nil {
			user := &users.User{ID: sessionState.AuthenticatedUser.ID}
			json, _ := json.Marshal(user)
			r.Header.Set("X-User", string(json))
		} else {
			r.Header.Del("X-User")
		}
	}

	// reverse proxies
	microserviceProxy := &httputil.ReverseProxy{Director: microserviceDirector}

	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)

	mux.Handle("/v1/courses", microserviceProxy)
	mux.Handle("/v1/courses/", microserviceProxy)
	mux.Handle("/v1/evaluations", microserviceProxy)
	mux.Handle("/v1/evaluations/", microserviceProxy)

	//wrap this new mux with your CORS middleware handler
	wrappedMux := handlers.NewCorsMiddleware(mux)

	log.Printf("server is listening at http://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, wrappedMux))

}
