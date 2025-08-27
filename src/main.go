package main

import(
	"fmt",
	"net/url",
	"net/http/httputil",
)

type simpleServer struct{
	addr string
	proxy httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleServer{
  serverUrl,err := url.Parse(addr)
	handleErr(err)
	return &simpleserver{
		addr: addr,
		proxy: *httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func handleErr(err error){
	if err != nil{
		fmt.Printf("Error: %v/n", err)
	}
	os.Exit(1)
}