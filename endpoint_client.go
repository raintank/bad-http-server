package main

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Ip   string
	Seen time.Time
	Good uint64
	Bad  uint64
}

type Clients struct {
	Clients []*Client
}

func (cl *Clients) Add(ip string, now time.Time) *Client {
	client := &Client{
		Ip:   ip,
		Seen: now,
	}
	for i, c := range cl.Clients {
		if c == nil {
			cl.Clients[i] = client
			return client
		}
	}
	cl.Clients = append(cl.Clients, client)
	return client
}

func (cl *Clients) Update(now, cutoff time.Time, ip string) *Client {
	var client *Client
	for i, c := range cl.Clients {
		if c == nil {
			continue
		}
		if c.Ip == ip {
			c.Seen = now
			client = c
		}
		if c.Seen.Before(cutoff) {
			cl.Clients[i] = nil
		}
	}
	return client
}

// EndpointClient keeps 2 buckets of active clients by response they will get.
// the size of the buckets matches the ratio as closely as possible.
// note: if we keep the Clients sorted, than results would be more consistent in light of pool changes
// but we don't need to care about that for now.
type EndpointClient struct {
	sync.Mutex
	Ratio int // value could be 0, 100 or anything in between
	Good  Clients
	Bad   Clients
}

// prunes out any stale ips if needed
// assures the ip is in the pool.
// assures the ratio is as close to the ideal as possible. (TODO: for now we only do this for new clients, should maintain this at all times)
// returns the Client object and whether it's in the fail group
func (ec *EndpointClient) Update(ip string) (client *Client, fail bool) {
	now := time.Now()
	cutoff := now.Add(-time.Duration(5) * time.Minute)
	clientGood := ec.Good.Update(now, cutoff, ip)
	clientBad := ec.Bad.Update(now, cutoff, ip)
	if clientGood != nil && clientBad != nil {
		panic("client was found both in good and bad")
	}
	if clientBad != nil {
		fail = true
		client = clientBad
	}
	if clientGood != nil {
		client = clientGood
	}
	if client == nil {
		if closestRatio(float64(ec.Ratio)/100, float64(len(ec.Bad.Clients)), float64(len(ec.Good.Clients))) {
			client = ec.Bad.Add(ip, now)
			fail = true
		} else {
			client = ec.Good.Add(ip, now)
		}
	}
	return client, fail
}

func (e *EndpointClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.Lock()
	defer e.Unlock()
	i := strings.LastIndex(r.RemoteAddr, ":")
	if i == -1 {
		panic("unrecognized remote addr " + r.RemoteAddr)
	}
	ip := r.RemoteAddr[:i]
	client, fail := e.Update(ip)
	if fail {
		client.Bad += 1
		http.Error(w, "panic.", http.StatusInternalServerError)
	} else {
		client.Good += 1
		w.Write([]byte("ok\n"))
	}
}

func NewEndpointClient(ratio int) Endpoint {
	e := &EndpointClient{
		Ratio: ratio,
	}
	return e
}

func (e *EndpointClient) Serve(w http.ResponseWriter, r *http.Request) Endpoint {
	e.ServeHTTP(w, r)
	return e
}
