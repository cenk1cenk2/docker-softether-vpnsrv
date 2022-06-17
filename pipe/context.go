package pipe

import (
	"net"
	"time"
)

type Ctx struct {
	Health struct {
		Duration     time.Duration
		SoftEtherPID int
		DnsMasqPID   int
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
