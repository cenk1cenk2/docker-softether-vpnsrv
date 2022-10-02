package pipe

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
	"github.com/mitchellh/go-ps"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func HealthCheck(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "parent").
		SetJobWrapper(func(job Job) Job {
			return TL.JobSequence(
				job,
				TL.JobParallel(
					HealthCheckPing(tl).Job(),
					HealthCheckSoftEther(tl).Job(),
					HealthCheckDhcpServer(tl).Job(),
				),
			)
		}).
		Set(func(t *Task[Pipe]) error {
			processes, err := ps.Processes()

			if err != nil {
				return err
			}

			for _, process := range processes {
				switch process.Executable() {
				case "vpnserver":
					t.Pipe.Ctx.Health.SoftEtherPIDs = append(t.Pipe.Ctx.Health.SoftEtherPIDs, process.Pid())
				case "dnsmasq":
					t.Pipe.Ctx.Health.DnsMasqPIDs = append(t.Pipe.Ctx.Health.DnsMasqPIDs, process.Pid())
				}
			}

			t.Log.Debugf("SoftEtherVPN server PIDs set: %s", t.Pipe.Ctx.Health.SoftEtherPIDs)
			t.Log.Debugf("DNSMASQ server PIDs set: %s", t.Pipe.Ctx.Health.DnsMasqPIDs)

			return nil
		})
}

func HealthCheckPing(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "ping").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Health.EnablePing
		}).
		SetJobWrapper(func(job Job) Job {
			return tl.JobBackground(tl.JobLoopWithWaitAfter(job, tl.Pipe.Ctx.Health.Duration))
		}).
		SetMarks(MARK_ROUTINE).
		Set(func(t *Task[Pipe]) error {
			pinger, err := ping.NewPinger(t.Pipe.Health.DhcpServerAddress)
			pinger.Count = 3
			pinger.Timeout = time.Second * 10

			if err != nil {
				return err
			}

			if err := pinger.Run(); err != nil {
				return err
			}

			stats := pinger.Statistics()

			if stats.PacketLoss == 100 {
				return fmt.Errorf(
					"Can not ping the upstream DHCP server: %s",
					t.Pipe.Health.DhcpServerAddress,
				)
			}

			t.Log.Debugf("Ping health check to %s in avg %s.", stats.IPAddr.String(), stats.AvgRtt)

			t.Log.Debugf("Next ping health check in: %s", t.Pipe.Ctx.Health.Duration.String())

			return nil
		})
}

func HealthCheckSoftEther(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "softether").
		SetJobWrapper(func(job Job) Job {
			return tl.JobBackground(tl.JobLoopWithWaitAfter(job, tl.Pipe.Ctx.Health.Duration))
		}).
		SetMarks(MARK_ROUTINE).
		Set(func(t *Task[Pipe]) error {
			for _, pid := range t.Pipe.Ctx.Health.SoftEtherPIDs {
				process, err := ps.FindProcess(pid)

				if err != nil {
					t.Log.Debugln(err)
				}

				if process == nil {
					return fmt.Errorf("SoftEther process is not alive: PID: %d", pid)
				}
			}

			t.Log.Debugf(
				"Next SoftEther process health check in: %s",
				t.Pipe.Ctx.Health.Duration.String(),
			)

			return nil
		})
}

func HealthCheckDhcpServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "dnsmasq").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_DHCP
		}).
		SetJobWrapper(func(job Job) Job {
			return tl.JobBackground(tl.JobLoopWithWaitAfter(job, tl.Pipe.Ctx.Health.Duration))
		}).
		SetMarks(MARK_ROUTINE).
		Set(func(t *Task[Pipe]) error {
			for _, pid := range t.Pipe.Ctx.Health.DnsMasqPIDs {
				process, err := ps.FindProcess(pid)

				if err != nil {
					t.Log.Debugln(err)
				}

				if process == nil {
					return fmt.Errorf("DNSMASQ process is not alive.")
				}
			}

			t.Log.Debugf(
				"Next DNSMASQ process health check in: %s",
				t.Pipe.Ctx.Health.Duration.String(),
			)

			return nil
		})
}
