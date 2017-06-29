package main

import "net/http"

type Endpoint interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
