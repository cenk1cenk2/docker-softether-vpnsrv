package pipe

import (
	"fmt"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

// TODO: idk why but this can not be adapted to plumber v4, maybe try again later

func Terminate(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate").
		Set(func(t *Task[Pipe]) error {
			t.SetSubtask(
				tl.JobParallel(
					TerminateSoftEther(tl).Job(),
					TerminateDhcpServer(tl).Job(),
					TerminateTapInterface(tl).Job(),
					TerminateBridgeInterface(tl).Job(),
				),
			)

			return nil
		}).
		EnableTerminator().
		SetOnTerminator(func(t *Task[Pipe]) error {
			tl.Log.Warnln("Running termination tasks...")

			if err := t.RunSubtasks(); err != nil {
				return err
			}

			t.Control.Cancel(fmt.Errorf("Trying to terminate..."))

			t.Log.Infoln("Graceful termination finished.")

			return nil
		})
}

func TerminateSoftEther(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate", "softether").
		SetJobWrapper(func(job Job) Job {
			return tl.GuardAlways(job)
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("softether-vpnsrv", "stop").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func TerminateDhcpServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate", "dnsmasq").
		SetJobWrapper(func(job Job) Job {
			return tl.GuardAlways(job)
		}).
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
			return t.RunCommandJobAsJobSequence()
		})
}

func TerminateTapInterface(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate", "interface", "tap").
		SetJobWrapper(func(job Job) Job {
			return tl.GuardAlways(job)
		}).
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
			return t.RunCommandJobAsJobSequence()
		})
}

func TerminateBridgeInterface(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate", "interface", "bridge").
		SetJobWrapper(func(job Job) Job {
			return tl.GuardAlways(job)
		}).
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
			return t.RunCommandJobAsJobSequence()
		})
}
