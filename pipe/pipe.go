package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Health struct {
		CheckInterval     string
		DhcpServerAddress string
		EnablePing        bool
	}

	DhcpServer struct {
		Template       string `validate:"omitempty,file"`
		Lease          string
		Gateway        string `validate:"omitempty,ip"`
		SendGateway    bool
		ForwardingZone cli.StringSlice `validate:"omitempty,ip"`
	}

	LinuxBridge struct {
		BridgeInterface   string
		UpstreamInterface string
		UseDhcp           bool
		StaticIp          string `validate:"omitempty,ip"`
	}

	SoftEther struct {
		Template     string `validate:"file"`
		TapInterface string
		DefaultHub   string
	}

	Server struct {
		Mode        string `validate:"oneof=dhcp bridge"`
		CidrAddress string `validate:"cidrv4"`
	}

	Pipe struct {
		Ctx

		Health
		DhcpServer
		LinuxBridge
		SoftEther
		Server
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		Set(
			func(tl *TaskList[Pipe]) Job {
				return tl.JobSequence(
					Tasks(tl).Job(),
					Services(tl).Job(),
					HealthCheck(tl).Job(),
					tl.JobWaitForTerminator(),
				)
			})
}
