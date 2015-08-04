package main

import (
	"encoding/json"
	"sync"
)

// since there's often intricate relations between Gets and Puts, this
// must be locked/unlocked by caller
type Endpoints struct {
	sync.Mutex
	endpoints map[string]Endpoint
}

func NewEndpoints() *Endpoints {
	return &Endpoints{
		endpoints: make(map[string]Endpoint),
	}
}

func (e *Endpoints) Json() ([]byte, error) {
	return json.Marshal(e.endpoints)
}
func (e *Endpoints) Get(key string) (Endpoint, bool) {
	r, ok := e.endpoints[key]
	return r, ok
}
func (e *Endpoints) Set(key string, endp Endpoint) {
	e.endpoints[key] = endp
}
