package main

import (
	"math"
	"net/http"
)

type Endpoint struct {
	Ratio int // value could be 0, 100 or anything in between
	Good  uint64
	Bad   uint64
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if we serve an ok the ratio would be:
	a := float64(e.Bad) / float64(e.Good+1+e.Bad)
	// if we serve an error the ratio would be:
	b := float64(e.Bad+1) / float64(e.Good+e.Bad+1)

	ratioNorm := float64(e.Ratio) / 100

	if math.Abs(b-ratioNorm) < math.Abs(a-ratioNorm) {
		e.Bad += 1
		http.Error(w, "panic.", http.StatusInternalServerError)
	} else {
		e.Good += 1
		w.Write([]byte("ok\n"))
	}
}

func New(ratio int) *Endpoint {
	e := &Endpoint{
		Ratio: ratio,
	}
	return e
}

func (e *Endpoint) Serve(w http.ResponseWriter, r *http.Request) *Endpoint {
	e.ServeHTTP(w, r)
	return e
}
