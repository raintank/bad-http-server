package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var path string
var addr string

var lock sync.Mutex
var ratio int // value could be 0, 100 or anything in between

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s <addr>\n", os.Args[0])
		os.Exit(2)
	}
	addr = os.Args[1]
	fmt.Println("will listen for http traffic on", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		cutoff := strings.LastIndex(r.URL.Path, "/")
		end := r.URL.Path[cutoff+1:]
		if end != "" {
			i, err := strconv.Atoi(end)
			if err == nil && ratio >= 0 && ratio <= 100 {
				ratio = i
				w.Write([]byte("updated\n"))
				return
			}
		}
		if int64(ratio) > time.Now().UnixNano()%100 {
			http.Error(w, "panic.", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok\n"))
	})
	http.ListenAndServe(addr, nil)
}
