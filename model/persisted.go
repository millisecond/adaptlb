package model

import (
	"time"
	"sync"
)

type LBType string
const (
	LBTypeHTTP LBType = "http"
	LBTypeTCP LBType = "tcp"
	LBTypeUDP LBType = "udp"
)

type LBStrategy string
const (
	LBStrategyRandom LBStrategy = "random"
	LBStrategyRoundRobin LBStrategy = "round_robin"
)

type LBStickySessionType string
const (
	LBStickySessionTypeCookie LBStickySessionType = "cookie"
	LBStickySessionTypeSrcIP LBStickySessionType = "src_ip"
	LBStickySessionTypeDstIP LBStickySessionType = "dst_ip"
	LBStickySessionTypeSrcPort LBStickySessionType = "src_port"
	LBStickySessionTypeDstPort LBStickySessionType = "dst_port"
	LBStickySessionTypeSrcIPPort LBStickySessionType = "src_ipport"
	LBStickySessionTypeDstIPPort LBStickySessionType = "dst_ipport"
)

type LBBackendType string
const (
	LBBackendTypeIndividual LBBackendType = "individual"
	LBBackendTypeASG LBBackendType = "asg"
	LBBackendTypeTG LBBackendType = "tg"
)

type LBHealthcheckType string
const (
	LBHealthcheckTypeHTTP LBHealthcheckType = "http"
	LBHealthcheckTypeTCP LBHealthcheckType = "tcp"
	LBHealthcheckTypeICMP LBHealthcheckType = "icmp"
)

type DNSType string
const (
	DNSTypeA DNSType = "A"
	DNSTypeAAAA DNSType = "AAAA"
	DNSTypeBothA DNSType = "BOTHA"
)

type Frontend struct {
	RowID string `json:"rowId,omitempty"`

	Type LBType `json:"type,omitempty"` // "http", "tcp", or"udp": HTTP load balancers can share ports, TCP/UDP are exclusive

	// HTTP can have many server pools balancing paths, TCP/UDP only have one
	ServerPools []*ServerPool `serverPools:"type,omitempty"`

	// HTTP-only
	FQDN string `json:"fqdn,omitempty"` // All HTTP domains use hostname routing, TCP and UDP use for updating A records
	//TLSHeader        string `json:"tlsHeader,omitempty"`
	//TLSHeaderValue   string `json:"tlsHeaderValue,omitempty"`
	//GZIPContentTypes string `json:"gzipContentTypes,omitempty"`
	NoRouteStatus int `json:"noRouteStatus,omitempty"`

	//All
	Ports          string     `json:"ports,omitempty"`      // "80" or "80,443" or "80-8000"
	SecurePorts    string     `json:"ports,omitempty"`      // "443" or "443,8443" or "443-8443"
	DNSRecords     DNSRecords `json:"dnsRecords,omitempty"` // Type of upstream DNS provider to update ("route53") or "" to disable updates
	ShutdownWaitMS int        `json:"shutdownWaitMS,omitempty"`

	Listeners *map[int]*Listener `json:"-"`
}

type ServerPool struct {
	RowID string `json:"rowId,omitempty"`

	Backends []Backend `json:"backends,omitempty"` // Can be multiple to facilitate IP Lists and blue/green deploys
	//CircuitBreaker *CircuitBreaker    `json:"circuitBreaker,omitempty"`

	Strategy      LBStrategy `json:"strategy,omitempty"`      // "random", "roundrobin" - picks servers, sticky will override
	StickySessionType LBStickySessionType `json:"stickySessionType,omitempty"` // "cookie", "src_ip", "src_port", "dst_ip", "dst_port"

	// HTTP/TCP-only
	TLSBackend bool `json:"tlsBackend,omitempty"` // what version of proto to use on back-side

	// HTTP-only
	Path string `json:"route,omitempty"` // "" == default, "/admin", "/checkout"

	// Net connection settings
	MaxIdle               int           `json:"maxIdle,omitempty"`
	MaxIdlePerHost        int           `json:"maxIdlePerHost,omitempty"`
	HealthCheck           *HealthCheck  `json:"healthCheck,omitempty"`
	DialTimeout           time.Duration `json:"dialTimeout,omitempty"`
	ResponseHeaderTimeout time.Duration `json:"responseHeaderTimeout,omitempty"`
	KeepAliveTimeout      time.Duration `json:"keepAliveTimeout,omitempty"`
	FlushInterval         time.Duration `json:"flushInterval,omitempty"`

	// In-memory state, don't persist
	LiveServers []*LiveServer  `json:"-"`
	LiveServerMutex *sync.RWMutex  `json:"-"`
	SharedLBState *SharedLBState `json:"-"`
}

// DoS prevention: if one of these conditions is triggered for a node, it's no longer available as a target.
// NOTE: this is true even if it's the target of a sticky session, if sticky must not be broken don't use this.
//type CircuitBreaker struct {
//	MaxConn        int `json:"maxConn"`
//	NodeOverweight int `json:"nodeOverweight"` // % over average node conns allowed to prevent over-sticky pounding
//}

type Backend struct {
	RowID string `json:"rowId,omitempty"`

	Type    LBBackendType `json:"type,omitempty"` // "individual", "asg", "targetgroup", "tagged"
	Address string `json:"address,omitempty"`   // ip/hostname, ASG-name, targetgroup name, tag
	Port    int    `json:"port"`           // port to use when connecting, invalid for "targetgroup"
	Weight  int    `json:"weight"`         // portion of traffic to send

	Connections *[]*LiveConnection
}

type HealthCheck struct {
	RowID string `json:"rowId,omitempty"`

	Type LBHealthcheckType // "http", "tcp", "icmp"

	// "http"-specific fields
	HTTPPath                    string `json:"httpPath,omitempty"`
	HTTPSuccessCodes            []int  `json:"httpSuccessCodes,omitempty"`
	HTTPSuccessResponseContains string `json:"httpSuccessResponseContains,omitempty"`

	IntervalSeconds           int `json:"intervalSeconds,omitempty"`
	Timeout            int    `json:"timeout,omitempty"`
	HealthyThreshold   int    `json:"healthyThreshold,omitempty"`
	UnhealthyThreshold int    `json:"unhealthyThreshold,omitempty"`
}

type DNSRecords struct {
	RowID string `json:"rowId,omitempty"`

	Enabled     bool
	ZoneName    string
	RecordSetID string
	Type        DNSType
	Records     []string
}
