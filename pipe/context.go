package pipe

import (
	"net"
	"time"
)

type Ctx struct {
	Health struct {
		Duration      time.Duration
		SoftEtherPIDs []int
		DnsMasqPIDs   []int
	}

	Server struct {
		Network    *net.IPNet
		RangeStart net.IP
		RangeEnd   net.IP
	}

	DhcpServer struct {
		Options map[string]string
	}
}
