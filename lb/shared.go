package lb

import (
	"errors"
	"github.com/millisecond/adaptlb/config"
	"github.com/millisecond/adaptlb/model"
	"net"
	"strconv"
	"strings"
	"sync"
)

var activationMutex = &sync.Mutex{}
var activeConfig *config.Config

func Activate(cfg *config.Config) error {
	activationMutex.Lock()
	defer activationMutex.Unlock()

	if activeConfig == nil {
		for _, frontend := range cfg.Frontends {
			err := addListener(frontend)
			if err != nil {
				return err
			}
		}
		activeConfig = cfg
		return nil
	}

	// we have an existing config, need to de-dupe
	removedFrontEnds := []*model.Frontend{}
	addedFrontends := activeConfig.Frontends

	// Update matches and collect FE's to remove
	for _, existing := range activeConfig.Frontends {
		found := false
		for _, toAdd := range cfg.Frontends {
			if toAdd.RowID == existing.RowID {
				//toAdd.Listeners = existing.Listeners
				found = true
			}
		}
		if !found {
			removedFrontEnds = append(removedFrontEnds, existing)
		}
	}

	// Find new FE's
	for _, toAdd := range cfg.Frontends {
		found := false
		for _, existing := range activeConfig.Frontends {
			if toAdd.RowID == existing.RowID {
				// TODO changed ports
				//toAdd.Listeners = existing.Listeners
				found = true
			}
		}
		if !found {
			addedFrontends = append(addedFrontends, toAdd)
		}
	}

	// Stop old ones
	//for _, fe := range removedFrontEnds {
	//	for _, listener := range *fe.Listeners {
	//		(*listener).Stop()
	//	}
	//}

	// Start listening on new ones
	//for _, fe := range addedFrontends {
	//	addListener(fe)
	//}

	activeConfig.Frontends = cfg.Frontends

	return nil
}

func addListener(frontend *model.Frontend) error {
	switch frontend.Type {
	case "http":
	case "tcp":
		err := AddTCPPort(frontend)
		if err != nil {
			return err
		}
	case "udp":
	default:
		return errors.New("Unknown config type: " + frontend.Type)
	}
	return nil
}

func portFromConn(c net.Conn) (int, error) {
	_, portS, err := net.SplitHostPort(c.LocalAddr().String())
	if err != nil {
		return -1, err
	}
	port, err := strconv.Atoi(portS)
	if err != nil {
		return -1, err
	}
	return port, nil
}

func parsePorts(portString string) ([]int, error) {
	ports := []int{}
	parts := strings.Split(portString, ",")
	for _, part := range parts {
		port, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}
	return ports, nil
}
