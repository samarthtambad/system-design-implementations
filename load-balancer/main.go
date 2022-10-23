package main

import (
	"net"
)

type IncomingReq struct {
	srcConn net.Conn
	reqId   string
}

func initLoadBalancer() *LoadBalancer {
	servers := []*Server{
		{Host: "localhost", Port: 9091, IsHealthy: true},
		{Host: "localhost", Port: 9092, IsHealthy: true},
		{Host: "localhost", Port: 9093, IsHealthy: true},
		{Host: "localhost", Port: 9094, IsHealthy: true},
	}

	lb := &LoadBalancer{
		strategy: newRoundRobinStrategy(servers),
		servers:  servers,
		host:     "",
		port:     9090,
	}

	return lb
}

func main() {
	lb := initLoadBalancer()
	lb.Run()
}
