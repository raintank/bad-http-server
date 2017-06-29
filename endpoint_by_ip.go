package main

import (
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Ip      string
	Updated time.Time
	Good    uint64
	Bad     uint64
}

func NewClient(ip string) *Client {
	return &Client{
		Ip:      ip,
		Updated: time.Now(),
	}
}

// EndpointByIp keeps 2 buckets of active clients by response they will get.
// the size of the buckets matches the ratio as closely as possible.
// note: if we keep the Clients sorted, than results would be more consistent in light of pool changes
// but we don't need to care about that for now.
type EndpointByIp struct {
	sync.Mutex
	Ratio          int // value could be 0, 100 or anything in between
	updated        time.Time
	Bad            []*Client
	Good           []*Client
	DroppedClients uint64
}

func NewEndpointByIp(ratio int) Endpoint {
	e := &EndpointByIp{
		Ratio:   ratio,
		updated: time.Now(),
	}
	go e.clean()
	return e
}

// clean will rebalance the good/bad buckets if these 2 conditions apply:
// * any of the clients hasn't been seen in an hour
// * the ratio hasn't been updated in an hour
// the rebalance aims to be stable (don't needlessly change assignments)
func (e *EndpointByIp) clean() {
	for range time.Tick(time.Hour) {
		e.Lock()
		if time.Now().Sub(e.updated) > time.Hour {
			var good []*Client
			var bad []*Client
			var dropped uint64
			for _, c := range e.Good {
				if time.Now().Sub(c.Updated) < time.Hour {
					good = append(good, c)
				} else {
					dropped += 1
				}
				for _, c := range e.Bad {
					if time.Now().Sub(c.Updated) < time.Hour {
						bad = append(bad, c)
					} else {
						dropped += 1
					}
				}
				if dropped != 0 {
					e.DroppedClients += dropped
					e.rebalance(good, bad)
				}
			}
		}
		e.Unlock()
	}
}

// rebalance will rebalance buckets with given starting points for good and bad
// caller must hold lock
func (e *EndpointByIp) rebalance(good, bad []*Client) {
	total := len(good) + len(bad)
	// numBad can be anywhere from 0 to total (inclusive)
	numBad := int(math.Floor((float64(total) * float64(e.Ratio) / 100) + 0.5))
	//fmt.Printf("rebalancing: good %v - bad %v . numBad: %d\n", good, bad, numBad)

	// move clients from good to bad as needed
	for len(bad) < numBad {
		//fmt.Println("moving to bad")
		bad = append(bad, good[len(good)-1])
		good = good[:len(good)-1]
	}

	// move clients from bad to good as needed
	for len(bad) > numBad {
		//fmt.Println("moving to good")
		good = append(good, bad[len(bad)-1])
		bad = bad[:len(bad)-1]
	}

	e.Good = good
	e.Bad = bad
}

func (ec *EndpointByIp) addOrUpdate(ip string) (client *Client, fail bool) {
	for _, c := range ec.Good {
		if c.Ip == ip {
			c.Updated = time.Now()
			return c, false
		}

	}
	for _, c := range ec.Bad {
		if c.Ip == ip {
			c.Updated = time.Now()
			return c, true
		}

	}
	c := NewClient(ip)
	if closestRatio(float64(ec.Ratio)/100, float64(len(ec.Bad)), float64(len(ec.Good))) {
		ec.Bad = append(ec.Bad, c)
		return c, true
	}
	ec.Good = append(ec.Good, c)
	return c, false
}

func (e *EndpointByIp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := strings.LastIndex(r.RemoteAddr, ":")
	if i == -1 {
		panic("unrecognized remote addr " + r.RemoteAddr)
	}
	ip := r.RemoteAddr[:i]
	e.Lock()
	defer e.Unlock()
	client, fail := e.addOrUpdate(ip)
	if fail {
		client.Bad += 1
		http.Error(w, "panic.", http.StatusInternalServerError)
	} else {
		client.Good += 1
		w.Write([]byte("ok\n"))
	}
}

func (e *EndpointByIp) Update(ratio int) {
	e.Lock()
	e.Ratio = ratio
	e.rebalance(e.Good, e.Bad)
	e.updated = time.Now()
	e.Unlock()
}
