package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
)

type Server struct {
	Host        string
	Port        int
	IsHealthy   bool
	NumRequests int
}

type Event struct {
	EventName string
	Data      interface{}
}

type LoadBalancer struct {
	strategy BalancingStrategy
	servers  []*Server
	events   chan Event
	host     string
	port     int
}

func (lb *LoadBalancer) Run() {
	fmt.Println("Starting load balancer ...")

	// Listen for incoming connections.
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", lb.host, lb.port))
	if err != nil {
		panic(err)
	}

	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Printf("Load balancer listening on port %d ...\n", lb.port)

	// handle commands
	go func() {
		for {
			select {
			case event := <-lb.events:
				switch event.EventName {
				case "CMD_exit":
					log.Println("Terminating gracefully ...")
					return
				case "CMD_addServer":
					server, success := event.Data.(Server)
					if !success {
						panic(err)
					}
					log.Printf("Adding server: %v", server)
				case "CMD_changeStrategy":
					strategyName, success := event.Data.(string)
					if !success {
						panic(err)
					}
					switch strategyName {
					case "round_robin":
						log.Println("round robin strategy")
					}
				}
			}
		}
	}()

	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			panic(err)
		}

		// Handle connections in a new goroutine.
		go lb.proxy(IncomingReq{
			srcConn: conn,
			reqId:   uuid.NewString(),
		})
	}
}

func (lb *LoadBalancer) proxy(req IncomingReq) {
	// get the server to route request
	server := lb.strategy.nextServer()
	log.Printf("Routing request %s to server %s:%d", req.reqId, server.Host, server.Port)

	// setup connection to the server
	backendConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))
	if err != nil {
		log.Printf("Error connecting to server %s - %s", fmt.Sprintf("%s:%d", server.Host, server.Port), err.Error())

		// send back error to source
		req.srcConn.Write([]byte("server not available"))
		req.srcConn.Close()
		panic(err)
	}

	server.NumRequests++
	go io.Copy(backendConn, req.srcConn)
	go io.Copy(req.srcConn, backendConn)
}
