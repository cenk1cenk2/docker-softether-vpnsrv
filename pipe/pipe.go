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
		DhcpServer
		LinuxBridge
		SoftEther
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

					TL.JobSequence(
						HealthCheckSetup(&TL).Job(),
						TL.JobParallel(
							TL.JobBackground(TL.JobLoop(HealthCheckPing(&TL).Job())),
							TL.JobBackground(TL.JobLoop(HealthCheckSoftEther(&TL).Job())),
							TL.JobBackground(TL.JobLoop(HealthCheckDhcpServer(&TL).Job())),
						),
					),

					KeepAlive(&TL).Job(),
				),

				// terminate handler
				TL.JobBackground(
					TL.JobIf(
						TerminatePredicate(&TL),
						TL.JobSequence(
							TL.GuardIgnorePanic(
								TL.JobSequence(
									TL.GuardResume(
										TerminateSoftEther(&TL).Job(),
										TASK_ANY,
									),
									TL.GuardResume(
										TerminateDhcpServer(&TL).Job(),
										TASK_ANY,
									),
									TL.GuardResume(
										TerminateTapInterface(&TL).Job(),
										TASK_ANY,
									),
									TL.GuardResume(
										TerminateBridgeInterface(&TL).Job(),
										TASK_ANY,
									),
								),
							),
							TL.GuardResume(Terminated(&TL).Job(), TASK_ANY),
						),
					),
				),
			),
		)
}
