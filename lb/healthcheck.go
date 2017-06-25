package lb

import (
	"github.com/millisecond/adaptlb/model"
	"strconv"
	"github.com/containous/traefik/log"
)

func healthcheck(frontend *model.Frontend) {
	// TODO: real healthcheck channel w/Stop and Update
	// For now, just copy all backends into live servers
	for _, pool := range frontend.ServerPools {
		pool.LiveServers = []*model.LiveServer{}
		for _, backend := range pool.Backends {
			switch backend.Type {
			case model.LBBackendTypeIndividual:
				pool.LiveServers = append(pool.LiveServers, &model.LiveServer{
					Address: backend.Address + ":" + strconv.Itoa(backend.Port),
				})
			default:
				log.Println("Unsupported backend type: " + string(backend.Type))
			}
		}
	}
	log.Println("Finished healthcheck")
}
