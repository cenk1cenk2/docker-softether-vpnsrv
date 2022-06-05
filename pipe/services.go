package pipe

import (
	"time"

	"github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func RunDnsServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("dnsmasq").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_DHCP
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("dnsmasq").
				SetLogLevel(logrus.DebugLevel, 0, logrus.DebugLevel).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			err := t.RunCommandJobAsJobSequence()

			t.Log.Infoln("Started DNSMASQ DHCP server.")

			return err
		})
}

func RunSoftEtherVpnServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("softether").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("softether-vpnsrv", "start").
				SetLogLevel(logrus.DebugLevel, 0, logrus.DebugLevel).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			err := t.RunCommandJobAsJobSequence()

			t.Log.Infoln("Started SoftEtherVPN server.")

			return err
		})
}

func HealthCheckPing(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health-ping").
		Set(func(t *Task[Pipe]) error {
			pinger, err := ping.NewPinger(t.Pipe.Health.DhcpServerAddress)
			pinger.Count = 3

			if err != nil {
				return err
			}

			if err := pinger.Run(); err != nil {
				return err
			}

			stats := pinger.Statistics()

			t.Log.Debugf("Ping health check to %s in avg %s.", stats.IPAddr.String(), stats.AvgRtt)

			t.Log.Debugf("Next health check in: %s", t.Pipe.Ctx.Health.Duration.String())
			time.Sleep(t.Pipe.Ctx.Health.Duration)

			return nil
		})
}
