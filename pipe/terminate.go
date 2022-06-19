package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func TerminatePredicate(tl *TaskList[Pipe]) JobPredicate {
	return tl.Predicate(func(tl *TaskList[Pipe]) bool {
		tl.Log.Debugln("Registered terminate listener.")

		a := <-tl.Pipe.Terminator.ShouldTerminate

		tl.Log.Warnln("Running termination tasks...")

		return a
	})
}

func Terminate(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate").
		SetJobWrapper(func(job Job) Job {
			return TL.GuardAlways(job)
		}).
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
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunSubtasks(); err != nil {
				return err
			}

			t.Log.Infoln("Graceful termination finished.")

			t.Pipe.Terminator.Terminated <- true

			close(t.Pipe.Terminator.ShouldTerminate)
			close(t.Pipe.Terminator.Terminated)

			return nil
		})
}

func TerminateSoftEther(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate:softether").
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
	return tl.CreateTask("terminate:dnsmasq").
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
	return tl.CreateTask("terminate:interface:tap").
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
	return tl.CreateTask("terminate:interface:bridge").
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
