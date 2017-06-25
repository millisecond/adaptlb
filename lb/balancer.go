package lb

import (
	"github.com/millisecond/adaptlb/model"
	"github.com/docker/docker/pkg/random"
)

// Sets LiveServer and other state on req
func LoadBalanceHTTP(req *model.LBRequest) {
	// TODO: special handling of req/resp, but use L4 to assign a server
	LoadBalanceL4(req)
}

func LoadBalanceL4(req *model.LBRequest) bool {
	// L4 FE's must have exactly one server pool
	req.ServerPool = req.Frontend.ServerPools[0]
	req.ServerPool.LiveServerMutex.RLock()
	serverCount := uint64(len(req.ServerPool.LiveServers))
	if serverCount == 0 {
		req.ServerPool.LiveServerMutex.RUnlock()
		return false
	}
	switch req.ServerPool.Strategy {
	case model.LBStrategyRoundRobin:
		req.LiveServer = req.ServerPool.LiveServers[int(req.ServerPool.SharedLBState.IncrAndGetRequests() %serverCount)]
	case model.LBStrategyRandom:
		fallthrough
	default:
		// Random, doesn't need to be cryptographically defensible, just spread out
		req.LiveServer = req.ServerPool.LiveServers[int(random.Rand.Uint64() % serverCount)]
	}
	req.ServerPool.LiveServerMutex.RUnlock()
	return true
}
