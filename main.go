package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var path string
var addr string

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s <addr>\n", os.Args[0])
		os.Exit(2)
	}
	addr = os.Args[1]
	fmt.Println("will listen for http traffic on", addr)

	endpoints := NewEndpoints()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "" {
			http.Error(w, "empty path", http.StatusBadRequest)
		}
		endpoints.Lock()
		defer endpoints.Unlock()
		js, err := endpoints.Json()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		key, ratio, err := parseKeyRatio(r.URL.Path, "/static/ratio")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		endpoints.Lock()
		defer endpoints.Unlock()
		e, ok := endpoints.Get(key)
		if !ok {
			e = NewEndpointBasic(ratio)
			endpoints.Set(key, e)
		}
		e.ServeHTTP(w, r)
	})
	http.HandleFunc("/static-by-ip/", func(w http.ResponseWriter, r *http.Request) {
		key, ratio, err := parseKeyRatio(r.URL.Path, "/static-by-ip/ratio")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		endpoints.Lock()
		defer endpoints.Unlock()
		e, ok := endpoints.Get(key)
		if !ok {
			e = NewEndpointByIp(ratio)
			endpoints.Set(key, e)
		}
		e.ServeHTTP(w, r)
	})
	http.HandleFunc("/dynamic/", func(w http.ResponseWriter, r *http.Request) {
		key, ratio, err := parseDynamicKeyRatio(r.URL.Path, "/dynamic/key[/ratio]")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		endpoints.Lock()
		defer endpoints.Unlock()

		// user wants to update
		if ratio != -1 {
			e, ok := endpoints.Get(key)
			if !ok {
				http.Error(w, "not found", http.StatusBadRequest)
			} else {
				e.Update(ratio)
				w.Write([]byte("updated\n"))
			}
			return
		}

		e, ok := endpoints.Get(key)
		if !ok {
			e = NewEndpointBasic(0)
			endpoints.Set(key, e)
		}
		e.ServeHTTP(w, r)
	})

	http.HandleFunc("/dynamic-by-ip/", func(w http.ResponseWriter, r *http.Request) {
		key, ratio, err := parseDynamicKeyRatio(r.URL.Path, "/dynamic-by-ip/key[/ratio]")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		endpoints.Lock()
		defer endpoints.Unlock()

		// user wants to update
		if ratio != -1 {
			e, ok := endpoints.Get(key)
			if !ok {
				http.Error(w, "not found", http.StatusBadRequest)
			} else {
				e.Update(ratio)
				w.Write([]byte("updated\n"))
			}
			return
		}

		e, ok := endpoints.Get(key)
		if !ok {
			e = NewEndpointByIp(0)
			endpoints.Set(key, e)
		}
		e.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
