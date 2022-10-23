package main

type BalancingStrategy interface {
	registerServer(Server)
	nextServer() Server
}

// RoundRobin balancing strategy
type RoundRobin struct {
	servers       []*Server
	prevServerIdx int
}

func (rr *RoundRobin) registerServer(server Server) {
	rr.servers = append(rr.servers, &server)

}

func (rr *RoundRobin) nextServer() Server {
	curIdx := rr.prevServerIdx + 1
	if curIdx == len(rr.servers) {
		curIdx = 0
	}

	server := *rr.servers[curIdx]
	rr.prevServerIdx = curIdx
	return server
}

func newRoundRobinStrategy(servers []*Server) BalancingStrategy {
	roundRobin := RoundRobin{
		servers:       servers,
		prevServerIdx: -1,
	}
	return &roundRobin
}
