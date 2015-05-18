package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var path string
var addr string

var lock sync.Mutex
var ratio int // value could be 0, 100 or anything in between
var good uint64
var bad uint64

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
				good = 0
				bad = 0
				w.Write([]byte("updated\n"))
				return
			}
		}
		// if we serve an ok the ratio would be:
		a := float64(bad) / float64(good+1+bad)
		// if we serve an error the ratio would be:
		b := float64(bad+1) / float64(good+bad+1)

		ratioNorm := float64(ratio) / 100

		if math.Abs(b-ratioNorm) < math.Abs(a-ratioNorm) {
			bad += 1
			http.Error(w, "panic.", http.StatusInternalServerError)
			return
		}
		good += 1
		w.Write([]byte("ok\n"))
	})
	http.ListenAndServe(addr, nil)
}
