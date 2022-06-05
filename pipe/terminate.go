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

func Terminated(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminated").
		Set(func(t *Task[Pipe]) error {
			t.Log.Infoln("Graceful termination finished.")

			t.Pipe.Terminator.Terminated <- true

			close(t.Pipe.Terminator.ShouldTerminate)
			close(t.Pipe.Terminator.Terminated)

			return nil
		})
}

func TerminateSoftEther(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate-softether").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("softether-vpnsrv", "stop").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			t.Channel.Err <- t.RunCommandJobAsJobSequenceWithExtension(func(job Job) Job {
				return tl.GuardResume(job, TASK_ANY)
			})

			return nil
		})
}

func TerminateDhcpServer(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate-dnsmasq").
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
			t.Channel.Err <- t.RunCommandJobAsJobSequenceWithExtension(func(job Job) Job {
				return tl.GuardResume(job, TASK_ANY)
			})

			return nil
		})
}

func TerminateTapInterface(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate-interface-tap").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_DHCP
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"ifconfig",
				t.Pipe.DhcpServer.TapInterface,
				"down",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"ip",
				"link",
				"delete",
				t.Pipe.DhcpServer.TapInterface,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"ip",
				"tuntap",
				"del",
				"dev",
				t.Pipe.DhcpServer.TapInterface,
				"mode",
				"tap",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			t.Channel.Err <- t.RunCommandJobAsJobSequenceWithExtension(func(job Job) Job {
				return tl.GuardResume(job, TASK_ANY)
			})

			return nil
		})
}
