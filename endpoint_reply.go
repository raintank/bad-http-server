package main

import (
	"net/http"
	"sync"
)

// EndpointReply aims for the given ratio by counting individual replies.
type EndpointReply struct {
	sync.Mutex
	Ratio int // value could be 0, 100 or anything in between
	Good  uint64
	Bad   uint64
}

func (e *EndpointReply) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func NewEndpointReply(ratio int) Endpoint {
	e := &EndpointReply{
		Ratio: ratio,
	}
	return e
}
