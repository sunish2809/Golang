package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type LoadBalancer struct {
	port            string
	roundRobinCount int64
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%int64(len(lb.servers))]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%int64(len(lb.servers))]
	}
	lb.roundRobinCount++
	return server

}

func (lb *LoadBalancer) serverProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Println("forwarding request to address", targetServer.Address())
	targetServer.Serve(w, r)
}
func (s *simpleServer) Address() string {
	return s.addr
}
func (s *simpleServer) IsAlive() bool {
	return true
}
func (s *simpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.bing.com"),
		newSimpleServer("http://www.duckduckgo.com"),
	}

	lb := NewLoadBalancer("8000", servers)
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serverProxy(w, r)
	}
	http.HandleFunc("/", handleRedirect)
	fmt.Println("serving request at localhost ", lb.port)
	http.ListenAndServe(":"+lb.port, nil)

}
