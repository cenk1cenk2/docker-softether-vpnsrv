package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func TerminateSoftEther(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate", "softether").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"softether-vpnsrv",
				"stop",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJob(func(t *Task[Pipe]) Job {
				return tl.GuardAlways(t.GetCommandJobAsJobSequence())
			})
		})
}

func TerminateDhcpServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate", "dnsmasq").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_DHCP
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"pkill",
				"dnsmasq",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJob(func(t *Task[Pipe]) Job {
				return tl.GuardAlways(t.GetCommandJobAsJobSequence())
			})
		})
}

func TerminateTapInterface(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate", "interface", "tap").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"ifconfig",
				t.Pipe.SoftEther.TapInterface,
				"down",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"ip",
				"link",
				"delete",
				t.Pipe.SoftEther.TapInterface,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"ip",
				"tuntap",
				"del",
				"dev",
				t.Pipe.SoftEther.TapInterface,
				"mode",
				"tap",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJob(func(t *Task[Pipe]) Job {
				return tl.GuardAlways(t.GetCommandJobAsJobSequence())
			})
		})
}

func TerminateBridgeInterface(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate", "interface", "bridge").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_BRIDGE
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"ifconfig",
				t.Pipe.LinuxBridge.BridgeInterface,
				"down",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"brctl",
				"delbr",
				t.Pipe.LinuxBridge.BridgeInterface,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJob(func(t *Task[Pipe]) Job {
				return tl.GuardAlways(t.GetCommandJobAsJobSequence())
			})
		})
}
