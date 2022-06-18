package pipe

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
	"github.com/mitchellh/go-ps"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func HealthCheckSetup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health:init").
		Set(func(t *Task[Pipe]) error {
			processes, err := ps.Processes()

			if err != nil {
				return err
			}

			for _, process := range processes {
				switch process.Executable() {
				case "vpnserver":
					t.Pipe.Ctx.Health.SoftEtherPID = process.Pid()
					t.Log.Debugf("SoftEtherVPN server PID set: %d", t.Pipe.Ctx.Health.SoftEtherPID)
				case "dnsmasq":
					t.Pipe.Ctx.Health.DnsMasqPID = process.Pid()
					t.Log.Debugf("DNSMASQ server PID set: %d", t.Pipe.Ctx.Health.DnsMasqPID)
				}
			}

			return nil
		})
}

func HealthCheckPing(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health:ping").
		Set(func(t *Task[Pipe]) error {
			pinger, err := ping.NewPinger(t.Pipe.Health.DhcpServerAddress)
			pinger.Count = 3
			pinger.Timeout = time.Second * 10

			if err != nil {
				t.Channel.Fatal <- err
			}

			if err := pinger.Run(); err != nil {
				t.Channel.Fatal <- err
			}

			stats := pinger.Statistics()

			if stats.PacketLoss == 100 {
				return fmt.Errorf("Can not reach the upstream DHCP server.")
			}

			t.Log.Debugf("Ping health check to %s in avg %s.", stats.IPAddr.String(), stats.AvgRtt)

			t.Log.Debugf("Next ping health check in: %s", t.Pipe.Ctx.Health.Duration.String())
			time.Sleep(t.Pipe.Ctx.Health.Duration)

			return nil
		})
}

func HealthCheckSoftEther(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health:softether").
		Set(func(t *Task[Pipe]) error {
			process, err := ps.FindProcess(t.Pipe.Ctx.Health.SoftEtherPID)

			if err != nil {
				t.Log.Debugln(err)
			}

			if process == nil {
				t.Channel.Fatal <- fmt.Errorf("SoftEther process is not alive.")
			}

			t.Log.Debugf(
				"Next SoftEther process health check in: %s",
				t.Pipe.Ctx.Health.Duration.String(),
			)

			time.Sleep(t.Pipe.Ctx.Health.Duration)

			return nil
		})
}

func HealthCheckDhcpServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health:dnsmasq").
		Set(func(t *Task[Pipe]) error {
			process, err := ps.FindProcess(t.Pipe.Ctx.Health.DnsMasqPID)

			if err != nil {
				t.Log.Debugln(err)
			}

			if process == nil {
				t.Channel.Fatal <- fmt.Errorf("SoftEther process is not alive.")
			}

			t.Log.Debugf(
				"Next DNSMASQ process health check in: %s",
				t.Pipe.Ctx.Health.Duration.String(),
			)

			time.Sleep(t.Pipe.Ctx.Health.Duration)

			return nil
		})
}
