package main

import (
	"fmt"
	"net/http"

	"loadbalancer_go/src/balancer"
	"loadbalancer_go/src/server"
)

func main() {
	servers := []server.Server{
		server.NewSimpleServer("https://www.facebook.com"),
		server.NewSimpleServer("https://www.bing.com"),
		server.NewSimpleServer("https://www.duckduckgo.com"),
	}

	lb := balancer.NewLoadBalancer("8000", servers)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(w, r)
	})

	fmt.Printf("Load balancer started on port %s\n", lb.Port())
	if err := http.ListenAndServe(":"+lb.Port(), nil); err != nil {
		fmt.Println("Error:", err)
	}
}
