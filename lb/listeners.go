package lb

import (
	"github.com/millisecond/adaptlb/model"
)

// HTTP types share lb, others don't
var sharedListeners = map[int]*Listener{}
var uniqueListeners = map[int]*Listener{}

