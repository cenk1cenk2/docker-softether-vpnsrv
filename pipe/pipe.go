package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type (
	Health struct {
		CheckInterval     string
		DhcpServerAddress string
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
		SetTasks(
			TL.JobSequence(
				// terminate handler
				Terminate(&TL).Job(),
				Setup(&TL).Job(),

				TL.JobParallel(
					CreatePostroutingRules(&TL).Job(),
					GenerateDhcpServerConfiguration(&TL).Job(),
					GenerateSoftEtherServerConfiguration(&TL).Job(),
				),

				TL.JobSequence(
					CreateTapDevice(&TL).Job(),
					CreateBridgeDevice(&TL).Job(),
				),

				TL.JobParallel(
					RunDnsServer(&TL).Job(),
					RunSoftEtherVpnServer(&TL).Job(),
				),

				HealthCheck(&TL).Job(),

				TL.JobWaitForTerminator(),
			),
		)
}
