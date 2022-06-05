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

	Logs struct {
		CleanServerLogs   bool
		CleanPacketLogs   bool
		CleanSecurityLogs bool
		CleanLogsInterval string
	}

	DhcpServer struct {
		Template       string `validate:"omitempty,file"`
		TapInterface   string
		Lease          string
		Gateway        string `validate:"omitempty,ip"`
		SendGateway    bool
		ForwardingZone cli.StringSlice `validate:"omitempty,ip"`
		Options        string          `validate:"omitempty,json"`
	}

	LinuxBridge struct {
		BridgeInterface   string
		UpstreamInterface string
		TapInterface      string
	}

	Server struct {
		Mode        string `validate:"oneof=dhcp bridge"`
		CidrAddress string `validate:"cidrv4"`
	}

	Terminator struct {
		ShouldTerminate chan bool
		Terminated      chan bool
	}

	Pipe struct {
		Ctx
		Terminator

		Health
		Logs
		DhcpServer
		LinuxBridge
		Server
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber, ctx *cli.Context) *TaskList[Pipe] {
	return TL.New(p).
		SetCliContext(ctx).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			tl.Pipe.Terminator.ShouldTerminate = make(chan bool, 1)
			tl.Pipe.Terminator.Terminated = make(chan bool, 1)

			return nil
		}).
		SetTasks(
			TL.JobParallel(
				TL.JobSequence(
					Setup(&TL).Job(),

					TL.JobParallel(
						CreatePostroutingRules(&TL).Job(),
						GenerateDhcpServerConfiguration(&TL).Job(),
					),

					TL.JobParallel(
						CreateTapDevice(&TL).Job(),
					),

					TL.JobParallel(
						RunDnsServer(&TL).Job(),
						RunSoftEtherVpnServer(&TL).Job(),
					),

					TL.JobLoop(
						TL.JobParallel(
							HealthCheckPing(&TL).Job(),
						),
					),
				),

				TL.JobBackground(
					TL.JobIf(
						TerminatePredicate(&TL),
						TL.GuardResume(
							TL.JobSequence(
								TL.JobParallel(
									TerminateSoftEther(&TL).Job(),
									TerminateDhcpServer(&TL).Job(),
									TerminateInterfaces(&TL).Job(),
								),
								Terminated(&TL).Job(),
							),
							TASK_ANY,
						),
					),
				),
			),
		)
}
