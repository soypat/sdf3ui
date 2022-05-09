package main

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type baseMiddleware struct {
	handler http.Handler
}

func newBaseMiddleware(handlerToWrap http.Handler) *baseMiddleware {
	return &baseMiddleware{
		handler: handlerToWrap,
	}
}

func (rh *baseMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// TODO
	aborted := corsMW(rw, r, "GET", "POST")
	if aborted {
		return
	}
	start := time.Now()
	// log.Printf("%s receive\n", r.URL.String())

	rh.handler.ServeHTTP(rw, r)
	log.Printf("%s elapsed %s \n", r.URL.String(), time.Since(start))
}

var (
	corsMWMethods string
	onceMethodSet sync.Once
)

// CORS middleware
func corsMW(rw http.ResponseWriter, r *http.Request, methods ...string) (aborted bool) {
	onceMethodSet.Do(func() {
		corsMWMethods = strings.Join(append(methods, "OPTIONS"), ",")
	})
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Max-Age", "999999")            // Allow javascript requests from localhost
	rw.Header().Set("Access-Control-Allow-Methods", corsMWMethods) // We will use GET and POST methods
	// rw.Header().Set("Access-Control-Allow-Headers", "Content-Type") // allow JSON interchange
	rw.Header().Set("Access-Control-Allow-Headers", "content-type")
	rw.Header().Set("Access-Control-Allow-Origin", "*") // Allow javascript requests from localhost
	if r.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}
