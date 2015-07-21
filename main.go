package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var path string
var addr string

var lock sync.Mutex

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s <addr>\n", os.Args[0])
		os.Exit(2)
	}
	addr = os.Args[1]
	fmt.Println("will listen for http traffic on", addr)

	endpoints := make(map[string]*Endpoint)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "" {
			http.Error(w, "empty path", http.StatusBadRequest)
		}
		lock.Lock()
		defer lock.Unlock()
		js, err := json.Marshal(endpoints)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
	http.HandleFunc("/bad/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) == 5 {
			http.Error(w, "empty key", http.StatusBadRequest)
			return
		}
		e, ok := endpoints[r.URL.Path]
		if ok {
			e.ServeHTTP(w, r)
			return
		}
		badRatio, err := strconv.Atoi(r.URL.Path[5:])
		if err != nil || badRatio < 0 || badRatio > 100 {
			http.Error(w, "bad ratio (should be a percentage between 0 and 100, inclusive)", http.StatusBadRequest)
			return
		}
		endpoints[r.URL.Path] = New(badRatio).Serve(w, r)
	})
	http.HandleFunc("/custom/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) == 8 {
			http.Error(w, "empty key", http.StatusBadRequest)
			return
		}
		e, ok := endpoints[r.URL.Path]
		if ok {
			e.ServeHTTP(w, r)
			return
		}
		remainder := r.URL.Path[8:]
		if strings.Count(remainder, "/") > 1 {
			http.Error(w, "too many slashes", http.StatusBadRequest)
			return
		}
		pos := strings.LastIndex(remainder, "/")
		if pos == -1 {
			endpoints[r.URL.Path] = New(0).Serve(w, r)
			return
		}
		badRatio, err := strconv.Atoi(remainder[pos+1:])
		if err != nil || badRatio < 0 || badRatio > 100 {
			http.Error(w, "bad ratio (should be a percentage between 0 and 100, inclusive)", http.StatusBadRequest)
			return
		}
		key := "/custom/" + remainder[:pos]
		e, ok = endpoints[key]
		if ok {
			e.Ratio = badRatio
		} else {
			endpoints[key] = New(badRatio)
		}
		w.Write([]byte("updated\n"))
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
