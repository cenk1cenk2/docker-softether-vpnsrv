package pipe

import (
	"net"
)

type Ctx struct {
	Health struct {
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
