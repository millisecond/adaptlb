package lb

import (
	"context"
	"github.com/millisecond/adaptlb/model"
	"strconv"
	"log"
)

func healthcheck(frontend *model.Frontend) {
	// TODO: real healthcheck channel w/Stop and Update
	// For now, just copy all backends into live servers
	ctx := context.Background()
	for _, pool := range frontend.ServerPools {
		func() {
			pool.LiveServers = []*model.LiveServer{}
			pool.LiveServerMutex.Lock(ctx)
			defer pool.LiveServerMutex.Unlock(ctx)
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
		}()
	}
	log.Println("Finished healthcheck")
}
