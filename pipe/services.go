package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func RunDnsServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("dnsmasq").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_DHCP
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("dnsmasq").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				return err
			}

			t.Log.Infoln("Started DNSMASQ DHCP server.")

			return nil
		})
}

func RunSoftEtherVpnServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("softether").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("softether-vpnsrv", "start").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				return err
			}

			t.Log.Infoln("Started SoftEtherVPN server.")

			return nil
		})
}

func KeepAlive(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("keep-alive").
		Set(func(t *Task[Pipe]) error {
			<-t.Plumber.Terminator.Terminated

			return nil
		})
}
