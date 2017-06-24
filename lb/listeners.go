package lb

import "github.com/millisecond/adaptlb/model"

// HTTP types share lb, others don't
var sharedListeners = map[int]*model.Listener{}
var uniqueListeners = map[int]*model.Listener{}

