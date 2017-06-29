package main

import (
	"net/http"
	"sync"
)

// EndpointBasic aims for the given ratio by counting individual replies.
type EndpointBasic struct {
	sync.Mutex
	Ratio int // value could be 0, 100 or anything in between
	Good  uint64
	Bad   uint64
}

func NewEndpointBasic(ratio int) Endpoint {
	return &EndpointBasic{
		Ratio: ratio,
	}
}

func (e *EndpointBasic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.Lock()
	defer e.Unlock()

	if closestRatio(float64(e.Ratio)/100, float64(e.Bad), float64(e.Good)) {
		e.Bad += 1
		http.Error(w, "panic.", http.StatusInternalServerError)
	} else {
		e.Good += 1
		w.Write([]byte("ok\n"))
	}
}

func (e *EndpointBasic) Update(ratio int) {
	e.Lock()
	e.Ratio = ratio
	e.Unlock()
}
