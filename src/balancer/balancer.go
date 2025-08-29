package balancer

import (
	"fmt"
	"loadbalancer_go/src/server"
	"net/http"
	"sync"
	"time"
)

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []server.Server
}

func NewLoadBalancer(port string, servers []server.Server) *LoadBalancer {
	lb :=  &LoadBalancer{
		port:    port,
		servers: servers,
	}
	go lb.healthCheck()

	return lb

}

func (lb *LoadBalancer) getNextAvailableServer() server.Server {
	for {
		srv := lb.servers[lb.roundRobinCount%len(lb.servers)]
		lb.roundRobinCount++

		if srv.IsAlive() {
			return srv
		}
	}
}

func (lb *LoadBalancer) ServeProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to %q\n", targetServer.Address())
	targetServer.Serve(w, r)
}

func (lb *LoadBalancer) Port() string {
	return lb.port
}

func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(10 * time.Second) // check every 10s
	defer ticker.Stop()

	for {
		<-ticker.C

		var wg sync.WaitGroup
		for _, srv := range lb.servers {
			wg.Add(1)
			go func(s server.Server) {
				defer wg.Done()
				resp, err := http.Get(s.Address())
				if err != nil || resp.StatusCode != http.StatusOK {
					fmt.Printf("Server %s is DOWN\n", s.Address())
					s.SetAlive(false)
				} else {
					fmt.Printf("Server %s is UP\n", s.Address())
					s.SetAlive(true)
					resp.Body.Close()
				}
			}(srv)
		}
		wg.Wait() // wait until all checks finish before next cycle
	}
}


